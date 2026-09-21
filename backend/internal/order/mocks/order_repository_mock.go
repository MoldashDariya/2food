package mocks

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"gitlab.com/2foood/backend/internal/domain"
)

type OrderRepositoryMock struct {
	mock.Mock
}

func (m *OrderRepositoryMock) Create(ctx context.Context, order *domain.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *OrderRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	args := m.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.Order), args.Error(1)
}

func (m *OrderRepositoryMock) GetByUser(ctx context.Context, userID uuid.UUID) ([]*domain.Order, error) {
	args := m.Called(ctx, userID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*domain.Order), args.Error(1)
}

func (m *OrderRepositoryMock) GetByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]*domain.Order, error) {
	args := m.Called(ctx, restaurantID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*domain.Order), args.Error(1)
}

func (m *OrderRepositoryMock) GetActiveByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]*domain.Order, error) {
	args := m.Called(ctx, restaurantID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*domain.Order), args.Error(1)
}

func (m *OrderRepositoryMock) GetByPickupCode(ctx context.Context, restaurantID uuid.UUID, pickupCode string) (*domain.Order, error) {
	args := m.Called(ctx, restaurantID, pickupCode)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.Order), args.Error(1)
}

func (m *OrderRepositoryMock) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.OrderStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *OrderRepositoryMock) AddStatusHistory(ctx context.Context, entry *domain.OrderStatusHistoryEntry) error {
	args := m.Called(ctx, entry)
	return args.Error(0)
}

func (m *OrderRepositoryMock) GetStatusHistory(ctx context.Context, orderID uuid.UUID) ([]*domain.OrderStatusHistoryEntry, error) {
	args := m.Called(ctx, orderID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*domain.OrderStatusHistoryEntry), args.Error(1)
}

func (m *OrderRepositoryMock) GetExpiredPendingOrders(ctx context.Context, threshold time.Time) ([]*domain.Order, error) {
	args := m.Called(ctx, threshold)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*domain.Order), args.Error(1)
}

func (m *OrderRepositoryMock) AddPaymentHistory(ctx context.Context, entry *domain.PaymentHistoryEntry) error {
	args := m.Called(ctx, entry)
	return args.Error(0)
}

func (m *OrderRepositoryMock) GetPaymentHistory(ctx context.Context, orderID uuid.UUID) ([]*domain.PaymentHistoryEntry, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.PaymentHistoryEntry), args.Error(1)
}
func (m *OrderRepositoryMock) GetPaymentHistoryByUser(ctx context.Context, userID uuid.UUID) ([]*domain.PaymentHistoryEntry, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.PaymentHistoryEntry), args.Error(1)
}
