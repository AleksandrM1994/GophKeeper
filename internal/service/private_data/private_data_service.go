package private_data

import (
	"context"

	"github.com/GophKeeper/internal/service/private_data/dto"
)

type PrivateDataService interface {
	SavePrivateData(ctx context.Context, req *dto.SavePrivateDataRequest) error
}
