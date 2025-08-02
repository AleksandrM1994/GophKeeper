package kafka

import (
	"context"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/segmentio/kafka-go"

	"github.com/GophKeeper/internal/repository"
	"github.com/GophKeeper/internal/service"
	"github.com/GophKeeper/internal/storage/bbolt"
	api "github.com/GophKeeper/pkg/api"
)

func (c *Controller) ReadMessage(ctx context.Context, topic string) error {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{c.kafkaHost},
		Topic:     topic,
		Partition: 0,
		MinBytes:  10e3,
		MaxBytes:  10e6,
	})

	msg, errReadMessage := reader.ReadMessage(ctx)
	if errReadMessage != nil {
		return fmt.Errorf("read message: %w", errReadMessage)
	}

	c.lg.Infow("Получено сообщение из топика", topic, string(msg.Value))

	var req *api.PrivateDataSaved
	switch topic {
	case GophKeeperPrivateDataSavedTopic:
		errUnmarshal := json.Unmarshal(msg.Value, &req)
		if errUnmarshal != nil {
			return fmt.Errorf("unmarshal data: %w", errUnmarshal)
		}

		createdAt := req.CreatedAt.AsTime()
		updatedAt := req.UpdatedAt.AsTime()
		errSavePrivateData := c.bboltService.SavePrivateData(req.Login, &bbolt.PrivateData{
			ID:        req.Id,
			Type:      repository.PrivateDataType(req.Type),
			Data:      req.Data,
			CreatedAt: service.DatePtr(createdAt),
			UpdatedAt: service.DatePtr(updatedAt),
		})
		if errSavePrivateData != nil {
			return fmt.Errorf("save private data: %w", errSavePrivateData)
		}
	}

	errClose := reader.Close()
	if errClose != nil {
		return fmt.Errorf("close reader: %w", errClose)
	}

	return nil
}
