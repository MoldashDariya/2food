package mocks

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type PaymentDispatcherMock struct {
	mock.Mock
}

func (m *PaymentDispatcherMock) ConfirmAsync(ctx context.Context, paymentID string, orderID uuid.UUID) error {
	args := m.Called(ctx, paymentID, orderID)
	return args.Error(0)
}

func (m *PaymentDispatcherMock) CancelAsync(ctx context.Context, paymentID string, orderID uuid.UUID) error {
	args := m.Called(ctx, paymentID, orderID)
	return args.Error(0)
}

func (m *PaymentDispatcherMock) RefundAsync(ctx context.Context, paymentID string, orderID uuid.UUID) error {
	args := m.Called(ctx, paymentID, orderID)
	return args.Error(0)
}

func (m *PaymentDispatcherMock) AuditAsync(ctx context.Context, paymentID string, orderID uuid.UUID, delay time.Duration) error {
	args := m.Called(ctx, paymentID, orderID, delay)
	return args.Error(0)
}
