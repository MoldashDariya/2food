package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type OTPManagerMock struct {
	mock.Mock
}

func (m *OTPManagerMock) Send(ctx context.Context, to string) error {
	args := m.Called(ctx, to)
	return args.Error(0)
}

func (m *OTPManagerMock) Verify(ctx context.Context, to, code string) error {
	args := m.Called(ctx, to, code)
	return args.Error(0)
}
