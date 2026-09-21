package http

import (
	"context"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gitlab.com/2foood/backend/internal/domain"
)

type ProductService interface {
	Create(ctx context.Context, input domain.CreateProductInput) (*domain.Product, error)
	GetByRestaurant(ctx context.Context, actor *domain.RestaurantStaff, restID uuid.UUID) ([]*domain.Product, error)
	GetActiveByRestaurant(ctx context.Context, restID uuid.UUID) ([]*domain.Product, error)
	GetByID(ctx context.Context, actor *domain.RestaurantStaff, id uuid.UUID) (*domain.Product, error)
	GetCategoriesByRestaurant(ctx context.Context, restID uuid.UUID) ([]*domain.ProductCategory, error)
	Update(ctx context.Context, actor *domain.RestaurantStaff, id uuid.UUID, input domain.UpdateProductInput) (*domain.Product, error)
	UpdateStatus(ctx context.Context, actor *domain.RestaurantStaff, id uuid.UUID, isActive bool) error
	UpdateStock(ctx context.Context, actor *domain.RestaurantStaff, id uuid.UUID, quantity int) error
	Delete(ctx context.Context, actor *domain.RestaurantStaff, id uuid.UUID) error
}

type Handler struct {
	log            *slog.Logger
	productService ProductService
	validate       *validator.Validate
}

func NewHandler(log *slog.Logger, p ProductService, v *validator.Validate) *Handler {
	return &Handler{
		log:            log,
		productService: p,
		validate:       v,
	}
}
