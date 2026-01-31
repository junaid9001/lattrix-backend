package worker

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptrace"
	"slices"
	"strings"
	"time"

	"github.com/junaid9001/lattrix-backend/internal/consumer"
	"github.com/junaid9001/lattrix-backend/internal/domain/models"
	"github.com/junaid9001/lattrix-backend/internal/http/dto"
	"github.com/junaid9001/lattrix-backend/internal/services"
	"github.com/junaid9001/lattrix-backend/internal/utils"
	"gorm.io/gorm"
)

type Consumer struct {
	jobConsumer     consumer.JobConsumer
	workNotiService *services.WorkspaceNotiService
	db              *gorm.DB
}

func NewConsumer(jobConsumer consumer.JobConsumer, workNotiRepo *services.WorkspaceNotiService, db *gorm.DB) *Consumer {
	return &Consumer{jobConsumer: jobConsumer, workNotiService: workNotiRepo, db: db}
}

func (c *Consumer) HandleCheckJob(data []byte) error {

	var job dto.CheckJob
	if err := json.Unmarshal(data, &job); err != nil {
		return err
	}

	var api models.API

	if err := c.db.First(&api, job.APIID).Error; err != nil {
		return err
	}
	downCount := api.DownCount

	var bodyReader io.Reader
	if strBody, ok := job.Body.(string); ok && strBody != "" {
		bodyReader = strings.NewReader(strBody)
	}
	cleanUrl := strings.TrimSpace(job.URL)
	req, err := http.NewRequest(job.Method, cleanUrl, bodyReader)
	if err != nil {
		return err
	}
	req.Close = true

	req.Header.Set("User-Agent", "Lattrix-Monitor/1.0")
	if api.Method == "POST" {
		req.Header.Set("Content-Type", "application/json")
	}

	//headers

	for k, v := range job.Headers {
		req.Header.Set(k, v)
	}

	if job.AuthType != "NONE" && job.AuthValue != nil {
		val := *job.AuthValue
		key := "Authorization"
		if job.AuthKey != nil {
			key = *job.AuthKey
		}

		if job.AuthIn != nil && *job.AuthIn == "QUERY" {
			q := req.URL.Query()
			q.Set(key, val)
			req.URL.RawQuery = q.Encode()
		} else {
			if job.AuthType == "BEARER" {
				req.Header.Set("Authorization", "Bearer "+val)
			} else {
				req.Header.Set(key, val)
			}
		}
	}

	//latency decomposition

	var dnsStart, dnsDone, connStart, connDone, tlsStart, tlsDone, firstByte time.Time
	trace := &httptrace.ClientTrace{
		DNSStart: func(di httptrace.DNSStartInfo) { dnsStart = time.Now() },
		DNSDone:  func(di httptrace.DNSDoneInfo) { dnsDone = time.Now() },
		ConnectStart: func(network, addr string) {
			if connStart.IsZero() {
				connStart = time.Now()
			}

		},
		ConnectDone:          func(network, addr string, err error) { connDone = time.Now() },
		TLSHandshakeStart:    func() { tlsStart = time.Now() },
		TLSHandshakeDone:     func(cs tls.ConnectionState, err error) { tlsDone = time.Now() },
		GotFirstResponseByte: func() { firstByte = time.Now() },
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	//
	client := &http.Client{
		Timeout: time.Duration(job.TimeoutMs) * time.Millisecond,
	}

	start := time.Now()
	resp, err := client.Do(req)
	latency := int(time.Since(start).Milliseconds())

	var dnsMs, tcpMs, tlsMs, processingMs, transferMs int

	var StatusCode *int

	var failureReasons []string

	isStatusCodeMatch := false
	isBodyMatch := true
	isResponseTimeMatch := true

	var sslExpiry *time.Time
	var sslDaysRemaining *int

	if err != nil {
		failureReasons = append(failureReasons, err.Error())
		isStatusCodeMatch = false
		isBodyMatch = false
	} else {
		defer resp.Body.Close()

		if !dnsStart.IsZero() && !dnsDone.IsZero() {
			dnsMs = int(dnsDone.Sub(dnsStart).Milliseconds())
		}

		if !connStart.IsZero() && !connDone.IsZero() {
			tcpMs = int(connDone.Sub(connStart).Milliseconds())
		}

		if !tlsStart.IsZero() && !tlsDone.IsZero() {
			tlsMs = int(tlsDone.Sub(tlsStart).Milliseconds())
		}

		var connectionFinished time.Time
		if !tlsDone.IsZero() {
			connectionFinished = tlsDone
		} else if !connDone.IsZero() {
			connectionFinished = connDone
		} else {
			connectionFinished = dnsDone
		}

		if !firstByte.IsZero() && !connectionFinished.IsZero() {
			processingMs = int(firstByte.Sub(connectionFinished).Milliseconds())
		}

		code := resp.StatusCode
		StatusCode = &code

		//ssl exp
		if resp.TLS != nil && len(resp.TLS.PeerCertificates) > 0 {
			cert := resp.TLS.PeerCertificates[0]
			exp := cert.NotAfter
			sslExpiry = &exp
			days := int(time.Until(exp).Hours() / 24)
			sslDaysRemaining = &days
		}

		//status code check
		if len(job.ExpectedStatusCodes) > 0 {
			if slices.Contains(job.ExpectedStatusCodes, code) {
				isStatusCodeMatch = true
			}

		} else {
			if code >= 200 && code < 400 {
				isStatusCodeMatch = true
			}
		}

		//latency check
		if job.ExpectedResponseTimeMs != nil {
			if latency > *job.ExpectedResponseTimeMs {
				isResponseTimeMatch = false
				e := fmt.Sprintf("Latency %dms exceeded limit %dms", latency, *job.ExpectedResponseTimeMs)
				failureReasons = append(failureReasons, e)
			}
		}

		if !isStatusCodeMatch {
			e := fmt.Sprintf("Unexpected Status Code: %d", code)
			failureReasons = append(failureReasons, e)
		}

		bodyBytes, _ := io.ReadAll(resp.Body)
		if !firstByte.IsZero() {
			transferMs = int(time.Since(firstByte).Milliseconds())
		}
		// body check
		if job.ExpectedBodyContains != nil && *job.ExpectedBodyContains != "" {

			bodyString := string(bodyBytes)
			if !strings.Contains(bodyString, *job.ExpectedBodyContains) {
				isBodyMatch = false
				e := fmt.Sprintf("Body does not contain: '%s'", *job.ExpectedBodyContains)
				failureReasons = append(failureReasons, e)
			}
		}

	}

	isSuccess := len(failureReasons) == 0
	statusStr := "DOWN"
	if isSuccess {
		statusStr = "UP"
		downCount = 0
	} else {
		downCount += 1
	}

	var finalErrMsg *string
	if len(failureReasons) > 0 {
		joinedMsg := strings.Join(failureReasons, "; ")
		finalErrMsg = &joinedMsg
	}

	if downCount == api.NotifyAfterFailures {

		reason := "Unknown Error"
		if finalErrMsg != nil {
			reason = *finalErrMsg
		}

		alertTitle := fmt.Sprintf("🚨 Alert: %v is DOWN", api.Name)
		alertBody := fmt.Sprintf("Monitor: %s\nURL: %s\nTime: %s\n\nError: %s",
			api.Name, api.URL, time.Now().Format(time.RFC822), reason)

		var result struct {
			CreatedByEmail string
		}

		err := c.db.Model(&models.ApiGroup{}).Select("created_by_email").
			Where("id = ?", api.ApiGroupID).
			First(&result).Error
		if err != nil {
			return err
		}

		go func(to, sub, body string) {
			if err := utils.SendEmail(to, sub, body); err != nil {
				log.Printf("Failed to send email alert: %v", err)
			}
		}(result.CreatedByEmail, alertTitle, alertBody)

		c.workNotiService.Create(api.WorkspaceID, reason, alertTitle)
		downCount = 0
	}

	c.db.Model(&models.API{}).Where("id = ?", job.APIID).Updates(map[string]interface{}{
		"last_checked_at":       time.Now(),
		"next_check_at":         time.Now().Add(time.Duration(api.IntervalSeconds) * time.Second),
		"last_status":           statusStr,
		"last_response_time_ms": int(latency),
		"last_error_message":    finalErrMsg,
		"down_count":            downCount,
	})

	// Log History
	intLatency := int(latency)
	history := &models.ApiCheckResult{
		APIID:               job.APIID,
		ApiName:             job.ApiName,
		CheckedAt:           time.Now(),
		Status:              statusStr,
		StatusCode:          StatusCode,
		ResponseTimeMs:      &intLatency,
		ErrorMessage:        finalErrMsg,
		Success:             isSuccess,
		DnsMs:               dnsMs,
		TcpMs:               tcpMs,
		TlsMs:               tlsMs,
		ProcessingMs:        processingMs,
		TransferMs:          transferMs,
		SslExpiry:           sslExpiry,
		SslDaysRemaining:    sslDaysRemaining,
		IsStatusCodeMatch:   isStatusCodeMatch,
		IsBodyMatch:         isBodyMatch,
		IsResponseTimeMatch: isResponseTimeMatch,
	}
	if err := c.db.Create(history).Error; err != nil {
		return err
	}
	return nil

}

func (c *Consumer) Start(ctx context.Context) error {
	return c.jobConsumer.Consume(ctx, c.HandleCheckJob)
}
