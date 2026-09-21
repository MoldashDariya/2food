package upload_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/2foood/backend/internal/upload"
	"gitlab.com/2foood/backend/internal/upload/mocks"
)

func setupTest() (*mocks.StorageManagerMock, *upload.Service) {
	storage := &mocks.StorageManagerMock{}
	svc := upload.NewService(storage)
	return storage, svc
}

func TestUploadService_PresignUpload(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name        string
		folder      string
		contentType string
		mockFn      func(storage *mocks.StorageManagerMock)
		wantErr     bool
	}{
		{
			name:        "storage error",
			folder:      "avatars",
			contentType: "image/jpeg",
			mockFn: func(storage *mocks.StorageManagerMock) {
				storage.On("PresignUpload", mock.Anything, "avatars", userID, "image/jpeg").
					Return("", "", errors.New("storage error"))
			},
			wantErr: true,
		},
		{
			name:        "success",
			folder:      "avatars",
			contentType: "image/jpeg",
			mockFn: func(storage *mocks.StorageManagerMock) {
				storage.On("PresignUpload", mock.Anything, "avatars", userID, "image/jpeg").
					Return("https://upload.example.com/presign", "https://cdn.example.com/file.jpg", nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			storage, svc := setupTest()
			tt.mockFn(storage)

			uploadURL, fileURL, err := svc.PresignUpload(context.Background(), tt.folder, userID, tt.contentType)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, uploadURL)
				assert.Empty(t, fileURL)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, uploadURL)
				assert.NotEmpty(t, fileURL)
			}
			storage.AssertExpectations(t)
		})
	}
}

func TestUploadService_Delete(t *testing.T) {
	fileURL := "https://cdn.example.com/avatars/file.jpg"

	tests := []struct {
		name    string
		mockFn  func(storage *mocks.StorageManagerMock)
		wantErr bool
	}{
		{
			name: "storage error",
			mockFn: func(storage *mocks.StorageManagerMock) {
				storage.On("Delete", mock.Anything, fileURL).Return(errors.New("storage error"))
			},
			wantErr: true,
		},
		{
			name: "success",
			mockFn: func(storage *mocks.StorageManagerMock) {
				storage.On("Delete", mock.Anything, fileURL).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			storage, svc := setupTest()
			tt.mockFn(storage)

			err := svc.Delete(context.Background(), fileURL)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			storage.AssertExpectations(t)
		})
	}
}