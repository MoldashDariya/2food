package http

import (
	"context"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gitlab.com/2foood/backend/internal/domain"
)

type RestaurantService interface {
	Create(ctx context.Context, ownerID uuid.UUID, input domain.CreateRestaurantInput) (*domain.Restaurant, error)
	GetAll(ctx context.Context) ([]*domain.Restaurant, error)
	GetByPartner(ctx context.Context, userID uuid.UUID) ([]*domain.Restaurant, error)
	GetByID(ctx context.Context, restID uuid.UUID) (*domain.Restaurant, error)
	UpdateInfo(ctx context.Context, actor *domain.RestaurantStaff, restID uuid.UUID, input domain.UpdateRestaurantInput) (*domain.Restaurant, error)
	UpdateCuisines(ctx context.Context, actor *domain.RestaurantStaff, restID uuid.UUID, cuisineIDs []int) error
	Delete(ctx context.Context, restID uuid.UUID) error
}

type StaffService interface {
	AddStaff(ctx context.Context, actor *domain.RestaurantStaff, restID, targetID uuid.UUID, role domain.StaffRole) error
	GetStaff(ctx context.Context, restID uuid.UUID) ([]*domain.RestaurantStaff, error)
	GetStaffMember(ctx context.Context, restID, userID uuid.UUID) (*domain.RestaurantStaff, error)
	UpdateStaffRole(ctx context.Context, actor *domain.RestaurantStaff, restID, targetID uuid.UUID, role domain.StaffRole) error
	RemoveStaff(ctx context.Context, actor *domain.RestaurantStaff, restID, targetID uuid.UUID) error
	TransferOwnership(ctx context.Context, actor *domain.RestaurantStaff, restID, targetID uuid.UUID) error
	AssignOwner(ctx context.Context, restID uuid.UUID, targetID uuid.UUID) error
}

type Handler struct {
	log          *slog.Logger
	restService  RestaurantService
	staffService StaffService
	validate     *validator.Validate
}

func NewHandler(log *slog.Logger, r RestaurantService, s StaffService, v *validator.Validate) *Handler {
	return &Handler{
		log:          log,
		restService:  r,
		staffService: s,
		validate:     v,
	}
}
