package middleware

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"gitlab.com/2foood/backend/internal/domain"
)

type TokenValidator interface {
	ValidateToken(tokenString string) (*domain.Claims, error)
}

type StaffRepo interface {
	GetMember(ctx context.Context, restaurantID, userID uuid.UUID) (*domain.RestaurantStaff, error)
}

type Manager struct {
	log       *slog.Logger
	tokenValidator TokenValidator
	staffRepo StaffRepo
}

func NewManager(log *slog.Logger, tokenValidator TokenValidator, staffRepo StaffRepo) *Manager {
	return &Manager{
		log:       log,
		tokenValidator: tokenValidator,
		staffRepo: staffRepo,
	}
}
