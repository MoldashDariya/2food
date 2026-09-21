package provider

import (
	"context"

	"gitlab.com/2foood/backend/internal/domain"
	"google.golang.org/api/idtoken"
)

type GoogleConfig struct {
	ClientID string `envconfig:"GOOGLE_CLIENT_ID" required:"true"`
}

type Google struct {
	clientID string
}

func NewGoogle(cfg GoogleConfig) *Google {
	return &Google{
		clientID: cfg.ClientID,
	}
}

func (g *Google) VerifyID(ctx context.Context, idToken string) (*domain.OAuthInfo, error) {
	payload, err := idtoken.Validate(ctx, idToken, g.clientID)
	if err != nil {
		return nil, domain.ErrInvalidOAuthToken
	}

	claims := payload.Claims

	oauthClaims := domain.OAuthInfo{
		OAuthID:       getClaim[string](claims, "sub"),
		Email:         getClaim[string](claims, "email"),
		EmailVerified: getClaim[bool](claims, "email_verified"),
		Name:          getClaim[string](claims, "name"),
		AvatarURL:     getClaim[string](claims, "picture"),
	}
	return &oauthClaims, nil
}

func getClaim[T any](claims map[string]any, key string) T {
	var zero T
	if v, ok := claims[key]; ok {
		if t, ok := v.(T); ok {
			return t
		}
	}

	return zero
}
