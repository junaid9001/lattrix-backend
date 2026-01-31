package publisher

import "github.com/segmentio/kafka-go"

type Kafkaconfig struct {
	Brokers []string
	Topic   string
}

func NewKafkaWriter(cfg Kafkaconfig) *kafka.Writer {
	return &kafka.Writer{
		Addr:  kafka.TCP(cfg.Brokers...),
		Topic: cfg.Topic,
	}
}
