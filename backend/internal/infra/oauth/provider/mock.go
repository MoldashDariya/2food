package provider

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"strings"

	"gitlab.com/2foood/backend/internal/domain"
)

type Mock struct {
	log *slog.Logger
}

func NewMock(log *slog.Logger) *Mock {
	return &Mock{log: log}
}

func (p *Mock) VerifyID(ctx context.Context, idToken string) (*domain.OAuthInfo, error) {
	if idToken == "" {
		idToken = "default_user"
	}

	seed := idToken
	if len(seed) > 30 {
		hash := sha256.Sum256([]byte(seed))
		seed = hex.EncodeToString(hash[:])[:10]
	}

	seed = strings.TrimSpace(seed)

	p.log.Info("MOCK OAuth verify",
		slog.String("input_token", idToken),
		slog.String("generated_user", seed),
	)

	return &domain.OAuthInfo{
		OAuthID:   fmt.Sprintf("mock-oauth-id-%s", seed),
		Email:     fmt.Sprintf("%s@mock.test", seed),
		Name:      fmt.Sprintf("Mock User %s", seed),
		AvatarURL: fmt.Sprintf("https://robohash.org/%s?set=set4", seed),
	}, nil
}
