package restaurant

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gitlab.com/2foood/backend/internal/domain"
)

func (s *Service) AddStaff(ctx context.Context, actor *domain.RestaurantStaff, restID, targetID uuid.UUID, role domain.StaffRole) error {
	const op = "restaurant.Service.AddStaff"

	// check if actor is owner of the restaurant
	if actor.Role != domain.StaffRoleOwner {
		return fmt.Errorf("%s: %w", op, domain.ErrPermissionDenied)
	}

	// check if target user is already a staff member
	_, err := s.staffRepo.GetMember(ctx, restID, targetID)
	if err == nil {
		return fmt.Errorf("%s: %w", op, domain.ErrAlreadyExists)
	}

	// owner cannot be added as staff
	if role == domain.StaffRoleOwner {
		return fmt.Errorf("%s: %w", op, domain.ErrPermissionDenied)
	}

	err = s.tm.Do(ctx, func(txCtx context.Context) error {
		staff := domain.NewRestaurantStaff(restID, targetID, role)

		if err := s.staffRepo.Add(txCtx, staff); err != nil {
			return fmt.Errorf("failed to add staff: %w", err)
		}

		targetUser, err := s.userRepo.GetByID(txCtx, targetID)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}

		if targetUser.Role == domain.UserRoleAdmin {
			return nil // skip setting role to partner for admin
		}

		if err := s.userRepo.SetRole(txCtx, targetID, domain.UserRolePartner); err != nil {
			return fmt.Errorf("failed to set user role to partner: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) GetStaff(ctx context.Context, restID uuid.UUID) ([]*domain.RestaurantStaff, error) {
	const op = "restaurant.Service.GetStaff"

	staff, err := s.staffRepo.GetByRestaurant(ctx, restID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return staff, nil
}

func (s *Service) GetStaffMember(ctx context.Context, restID, userID uuid.UUID) (*domain.RestaurantStaff, error) {
	const op = "restaurant.Service.GetStaffMember"

	staff, err := s.staffRepo.GetMember(ctx, restID, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return staff, nil
}

func (s *Service) UpdateStaffRole(ctx context.Context, actor *domain.RestaurantStaff, restID, targetID uuid.UUID, role domain.StaffRole) error {
	const op = "restaurant.Service.UpdateStaffRole"

	// check if actor is owner of the restaurant
	if actor.Role != domain.StaffRoleOwner {
		return fmt.Errorf("%s: %w", op, domain.ErrPermissionDenied)
	}

	// check if target user is a staff member
	target, err := s.staffRepo.GetMember(ctx, restID, targetID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	// owner cannot be updated
	if target.Role == domain.StaffRoleOwner {
		return fmt.Errorf("%s: %w", op, domain.ErrPermissionDenied)
	}

	// update role to owner is not allowed
	if role == domain.StaffRoleOwner {
		return fmt.Errorf("%s: %w", op, domain.ErrPermissionDenied)
	}

	if err := s.staffRepo.SetRole(ctx, restID, targetID, role); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) TransferOwnership(ctx context.Context, actor *domain.RestaurantStaff, restID, targetID uuid.UUID) error {
	const op = "restaurant.Service.TransferOwnership"

	// check if actor is admin
	actorUser, err := s.userRepo.GetByID(ctx, actor.UserID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// admin must use assign owner endpoint
	if actorUser.Role == domain.UserRoleAdmin {
		return fmt.Errorf("%s: %w", op, domain.ErrPermissionDenied)
	}

	// check if actor is owner of the restaurant
	if actor.Role != domain.StaffRoleOwner {
		return fmt.Errorf("%s: %w", op, domain.ErrPermissionDenied)
	}

	if actor.UserID == targetID {
		return fmt.Errorf("%s: %w", op, domain.ErrPermissionDenied)
	}

	err = s.tm.Do(ctx, func(txCtx context.Context) error {
		if err := s.staffRepo.SetRole(txCtx, restID, actor.UserID, domain.StaffRoleManager); err != nil {
			return fmt.Errorf("failed to demote current owner: %w", err)
		}

		err := s.staffRepo.SetRole(txCtx, restID, targetID, domain.StaffRoleOwner)

		// if target user is not a staff member, add them as owner
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				newOwner := domain.NewRestaurantStaff(restID, targetID, domain.StaffRoleOwner)
				if err := s.staffRepo.Add(txCtx, newOwner); err != nil {
					return fmt.Errorf("failed to add new owner: %w", err)
				}
			} else {
				return fmt.Errorf("failed to set new owner role: %w", err)
			}
		}

		targetUser, err := s.userRepo.GetByID(txCtx, targetID)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}

		if targetUser.Role == domain.UserRoleAdmin {
			return nil // skip setting role to partner for admin
		}

		if err := s.userRepo.SetRole(txCtx, targetID, domain.UserRolePartner); err != nil {
			return fmt.Errorf("failed to set user role to partner: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) AssignOwner(ctx context.Context, restID uuid.UUID, targetID uuid.UUID) error {
	const op = "restaurant.Service.AssignOwner"

	err := s.tm.Do(ctx, func(txCtx context.Context) error {
		// check if restaurant exists
		if _, err := s.restRepo.GetByID(txCtx, restID); err != nil {
			return fmt.Errorf("failed to get restaurant: %w", err)
		}

		// find current owner and demote if necessary
		owner, err := s.staffRepo.GetOwner(txCtx, restID)
		if err != nil && !errors.Is(err, domain.ErrNotFound) {
			return fmt.Errorf("failed to get owner: %w", err)
		}
		if err == nil && owner.UserID != targetID {
			if err := s.staffRepo.SetRole(txCtx, restID, owner.UserID, domain.StaffRoleManager); err != nil {
				return fmt.Errorf("failed to demote current owner: %w", err)
			}
		}

		// assign new owner
		if _, err := s.staffRepo.GetMember(txCtx, restID, targetID); err == nil {
			if err := s.staffRepo.SetRole(txCtx, restID, targetID, domain.StaffRoleOwner); err != nil {
				return fmt.Errorf("failed to set new owner: %w", err)
			}
		} else {
			if !errors.Is(err, domain.ErrNotFound) {
				return fmt.Errorf("failed to check target membership: %w", err)
			}
			newOwner := domain.NewRestaurantStaff(restID, targetID, domain.StaffRoleOwner)
			if err := s.staffRepo.Add(txCtx, newOwner); err != nil {
				return fmt.Errorf("failed to add new owner: %w", err)
			}
		}

		// set user role to partner if necessary
		targetUser, err := s.userRepo.GetByID(txCtx, targetID)
		if err != nil {
			return fmt.Errorf("failed to get target user: %w", err)
		}

		if targetUser.Role != domain.UserRoleAdmin && targetUser.Role != domain.UserRolePartner {
			if err := s.userRepo.SetRole(txCtx, targetID, domain.UserRolePartner); err != nil {
				return fmt.Errorf("failed to set user role to partner: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) RemoveStaff(ctx context.Context, actor *domain.RestaurantStaff, restID, targetID uuid.UUID) error {
	const op = "restaurant.Service.RemoveStaff"

	// check if actor is owner of the restaurant
	if actor.Role != domain.StaffRoleOwner {
		return fmt.Errorf("%s: %w", op, domain.ErrPermissionDenied)
	}

	// check if target user is a staff member
	target, err := s.staffRepo.GetMember(ctx, restID, targetID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	// owner cannot be removed
	if target.Role == domain.StaffRoleOwner {
		return fmt.Errorf("%s: %w", op, domain.ErrPermissionDenied)
	}

	err = s.tm.Do(ctx, func(txCtx context.Context) error {
		if err := s.staffRepo.Remove(txCtx, restID, targetID); err != nil {
			return fmt.Errorf("failed to remove staff: %w", err)
		}

		count, err := s.staffRepo.CountRestaurants(txCtx, targetID)
		if err != nil {
			return fmt.Errorf("failed to count restaurants: %w", err)
		}

		if count == 0 {
			targetUser, err := s.userRepo.GetByID(txCtx, targetID)
			if err != nil {
				return fmt.Errorf("failed to get user: %w", err)
			}

			if targetUser.Role == domain.UserRoleAdmin {
				return nil // skip setting role to customer for admin
			}

			if err := s.userRepo.SetRole(txCtx, targetID, domain.UserRoleCustomer); err != nil {
				return fmt.Errorf("failed to set user role to customer: %w", err)
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
