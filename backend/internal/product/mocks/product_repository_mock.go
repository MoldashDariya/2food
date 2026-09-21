package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"gitlab.com/2foood/backend/internal/domain"
)

type ProductRepositoryMock struct {
	mock.Mock
}

func (m *ProductRepositoryMock) Create(ctx context.Context, product *domain.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *ProductRepositoryMock) GetByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]*domain.Product, error) {
	args := m.Called(ctx, restaurantID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*domain.Product), args.Error(1)
}

func (m *ProductRepositoryMock) GetActiveByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]*domain.Product, error) {
	args := m.Called(ctx, restaurantID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*domain.Product), args.Error(1)
}

func (m *ProductRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	args := m.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *ProductRepositoryMock) GetCategoriesByRestaurant(ctx context.Context, restID uuid.UUID) ([]*domain.ProductCategory, error) {
	args := m.Called(ctx, restID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*domain.ProductCategory), args.Error(1)
}

func (m *ProductRepositoryMock) Update(ctx context.Context, product *domain.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *ProductRepositoryMock) UpdateStatus(ctx context.Context, id uuid.UUID, isActive bool) error {
	args := m.Called(ctx, id, isActive)
	return args.Error(0)
}

func (m *ProductRepositoryMock) UpdateStock(ctx context.Context, id uuid.UUID, quantity int) error {
	args := m.Called(ctx, id, quantity)
	return args.Error(0)
}

func (m *ProductRepositoryMock) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
