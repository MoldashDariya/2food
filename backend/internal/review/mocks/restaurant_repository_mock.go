package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"gitlab.com/2foood/backend/internal/domain"
)

type RestaurantRepositoryMock struct {
	mock.Mock
}

func (m *RestaurantRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*domain.Restaurant, error) {
	args := m.Called(ctx, id)

	if args.Get(0) != nil {
		return args.Get(0).(*domain.Restaurant), args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *RestaurantRepositoryMock) IncrementRating(ctx context.Context, restID uuid.UUID, rating int8) error {
	args := m.Called(ctx, restID, rating)
	return args.Error(0)
}
