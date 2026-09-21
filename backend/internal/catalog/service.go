package catalog

import (
	"context"

	"gitlab.com/2foood/backend/internal/domain"
)

type CityRepository interface {
	Create(ctx context.Context, city *domain.City) error
	GetAll(ctx context.Context) ([]*domain.City, error)
	GetByID(ctx context.Context, id int) (*domain.City, error)
	GetByName(ctx context.Context, name string) (*domain.City, error)
	Delete(ctx context.Context, id int) error
}

type CuisineTypeRepository interface {
	Create(ctx context.Context, cuisine *domain.CuisineType) error
	GetAll(ctx context.Context) ([]*domain.CuisineType, error)
	GetByID(ctx context.Context, id int) (*domain.CuisineType, error)
	GetByName(ctx context.Context, name string) (*domain.CuisineType, error)
	Delete(ctx context.Context, id int) error
}

type ProductCategoryRepository interface {
	Create(ctx context.Context, category *domain.ProductCategory) error
	GetAll(ctx context.Context) ([]*domain.ProductCategory, error)
	GetByID(ctx context.Context, id int) (*domain.ProductCategory, error)
	GetByName(ctx context.Context, name string) (*domain.ProductCategory, error)
	Delete(ctx context.Context, id int) error
}

type Service struct {
	cityRepo            CityRepository
	cuisineTypeRepo     CuisineTypeRepository
	productCategoryRepo ProductCategoryRepository
}

func NewService(cityRepo CityRepository, cuisineTypeRepo CuisineTypeRepository, productCategoryRepo ProductCategoryRepository) *Service {
	return &Service{cityRepo: cityRepo, cuisineTypeRepo: cuisineTypeRepo, productCategoryRepo: productCategoryRepo}
}
