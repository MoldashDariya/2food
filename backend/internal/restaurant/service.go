package restaurant

import (
	"context"

	"github.com/google/uuid"
	"gitlab.com/2foood/backend/internal/domain"
)

type TransactionManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	SetRole(ctx context.Context, id uuid.UUID, role domain.UserRole) error
}

type RestaurantRepository interface {
	Create(ctx context.Context, rest *domain.Restaurant) error
	GetAll(ctx context.Context) ([]*domain.Restaurant, error)
	GetByPartner(ctx context.Context, userID uuid.UUID) ([]*domain.Restaurant, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Restaurant, error)
	UpdateInfo(ctx context.Context, rest *domain.Restaurant) error
	UpdateCuisines(ctx context.Context, restID uuid.UUID, cuisineIDs []int) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type StaffRepository interface {
	Add(ctx context.Context, staff *domain.RestaurantStaff) error
	GetByRestaurant(ctx context.Context, restID uuid.UUID) ([]*domain.RestaurantStaff, error)
	GetMember(ctx context.Context, restID, userID uuid.UUID) (*domain.RestaurantStaff, error)
	GetOwner(ctx context.Context, restID uuid.UUID) (*domain.RestaurantStaff, error)
	CountRestaurants(ctx context.Context, userID uuid.UUID) (int, error)
	SetRole(ctx context.Context, restID, userID uuid.UUID, role domain.StaffRole) error
	Remove(ctx context.Context, restID, userID uuid.UUID) error
}

type Service struct {
	tm        TransactionManager
	userRepo  UserRepository
	restRepo  RestaurantRepository
	staffRepo StaffRepository
}

func NewService(
	tm TransactionManager, userRepo UserRepository, restRepo RestaurantRepository, staffRepo StaffRepository) *Service {
	return &Service{
		tm:        tm,
		userRepo:  userRepo,
		restRepo:  restRepo,
		staffRepo: staffRepo,
	}
}
