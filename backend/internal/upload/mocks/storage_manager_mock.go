package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type StorageManagerMock struct {
	mock.Mock
}

func (m *StorageManagerMock) PresignUpload(ctx context.Context, folder string, userID uuid.UUID, contentType string) (string, string, error) {
	args := m.Called(ctx, folder, userID, contentType)

	if args.Get(0) == nil {
		return "", "", args.Error(2)
	}

	if args.Get(1) == nil {
		return "", "", args.Error(2)
	}

	return args.Get(0).(string), args.Get(1).(string), args.Error(2)
}

func (m *StorageManagerMock) Delete(ctx context.Context, fileURL string) error {
	args := m.Called(ctx, fileURL)
	return args.Error(0)
}
