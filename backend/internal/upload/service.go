package upload

import (
	"context"

	"github.com/google/uuid"
)

type StorageManager interface {
	PresignUpload(ctx context.Context, folder string, userID uuid.UUID, contentType string) (uploadURL string, fileURL string, err error)
	Delete(ctx context.Context, fileURL string) error
}

type Service struct {
	storage StorageManager
}

func NewService(storage StorageManager) *Service {
	return &Service{storage: storage}
}

func (s *Service) PresignUpload(ctx context.Context, folder string, userID uuid.UUID, contentType string) (string, string, error) {
	return s.storage.PresignUpload(ctx, folder, userID, contentType)
}

func (s *Service) Delete(ctx context.Context, fileURL string) error {
	return s.storage.Delete(ctx, fileURL)
}
