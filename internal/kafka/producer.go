package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

func SendMessage(ctx context.Context, kafkaHost string, topic string, message []byte) error {
	w := &kafka.Writer{
		Addr:     kafka.TCP(kafkaHost),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	err := w.WriteMessages(ctx, kafka.Message{
		Value: message,
	})
	if err != nil {
		return fmt.Errorf("write messages: %w", err)
	}
	return nil
}
