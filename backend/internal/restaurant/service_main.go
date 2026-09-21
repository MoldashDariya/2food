package restaurant

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gitlab.com/2foood/backend/internal/domain"
)

func (s *Service) Create(ctx context.Context, ownerID uuid.UUID, input domain.CreateRestaurantInput) (*domain.Restaurant, error) {
	const op = "restaurant.Service.Create"

	rest := domain.NewRestaurant(input)

	// create restaurant and add owner to staff
	err := s.tm.Do(ctx, func(txCtx context.Context) error {
		if err := s.restRepo.Create(txCtx, rest); err != nil {
			return fmt.Errorf("failed to create restaurant: %w", err)
		}

		ownerStaff := domain.NewRestaurantStaff(rest.ID, ownerID, domain.StaffRoleOwner)
		if err := s.staffRepo.Add(txCtx, ownerStaff); err != nil {
			return fmt.Errorf("failed to add owner to staff: %w", err)
		}

		ownerUser, err := s.userRepo.GetByID(txCtx, ownerID)
		if err != nil {
			return fmt.Errorf("failed to get owner user: %w", err)
		}

		if ownerUser.Role == domain.UserRoleAdmin {
			return nil // skip setting role to partner for admin
		}

		if ownerUser.Role == domain.UserRolePartner {
			return nil // skip setting role to partner for partner
		}

		if err := s.userRepo.SetRole(txCtx, ownerID, domain.UserRolePartner); err != nil {
			return fmt.Errorf("failed to set owner role to partner: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return rest, nil
}

func (s *Service) GetAll(ctx context.Context) ([]*domain.Restaurant, error) {
	const op = "restaurant.Service.GetAll"

	rests, err := s.restRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return rests, nil
}

func (s *Service) GetByPartner(ctx context.Context, userID uuid.UUID) ([]*domain.Restaurant, error) {
	const op = "restaurant.Service.GetByPartner"

	rests, err := s.restRepo.GetByPartner(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return rests, nil
}

func (s *Service) GetByID(ctx context.Context, restID uuid.UUID) (*domain.Restaurant, error) {
	const op = "restaurant.Service.GetByID"

	rest, err := s.restRepo.GetByID(ctx, restID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return rest, nil
}

func (s *Service) UpdateInfo(ctx context.Context, actor *domain.RestaurantStaff, restID uuid.UUID, input domain.UpdateRestaurantInput) (*domain.Restaurant, error) {
	const op = "restaurant.Service.UpdateInfo"

	// check if actor is owner of the restaurant
	if actor.Role != domain.StaffRoleOwner {
		return nil, fmt.Errorf("%s: %w", op, domain.ErrPermissionDenied)
	}

	var updatedRest *domain.Restaurant

	// update restaurant info
	err := s.tm.Do(ctx, func(txCtx context.Context) error {
		rest, err := s.restRepo.GetByID(txCtx, restID)
		if err != nil {
			return fmt.Errorf("failed to get restaurant: %w", err)
		}

		rest.Update(input)

		if err := s.restRepo.UpdateInfo(txCtx, rest); err != nil {
			return fmt.Errorf("failed to update restaurant info: %w", err)
		}

		updatedRest = rest
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return updatedRest, nil
}

func (s *Service) UpdateCuisines(ctx context.Context, actor *domain.RestaurantStaff, restID uuid.UUID, cuisineIDs []int) error {
	const op = "restaurant.Service.UpdateCuisines"

	// only owner can update cuisines
	if actor.Role != domain.StaffRoleOwner {
		return fmt.Errorf("%s: %w", op, domain.ErrPermissionDenied)
	}

	err := s.tm.Do(ctx, func(txCtx context.Context) error {
		// transactionally update cuisines
		if err := s.restRepo.UpdateCuisines(txCtx, restID, cuisineIDs); err != nil {
			return fmt.Errorf("failed to update cuisines: %w", err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, restID uuid.UUID) error {
	const op = "restaurant.Service.Delete"

	if err := s.restRepo.Delete(ctx, restID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
