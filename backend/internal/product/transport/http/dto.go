package http

import (
	"gitlab.com/2foood/backend/internal/domain"
)

type ProductResponse struct {
	Product *domain.Product `json:"product"`
}

type ProductsResponse struct {
	Products []*domain.Product `json:"products"`
}

type CategoriesResponse struct {
	Categories []*domain.ProductCategory `json:"categories"`
}

type CreateProductRequest struct {
	CategoryID        *int    `json:"category_id"`
	Name              string  `json:"name" validate:"required,min=2,max=255"`
	Description       *string `json:"description"`
	ImageURL          *string `json:"image_url" validate:"omitempty,url"`
	OriginalPrice     float64 `json:"original_price" validate:"required,gt=0"`
	DiscountedPrice   float64 `json:"discounted_price" validate:"required,gt=0"`
	PickupTimeStart   string  `json:"pickup_time_start" validate:"required,datetime=15:04"`
	PickupTimeEnd     string  `json:"pickup_time_end" validate:"required,datetime=15:04"`
	QuantityAvailable int     `json:"quantity_available" validate:"required,gte=0"`
}

type UpdateProductRequest struct {
	CategoryID      *int     `json:"category_id"`
	Name            *string  `json:"name" validate:"omitempty,min=2,max=255"`
	Description     *string  `json:"description"`
	ImageURL        *string  `json:"image_url" validate:"omitempty,url"`
	OriginalPrice   *float64 `json:"original_price" validate:"omitempty,gt=0"`
	DiscountedPrice *float64 `json:"discounted_price" validate:"omitempty,gt=0"`
	PickupTimeStart *string  `json:"pickup_time_start" validate:"omitempty,datetime=15:04"`
	PickupTimeEnd   *string  `json:"pickup_time_end" validate:"omitempty,datetime=15:04"`
}

type UpdateStatusRequest struct {
	IsActive *bool `json:"is_active" validate:"required,boolean"`
}

type UpdateStockRequest struct {
	Quantity *int `json:"quantity" validate:"required,gte=0"`
}
