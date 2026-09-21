package push

import (
	"context"

	"github.com/google/uuid"
	"gitlab.com/2foood/backend/internal/domain"
)

type Provider interface {
	Send(ctx context.Context, token string, title, body string, data map[string]any) error
	SendMulti(ctx context.Context, tokens []string, title, body string, data map[string]any) (staleTokens []string, err error)
}

type DeviceTokenRepository interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.DeviceToken, error)
	Delete(ctx context.Context, userID uuid.UUID, token string) error
}

type Manager struct {
	provider        Provider
	deviceTokenRepo DeviceTokenRepository
}

func NewManager(provider Provider, deviceTokenRepo DeviceTokenRepository) *Manager {
	return &Manager{provider: provider, deviceTokenRepo: deviceTokenRepo}
}

func (m *Manager) SendToUser(ctx context.Context, userID uuid.UUID, title, body string, data map[string]any) (int, error) {
	tokens, err := m.deviceTokenRepo.GetByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}

	if len(tokens) == 0 {
		return 0, nil // no tokens to send to user
	}

	rawTokens := make([]string, len(tokens))
	for i, t := range tokens {
		rawTokens[i] = t.Token
	}

	stale, err := m.provider.SendMulti(ctx, rawTokens, title, body, data)
	if err != nil {
		return 0, err
	}

	// delete stale tokens
	for _, token := range stale {
		_ = m.deviceTokenRepo.Delete(ctx, userID, token)
	}

	return len(tokens) - len(stale), nil // number of tokens sent successfully
}
