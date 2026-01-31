package consumer

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type JobConsumer interface {
	Consume(ctx context.Context, handler func([]byte) error) error
}

type KafkaJobConsumer struct {
	Reader *kafka.Reader
}

func NewKafkaJobConsumer(reader *kafka.Reader) JobConsumer {
	return &KafkaJobConsumer{Reader: reader}
}

func (c *KafkaJobConsumer) Consume(ctx context.Context, handler func([]byte) error) error {
	for {
		msg, err := c.Reader.ReadMessage(ctx)
		if err != nil {
			return err
		}
		if err := handler(msg.Value); err != nil {

		}
	}
}
