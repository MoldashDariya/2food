package provider

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

const mockBaseURL = "http://localhost:9000/mock-bucket"

type Mock struct {
	log *slog.Logger
}

func NewMock(log *slog.Logger) *Mock {
	return &Mock{log: log}
}

func (m *Mock) PresignPut(ctx context.Context, key string, contentType string, ttl time.Duration) (string, error) {
	uploadURL := fmt.Sprintf("%s/%s?mock=true&ttl=%s", mockBaseURL, key, ttl)

	m.log.Info("mock presign put",
		slog.String("key", key),
		slog.String("content_type", contentType),
		slog.String("upload_url", uploadURL),
	)

	return uploadURL, nil
}

func (m *Mock) Delete(ctx context.Context, key string) error {
	m.log.Info("mock delete", slog.String("key", key))
	return nil
}

func (m *Mock) PublicURL(key string) string {
	return fmt.Sprintf("%s/%s", mockBaseURL, key)
}

func (m *Mock) KeyFromURL(fileURL string) string {
	prefix := mockBaseURL + "/"
	if len(fileURL) > len(prefix) {
		return fileURL[len(prefix):]
	}
	return fileURL
}
