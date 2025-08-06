package kafka

import (
	"context"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/segmentio/kafka-go"

	"github.com/GophKeeper/internal/repository"
	"github.com/GophKeeper/internal/service"
	"github.com/GophKeeper/internal/storage/sqlite"
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
		errSavePrivateData := c.sqliteService.SavePrivateData(&sqlite.PrivateData{
			ID:        req.Id,
			Type:      FromProto(req.Type),
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
