package consumer

import "github.com/segmentio/kafka-go"

type ReaderConfig struct {
	Brokers []string
	Topic   string
	GroupID string
}

func NewKafkaConsumer(cfg ReaderConfig) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Brokers,
		Topic:   cfg.Topic,
		GroupID: cfg.GroupID,
	})
}
