package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gitlab.com/2foood/backend/internal/domain"
)

type TransactionManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type UserRepository interface {
	GetByID(ctx context.Context, userID uuid.UUID) (*domain.User, error)
	GetByPhone(ctx context.Context, phone string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	SetPhone(ctx context.Context, userID uuid.UUID, phone string) error
	SetEmail(ctx context.Context, userID uuid.UUID, email string) error
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, userID uuid.UUID) error
}

type DeviceTokenRepository interface {
	Register(ctx context.Context, dt *domain.DeviceToken) error
	Delete(ctx context.Context, userID uuid.UUID, token string) error
}

type OTPManager interface {
	Send(ctx context.Context, to string) error
	Verify(ctx context.Context, to, code string) error
}

type Service struct {
	tm              TransactionManager
	userRepo        UserRepository
	deviceTokenRepo DeviceTokenRepository
	phoneOTP        OTPManager
	emailOTP        OTPManager
}

func NewService(tm TransactionManager, userRepo UserRepository, deviceTokenRepo DeviceTokenRepository, phoneOTP OTPManager, emailOTP OTPManager) *Service {
	return &Service{
		tm:              tm,
		userRepo:        userRepo,
		deviceTokenRepo: deviceTokenRepo,
		phoneOTP:        phoneOTP,
		emailOTP:        emailOTP,
	}
}

func (s *Service) Get(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	const op = "user.Service.Get"

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Service) Update(ctx context.Context, userID uuid.UUID, input domain.UpdateUserInput) (*domain.User, error) {
	const op = "user.Service.Update"

	var updatedUser *domain.User

	err := s.tm.Do(ctx, func(txCtx context.Context) error {
		user, err := s.userRepo.GetByID(txCtx, userID)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}

		user.Update(input)

		if err := s.userRepo.Update(txCtx, user); err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}

		updatedUser = user
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return updatedUser, nil
}

func (s *Service) UpdatePhone(ctx context.Context, userID uuid.UUID, newPhone string, code string) error {
	const op = "user.Service.UpdatePhone"

	if err := s.phoneOTP.Verify(ctx, newPhone, code); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err := s.userRepo.GetByPhone(ctx, newPhone)
	if err == nil {
		return fmt.Errorf("%s: %w", op, domain.ErrAlreadyExists)
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := s.userRepo.SetPhone(ctx, userID, newPhone); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) UpdateEmail(ctx context.Context, userID uuid.UUID, email string, code string) error {
	const op = "user.Service.UpdateEmail"

	if err := s.emailOTP.Verify(ctx, email, code); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil {
		return fmt.Errorf("%s: %w", op, domain.ErrAlreadyExists)
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := s.userRepo.SetEmail(ctx, userID, email); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, userID uuid.UUID) error {
	const op = "user.Service.Delete"

	if err := s.userRepo.Delete(ctx, userID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
