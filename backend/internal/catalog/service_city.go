package catalog

import (
	"context"
	"fmt"

	"gitlab.com/2foood/backend/internal/domain"
)

func (s *Service) CreateCity(ctx context.Context, input domain.CreateCityInput) (*domain.City, error) {
	const op = "catalog.Service.CreateCity"

	city := domain.NewCity(input)

	if err := s.cityRepo.Create(ctx, city); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return city, nil
}

func (s *Service) GetCities(ctx context.Context) ([]*domain.City, error) {
	const op = "catalog.Service.GetCities"

	cities, err := s.cityRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return cities, nil
}

func (s *Service) GetCityByID(ctx context.Context, id int) (*domain.City, error) {
	const op = "catalog.Service.GetCityByID"

	city, err := s.cityRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return city, nil
}

func (s *Service) GetCityByName(ctx context.Context, name string) (*domain.City, error) {
	const op = "catalog.Service.GetCityByName"

	city, err := s.cityRepo.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return city, nil
}

func (s *Service) DeleteCity(ctx context.Context, id int) error {
	const op = "catalog.Service.DeleteCity"

	if err := s.cityRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
