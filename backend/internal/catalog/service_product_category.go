package catalog

import (
	"context"
	"fmt"

	"gitlab.com/2foood/backend/internal/domain"
)

func (s *Service) CreateProductCategory(ctx context.Context, input domain.CreateProductCategoryInput) (*domain.ProductCategory, error) {
	const op = "catalog.Service.CreateProductCategory"

	category := domain.NewProductCategory(input)

	if err := s.productCategoryRepo.Create(ctx, category); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return category, nil
}

func (s *Service) GetProductCategories(ctx context.Context) ([]*domain.ProductCategory, error) {
	const op = "catalog.Service.GetProductCategories"

	categories, err := s.productCategoryRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return categories, nil
}

func (s *Service) GetProductCategoryByID(ctx context.Context, id int) (*domain.ProductCategory, error) {
	const op = "catalog.Service.GetProductCategoryByID"

	category, err := s.productCategoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return category, nil
}

func (s *Service) GetProductCategoryByName(ctx context.Context, name string) (*domain.ProductCategory, error) {
	const op = "catalog.Service.GetProductCategoryByName"

	category, err := s.productCategoryRepo.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return category, nil
}

func (s *Service) DeleteProductCategory(ctx context.Context, id int) error {
	const op = "catalog.Service.DeleteProductCategory"

	if err := s.productCategoryRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
