package private_data

import (
	"context"
	"fmt"
	"time"

	"github.com/GophKeeper/internal/repository"
	"github.com/GophKeeper/internal/service"
	"github.com/GophKeeper/internal/service/private_data/dto"
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

	err := s.privateDataRepo.CreatePrivateData(ctx, &repository.PrivateData{
		Data:      req.Data,
		Type:      req.Type,
		CreatedAt: service.DatePtr(timeNow),
		UpdatedAt: service.DatePtr(timeNow),
		UserID:    req.UserID,
	})
	if err != nil {
		return fmt.Errorf("save private data error: %v", err)
	}
	return nil
}
