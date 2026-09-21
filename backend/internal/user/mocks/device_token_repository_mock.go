package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"gitlab.com/2foood/backend/internal/domain"
)

type DeviceTokenRepositoryMock struct {
	mock.Mock
}

func (m *DeviceTokenRepositoryMock) Register(ctx context.Context, dt *domain.DeviceToken) error {
	args := m.Called(ctx, dt)
	return args.Error(0)
}

func (m *DeviceTokenRepositoryMock) Delete(ctx context.Context, userID uuid.UUID, token string) error {
	args := m.Called(ctx, userID, token)
	return args.Error(0)
}
