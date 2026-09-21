package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"gitlab.com/2foood/backend/internal/domain"
)

type ReviewRepositoryMock struct {
	mock.Mock
}

func (m *ReviewRepositoryMock) Create(ctx context.Context, review *domain.Review) error {
	args := m.Called(ctx, review)
	return args.Error(0)
}

func (m *ReviewRepositoryMock) GetAll(ctx context.Context) ([]*domain.Review, error) {
	args := m.Called(ctx)

	if args.Get(0) != nil {
		return args.Get(0).([]*domain.Review), args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *ReviewRepositoryMock) GetByRestaurant(ctx context.Context, restID uuid.UUID) ([]*domain.Review, error) {
	args := m.Called(ctx, restID)

	if args.Get(0) != nil {
		return args.Get(0).([]*domain.Review), args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *ReviewRepositoryMock) GetByRestaurantAll(ctx context.Context, restID uuid.UUID) ([]*domain.Review, error) {
	args := m.Called(ctx, restID)

	if args.Get(0) != nil {
		return args.Get(0).([]*domain.Review), args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *ReviewRepositoryMock) GetByUser(ctx context.Context, userID uuid.UUID) ([]*domain.Review, error) {
	args := m.Called(ctx, userID)

	if args.Get(0) != nil {
		return args.Get(0).([]*domain.Review), args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *ReviewRepositoryMock) GetByOrder(ctx context.Context, orderID uuid.UUID) (*domain.Review, error) {
	args := m.Called(ctx, orderID)

	if args.Get(0) != nil {
		return args.Get(0).(*domain.Review), args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *ReviewRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*domain.Review, error) {
	args := m.Called(ctx, id)

	if args.Get(0) != nil {
		return args.Get(0).(*domain.Review), args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *ReviewRepositoryMock) ExistsByOrder(ctx context.Context, orderID uuid.UUID) (bool, error) {
	args := m.Called(ctx, orderID)
	return args.Bool(0), args.Error(1)
}

func (m *ReviewRepositoryMock) GetPhotosByRestaurant(ctx context.Context, restID uuid.UUID) ([]*domain.ReviewPhoto, error) {
	args := m.Called(ctx, restID)

	if args.Get(0) != nil {
		return args.Get(0).([]*domain.ReviewPhoto), args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *ReviewRepositoryMock) SetReply(ctx context.Context, review *domain.Review) error {
	args := m.Called(ctx, review)
	return args.Error(0)
}

func (m *ReviewRepositoryMock) SetVisibility(ctx context.Context, id uuid.UUID, isHidden bool) error {
	args := m.Called(ctx, id, isHidden)
	return args.Error(0)
}
