package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"gitlab.com/2foood/backend/internal/domain"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *UserRepositoryMock) SetRole(ctx context.Context, id uuid.UUID, role domain.UserRole) error {
	args := m.Called(ctx, id, role)
	return args.Error(0)
}
