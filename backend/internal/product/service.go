package product

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gitlab.com/2foood/backend/internal/domain"
)

type TransactionManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	GetByRestaurant(ctx context.Context, restID uuid.UUID) ([]*domain.Product, error)
	GetActiveByRestaurant(ctx context.Context, restID uuid.UUID) ([]*domain.Product, error)
	GetCategoriesByRestaurant(ctx context.Context, restID uuid.UUID) ([]*domain.ProductCategory, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	Update(ctx context.Context, product *domain.Product) error
	UpdateStatus(ctx context.Context, id uuid.UUID, isActive bool) error
	UpdateStock(ctx context.Context, id uuid.UUID, quantity int) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type Service struct {
	tm          TransactionManager
	productRepo ProductRepository
}

func NewService(tm TransactionManager, productRepo ProductRepository) *Service {
	return &Service{
		tm:          tm,
		productRepo: productRepo,
	}
}

func (s *Service) Create(ctx context.Context, input domain.CreateProductInput) (*domain.Product, error) {
	const op = "product.Service.Create"

	product := domain.NewProduct(input)

	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return product, nil
}

func (s *Service) GetByRestaurant(ctx context.Context, actor *domain.RestaurantStaff, restID uuid.UUID) ([]*domain.Product, error) {
	const op = "product.Service.GetByRestaurant"

	if actor.RestaurantID != restID {
		return nil, fmt.Errorf("%s: %w", op, domain.ErrPermissionDenied)
	}

	products, err := s.productRepo.GetByRestaurant(ctx, restID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return products, nil
}

func (s *Service) GetActiveByRestaurant(ctx context.Context, restID uuid.UUID) ([]*domain.Product, error) {
	const op = "product.Service.GetActiveByRestaurant"

	products, err := s.productRepo.GetActiveByRestaurant(ctx, restID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return products, nil
}

func (s *Service) GetByID(ctx context.Context, actor *domain.RestaurantStaff, productID uuid.UUID) (*domain.Product, error) {
	const op = "product.Service.GetByID"

	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if actor.RestaurantID != product.RestaurantID {
		return nil, fmt.Errorf("%s: %w", op, domain.ErrPermissionDenied)
	}

	return product, nil
}

func (s *Service) GetCategoriesByRestaurant(ctx context.Context, restID uuid.UUID) ([]*domain.ProductCategory, error) {
	const op = "product.Service.GetCategoriesByRestaurant"

	categories, err := s.productRepo.GetCategoriesByRestaurant(ctx, restID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return categories, nil
}

func (s *Service) Update(ctx context.Context, actor *domain.RestaurantStaff, productID uuid.UUID, input domain.UpdateProductInput) (*domain.Product, error) {
	const op = "product.Service.Update"

	var updatedProduct *domain.Product

	err := s.tm.Do(ctx, func(txCtx context.Context) error {
		product, err := s.productRepo.GetByID(txCtx, productID)
		if err != nil {
			return fmt.Errorf("failed to get product: %w", err)
		}

		if actor.RestaurantID != product.RestaurantID {
			return fmt.Errorf("product does not belong to restaurant: %w", domain.ErrPermissionDenied)
		}

		product.Update(input)

		if err := s.productRepo.Update(txCtx, product); err != nil {
			return fmt.Errorf("failed to update product: %w", err)
		}

		updatedProduct = product
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return updatedProduct, nil
}

func (s *Service) UpdateStatus(ctx context.Context, actor *domain.RestaurantStaff, productID uuid.UUID, isActive bool) error {
	const op = "product.Service.UpdateStatus"

	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}

	if actor.RestaurantID != product.RestaurantID {
		return fmt.Errorf("product does not belong to restaurant: %w", domain.ErrPermissionDenied)
	}

	if err := s.productRepo.UpdateStatus(ctx, productID, isActive); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) UpdateStock(ctx context.Context, actor *domain.RestaurantStaff, productID uuid.UUID, quantity int) error {
	const op = "product.Service.UpdateStock"

	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}

	if actor.RestaurantID != product.RestaurantID {
		return fmt.Errorf("product does not belong to restaurant: %w", domain.ErrPermissionDenied)
	}

	if err := s.productRepo.UpdateStock(ctx, productID, quantity); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, actor *domain.RestaurantStaff, productID uuid.UUID) error {
	const op = "product.Service.Delete"

	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}

	if actor.RestaurantID != product.RestaurantID {
		return fmt.Errorf("product does not belong to restaurant: %w", domain.ErrPermissionDenied)
	}

	if err := s.productRepo.Delete(ctx, productID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
