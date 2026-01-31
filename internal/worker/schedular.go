package worker

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/junaid9001/lattrix-backend/internal/domain/repository"
	"github.com/junaid9001/lattrix-backend/internal/http/dto"
	"github.com/junaid9001/lattrix-backend/internal/publisher"
)

type Schedular struct {
	apiRepo   repository.ApiRepository
	publisher publisher.Publisher
}

func NewSchedular(apiRepo repository.ApiRepository, publisher publisher.Publisher) *Schedular {
	return &Schedular{apiRepo: apiRepo, publisher: publisher}
}
func (s *Schedular) Start(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {

		select {

		case <-ctx.Done():
			return

		case <-ticker.C:
			apis, err := s.apiRepo.ListDueForCheck(time.Now())
			if err != nil {
				continue
			}

			for _, api := range apis {

				var headers map[string]string
				var expectedStatusCodes []int

				if len(api.Headers) > 0 {
					_ = json.Unmarshal(api.Headers, &headers)
				}
				if len(api.ExpectedStatusCodes) > 0 {
					_ = json.Unmarshal(api.ExpectedStatusCodes, &expectedStatusCodes)
				}

				var bodyStr string
				if len(api.Body) > 0 {
					bodyStr = string(api.Body)
				}
				job := dto.CheckJob{
					APIID:   api.ID,
					ApiName: api.Name,
					URL:     api.URL,
					Method:  api.Method,

					AuthType:  api.AuthType,
					AuthIn:    api.AuthIn,
					AuthKey:   api.AuthKey,
					AuthValue: api.AuthValue,

					Headers:  headers,
					BodyType: api.BodyType,
					Body:     bodyStr,

					TimeoutMs:              api.TimeoutMs,
					ExpectedStatusCodes:    expectedStatusCodes,
					ExpectedResponseTimeMs: api.ExpectedResponseTimeMs,
					ExpectedBodyContains:   api.ExpectedBodyContains,
				}
				data, err := json.Marshal(job)
				if err != nil {
					log.Printf("failed to marshal job: %v", err)
					continue
				}

				//
				// _, _ = s.apiRepo.Update(api.ID, api.ApiGroupID, map[string]any{
				// 	"last_checked_at": time.Now(),
				// 	"next_check_at":   time.Now().Add(time.Duration(api.IntervalSeconds) * time.Second),
				// }, api.WorkspaceID)

				if err := s.publisher.Publish(ctx, data); err != nil {
					log.Printf("failed to publish job: %v", err)
					continue
				}

			}
		}
	}

}
