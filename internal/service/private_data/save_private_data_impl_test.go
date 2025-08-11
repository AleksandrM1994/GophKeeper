package private_data

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/GophKeeper/config"
	"github.com/GophKeeper/internal/kafka"
	"github.com/GophKeeper/internal/mocks"
	"github.com/GophKeeper/internal/repository"
	"github.com/GophKeeper/internal/service/private_data/dto"
)

func TestSavePrivateData_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPrivateDataRepository(ctrl)
	mockKafka := mocks.NewMockKafkaService(ctrl)
	lg := zap.NewExample().Sugar()
	cfg := config.Config{KafkaHost: "kafka:9092"}
	svc := NewService(lg, cfg, mockRepo, mockKafka)

	req := &dto.SavePrivateDataRequest{
		Data:   []byte("secret"),
		Type:   repository.PrivateDataTypeText,
		UserID: "user1",
		Nonce:  []byte("nonce123"),
		Login:  "login1",
	}

	mockRepo.
		EXPECT().
		CreatePrivateData(gomock.Any(), gomock.AssignableToTypeOf(&repository.PrivateData{})).
		DoAndReturn(func(_ context.Context, pd *repository.PrivateData) error {
			assert.Equal(t, req.Data, pd.Data)
			assert.Equal(t, req.Type, pd.Type)
			assert.Equal(t, req.UserID, pd.UserID)
			assert.Equal(t, req.Nonce, pd.Nonce)
			return nil
		})

	mockKafka.
		EXPECT().
		SendMessage(gomock.Any(), gomock.Any(), kafka.GophKeeperPrivateDataSavedTopic, gomock.Any()).
		Return(nil)

	err := svc.SavePrivateData(context.Background(), req)

	assert.NoError(t, err)
}

func TestSavePrivateData_KafkaError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPrivateDataRepository(ctrl)
	mockKafka := mocks.NewMockKafkaService(ctrl)
	lg := zap.NewExample().Sugar()
	cfg := config.Config{KafkaHost: "kafka:9092"}
	svc := NewService(lg, cfg, mockRepo, mockKafka)

	req := &dto.SavePrivateDataRequest{
		Data:   []byte("secret"),
		Type:   repository.PrivateDataTypeText,
		UserID: "user1",
		Nonce:  []byte("nonce123"),
		Login:  "login1",
	}

	mockRepo.
		EXPECT().
		CreatePrivateData(gomock.Any(), gomock.Any()).
		Return(nil)

	// Возвращаем ошибку при Send
	kafkaErr := errors.New("kafka down")
	mockKafka.
		EXPECT().
		SendMessage(gomock.Any(), cfg.KafkaHost, kafka.GophKeeperPrivateDataSavedTopic, gomock.Any()).
		Return(kafkaErr)

	err := svc.SavePrivateData(context.Background(), req)
	assert.EqualError(t, err, fmt.Sprintf("send message error: %v", kafkaErr))
}
