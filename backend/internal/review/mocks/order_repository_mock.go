package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"gitlab.com/2foood/backend/internal/domain"
)

type OrderRepositoryMock struct {
	mock.Mock
}

func (m *OrderRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	args := m.Called(ctx, id)

	if args.Get(0) != nil {
		return args.Get(0).(*domain.Order), args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *OrderRepositoryMock) SetReviewed(ctx context.Context, orderID uuid.UUID) error {
	args := m.Called(ctx, orderID)
	return args.Error(0)
}
