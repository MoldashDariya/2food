package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"gitlab.com/2foood/backend/internal/domain"
)

type PushDispatcherMock struct {
	mock.Mock
}

func (m *PushDispatcherMock) SendAsync(ctx context.Context, userID uuid.UUID, message *domain.PushMessage) error {
	args := m.Called(ctx, userID, message)
	return args.Error(0)
}
