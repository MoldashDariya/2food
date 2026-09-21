package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"gitlab.com/2foood/backend/internal/domain"
)

type OAuthManagerMock struct {
	mock.Mock
}

func (m *OAuthManagerMock) VerifyID(ctx context.Context, token string) (*domain.OAuthInfo, error) {
	args := m.Called(ctx, token)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.OAuthInfo), args.Error(1)
}

func (m *OAuthManagerMock) CreateSession(ctx context.Context, info *domain.OAuthInfo) (string, error) {
	args := m.Called(ctx, info)

	if args.Get(0) == nil {
		return "", args.Error(1)
	}

	return args.Get(0).(string), args.Error(1)
}

func (m *OAuthManagerMock) GetSession(ctx context.Context, sessionID string) (*domain.OAuthInfo, error) {
	args := m.Called(ctx, sessionID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.OAuthInfo), args.Error(1)
}

func (m *OAuthManagerMock) DeleteSession(ctx context.Context, sessionID string) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}
