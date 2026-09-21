package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gitlab.com/2foood/backend/internal/domain"
)


func (s *Service) RegisterDeviceToken(ctx context.Context, input domain.RegisterDeviceTokenInput) error {
	const op = "user.Service.RegisterDeviceToken"

	dt := domain.NewDeviceToken(input)

	if err := s.deviceTokenRepo.Register(ctx, dt); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) DeleteDeviceToken(ctx context.Context, userID uuid.UUID, token string) error {
	const op = "user.Service.DeleteDeviceToken"

	if err := s.deviceTokenRepo.Delete(ctx, userID, token); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
