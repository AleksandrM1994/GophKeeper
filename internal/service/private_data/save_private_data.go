package private_data

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/GophKeeper/internal/kafka"
	"github.com/GophKeeper/internal/repository"
	"github.com/GophKeeper/internal/service"
	"github.com/GophKeeper/internal/service/private_data/dto"
	api "github.com/GophKeeper/pkg/api"
)

func (s *PrivateDataServiceImpl) SavePrivateData(ctx context.Context, req *dto.SavePrivateDataRequest) error {
	errValidate := req.Validate()
	if errValidate != nil {
		return fmt.Errorf("validate error: %v", errValidate)
	}

	mosLoc, errLoadLocation := time.LoadLocation("Europe/Moscow")
	if errLoadLocation != nil {
		return fmt.Errorf("time.LoadLocation:%w", errLoadLocation)
	}
	timeNowString := time.Now().In(mosLoc).Format(time.RFC3339)
	timeNow, errTimeParse := time.Parse(time.RFC3339, timeNowString)
	if errTimeParse != nil {
		return fmt.Errorf("time.Parse:%w", errTimeParse)
	}

	id := uuid.New().String()
	err := s.privateDataRepo.CreatePrivateData(ctx, &repository.PrivateData{
		ID:        id,
		Data:      req.Data,
		Type:      req.Type,
		CreatedAt: service.DatePtr(timeNow),
		UpdatedAt: service.DatePtr(timeNow),
		UserID:    req.UserID,
		Nonce:     req.Nonce,
	})
	if err != nil {
		return fmt.Errorf("save private data error: %v", err)
	}

	data := &api.PrivateDataSaved{
		Id:        id,
		Type:      ToProto(req.Type),
		Data:      req.Data,
		CreatedAt: timestamppb.New(timeNow),
		UpdatedAt: timestamppb.New(timeNow),
		Nonce:     req.Nonce,
		Login:     req.Login,
	}

	dataBytes, errMarshal := json.Marshal(data)
	if errMarshal != nil {
		return fmt.Errorf("marshal data error: %v", errMarshal)
	}

	errSendMessage := s.kafkaService.SendMessage(
		ctx,
		s.cfg.KafkaHost,
		kafka.GophKeeperPrivateDataSavedTopic,
		dataBytes,
	)
	if errSendMessage != nil {
		return fmt.Errorf("send message error: %v", errSendMessage)
	}
	return nil
}

func ToProto(in repository.PrivateDataType) api.PrivateDataSaved_PrivateDataType {
	switch in {
	case repository.PrivateDataTypeText:
		return api.PrivateDataSaved_TEXT_TYPE
	case repository.PrivateDataTypeFile:
		return api.PrivateDataSaved_FILE_TYPE
	case repository.PrivateDataTypeAuth:
		return api.PrivateDataSaved_AUTH_TYPE
	case repository.PrivateDataTypeBank:
		return api.PrivateDataSaved_BANK_TYPE
	default:
		return api.PrivateDataSaved_UNKNOWN_TYPE
	}
}
