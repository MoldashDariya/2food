package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/2foood/backend/internal/auth"
	"gitlab.com/2foood/backend/internal/auth/mocks"
	"gitlab.com/2foood/backend/internal/domain"
)

func setupTest() (*mocks.UserRepositoryMock, *mocks.OTPManagerMock, *mocks.OTPManagerMock, *mocks.OAuthManagerMock, *mocks.TokenManagerMock, *auth.Service) {
	userRepo := &mocks.UserRepositoryMock{}
	phoneOTP := &mocks.OTPManagerMock{}
	emailOTP := &mocks.OTPManagerMock{}
	oauth := &mocks.OAuthManagerMock{}
	tokens := &mocks.TokenManagerMock{}
	svc := auth.NewService(userRepo, phoneOTP, emailOTP, oauth, tokens)
	return userRepo, phoneOTP, emailOTP, oauth, tokens, svc
}

func TestAuthService_SendPhoneOTP(t *testing.T) {
	tests := []struct {
		name    string
		phone   string
		mockFn  func(otp *mocks.OTPManagerMock)
		wantErr bool
	}{
		{
			name:  "success",
			phone: "+79999999999",
			mockFn: func(otp *mocks.OTPManagerMock) {
				otp.On("Send", mock.Anything, "+79999999999").Return(nil)
			},
			wantErr: false,
		},
		{
			name:  "otp manager error",
			phone: "+79999999999",
			mockFn: func(otp *mocks.OTPManagerMock) {
				otp.On("Send", mock.Anything, "+79999999999").Return(errors.New("otp error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, phoneOTP, _, _, _, svc := setupTest()
			tt.mockFn(phoneOTP)

			err := svc.SendPhoneOTP(context.Background(), tt.phone)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			phoneOTP.AssertExpectations(t)
		})
	}
}

func TestAuthService_SendEmailOTP(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		mockFn  func(otp *mocks.OTPManagerMock)
		wantErr bool
	}{
		{
			name:  "success",
			email: "user@example.com",
			mockFn: func(otp *mocks.OTPManagerMock) {
				otp.On("Send", mock.Anything, "user@example.com").Return(nil)
			},
			wantErr: false,
		},
		{
			name:  "otp manager error",
			email: "user@example.com",
			mockFn: func(otp *mocks.OTPManagerMock) {
				otp.On("Send", mock.Anything, "user@example.com").Return(errors.New("otp error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, _, emailOTP, _, _, svc := setupTest()
			tt.mockFn(emailOTP)

			err := svc.SendEmailOTP(context.Background(), tt.email)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			emailOTP.AssertExpectations(t)
		})
	}
}

func TestAuthService_LoginWithPhone(t *testing.T) {
	tests := []struct {
		name        string
		phone       string
		code        string
		mockFn      func(userRepo *mocks.UserRepositoryMock, otp *mocks.OTPManagerMock, tokens *mocks.TokenManagerMock)
		wantErr     bool
		wantNewUser bool
	}{
		{
			name:  "otp verification error",
			phone: "+79999999999",
			code:  "123456",
			mockFn: func(userRepo *mocks.UserRepositoryMock, otp *mocks.OTPManagerMock, tokens *mocks.TokenManagerMock) {
				otp.On("Verify", mock.Anything, "+79999999999", "123456").Return(errors.New("invalid code"))
			},
			wantErr: true,
		},
		{
			name:  "user not found, created, tokens generated",
			phone: "+79999999999",
			code:  "123456",
			mockFn: func(userRepo *mocks.UserRepositoryMock, otp *mocks.OTPManagerMock, tokens *mocks.TokenManagerMock) {
				otp.On("Verify", mock.Anything, "+79999999999", "123456").Return(nil)
				userRepo.On("GetByPhone", mock.Anything, "+79999999999").Return(nil, domain.ErrNotFound)
				userRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)
				tokens.On("GenerateTokenPair", mock.AnythingOfType("uuid.UUID"), mock.Anything).
					Return(&domain.TokenPair{AccessToken: "acc", RefreshToken: "ref"}, nil)
			},
			wantErr:     false,
			wantNewUser: true,
		},
		{
			name:  "existing user, tokens generated",
			phone: "+79999999999",
			code:  "123456",
			mockFn: func(userRepo *mocks.UserRepositoryMock, otp *mocks.OTPManagerMock, tokens *mocks.TokenManagerMock) {
				otp.On("Verify", mock.Anything, "+79999999999", "123456").Return(nil)
				u := domain.NewUserWithPhone("+79999999999", domain.LanguageEnglish)
				userRepo.On("GetByPhone", mock.Anything, "+79999999999").Return(u, nil)
				tokens.On("GenerateTokenPair", u.ID, u.Role).
					Return(&domain.TokenPair{AccessToken: "acc", RefreshToken: "ref"}, nil)
			},
			wantErr:     false,
			wantNewUser: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			userRepo, phoneOTP, _, _, tokens, svc := setupTest()
			tt.mockFn(userRepo, phoneOTP, tokens)

			user, pair, err := svc.LoginWithPhone(context.Background(), tt.phone, tt.code, domain.LanguageEnglish)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.NotNil(t, pair)
			}
			if tt.wantNewUser {
				userRepo.AssertNumberOfCalls(t, "Create", 1)
			} else {
				userRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
			}
			userRepo.AssertExpectations(t)
			phoneOTP.AssertExpectations(t)
			tokens.AssertExpectations(t)
		})
	}
}

