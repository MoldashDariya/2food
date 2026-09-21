package oauth

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gitlab.com/2foood/backend/internal/domain"
)

type Provider interface {
	VerifyID(ctx context.Context, token string) (*domain.OAuthInfo, error)
}

type Cache interface {
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
}

type Manager struct {
	cfg      Config
	provider Provider
	cache    Cache
}

func NewManager(cfg Config, provider Provider, cache Cache) *Manager {
	return &Manager{
		cfg:      cfg,
		provider: provider,
		cache:    cache,
	}
}

func (m *Manager) VerifyID(ctx context.Context, token string) (*domain.OAuthInfo, error) {
	return m.provider.VerifyID(ctx, token)
}

func (m *Manager) CreateSession(ctx context.Context, info *domain.OAuthInfo) (string, error) {
	sessionID := uuid.New().String()
	key := m.cfg.KeyPrefix + sessionID

	data, err := json.Marshal(info)
	if err != nil {
		return "", err
	}

	if err := m.cache.Set(ctx, key, string(data), m.cfg.SessionTTL); err != nil {
		return "", err
	}

	return sessionID, nil
}

func (m *Manager) GetSession(ctx context.Context, sessionID string) (*domain.OAuthInfo, error) {
	key := m.cfg.KeyPrefix + sessionID

	data, err := m.cache.Get(ctx, key)
	if err != nil {
		return nil, domain.ErrSessionNotFound
	}

	var info domain.OAuthInfo
	if err := json.Unmarshal([]byte(data), &info); err != nil {
		return nil, err
	}

	return &info, nil
}

func (m *Manager) DeleteSession(ctx context.Context, sessionID string) error {
	key := m.cfg.KeyPrefix + sessionID
	return m.cache.Del(ctx, key)
}
