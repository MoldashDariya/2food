package http

import (
	"gitlab.com/2foood/backend/internal/domain"
)

const (
	StatusSuccess   = "success"
	StatusNeedPhone = "need_phone"
)

type PhoneSendOTPRequest struct {
	Phone string `json:"phone" validate:"required,e164"`
}

type EmailSendOTPRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type PhoneLoginRequest struct {
	Phone string `json:"phone" validate:"required,e164"`
	Code  string `json:"code" validate:"required,len=6"`
}

type GoogleLoginRequest struct {
	IDToken string `json:"id_token" validate:"required"`
}

type AuthResponse struct {
	Status string            `json:"status" example:"success"`
	User   *domain.User      `json:"user"`
	Tokens *domain.TokenPair `json:"tokens"`
}

type OAuthNeedPhoneResponse struct {
	Status    string `json:"status" example:"need_phone"`
	SessionID string `json:"session_id" example:"u-u-i-d"`
}

type CompleteGoogleLoginRequest struct {
	SessionID string `json:"session_id" validate:"required"`
	Phone     string `json:"phone" validate:"required,e164"`
	Code      string `json:"code" validate:"required,len=6"`
}

type LinkGoogleRequest struct {
	IDToken string `json:"id_token" validate:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RefreshResponse struct {
	Tokens *domain.TokenPair `json:"tokens"`
}
