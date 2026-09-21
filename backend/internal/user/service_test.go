package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/2foood/backend/internal/domain"
	"gitlab.com/2foood/backend/internal/user"
	"gitlab.com/2foood/backend/internal/user/mocks"
)

func setupTest() (*mocks.UserRepositoryMock, *mocks.DeviceTokenRepositoryMock, *mocks.OTPManagerMock, *mocks.OTPManagerMock, *user.Service) {
	tm := &mocks.TransactionManagerFake{}
	userRepo := &mocks.UserRepositoryMock{}
	deviceTokenRepo := &mocks.DeviceTokenRepositoryMock{}
	phoneOTP := &mocks.OTPManagerMock{}
	emailOTP := &mocks.OTPManagerMock{}

	svc := user.NewService(tm, userRepo, deviceTokenRepo, phoneOTP, emailOTP)

	return userRepo, deviceTokenRepo, phoneOTP, emailOTP, svc
}

func TestUserService_Get(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name    string
		mockFn  func(userRepo *mocks.UserRepositoryMock)
		wantErr bool
	}{
		{
			name: "not found",
			mockFn: func(userRepo *mocks.UserRepositoryMock) {
				userRepo.On("GetByID", mock.Anything, userID).Return(nil, domain.ErrNotFound)
			},
			wantErr: true,
		},
		{
			name: "success",
			mockFn: func(userRepo *mocks.UserRepositoryMock) {
				name := "Test User"
				userRepo.On("GetByID", mock.Anything, userID).
					Return(&domain.User{ID: userID, Name: &name}, nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			userRepo, _, _, _, svc := setupTest()
			tt.mockFn(userRepo)

			res, err := svc.Get(context.Background(), userID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, userID, res.ID)
			}
			userRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_Update(t *testing.T) {
	userID := uuid.New()
	inputName := "Updated Name"
	input := domain.UpdateUserInput{
		Name: &inputName,
	}

	tests := []struct {
		name    string
		mockFn  func(userRepo *mocks.UserRepositoryMock)
		wantErr bool
	}{
		{
			name: "user not found",
			mockFn: func(userRepo *mocks.UserRepositoryMock) {
				userRepo.On("GetByID", mock.Anything, userID).Return(nil, domain.ErrNotFound)
			},
			wantErr: true,
		},
		{
			name: "db update error",
			mockFn: func(userRepo *mocks.UserRepositoryMock) {
				oldName := "Old Name"
				userRepo.On("GetByID", mock.Anything, userID).
					Return(&domain.User{ID: userID, Name: &oldName}, nil)
				userRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.User")).
					Return(errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "success",
			mockFn: func(userRepo *mocks.UserRepositoryMock) {
				oldName := "Old Name"
				userRepo.On("GetByID", mock.Anything, userID).
					Return(&domain.User{ID: userID, Name: &oldName}, nil)
				userRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.User")).
					Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			userRepo, _, _, _, svc := setupTest()
			tt.mockFn(userRepo)

			res, err := svc.Update(context.Background(), userID, input)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, inputName, *res.Name)
			}
			userRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_Delete(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name    string
		mockFn  func(userRepo *mocks.UserRepositoryMock)
		wantErr bool
	}{
		{
			name: "db error",
			mockFn: func(userRepo *mocks.UserRepositoryMock) {
				userRepo.On("Delete", mock.Anything, userID).Return(errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "success",
			mockFn: func(userRepo *mocks.UserRepositoryMock) {
				userRepo.On("Delete", mock.Anything, userID).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			userRepo, _, _, _, svc := setupTest()
			tt.mockFn(userRepo)

			err := svc.Delete(context.Background(), userID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			userRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_UpdatePhone(t *testing.T) {
	userID := uuid.New()
	newPhone := "+79991234567"
	code := "123456"

	tests := []struct {
		name    string
		mockFn  func(userRepo *mocks.UserRepositoryMock, phoneOTP *mocks.OTPManagerMock)
		wantErr bool
	}{
		{
			name: "invalid OTP code",
			mockFn: func(u *mocks.UserRepositoryMock, otp *mocks.OTPManagerMock) {
				otp.On("Verify", mock.Anything, newPhone, code).Return(domain.ErrCodeInvalid)
			},
			wantErr: true,
		},
		{
			name: "expired OTP code",
			mockFn: func(u *mocks.UserRepositoryMock, otp *mocks.OTPManagerMock) {
				otp.On("Verify", mock.Anything, newPhone, code).Return(domain.ErrCodeExpired)
			},
			wantErr: true,
		},
		{
			name: "phone already taken",
			mockFn: func(u *mocks.UserRepositoryMock, otp *mocks.OTPManagerMock) {
				otp.On("Verify", mock.Anything, newPhone, code).Return(nil)
				u.On("GetByPhone", mock.Anything, newPhone).Return(&domain.User{}, nil)
			},
			wantErr: true,
		},
		{
			name: "db error on GetByPhone",
			mockFn: func(u *mocks.UserRepositoryMock, otp *mocks.OTPManagerMock) {
				otp.On("Verify", mock.Anything, newPhone, code).Return(nil)
				u.On("GetByPhone", mock.Anything, newPhone).Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "db error on SetPhone",
			mockFn: func(u *mocks.UserRepositoryMock, otp *mocks.OTPManagerMock) {
				otp.On("Verify", mock.Anything, newPhone, code).Return(nil)
				u.On("GetByPhone", mock.Anything, newPhone).Return(nil, domain.ErrNotFound)
				u.On("SetPhone", mock.Anything, userID, newPhone).Return(errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "success",
			mockFn: func(u *mocks.UserRepositoryMock, otp *mocks.OTPManagerMock) {
				otp.On("Verify", mock.Anything, newPhone, code).Return(nil)
				u.On("GetByPhone", mock.Anything, newPhone).Return(nil, domain.ErrNotFound)
				u.On("SetPhone", mock.Anything, userID, newPhone).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			userRepo, _, phoneOTP, _, svc := setupTest()
			tt.mockFn(userRepo, phoneOTP)

			err := svc.UpdatePhone(context.Background(), userID, newPhone, code)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			userRepo.AssertExpectations(t)
			phoneOTP.AssertExpectations(t)
		})
	}
}

func TestUserService_UpdateEmail(t *testing.T) {
	userID := uuid.New()
	newEmail := "user@example.com"
	code := "123456"

	tests := []struct {
		name    string
		mockFn  func(userRepo *mocks.UserRepositoryMock, emailOTP *mocks.OTPManagerMock)
		wantErr bool
	}{
		{
			name: "invalid OTP code",
			mockFn: func(u *mocks.UserRepositoryMock, otp *mocks.OTPManagerMock) {
				otp.On("Verify", mock.Anything, newEmail, code).Return(domain.ErrCodeInvalid)
			},
			wantErr: true,
		},
		{
			name: "expired OTP code",
			mockFn: func(u *mocks.UserRepositoryMock, otp *mocks.OTPManagerMock) {
				otp.On("Verify", mock.Anything, newEmail, code).Return(domain.ErrCodeExpired)
			},
			wantErr: true,
		},
		{
			name: "email already taken",
			mockFn: func(u *mocks.UserRepositoryMock, otp *mocks.OTPManagerMock) {
				otp.On("Verify", mock.Anything, newEmail, code).Return(nil)
				u.On("GetByEmail", mock.Anything, newEmail).Return(&domain.User{}, nil)
			},
			wantErr: true,
		},
		{
			name: "db error on GetByEmail",
			mockFn: func(u *mocks.UserRepositoryMock, otp *mocks.OTPManagerMock) {
				otp.On("Verify", mock.Anything, newEmail, code).Return(nil)
				u.On("GetByEmail", mock.Anything, newEmail).Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "db error on SetEmail",
			mockFn: func(u *mocks.UserRepositoryMock, otp *mocks.OTPManagerMock) {
				otp.On("Verify", mock.Anything, newEmail, code).Return(nil)
				u.On("GetByEmail", mock.Anything, newEmail).Return(nil, domain.ErrNotFound)
				u.On("SetEmail", mock.Anything, userID, newEmail).Return(errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "success",
			mockFn: func(u *mocks.UserRepositoryMock, otp *mocks.OTPManagerMock) {
				otp.On("Verify", mock.Anything, newEmail, code).Return(nil)
				u.On("GetByEmail", mock.Anything, newEmail).Return(nil, domain.ErrNotFound)
				u.On("SetEmail", mock.Anything, userID, newEmail).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			userRepo, _, _, emailOTP, svc := setupTest()
			tt.mockFn(userRepo, emailOTP)

			err := svc.UpdateEmail(context.Background(), userID, newEmail, code)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			userRepo.AssertExpectations(t)
			emailOTP.AssertExpectations(t)
		})
	}
}

func TestUserService_RegisterDeviceToken(t *testing.T) {
	input := domain.RegisterDeviceTokenInput{
		UserID: uuid.New(),
		Token:  "test-device-token-123",
	}

	tests := []struct {
		name    string
		mockFn  func(deviceRepo *mocks.DeviceTokenRepositoryMock)
		wantErr bool
	}{
		{
			name: "repo error",
			mockFn: func(deviceRepo *mocks.DeviceTokenRepositoryMock) {
				deviceRepo.On("Register", mock.Anything, mock.AnythingOfType("*domain.DeviceToken")).
					Return(errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "success",
			mockFn: func(deviceRepo *mocks.DeviceTokenRepositoryMock) {
				deviceRepo.On("Register", mock.Anything, mock.AnythingOfType("*domain.DeviceToken")).
					Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, deviceRepo, _, _, svc := setupTest()
			tt.mockFn(deviceRepo)

			err := svc.RegisterDeviceToken(context.Background(), input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			deviceRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_DeleteDeviceToken(t *testing.T) {
	userID := uuid.New()
	token := "test-device-token-123"

	tests := []struct {
		name    string
		mockFn  func(deviceRepo *mocks.DeviceTokenRepositoryMock)
		wantErr bool
	}{
		{
			name: "repo error",
			mockFn: func(deviceRepo *mocks.DeviceTokenRepositoryMock) {
				deviceRepo.On("Delete", mock.Anything, userID, token).
					Return(errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "success",
			mockFn: func(deviceRepo *mocks.DeviceTokenRepositoryMock) {
				deviceRepo.On("Delete", mock.Anything, userID, token).
					Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, deviceRepo, _, _, svc := setupTest()
			tt.mockFn(deviceRepo)

			err := svc.DeleteDeviceToken(context.Background(), userID, token)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			deviceRepo.AssertExpectations(t)
		})
	}
}
