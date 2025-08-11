package kafka

import (
	"context"
	"time"

	"github.com/goccy/go-json"
	"github.com/segmentio/kafka-go"

	"github.com/GophKeeper/internal/repository"
	"github.com/GophKeeper/internal/service"
	"github.com/GophKeeper/internal/storage/sqlite"
	api "github.com/GophKeeper/pkg/api"
)

func (c *KafkaServiceImpl) ReadMessage(ctx context.Context, topic string) error {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:         []string{c.kafkaHost},
		Topic:           topic,
		GroupID:         "gophkeeper-cli-group",
		Partition:       0,
		MinBytes:        10e3,
		MaxBytes:        10e6,
		ReadLagInterval: time.Second * 5, // Интервал обновления информации о лаге
		MaxWait:         time.Second * 5,
	})

	defer func() {
		if err := reader.Close(); err != nil {
			c.lg.Errorf("failed to close Kafka reader: %v", err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			c.lg.Infof("Consumer context cancelled, exiting for topic: %s", topic)
			return nil
		default:
		}

		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			c.lg.Errorf("failed to read message from topic %s: %v", topic, err)
			continue
		}

		c.lg.Infow("Получено сообщение из топика", "topic", topic, "message", string(msg.Value))

		var req *api.PrivateDataSaved
		err = json.Unmarshal(msg.Value, &req)
		if err != nil {
			c.lg.Errorf("failed to unmarshal message from topic %s: %v", topic, err)
			continue
		}

		createdAt := req.CreatedAt.AsTime()
		updatedAt := req.UpdatedAt.AsTime()

		err = c.sqliteService.SavePrivateData(ctx, &sqlite.PrivateData{
			ID:        req.Id,
			Type:      FromProto(req.Type),
			Data:      req.Data,
			CreatedAt: service.DatePtr(createdAt),
			UpdatedAt: service.DatePtr(updatedAt),
			Nonce:     req.Nonce,
			UserLogin: req.Login,
		})
		if err != nil {
			c.lg.Errorf("failed to save private data from topic %s: %v", topic, err)
		}

		// Фиксируем offset
		if err := reader.CommitMessages(ctx, msg); err != nil {
			c.lg.Errorf("failed to commit message offset in topic %s: %v", topic, err)
		}
	}
}

func FromProto(in api.PrivateDataSaved_PrivateDataType) repository.PrivateDataType {
	switch in {
	case api.PrivateDataSaved_TEXT_TYPE:
		return repository.PrivateDataTypeText
	case api.PrivateDataSaved_FILE_TYPE:
		return repository.PrivateDataTypeFile
	case api.PrivateDataSaved_AUTH_TYPE:
		return repository.PrivateDataTypeAuth
	case api.PrivateDataSaved_BANK_TYPE:
		return repository.PrivateDataTypeBank
	default:
		return repository.PrivateDataTypeUnknown
	}
}
