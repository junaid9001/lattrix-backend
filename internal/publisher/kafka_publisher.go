package publisher

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Publisher interface {
	Publish(ctx context.Context, data []byte) error
}

type KafkaPublisher struct {
	writer *kafka.Writer
}

func NewkafkaPublisher(writer *kafka.Writer) *KafkaPublisher {
	return &KafkaPublisher{writer: writer}
}

func (p *KafkaPublisher) Publish(ctx context.Context, data []byte) error {
	return p.writer.WriteMessages(ctx, kafka.Message{
		Value: data,
	})
}
