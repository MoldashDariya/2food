package domain

import "errors"

var (
	ErrNotFound                   = errors.New("resource not found")
	ErrAlreadyExists              = errors.New("resource already exists")
	ErrTooManyAttempts            = errors.New("too many attempts")
	ErrInvalidToken               = errors.New("invalid token")
	ErrExpiredToken               = errors.New("expired token")
	ErrCodeExpired                = errors.New("verification code expired")
	ErrCodeInvalid                = errors.New("verification code invalid")
	ErrInvalidAuthorizationHeader = errors.New("invalid authorization header")
	ErrInvalidJSON                = errors.New("invalid json format")
	ErrInvalidOAuthToken          = errors.New("invalid oauth token")
	ErrInternalServerError        = errors.New("internal server error")
	ErrPermissionDenied           = errors.New("permission denied")
	ErrInvalidParam               = errors.New("invalid param")
	ErrMissingParam               = errors.New("missing param")
	ErrProductUnavailable         = errors.New("product is not available")
	ErrInsufficientStock          = errors.New("insufficient stock")
	ErrInvalidStatusTransition    = errors.New("invalid status transition")
	ErrPickupWindowMismatch       = errors.New("pickup window mismatch")
	ErrIgnoredWebhook             = errors.New("ignored webhook event")
	ErrSessionNotFound            = errors.New("session not found")
	ErrUnsupportedContentType     = errors.New("unsupported content type")
	ErrOrderNotCompleted          = errors.New("order not completed")
	ErrInvalidPickupCode          = errors.New("invalid pickup code")
)
