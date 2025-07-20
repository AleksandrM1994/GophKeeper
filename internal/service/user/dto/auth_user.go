package dto

import (
	"fmt"

	custom_errs "github.com/GophKeeper/internal/errors"
)

type AuthUserRequest struct {
	Login    string
	Password string
}

type AuthUserResponse struct {
	JWT    string
	UserID string
}

func (r *AuthUserRequest) Validate() error {
	switch {
	case r.Login == "":
		return fmt.Errorf("invalid login:%w", custom_errs.ErrValidate)
	case r.Password == "":
		return fmt.Errorf("invalid password:%w", custom_errs.ErrValidate)
	}
	return nil
}
