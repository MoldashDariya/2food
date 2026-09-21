package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type StaffRole string

const (
	StaffRoleManager StaffRole = "manager"
	StaffRoleOwner   StaffRole = "owner"
)

type DayHours struct {
	Open  string `json:"open"`
	Close string `json:"close"`
}

type WorkingHours struct {
	Mon *DayHours `json:"mon,omitempty"`
	Tue *DayHours `json:"tue,omitempty"`
	Wed *DayHours `json:"wed,omitempty"`
	Thu *DayHours `json:"thu,omitempty"`
	Fri *DayHours `json:"fri,omitempty"`
	Sat *DayHours `json:"sat,omitempty"`
	Sun *DayHours `json:"sun,omitempty"`
}

type SocialLinks struct {
	Instagram *string `json:"instagram,omitempty"`
	Telegram  *string `json:"telegram,omitempty"`
}

type Restaurant struct {
	ID           uuid.UUID     `json:"id"`
	Name         string        `json:"name"`
	Description  *string       `json:"description,omitempty"`
	Address      *string       `json:"address,omitempty"`
	Phone        *string       `json:"phone,omitempty"`
	Website      *string       `json:"website,omitempty"`
	ImageURL     *string       `json:"image_url,omitempty"`
	Latitude     *float64      `json:"latitude,omitempty"`
	Longitude    *float64      `json:"longitude,omitempty"`
	WorkingHours *WorkingHours `json:"working_hours,omitempty"`
	SocialLinks  *SocialLinks  `json:"social_links,omitempty"`
	CityID       *int          `json:"city_id,omitempty"`
	City         *City         `json:"city,omitempty"`
	CuisineTypes []CuisineType `json:"cuisine_types,omitempty"`
	RatingAvg    float64       `json:"rating_avg"`
	RatingCount  int           `json:"rating_count"`
	IsActive     bool          `json:"is_active"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	Role         StaffRole     `json:"role,omitempty"`
}

type RestaurantStaff struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	RestaurantID uuid.UUID `json:"restaurant_id"`
	Role         StaffRole `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateRestaurantInput struct {
	Name         string
	Description  *string
	Address      *string
	Phone        *string
	Website      *string
	ImageURL     *string
	Latitude     *float64
	Longitude    *float64
	WorkingHours *WorkingHours
	SocialLinks  *SocialLinks
	CityID       *int
}

type UpdateRestaurantInput struct {
	Name         *string
	Description  *string
	Address      *string
	Phone        *string
	Website      *string
	ImageURL     *string
	Latitude     *float64
	Longitude    *float64
	WorkingHours *WorkingHours
	SocialLinks  *SocialLinks
	CityID       *int
}

func NewRestaurant(input CreateRestaurantInput) *Restaurant {
	now := time.Now().UTC()
	return &Restaurant{
		ID:           uuid.New(),
		Name:         input.Name,
		Description:  input.Description,
		Address:      input.Address,
		Phone:        input.Phone,
		Website:      input.Website,
		ImageURL:     input.ImageURL,
		Latitude:     input.Latitude,
		Longitude:    input.Longitude,
		WorkingHours: input.WorkingHours,
		SocialLinks:  input.SocialLinks,
		CityID:       input.CityID,
		IsActive:     true, // active by default
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func NewRestaurantStaff(restaurantID, userID uuid.UUID, role StaffRole) *RestaurantStaff {
	now := time.Now().UTC()
	if role == "" {
		role = StaffRoleManager // default role is manager
	}
	return &RestaurantStaff{
		ID:           uuid.New(),
		RestaurantID: restaurantID,
		UserID:       userID,
		Role:         role,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func (r *Restaurant) Update(input UpdateRestaurantInput) {
	if input.Name != nil {
		r.Name = *input.Name
	}
	if input.Description != nil {
		r.Description = input.Description
	}
	if input.Address != nil {
		r.Address = input.Address
	}
	if input.Phone != nil {
		r.Phone = input.Phone
	}
	if input.Website != nil {
		r.Website = input.Website
	}
	if input.ImageURL != nil {
		r.ImageURL = input.ImageURL
	}
	if input.Latitude != nil {
		r.Latitude = input.Latitude
	}
	if input.Longitude != nil {
		r.Longitude = input.Longitude
	}
	if input.WorkingHours != nil {
		r.WorkingHours = input.WorkingHours
	}
	if input.SocialLinks != nil {
		r.SocialLinks = input.SocialLinks
	}
	if input.CityID != nil {
		r.CityID = input.CityID
	}

	r.UpdatedAt = time.Now().UTC()
}

// MarshalWorkingHours serializes WorkingHours to JSON for storage
func (r *Restaurant) MarshalWorkingHours() ([]byte, error) {
	if r.WorkingHours == nil {
		return nil, nil
	}

	return json.Marshal(r.WorkingHours)
}

// MarshalSocialLinks serializes SocialLinks to JSON for storage
func (r *Restaurant) MarshalSocialLinks() ([]byte, error) {
	if r.SocialLinks == nil {
		return nil, nil
	}

	return json.Marshal(r.SocialLinks)
}
