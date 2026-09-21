package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"gitlab.com/2foood/backend/internal/domain"
)

type Manager struct {
	secret          []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewManager(cfg Config) *Manager {
	return &Manager{
		secret:          []byte(cfg.Secret),
		accessTokenTTL:  cfg.AccessTokenTTL,
		refreshTokenTTL: cfg.RefreshTokenTTL,
	}
}

// two tokens: short-lived access + long-lived refresh. same secret, the Scope claim tells them apart
func (m *Manager) GenerateTokenPair(userID uuid.UUID, role domain.UserRole) (*domain.TokenPair, error) {
	accessToken, err := m.generateToken(userID, m.accessTokenTTL, role, domain.ScopeAccessToken)
	if err != nil {
		return nil, err
	}

	refreshToken, err := m.generateToken(userID, m.refreshTokenTTL, role, domain.ScopeRefreshToken)
	if err != nil {
		return nil, err
	}

	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (m *Manager) ValidateToken(tokenString string) (*domain.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.Claims{}, func(token *jwt.Token) (interface{}, error) {
		// only accept HMAC, otherwise someone could send a token signed with another algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrInvalidToken
		}
		return m.secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, domain.ErrExpiredToken
		}
		return nil, domain.ErrInvalidToken
	}

	claims, ok := token.Claims.(*domain.Claims)
	if !ok || !token.Valid {
		return nil, domain.ErrInvalidToken
	}

	return claims, nil
}

func (m *Manager) generateToken(userID uuid.UUID, ttl time.Duration, role domain.UserRole, scope domain.Scope) (string, error) {
	now := time.Now()

	claims := domain.Claims{
		UserID: userID,
		Role:   role,
		Scope:  scope,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}
