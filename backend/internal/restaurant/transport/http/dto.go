package http

import (
	"github.com/google/uuid"
	"gitlab.com/2foood/backend/internal/domain"
)

type RestaurantResponse struct {
	Restaurant *domain.Restaurant `json:"restaurant"`
}

type RestaurantsResponse struct {
	Restaurants []*domain.Restaurant `json:"restaurants"`
}

type StaffResponse struct {
	Staff []*domain.RestaurantStaff `json:"staff"`
}

type StaffMemberResponse struct {
	StaffMember *domain.RestaurantStaff `json:"staff_member"`
}

type CreateRestaurantRequest struct {
	OwnerID      uuid.UUID            `json:"owner_id"    validate:"required"`
	Name         string               `json:"name"        validate:"required,min=2"`
	Description  *string              `json:"description"`
	Address      *string              `json:"address"`
	Phone        *string              `json:"phone"`
	Website      *string              `json:"website"     validate:"omitempty,url"`
	ImageURL     *string              `json:"image_url"   validate:"omitempty,url"`
	Latitude     *float64             `json:"latitude"    validate:"omitempty,min=-90,max=90"`
	Longitude    *float64             `json:"longitude"   validate:"omitempty,min=-180,max=180"`
	WorkingHours *domain.WorkingHours `json:"working_hours"`
	SocialLinks  *domain.SocialLinks  `json:"social_links"`
	CityID       *int                 `json:"city_id"`
}

type UpdateRestaurantRequest struct {
	Name         *string              `json:"name"        validate:"omitempty,min=2"`
	Description  *string              `json:"description"`
	Address      *string              `json:"address"`
	Phone        *string              `json:"phone"`
	Website      *string              `json:"website"     validate:"omitempty,url"`
	ImageURL     *string              `json:"image_url"   validate:"omitempty,url"`
	Latitude     *float64             `json:"latitude"    validate:"omitempty,min=-90,max=90"`
	Longitude    *float64             `json:"longitude"   validate:"omitempty,min=-180,max=180"`
	WorkingHours *domain.WorkingHours `json:"working_hours"`
	SocialLinks  *domain.SocialLinks  `json:"social_links"`
	CityID       *int                 `json:"city_id"`
}

type AddStaffRequest struct {
	UserID uuid.UUID        `json:"user_id" validate:"required"`
	Role   domain.StaffRole `json:"role" validate:"required,oneof=manager"`
}

type UpdateStaffRoleRequest struct {
	Role domain.StaffRole `json:"role" validate:"required,oneof=manager"`
}

type AssignOwnerRequest struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
}

type UpdateCuisinesRequest struct {
	CuisineIDs []int `json:"cuisine_ids" validate:"required"`
}
