package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gitlab.com/2foood/backend/internal/domain"
)

type Provider interface {
	PresignPut(ctx context.Context, key string, contentType string, ttl time.Duration) (uploadURL string, err error)
	Delete(ctx context.Context, key string) error
	PublicURL(key string) string
	KeyFromURL(url string) string
}

type Manager struct {
	cfg      Config
	provider Provider
}

func NewManager(cfg Config, provider Provider) *Manager {
	return &Manager{
		cfg:      cfg,
		provider: provider,
	}
}

func (m *Manager) PresignUpload(ctx context.Context, folder string, userID uuid.UUID, contentType string) (uploadURL string, fileURL string, err error) {
	ext, err := m.contentTypeToExt(contentType)
	if err != nil {
		return "", "", domain.ErrUnsupportedContentType
	}

	key := fmt.Sprintf("%s/%s/%s%s", folder, userID, uuid.New(), ext)

	uploadURL, err = m.provider.PresignPut(ctx, key, contentType, m.cfg.PresignTTL)
	if err != nil {
		return "", "", fmt.Errorf("failed to presign upload: %w", err)
	}

	fileURL = m.provider.PublicURL(key)
	return uploadURL, fileURL, nil
}

func (m *Manager) Delete(ctx context.Context, fileURL string) error {
	key := m.provider.KeyFromURL(fileURL)
	return m.provider.Delete(ctx, key)
}

func (m *Manager) contentTypeToExt(ct string) (string, error) {
	switch ct {
	case "image/jpeg":
		return ".jpg", nil
	case "image/png":
		return ".png", nil
	case "image/webp":
		return ".webp", nil
	default:
		return "", fmt.Errorf("unsupported content type: %s", ct)
	}
}
