package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"gitlab.com/2foood/backend/internal/domain"
)

type CuisineTypeRepositoryMock struct {
	mock.Mock
}

func (m *CuisineTypeRepositoryMock) Create(ctx context.Context, cuisine *domain.CuisineType) error {
	args := m.Called(ctx, cuisine)
	return args.Error(0)
}

func (m *CuisineTypeRepositoryMock) GetAll(ctx context.Context) ([]*domain.CuisineType, error) {
	args := m.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*domain.CuisineType), args.Error(1)
}

func (m *CuisineTypeRepositoryMock) GetByID(ctx context.Context, id int) (*domain.CuisineType, error) {
	args := m.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.CuisineType), args.Error(1)
}

func (m *CuisineTypeRepositoryMock) GetByName(ctx context.Context, name string) (*domain.CuisineType, error) {
	args := m.Called(ctx, name)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.CuisineType), args.Error(1)
}

func (m *CuisineTypeRepositoryMock) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
