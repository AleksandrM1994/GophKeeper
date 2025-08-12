package dto

import "github.com/GophKeeper/internal/repository"

type SavePrivateDataRequest struct {
	Data   []byte
	Type   repository.PrivateDataType
	UserID string
	Nonce  []byte
	Login  string
}

func (r *SavePrivateDataRequest) Validate() error {
	return nil
}

type SavePrivateDataResponse struct {
}
