package services

import (
	"context"
	"sync"
	"time"

	"github.com/junaid9001/lattrix-backend/internal/config"
	"github.com/junaid9001/lattrix-backend/internal/domain/models"
	"github.com/segmentio/kafka-go"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/event"
	"gorm.io/gorm"
)

type AdminService struct {
	db          *gorm.DB
	kafkaWriter *kafka.Writer
	cfg         *config.Config
}

func NewAdminService(db *gorm.DB, kafkaWriter *kafka.Writer, cfg *config.Config) *AdminService {
	stripe.Key = cfg.STRIPE_SECRET_KEY
	return &AdminService{db: db, kafkaWriter: kafkaWriter, cfg: cfg}
}

type DashboardStats struct {
	TotalUsers        int64   `json:"total_users"`
	ActiveWorkspaces  int64   `json:"active_workspaces"`
	TotalRevenue      float64 `json:"total_revenue"`
	MostSoldPlan      string  `json:"most_sold_plan"`
	MostSoldPlanCount int     `json:"most_sold_plan_count"`
	ProUsers          int64   `json:"pro_users"`
	AgencyUsers       int64   `json:"agency_users"`
}

type StripeActivity struct {
	Type    string `json:"type"`
	ID      string `json:"id"`
	Created string `json:"created"`
}

type KafkaMetrics struct {
	Status      string            `json:"status"`
	BrokerCount int               `json:"broker_count"`
	TopicCount  int               `json:"topic_count"`
	WriterStats WriterPerformance `json:"writer_performance"`
}

type WriterPerformance struct {
	MessagesSent int64   `json:"messages_sent"`
	AvgWriteTime float64 `json:"avg_write_time_ms"`
	Errors       int64   `json:"errors"`
}

type SystemHealth struct {
	Database string       `json:"database"`
	Kafka    KafkaMetrics `json:"kafka"`
}

func (s *AdminService) GetDashboardStats() (*DashboardStats, error) {
	var stats DashboardStats
	var proCount int64
	var agencyCount int64

	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		s.db.Model(&models.User{}).Count(&stats.TotalUsers)
	}()

	go func() {
		defer wg.Done()
		s.db.Model(&models.Workspace{}).Count(&stats.ActiveWorkspaces)
	}()

	go func() {
		defer wg.Done()
		s.db.Model(&models.User{}).Where("plan ILIKE ?", "Pro").Count(&proCount)
	}()

	go func() {
		defer wg.Done()
		s.db.Model(&models.User{}).Where("plan ILIKE ?", "Agency").Count(&agencyCount)
	}()

	wg.Wait()

	// Logic Calculation (CPU bound, practically instant)
	stats.ProUsers = proCount
	stats.AgencyUsers = agencyCount
	stats.TotalRevenue = (float64(proCount) * 15.00) + (float64(agencyCount) * 45.00)

	if agencyCount > proCount {
		stats.MostSoldPlan = "Agency Plan"
		stats.MostSoldPlanCount = int(agencyCount)
	} else if proCount > 0 {
		stats.MostSoldPlan = "Pro Plan"
		stats.MostSoldPlanCount = int(proCount)
	} else {
		stats.MostSoldPlan = "Free Plan"
		stats.MostSoldPlanCount = int(stats.TotalUsers - proCount - agencyCount)
	}

	return &stats, nil
}

func (s *AdminService) GetRecentStripeActivities() ([]StripeActivity, error) {
	params := &stripe.EventListParams{}
	params.Limit = stripe.Int64(5)

	iter := event.List(params)
	var activities []StripeActivity

	for iter.Next() {
		ev := iter.Event()
		activities = append(activities, StripeActivity{
			Type:    string(ev.Type),
			ID:      ev.ID,
			Created: time.Unix(ev.Created, 0).Format("2006-01-02 15:04"),
		})
	}

	if activities == nil {
		activities = []StripeActivity{}
	}
	return activities, nil
}

func (s *AdminService) GetSystemHealth() SystemHealth {
	health := SystemHealth{
		Database: "Checking...",
		Kafka:    KafkaMetrics{Status: "Checking..."},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	sqlDB, err := s.db.DB()
	if err != nil {
		health.Database = "Error"
	} else {
		if err := sqlDB.PingContext(ctx); err != nil {
			health.Database = "Offline"
		} else {
			health.Database = "Online"
		}
	}

	if s.kafkaWriter != nil {

		stats := s.kafkaWriter.Stats()
		health.Kafka.WriterStats = WriterPerformance{
			MessagesSent: stats.Messages,
			Errors:       stats.Errors,
			AvgWriteTime: float64(stats.WriteTime.Avg.Microseconds()) / 1000.0,
		}

		dialer := &kafka.Dialer{
			Timeout:   500 * time.Millisecond,
			DualStack: true,
		}

		conn, err := dialer.Dial("tcp", s.kafkaWriter.Addr.String())
		if err != nil {
			health.Kafka.Status = "Unreachable"
		} else {
			health.Kafka.Status = "Online"

			brokers, _ := conn.Brokers()
			health.Kafka.BrokerCount = len(brokers)

			partitions, _ := conn.ReadPartitions()
			topics := make(map[string]bool)
			for _, p := range partitions {
				topics[p.Topic] = true
			}
			health.Kafka.TopicCount = len(topics)

			conn.Close()
		}
	} else {
		health.Kafka.Status = "Not Configured"
	}

	return health
}

func (s *AdminService) GetAllUsers(page, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64
	offset := (page - 1) * limit
	s.db.Model(&models.User{}).Count(&total)
	err := s.db.Limit(limit).Offset(offset).Order("created_at desc").Find(&users).Error
	return users, total, err
}

func (s *AdminService) ToggleUserBan(userID int) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	user.IsActive = !user.IsActive
	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