func TestAuthService_LoginWithGoogle(t *testing.T) {
	tests := []struct {
		name          string
		idToken       string
		mockFn        func(userRepo *mocks.UserRepositoryMock, oauth *mocks.OAuthManagerMock, tokens *mocks.TokenManagerMock)
		wantErr       bool
		wantNeedPhone bool
	}{
		{
			name:    "verify id fails",
			idToken: "bad",
			mockFn: func(userRepo *mocks.UserRepositoryMock, oauth *mocks.OAuthManagerMock, tokens *mocks.TokenManagerMock) {
				oauth.On("VerifyID", mock.Anything, "bad").Return(nil, errors.New("invalid token"))
			},
			wantErr: true,
		},
		{
			name:    "user not found in db, session created, need phone",
			idToken: "tok-new-user",
			mockFn: func(userRepo *mocks.UserRepositoryMock, oauth *mocks.OAuthManagerMock, tokens *mocks.TokenManagerMock) {
				info := &domain.OAuthInfo{OAuthID: "google-new", Email: "new@example.com"}
				oauth.On("VerifyID", mock.Anything, "tok-new-user").Return(info, nil)
				userRepo.On("GetByGoogleID", mock.Anything, "google-new").Return(nil, domain.ErrNotFound)
				oauth.On("CreateSession", mock.Anything, info).Return("session-123", nil)
			},
			wantErr:       false,
			wantNeedPhone: true,
		},
		{
			name:    "existing user without phone, need phone",
			idToken: "tok-need-phone",
			mockFn: func(userRepo *mocks.UserRepositoryMock, oauth *mocks.OAuthManagerMock, tokens *mocks.TokenManagerMock) {
				info := &domain.OAuthInfo{OAuthID: "google-4", Email: "need@example.com"}
				oauth.On("VerifyID", mock.Anything, "tok-need-phone").Return(info, nil)
				u := domain.NewUserWithPhone("", domain.LanguageEnglish)
				u.Phone = nil
				u.BindGoogleData(domain.OAuthInfo{OAuthID: "google-4", Email: info.Email})
				userRepo.On("GetByGoogleID", mock.Anything, "google-4").Return(u, nil)
				oauth.On("CreateSession", mock.Anything, info).Return("session-456", nil)
			},
			wantErr:       false,
			wantNeedPhone: true,
		},
		{
			name:    "existing user with phone, success",
			idToken: "tok-success",
			mockFn: func(userRepo *mocks.UserRepositoryMock, oauth *mocks.OAuthManagerMock, tokens *mocks.TokenManagerMock) {
				info := &domain.OAuthInfo{OAuthID: "google-6", Email: "ok@example.com"}
				oauth.On("VerifyID", mock.Anything, "tok-success").Return(info, nil)
				u := domain.NewUserWithPhone("+79999999999", domain.LanguageEnglish)
				u.BindGoogleData(domain.OAuthInfo{OAuthID: "google-6", Email: info.Email})
				userRepo.On("GetByGoogleID", mock.Anything, "google-6").Return(u, nil)
				tokens.On("GenerateTokenPair", u.ID, u.Role).
					Return(&domain.TokenPair{AccessToken: "acc", RefreshToken: "ref"}, nil)
			},
			wantErr:       false,
			wantNeedPhone: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			userRepo, _, _, oauth, tokens, svc := setupTest()
			tt.mockFn(userRepo, oauth, tokens)

			res, err := svc.LoginWithGoogle(context.Background(), tt.idToken)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantNeedPhone, res.NeedPhone)
				if res.NeedPhone {
					assert.NotEmpty(t, res.SessionID)
				} else {
					assert.NotNil(t, res.Tokens)
				}
			}
			userRepo.AssertExpectations(t)
			oauth.AssertExpectations(t)
			tokens.AssertExpectations(t)
		})
	}
}

