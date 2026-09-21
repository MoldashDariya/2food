package http

import (
	"context"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gitlab.com/2foood/backend/internal/domain"
)

type ReviewService interface {
	Create(ctx context.Context, input domain.CreateReviewInput) (*domain.Review, error)
	GetAll(ctx context.Context) ([]*domain.Review, error)
	GetByRestaurant(ctx context.Context, restID uuid.UUID) ([]*domain.Review, error)
	GetByRestaurantAll(ctx context.Context, restID uuid.UUID) ([]*domain.Review, error)
	GetByUser(ctx context.Context, userID uuid.UUID) ([]*domain.Review, error)
	GetByOrder(ctx context.Context, actorID uuid.UUID, orderID uuid.UUID) (*domain.Review, error)
	GetPhotosByRestaurant(ctx context.Context, restID uuid.UUID) ([]*domain.ReviewPhoto, error)
	AddReply(ctx context.Context, actor *domain.RestaurantStaff, reviewID uuid.UUID, reply string) (*domain.Review, error)
	Hide(ctx context.Context, reviewID uuid.UUID) error
	Unhide(ctx context.Context, reviewID uuid.UUID) error
}

type Handler struct {
	log      *slog.Logger
	service  ReviewService
	validate *validator.Validate
}

func NewHandler(log *slog.Logger, service ReviewService, validate *validator.Validate) *Handler {
	return &Handler{
		log:      log,
		service:  service,
		validate: validate,
	}
}
