package domain

import (
	"time"

	"github.com/google/uuid"
)

type DevicePlatform string

const (
	DevicePlatformIOS     DevicePlatform = "ios"
	DevicePlatformAndroid DevicePlatform = "android"
)

type DeviceToken struct {
	ID        uuid.UUID      `json:"id"`
	UserID    uuid.UUID      `json:"user_id"`
	Token     string         `json:"token"`
	Platform  DevicePlatform `json:"platform"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type RegisterDeviceTokenInput struct {
	UserID   uuid.UUID
	Token    string
	Platform DevicePlatform
}

func NewDeviceToken(input RegisterDeviceTokenInput) *DeviceToken {
	now := time.Now().UTC()
	return &DeviceToken{
		ID:        uuid.New(),
		UserID:    input.UserID,
		Token:     input.Token,
		Platform:  input.Platform,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
