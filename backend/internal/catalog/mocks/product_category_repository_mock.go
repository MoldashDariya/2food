package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"gitlab.com/2foood/backend/internal/domain"
)

type ProductCategoryRepositoryMock struct {
	mock.Mock
}

func (m *ProductCategoryRepositoryMock) Create(ctx context.Context, category *domain.ProductCategory) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *ProductCategoryRepositoryMock) GetAll(ctx context.Context) ([]*domain.ProductCategory, error) {
	args := m.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*domain.ProductCategory), args.Error(1)
}

func (m *ProductCategoryRepositoryMock) GetByID(ctx context.Context, id int) (*domain.ProductCategory, error) {
	args := m.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.ProductCategory), args.Error(1)
}

func (m *ProductCategoryRepositoryMock) GetByName(ctx context.Context, name string) (*domain.ProductCategory, error) {
	args := m.Called(ctx, name)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.ProductCategory), args.Error(1)
}

func (m *ProductCategoryRepositoryMock) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
