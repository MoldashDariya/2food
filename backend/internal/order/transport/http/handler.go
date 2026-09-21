package http

import (
	"context"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gitlab.com/2foood/backend/internal/domain"
)

type OrderService interface {
	Create(ctx context.Context, input domain.CreateOrderInput) (*domain.Order, string, error)
	ConfirmPaymentWebhook(ctx context.Context, payload []byte, signature string) error
	GetByID(ctx context.Context, orderID uuid.UUID) (*domain.Order, error)
	GetByUser(ctx context.Context, userID uuid.UUID) ([]*domain.Order, error)
	GetByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]*domain.Order, error)
	GetActiveByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]*domain.Order, error)
	GetByPickupCode(ctx context.Context, restaurantID uuid.UUID, pickupCode string) (*domain.Order, error)
	MarkAsAccepted(ctx context.Context, actor *domain.RestaurantStaff, orderID uuid.UUID) error
	MarkAsReady(ctx context.Context, actor *domain.RestaurantStaff, orderID uuid.UUID) error
	MarkAsCollected(ctx context.Context, actor *domain.RestaurantStaff, orderID uuid.UUID, pickupCode string, force bool) error
	GetStatusHistoryByUser(ctx context.Context, userID uuid.UUID, orderID uuid.UUID) ([]*domain.OrderStatusHistoryEntry, error)
	GetStatusHistoryByRestaurant(ctx context.Context, actor *domain.RestaurantStaff, orderID uuid.UUID) ([]*domain.OrderStatusHistoryEntry, error)
	GetPaymentHistoryByOrder(ctx context.Context, userID uuid.UUID, orderID uuid.UUID) ([]*domain.PaymentHistoryEntry, error)
	GetPaymentHistoryByUser(ctx context.Context, userID uuid.UUID) ([]*domain.PaymentHistoryEntry, error)
	CancelByUser(ctx context.Context, userID uuid.UUID, orderID uuid.UUID, reason *domain.OrderReason) error
	CancelByRestaurant(ctx context.Context, actor *domain.RestaurantStaff, orderID uuid.UUID, reason *domain.OrderReason) error
}

type Handler struct {
	log      *slog.Logger
	service  OrderService
	validate *validator.Validate
}

func NewHandler(log *slog.Logger, s OrderService, v *validator.Validate) *Handler {
	return &Handler{
		log:      log,
		service:  s,
		validate: v,
	}
}
