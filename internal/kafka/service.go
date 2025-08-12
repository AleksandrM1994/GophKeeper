package kafka

import "context"

type KafkaService interface {
	SendMessage(ctx context.Context, host, topic string, payload []byte) error
	ReadMessage(ctx context.Context, topic string) error
}
