package http

import (
	"time"

	"github.com/google/uuid"
	"gitlab.com/2foood/backend/internal/domain"
)

type OrderResponse struct {
	Order *domain.Order `json:"order"`
}

type OrdersResponse struct {
	Orders []*domain.Order `json:"orders"`
}

type CreateOrderItemDTO struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,min=1"`
}

type CreateOrderRequest struct {
	RestaurantID      uuid.UUID            `json:"restaurant_id" validate:"required"`
	PickupWindowStart time.Time            `json:"pickup_window_start" validate:"required"`
	PickupWindowEnd   time.Time            `json:"pickup_window_end" validate:"required"`
	Items             []CreateOrderItemDTO `json:"items" validate:"required,min=1,dive"`
}

type CreateOrderResponse struct {
	Order        *domain.Order `json:"order"`
	ClientSecret string        `json:"client_secret"`
}

type ConfirmPaymentRequest struct {
	OrderID   uuid.UUID `json:"order_id" validate:"required"`
	PaymentID string    `json:"payment_id" validate:"required"`
}

type GetByPickupCodeRequest struct {
	PickupCode string `json:"pickup_code" validate:"required"`
}

type CancelByRestaurantRequest struct {
	Reason domain.OrderReason `json:"reason" validate:"required,min=10,max=500"`
}

type CancelByUserRequest struct {
	Reason *domain.OrderReason `json:"reason" validate:"omitempty,min=10,max=500"`
}

type StatusHistoryResponse struct {
	History []*domain.OrderStatusHistoryEntry `json:"history"`
}

type PaymentHistoryResponse struct {
	History []*domain.PaymentHistoryEntry `json:"history"`
}

type MarkAsCollectedRequest struct {
	PickupCode string `json:"pickup_code" validate:"omitempty"`
	Force      bool   `json:"force"` // flag for skipping pickup code validation
}
