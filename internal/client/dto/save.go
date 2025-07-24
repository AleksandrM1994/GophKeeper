package dto

import "github.com/GophKeeper/internal/repository"

type SavePrivateDataRequest struct {
	Data []byte
	Type repository.PrivateDataType
	JWT  string
}
