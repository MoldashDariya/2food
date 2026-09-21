package otp

import (
	"context"
	"crypto/rand"
	"math/big"
	"time"

	"gitlab.com/2foood/backend/internal/domain"
)

// provider = how the code is delivered (twilio sms, smtp email, or a mock for dev)
type Provider interface {
	Send(ctx context.Context, to string, code string) error
}

type Cache interface {
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
	Incr(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, ttl time.Duration) error
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

// code lives in redis with a ttl, key = prefix + phone/email
func (m *Manager) Send(ctx context.Context, to string) error {
	if err := m.checkRateLimit(ctx, to); err != nil {
		return err
	}

	code, err := m.generateCode(m.cfg.Length)
	if err != nil {
		return err
	}

	key := m.cfg.KeyPrefix + to
	if err := m.cache.Set(ctx, key, code, m.cfg.TTL); err != nil {
		return err
	}

	return m.provider.Send(ctx, to, code)
}

func (m *Manager) Verify(ctx context.Context, to, code string) error {
	key := m.cfg.KeyPrefix + to

	stored, err := m.cache.Get(ctx, key)
	if err != nil {
		return domain.ErrCodeExpired
	}

	if stored != code {
		return domain.ErrCodeInvalid
	}

	// codes are one-time, delete after a successful check
	_ = m.cache.Del(ctx, key)
	return nil
}

func (m *Manager) generateCode(length int) (string, error) {
	const digits = "0123456789"
	code := make([]byte, length)
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		code[i] = digits[n.Int64()]
	}
	return string(code), nil
}

// counts how many codes were sent to this phone/email in the current period.
// note: only sending is limited, wrong guesses in Verify aren't counted (only the global ip rate limit applies)
func (m *Manager) checkRateLimit(ctx context.Context, to string) error {
	key := m.cfg.KeyPrefix + m.cfg.RateKeyPrefix + to

	count, err := m.cache.Incr(ctx, key)
	if err != nil {
		return err
	}

	if count == 1 {
		_ = m.cache.Expire(ctx, key, m.cfg.RatePeriod)
	}

	if count > m.cfg.RateLimit {
		return domain.ErrTooManyAttempts
	}

	return nil
}
