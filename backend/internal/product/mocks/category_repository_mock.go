package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"gitlab.com/2foood/backend/internal/domain"
)

type CategoryRepositoryMock struct {
	mock.Mock
}

func (m *CategoryRepositoryMock) Create(ctx context.Context, category *domain.ProductCategory) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *CategoryRepositoryMock) GetByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]*domain.ProductCategory, error) {
	args := m.Called(ctx, restaurantID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*domain.ProductCategory), args.Error(1)
}

func (m *CategoryRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*domain.ProductCategory, error) {
	args := m.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.ProductCategory), args.Error(1)
}

func (m *CategoryRepositoryMock) Update(ctx context.Context, category *domain.ProductCategory) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *CategoryRepositoryMock) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
