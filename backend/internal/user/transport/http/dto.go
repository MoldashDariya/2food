package http

import "gitlab.com/2foood/backend/internal/domain"

type UserResponse struct {
	User *domain.User `json:"user"`
}

type UpdateUserRequest struct {
	Name                 *string                      `json:"name,omitempty"`
	AvatarURL            *string                      `json:"avatar_url,omitempty" validate:"omitempty,url"`
	CityID               *int                         `json:"city_id,omitempty"`
	Language             *domain.Language             `json:"language,omitempty" validate:"omitempty,oneof=kk en ru uk zh"`
	NotificationSettings *domain.NotificationSettings `json:"notification_settings,omitempty"`
	DietaryPreferences   []domain.DietaryTag          `json:"dietary_preferences,omitempty" validate:"omitempty,dive,oneof=vegetarian vegan gluten_free halal kosher"`
}

type UpdatePhoneRequest struct {
	Phone string `json:"phone" validate:"required,e164"`
	Code  string `json:"code" validate:"required,len=6"`
}

type UpdateEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required,len=6"`
}

type RegisterDeviceTokenRequest struct {
	Token    string                `json:"token" validate:"required"`
	Platform domain.DevicePlatform `json:"platform" validate:"required,oneof=ios android"`
}

type DeleteDeviceTokenRequest struct {
	Token string `json:"token" validate:"required"`
}
