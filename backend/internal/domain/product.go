package domain

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID           uuid.UUID        `json:"id"`
	RestaurantID uuid.UUID        `json:"restaurant_id"`
	CategoryID   *int             `json:"category_id,omitempty"`
	Category     *ProductCategory `json:"category,omitempty"`

	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	ImageURL    *string `json:"image_url,omitempty"`

	OriginalPrice   float64 `json:"original_price"`
	DiscountedPrice float64 `json:"discounted_price"`

	PickupTimeStart string `json:"pickup_time_start"`
	PickupTimeEnd   string `json:"pickup_time_end"`

	QuantityAvailable int  `json:"quantity_available"`
	IsActive          bool `json:"is_active"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateProductInput struct {
	RestaurantID      uuid.UUID
	CategoryID        *int
	Name              string
	Description       *string
	ImageURL          *string
	OriginalPrice     float64
	DiscountedPrice   float64
	PickupTimeStart   string
	PickupTimeEnd     string
	QuantityAvailable int
}

type UpdateProductInput struct {
	CategoryID      *int
	Name            *string
	Description     *string
	ImageURL        *string
	OriginalPrice   *float64
	DiscountedPrice *float64
	PickupTimeStart *string
	PickupTimeEnd   *string
}

func NewProduct(input CreateProductInput) *Product {
	now := time.Now().UTC()
	return &Product{
		ID:                uuid.New(),
		RestaurantID:      input.RestaurantID,
		CategoryID:        input.CategoryID,
		Name:              input.Name,
		Description:       input.Description,
		ImageURL:          input.ImageURL,
		OriginalPrice:     input.OriginalPrice,
		DiscountedPrice:   input.DiscountedPrice,
		PickupTimeStart:   input.PickupTimeStart,
		PickupTimeEnd:     input.PickupTimeEnd,
		QuantityAvailable: input.QuantityAvailable,
		IsActive:          true,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}

func (p *Product) Update(input UpdateProductInput) {
	if input.CategoryID != nil {
		p.CategoryID = input.CategoryID
	}
	if input.Name != nil {
		p.Name = *input.Name
	}
	if input.Description != nil {
		p.Description = input.Description
	}
	if input.ImageURL != nil {
		p.ImageURL = input.ImageURL
	}
	if input.OriginalPrice != nil {
		p.OriginalPrice = *input.OriginalPrice
	}
	if input.DiscountedPrice != nil {
		p.DiscountedPrice = *input.DiscountedPrice
	}
	if input.PickupTimeStart != nil {
		p.PickupTimeStart = *input.PickupTimeStart
	}
	if input.PickupTimeEnd != nil {
		p.PickupTimeEnd = *input.PickupTimeEnd
	}

	p.UpdatedAt = time.Now().UTC()
}