func TestAuthService_CompleteGoogleLogin(t *testing.T) {
	tests := []struct {
		name      string
		sessionID string
		phone     string
		code      string
		mockFn    func(userRepo *mocks.UserRepositoryMock, otp *mocks.OTPManagerMock, oauth *mocks.OAuthManagerMock, tokens *mocks.TokenManagerMock)
		wantErr   bool
	}{
		{
			name:      "invalid or expired session",
			sessionID: "bad-session",
			phone:     "+79999999999",
			code:      "123456",
			mockFn: func(userRepo *mocks.UserRepositoryMock, otp *mocks.OTPManagerMock, oauth *mocks.OAuthManagerMock, tokens *mocks.TokenManagerMock) {
				oauth.On("GetSession", mock.Anything, "bad-session").Return(nil, errors.New("session not found"))
			},
			wantErr: true,
		},
		{
			name:      "otp verification fails",
			sessionID: "valid-session",
			phone:     "+79999999999",
			code:      "wrong",
			mockFn: func(userRepo *mocks.UserRepositoryMock, otp *mocks.OTPManagerMock, oauth *mocks.OAuthManagerMock, tokens *mocks.TokenManagerMock) {
				info := &domain.OAuthInfo{OAuthID: "google-1", Email: "test@test.com"}
				oauth.On("GetSession", mock.Anything, "valid-session").Return(info, nil)
				otp.On("Verify", mock.Anything, "+79999999999", "wrong").Return(errors.New("invalid otp"))
			},
			wantErr: true,
		},
		{
			name:      "phone already linked to another google account",
			sessionID: "valid-session",
			phone:     "+79999999999",
			code:      "123456",
			mockFn: func(userRepo *mocks.UserRepositoryMock, otp *mocks.OTPManagerMock, oauth *mocks.OAuthManagerMock, tokens *mocks.TokenManagerMock) {
				info := &domain.OAuthInfo{OAuthID: "google-new", Email: "test@test.com"}
				oauth.On("GetSession", mock.Anything, "valid-session").Return(info, nil)
				otp.On("Verify", mock.Anything, "+79999999999", "123456").Return(nil)
				existingUser := domain.NewUserWithPhone("+79999999999", domain.LanguageEnglish)
				otherGoogleID := "google-old"
				existingUser.GoogleID = &otherGoogleID
				userRepo.On("GetByPhone", mock.Anything, "+79999999999").Return(existingUser, nil)
			},
			wantErr: true,
		},
		{
			name:      "success - merge existing user",
			sessionID: "valid-session",
			phone:     "+79999999999",
			code:      "123456",
			mockFn: func(userRepo *mocks.UserRepositoryMock, otp *mocks.OTPManagerMock, oauth *mocks.OAuthManagerMock, tokens *mocks.TokenManagerMock) {
				info := &domain.OAuthInfo{OAuthID: "google-merge", Email: "test@test.com"}
				oauth.On("GetSession", mock.Anything, "valid-session").Return(info, nil)
				otp.On("Verify", mock.Anything, "+79999999999", "123456").Return(nil)
				existingUser := domain.NewUserWithPhone("+79999999999", domain.LanguageEnglish)
				existingUser.GoogleID = nil
				userRepo.On("GetByPhone", mock.Anything, "+79999999999").Return(existingUser, nil)
				userRepo.On("LinkGoogle", mock.Anything, existingUser.ID, info).Return(nil)
				oauth.On("DeleteSession", mock.Anything, "valid-session").Return(nil)
				tokens.On("GenerateTokenPair", existingUser.ID, existingUser.Role).
					Return(&domain.TokenPair{AccessToken: "acc", RefreshToken: "ref"}, nil)
			},
			wantErr: false,
		},
		{
			name:      "success - create new user",
			sessionID: "valid-session",
			phone:     "+79999999999",
			code:      "123456",
			mockFn: func(userRepo *mocks.UserRepositoryMock, otp *mocks.OTPManagerMock, oauth *mocks.OAuthManagerMock, tokens *mocks.TokenManagerMock) {
				info := &domain.OAuthInfo{OAuthID: "google-new", Email: "test@test.com"}
				oauth.On("GetSession", mock.Anything, "valid-session").Return(info, nil)
				otp.On("Verify", mock.Anything, "+79999999999", "123456").Return(nil)
				userRepo.On("GetByPhone", mock.Anything, "+79999999999").Return(nil, domain.ErrNotFound)
				userRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)
				oauth.On("DeleteSession", mock.Anything, "valid-session").Return(nil)
				tokens.On("GenerateTokenPair", mock.AnythingOfType("uuid.UUID"), mock.Anything).
					Return(&domain.TokenPair{AccessToken: "acc", RefreshToken: "ref"}, nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			userRepo, phoneOTP, _, oauth, tokens, svc := setupTest()
			tt.mockFn(userRepo, phoneOTP, oauth, tokens)

			user, pair, err := svc.CompleteGoogleLogin(context.Background(), tt.sessionID, tt.phone, tt.code, domain.LanguageEnglish)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.NotNil(t, pair)
			}
			userRepo.AssertExpectations(t)
			phoneOTP.AssertExpectations(t)
			oauth.AssertExpectations(t)
			tokens.AssertExpectations(t)
		})
	}
}

