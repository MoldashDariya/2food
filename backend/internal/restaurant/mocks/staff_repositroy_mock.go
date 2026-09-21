package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"gitlab.com/2foood/backend/internal/domain"
)

type StaffRepositoryMock struct {
	mock.Mock
}

func (m *StaffRepositoryMock) Add(ctx context.Context, staff *domain.RestaurantStaff) error {
	args := m.Called(ctx, staff)
	return args.Error(0)
}

func (m *StaffRepositoryMock) GetByRestaurant(ctx context.Context, restID uuid.UUID) ([]*domain.RestaurantStaff, error) {
	args := m.Called(ctx, restID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*domain.RestaurantStaff), args.Error(1)
}

func (m *StaffRepositoryMock) GetMember(ctx context.Context, restID, userID uuid.UUID) (*domain.RestaurantStaff, error) {
	args := m.Called(ctx, restID, userID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.RestaurantStaff), args.Error(1)
}

func (m *StaffRepositoryMock) CountRestaurants(ctx context.Context, userID uuid.UUID) (int, error) {
	args := m.Called(ctx, userID)
	return args.Int(0), args.Error(1)
}

func (m *StaffRepositoryMock) SetRole(ctx context.Context, restID, userID uuid.UUID, role domain.StaffRole) error {
	args := m.Called(ctx, restID, userID, role)
	return args.Error(0)
}

func (m *StaffRepositoryMock) Remove(ctx context.Context, restID, userID uuid.UUID) error {
	args := m.Called(ctx, restID, userID)
	return args.Error(0)
}

func (m *StaffRepositoryMock) GetOwner(ctx context.Context, restID uuid.UUID) (*domain.RestaurantStaff, error) {
	args := m.Called(ctx, restID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.RestaurantStaff), args.Error(1)
}
