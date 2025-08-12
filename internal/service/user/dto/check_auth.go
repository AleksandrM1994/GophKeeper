package dto

type CheckAuthRequest struct {
	JWT string
}

type CheckAuthResponse struct {
	UserID string
}
