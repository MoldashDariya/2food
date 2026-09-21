package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"gitlab.com/2foood/backend/internal/domain"
)

type CityRepositoryMock struct {
	mock.Mock
}

func (m *CityRepositoryMock) Create(ctx context.Context, city *domain.City) error {
	args := m.Called(ctx, city)
	return args.Error(0)
}

func (m *CityRepositoryMock) GetAll(ctx context.Context) ([]*domain.City, error) {
	args := m.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*domain.City), args.Error(1)
}

func (m *CityRepositoryMock) GetByID(ctx context.Context, id int) (*domain.City, error) {
	args := m.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.City), args.Error(1)
}

func (m *CityRepositoryMock) GetByName(ctx context.Context, name string) (*domain.City, error) {
	args := m.Called(ctx, name)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.City), args.Error(1)
}

func (m *CityRepositoryMock) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}