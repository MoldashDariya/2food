package http

import (
	"context"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"gitlab.com/2foood/backend/internal/domain"
	"gitlab.com/2foood/backend/internal/pkg/translator"
)

type CatalogService interface {
	// City
	CreateCity(ctx context.Context, input domain.CreateCityInput) (*domain.City, error)
	GetCities(ctx context.Context) ([]*domain.City, error)
	GetCityByID(ctx context.Context, id int) (*domain.City, error)
	GetCityByName(ctx context.Context, name string) (*domain.City, error)
	DeleteCity(ctx context.Context, id int) error
	// Cuisine Type
	CreateCuisineType(ctx context.Context, input domain.CreateCuisineTypeInput) (*domain.CuisineType, error)
	GetCuisineTypes(ctx context.Context) ([]*domain.CuisineType, error)
	GetCuisineTypeByID(ctx context.Context, id int) (*domain.CuisineType, error)
	GetCuisineTypeByName(ctx context.Context, name string) (*domain.CuisineType, error)
	DeleteCuisineType(ctx context.Context, id int) error
	// Product Category
	CreateProductCategory(ctx context.Context, input domain.CreateProductCategoryInput) (*domain.ProductCategory, error)
	GetProductCategories(ctx context.Context) ([]*domain.ProductCategory, error)
	GetProductCategoryByID(ctx context.Context, id int) (*domain.ProductCategory, error)
	GetProductCategoryByName(ctx context.Context, name string) (*domain.ProductCategory, error)
	DeleteProductCategory(ctx context.Context, id int) error
}

type Handler struct {
	log      *slog.Logger
	service  CatalogService
	validate *validator.Validate
	tr       *translator.Translator
}

func NewHandler(log *slog.Logger, service CatalogService, validate *validator.Validate, tr *translator.Translator) *Handler {
	return &Handler{
		log:      log,
		service:  service,
		validate: validate,
		tr:       tr,
	}
}
