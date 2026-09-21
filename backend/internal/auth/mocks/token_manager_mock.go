package mocks

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"gitlab.com/2foood/backend/internal/domain"
)

type TokenManagerMock struct {
	mock.Mock
}

func (m *TokenManagerMock) GenerateTokenPair(userID uuid.UUID, role domain.UserRole) (*domain.TokenPair, error) {
	args := m.Called(userID, role)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.TokenPair), args.Error(1)
}

func (m *TokenManagerMock) ValidateToken(tokenString string) (*domain.Claims, error) {
	args := m.Called(tokenString)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.Claims), args.Error(1)
}