func TestAuthService_LinkGoogle(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name    string
		idToken string
		mockFn  func(userRepo *mocks.UserRepositoryMock, oauth *mocks.OAuthManagerMock)
		wantErr bool
	}{
		{
			name:    "google account already linked to another user",
			idToken: "valid-token",
			mockFn: func(userRepo *mocks.UserRepositoryMock, oauth *mocks.OAuthManagerMock) {
				info := &domain.OAuthInfo{OAuthID: "google-1"}
				oauth.On("VerifyID", mock.Anything, "valid-token").Return(info, nil)
				otherUser := &domain.User{ID: uuid.New()}
				userRepo.On("GetByGoogleID", mock.Anything, "google-1").Return(otherUser, nil)
			},
			wantErr: true,
		},
		{
			name:    "google account already linked to this user (idempotent success)",
			idToken: "valid-token",
			mockFn: func(userRepo *mocks.UserRepositoryMock, oauth *mocks.OAuthManagerMock) {
				info := &domain.OAuthInfo{OAuthID: "google-1"}
				oauth.On("VerifyID", mock.Anything, "valid-token").Return(info, nil)
				sameUser := &domain.User{ID: userID}
				userRepo.On("GetByGoogleID", mock.Anything, "google-1").Return(sameUser, nil)
			},
			wantErr: false,
		},
		{
			name:    "success linking new google account",
			idToken: "valid-token",
			mockFn: func(userRepo *mocks.UserRepositoryMock, oauth *mocks.OAuthManagerMock) {
				info := &domain.OAuthInfo{OAuthID: "google-new"}
				oauth.On("VerifyID", mock.Anything, "valid-token").Return(info, nil)
				userRepo.On("GetByGoogleID", mock.Anything, "google-new").Return(nil, domain.ErrNotFound)
				userRepo.On("LinkGoogle", mock.Anything, userID, info).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			userRepo, _, _, oauth, _, svc := setupTest()
			tt.mockFn(userRepo, oauth)

			err := svc.LinkGoogle(context.Background(), userID, tt.idToken)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			userRepo.AssertExpectations(t)
			oauth.AssertExpectations(t)
		})
	}
}

func TestAuthService_RefreshTokens(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name         string
		refreshToken string
		mockFn       func(userRepo *mocks.UserRepositoryMock, tokens *mocks.TokenManagerMock)
		wantErr      bool
	}{
		{
			name:         "invalid token scope",
			refreshToken: "access-token-string",
			mockFn: func(userRepo *mocks.UserRepositoryMock, tokens *mocks.TokenManagerMock) {
				tokens.On("ValidateToken", "access-token-string").
					Return(&domain.Claims{Scope: domain.ScopeAccessToken}, nil)
			},
			wantErr: true,
		},
		{
			name:         "user not found",
			refreshToken: "valid-refresh",
			mockFn: func(userRepo *mocks.UserRepositoryMock, tokens *mocks.TokenManagerMock) {
				tokens.On("ValidateToken", "valid-refresh").
					Return(&domain.Claims{UserID: userID, Scope: domain.ScopeRefreshToken}, nil)
				userRepo.On("GetByID", mock.Anything, userID).Return(nil, domain.ErrNotFound)
			},
			wantErr: true,
		},
		{
			name:         "success",
			refreshToken: "valid-refresh",
			mockFn: func(userRepo *mocks.UserRepositoryMock, tokens *mocks.TokenManagerMock) {
				tokens.On("ValidateToken", "valid-refresh").
					Return(&domain.Claims{UserID: userID, Scope: domain.ScopeRefreshToken}, nil)
				u := &domain.User{ID: userID, Role: domain.UserRoleCustomer}
				userRepo.On("GetByID", mock.Anything, userID).Return(u, nil)
				tokens.On("GenerateTokenPair", userID, u.Role).
					Return(&domain.TokenPair{AccessToken: "new-acc", RefreshToken: "new-ref"}, nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			userRepo, _, _, _, tokens, svc := setupTest()
			tt.mockFn(userRepo, tokens)

			pair, err := svc.RefreshTokens(context.Background(), tt.refreshToken)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, pair)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, pair)
				assert.Equal(t, "new-acc", pair.AccessToken)
			}
			userRepo.AssertExpectations(t)
			tokens.AssertExpectations(t)
		})
	}
}