package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	UserRoleCustomer UserRole = "customer"
	UserRolePartner  UserRole = "partner"
	UserRoleAdmin    UserRole = "admin"
)

type NotificationSettings struct {
	OrderUpdates bool `json:"order_updates"`
	Promotions   bool `json:"promotions"`
}

type User struct {
	ID                   uuid.UUID            `json:"id"`
	Role                 UserRole             `json:"role"`
	Phone                *string              `json:"phone,omitempty"`
	Email                *string              `json:"email,omitempty"`
	GoogleID             *string              `json:"-"`
	Name                 *string              `json:"name,omitempty"`
	AvatarURL            *string              `json:"avatar_url,omitempty"`
	CityID               *int                 `json:"city_id,omitempty"`
	Language             Language             `json:"language"`
	NotificationSettings NotificationSettings `json:"notification_settings"`
	DietaryPreferences   []DietaryTag         `json:"dietary_preferences"`
	CreatedAt            time.Time            `json:"created_at"`
	UpdatedAt            time.Time            `json:"updated_at"`
}

type UpdateUserInput struct {
	Name                 *string
	Email                *string
	AvatarURL            *string
	CityID               *int
	Language             *Language
	NotificationSettings *NotificationSettings
	DietaryPreferences   []DietaryTag
}

func NewUserWithPhone(phone string, lang Language) *User {
	now := time.Now().UTC()

	if lang == "" {
		lang = LanguageEnglish // default to English
	}

	return &User{
		ID:       uuid.New(),
		Role:     UserRoleCustomer,
		Phone:    &phone,
		Language: lang,
		NotificationSettings: NotificationSettings{
			OrderUpdates: true,
			Promotions:   true,
		},
		DietaryPreferences: []DietaryTag{},
		CreatedAt:          now,
		UpdatedAt:          now,
	}
}

func (u *User) BindGoogleData(info OAuthInfo) {
	u.GoogleID = &info.OAuthID

	if info.EmailVerified && info.Email != "" {
		u.Email = &info.Email
	}
	if info.Name != "" {
		u.Name = &info.Name
	}
	if info.AvatarURL != "" {
		u.AvatarURL = &info.AvatarURL
	}
}

func (u *User) Update(input UpdateUserInput) {
	if input.Name != nil {
		u.Name = input.Name
	}
	if input.AvatarURL != nil {
		u.AvatarURL = input.AvatarURL
	}
	if input.CityID != nil {
		u.CityID = input.CityID
	}
	if input.Language != nil {
		u.Language = *input.Language
	}
	if input.NotificationSettings != nil {
		u.NotificationSettings = *input.NotificationSettings
	}
	if input.DietaryPreferences != nil {
		u.DietaryPreferences = input.DietaryPreferences
	}

	u.UpdatedAt = time.Now().UTC()
}

func (u *User) IsAdmin() bool {
	return u.Role == UserRoleAdmin
}

func (u *User) IsPartner() bool {
	return u.Role == UserRolePartner
}

func (u *User) IsCustomer() bool {
	return u.Role == UserRoleCustomer
}
