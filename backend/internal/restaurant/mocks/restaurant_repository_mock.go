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

func (m *RestaurantRepositoryMock) Create(ctx context.Context, rest *domain.Restaurant) error {
	args := m.Called(ctx, rest)
	return args.Error(0)
}

func (m *RestaurantRepositoryMock) GetAll(ctx context.Context) ([]*domain.Restaurant, error) {
	args := m.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*domain.Restaurant), args.Error(1)
}

func (m *RestaurantRepositoryMock) GetByPartner(ctx context.Context, userID uuid.UUID) ([]*domain.Restaurant, error) {
	args := m.Called(ctx, userID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*domain.Restaurant), args.Error(1)
}

func (m *RestaurantRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*domain.Restaurant, error) {
	args := m.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.Restaurant), args.Error(1)
}

func (m *RestaurantRepositoryMock) UpdateInfo(ctx context.Context, rest *domain.Restaurant) error {
	args := m.Called(ctx, rest)
	return args.Error(0)
}

func (m *RestaurantRepositoryMock) UpdateCuisines(ctx context.Context, id uuid.UUID, cuisineTypes []int) error {
	args := m.Called(ctx, id, cuisineTypes)
	return args.Error(0)
}

func (m *RestaurantRepositoryMock) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
