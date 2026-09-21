package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gitlab.com/2foood/backend/internal/domain"
)

// small interfaces defined where they're used (mocks are in ./mocks)
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByPhone(ctx context.Context, phone string) (*domain.User, error)
	GetByGoogleID(ctx context.Context, googleID string) (*domain.User, error)
	SetPhone(ctx context.Context, id uuid.UUID, phone string) error
	LinkGoogle(ctx context.Context, userID uuid.UUID, info *domain.OAuthInfo) error
}

type OTPManager interface {
	Send(ctx context.Context, to string) error
	Verify(ctx context.Context, to, code string) error
}

type OAuthManager interface {
	VerifyID(ctx context.Context, token string) (*domain.OAuthInfo, error)
	CreateSession(ctx context.Context, info *domain.OAuthInfo) (string, error)
	GetSession(ctx context.Context, sessionID string) (*domain.OAuthInfo, error)
	DeleteSession(ctx context.Context, sessionID string) error
}

type TokenManager interface {
	GenerateTokenPair(userID uuid.UUID, role domain.UserRole) (*domain.TokenPair, error)
	ValidateToken(tokenString string) (*domain.Claims, error)
}

type Service struct {
	userRepo    UserRepository
	phoneOTP    OTPManager
	emailOTP    OTPManager
	googleOAuth OAuthManager
	tokens      TokenManager
}

func NewService(userRepo UserRepository, phoneOTP OTPManager, emailOTP OTPManager, oauth OAuthManager, tokens TokenManager) *Service {
	return &Service{
		userRepo:    userRepo,
		phoneOTP:    phoneOTP,
		emailOTP:    emailOTP,
		googleOAuth: oauth,
		tokens:      tokens,
	}
}

func (s *Service) SendPhoneOTP(ctx context.Context, phone string) error {
	const op = "auth.Service.SendPhoneOTP"

	err := s.phoneOTP.Send(ctx, phone)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) SendEmailOTP(ctx context.Context, email string) error {
	const op = "auth.Service.SendEmailOTP"

	err := s.emailOTP.Send(ctx, email)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// login and sign up are the same thing here: if nobody has this phone yet, we create the user
func (s *Service) LoginWithPhone(ctx context.Context, phone, code string, lang domain.Language) (*domain.User, *domain.TokenPair, error) {
	const op = "auth.Service.LoginWithPhone"

	if err := s.phoneOTP.Verify(ctx, phone, code); err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	user, err := s.userRepo.GetByPhone(ctx, phone)
	if err != nil {
		if !errors.Is(err, domain.ErrNotFound) {
			return nil, nil, fmt.Errorf("%s: %w", op, err)
		}

		user = domain.NewUserWithPhone(phone, lang)
		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	tokens, err := s.tokens.GenerateTokenPair(user.ID, user.Role)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, tokens, nil
}

// google login has two outcomes:
// 1) known user with a phone -> tokens right away
// 2) new user (or no phone yet) -> we need a phone number first. google info is kept
//    in a session (redis) and the app finishes with CompleteGoogleLogin
func (s *Service) LoginWithGoogle(ctx context.Context, idToken string) (*domain.OAuthResult, error) {
	const op = "auth.Service.LoginWithGoogle"

	info, err := s.googleOAuth.VerifyID(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	
	user, err := s.userRepo.GetByGoogleID(ctx, info.OAuthID)
	if err != nil {
		if !errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		// create session for new user
		sessionID, err := s.googleOAuth.CreateSession(ctx, info)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		return &domain.OAuthResult{
			SessionID: sessionID,
			NeedPhone: true,
		}, nil
	}

	// if user already exists but doesn't have a phone number
	if user.Phone == nil || *user.Phone == "" {
		sessionID, err := s.googleOAuth.CreateSession(ctx, info)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		return &domain.OAuthResult{
			User:      user,
			SessionID: sessionID,
			NeedPhone: true,
		}, nil
	}

	tokens, err := s.tokens.GenerateTokenPair(user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &domain.OAuthResult{
		User:      user,
		Tokens:    tokens,
		NeedPhone: false,
	}, nil
}

// second step of the google flow: user confirms the phone with an otp code,
// then we create / link the account
func (s *Service) CompleteGoogleLogin(ctx context.Context, sessionID, phone string, code string, lang domain.Language) (*domain.User, *domain.TokenPair, error) {
	const op = "auth.Service.CompleteGoogleLogin"

	info, err := s.googleOAuth.GetSession(ctx, sessionID)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := s.phoneOTP.Verify(ctx, phone, code); err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	user, err := s.userRepo.GetByPhone(ctx, phone)

	if err == nil {
		// if user already exists but has a different google ID, return an error
		if user.GoogleID != nil && *user.GoogleID != "" && *user.GoogleID != info.OAuthID {
			return nil, nil, fmt.Errorf("%s: %w", op, domain.ErrAlreadyExists)
		}
		// if user already exists and has the same google ID, link the phone number
		if err := s.userRepo.LinkGoogle(ctx, user.ID, info); err != nil {
			return nil, nil, fmt.Errorf("%s: %w", op, err)
		}
	} else if errors.Is(err, domain.ErrNotFound) {
		// if user doesn't exist, create a new user
		user = domain.NewUserWithPhone(phone, lang)
		user.BindGoogleData(*info)

		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, nil, fmt.Errorf("%s: %w", op, err)
		}
	} else {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	// session is one-time use
	_ = s.googleOAuth.DeleteSession(ctx, sessionID)

	tokens, err := s.tokens.GenerateTokenPair(user.ID, user.Role)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, tokens, nil
}

func (s *Service) LinkGoogle(ctx context.Context, userID uuid.UUID, idToken string) error {
	const op = "auth.Service.LinkGoogle"

	info, err := s.googleOAuth.VerifyID(ctx, idToken)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	existingUser, err := s.userRepo.GetByGoogleID(ctx, info.OAuthID)
	if err == nil {
		if existingUser.ID != userID {
			return fmt.Errorf("%s: %w", op, domain.ErrAlreadyExists)
		}
		return nil
	}

	if !errors.Is(err, domain.ErrNotFound) {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := s.userRepo.LinkGoogle(ctx, userID, info); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// the scope check stops someone from sending an access token here.
// role is re-read from the db so role changes apply after the next refresh
func (s *Service) RefreshTokens(ctx context.Context, refreshToken string) (*domain.TokenPair, error) {
	const op = "auth.Service.RefreshTokens"

	claims, err := s.tokens.ValidateToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if claims.Scope != domain.ScopeRefreshToken {
		return nil, fmt.Errorf("%s: %w", op, domain.ErrInvalidToken)
	}

	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	tokens, err := s.tokens.GenerateTokenPair(claims.UserID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return tokens, nil
}
