package catalog

import (
	"context"
	"fmt"

	"gitlab.com/2foood/backend/internal/domain"
)

func (s *Service) CreateCuisineType(ctx context.Context, input domain.CreateCuisineTypeInput) (*domain.CuisineType, error) {
	const op = "catalog.Service.CreateCuisineType"

	cuisine := domain.NewCuisineType(input)

	if err := s.cuisineTypeRepo.Create(ctx, cuisine); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return cuisine, nil
}

func (s *Service) GetCuisineTypes(ctx context.Context) ([]*domain.CuisineType, error) {
	const op = "catalog.Service.GetCuisineTypes"

	cuisines, err := s.cuisineTypeRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return cuisines, nil
}

func (s *Service) GetCuisineTypeByID(ctx context.Context, id int) (*domain.CuisineType, error) {
	const op = "catalog.Service.GetCuisineTypeByID"

	cuisine, err := s.cuisineTypeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return cuisine, nil
}

func (s *Service) GetCuisineTypeByName(ctx context.Context, name string) (*domain.CuisineType, error) {
	const op = "catalog.Service.GetCuisineTypeByName"

	cuisine, err := s.cuisineTypeRepo.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return cuisine, nil
}

func (s *Service) DeleteCuisineType(ctx context.Context, id int) error {
	const op = "catalog.Service.DeleteCuisineType"

	if err := s.cuisineTypeRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
