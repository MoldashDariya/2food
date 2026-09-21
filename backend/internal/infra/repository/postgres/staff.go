package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gitlab.com/2foood/backend/internal/domain"
)

type StaffRepo struct {
	BaseRepo
}

func NewStaffRepo(db *sql.DB) *StaffRepo {
	return &StaffRepo{BaseRepo: NewBaseRepo(db)}
}

func (r *StaffRepo) Add(ctx context.Context, staff *domain.RestaurantStaff) error {
	const op = "postgres.StaffRepo.Add"

	q := r.Queryer(ctx)

	query := `
		INSERT INTO restaurant_staff (id, restaurant_id, user_id, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := q.ExecContext(ctx, query,
		staff.ID,
		staff.RestaurantID,
		staff.UserID,
		staff.Role,
		staff.CreatedAt,
		staff.UpdatedAt,
	)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == pgUniqueErrCode {
			return fmt.Errorf("%s: %w", op, domain.ErrAlreadyExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *StaffRepo) GetByRestaurant(ctx context.Context, restID uuid.UUID) ([]*domain.RestaurantStaff, error) {
	const op = "postgres.StaffRepo.GetByRestaurant"

	q := r.Queryer(ctx)

	query := `
		SELECT id, restaurant_id, user_id, role, created_at, updated_at
		FROM restaurant_staff
		WHERE restaurant_id = $1
	`

	rows, err := q.QueryContext(ctx, query, restID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var list []*domain.RestaurantStaff
	for rows.Next() {
		s, err := r.scanRow(rows)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		list = append(list, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return list, nil
}

func (r *StaffRepo) GetMember(ctx context.Context, restID, userID uuid.UUID) (*domain.RestaurantStaff, error) {
	const op = "postgres.StaffRepo.GetMember"

	q := r.Queryer(ctx)

	query := `
		SELECT id, restaurant_id, user_id, role, created_at, updated_at
		FROM restaurant_staff
		WHERE restaurant_id = $1 AND user_id = $2
	`

	s, err := r.scanRow(q.QueryRowContext(ctx, query, restID, userID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return s, nil
}

func (r *StaffRepo) GetOwner(ctx context.Context, restID uuid.UUID) (*domain.RestaurantStaff, error) {
	const op = "postgres.StaffRepo.GetOwner"

	q := r.Queryer(ctx)

	query := `
		SELECT id, restaurant_id, user_id, role, created_at, updated_at
		FROM restaurant_staff
		WHERE restaurant_id = $1 AND role = 'owner'
	`

	s, err := r.scanRow(q.QueryRowContext(ctx, query, restID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return s, nil
}

func (r *StaffRepo) CountRestaurants(ctx context.Context, userID uuid.UUID) (int, error) {
	const op = "postgres.StaffRepo.CountRestaurants"

	q := r.Queryer(ctx)

	query := `
		SELECT COUNT(*) FROM restaurant_staff WHERE user_id = $1
	`

	var count int
	err := q.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return count, nil
}

func (r *StaffRepo) SetRole(ctx context.Context, restID, userID uuid.UUID, role domain.StaffRole) error {
	const op = "postgres.StaffRepo.UpdateRole"

	q := r.Queryer(ctx)

	query := `
		UPDATE restaurant_staff SET role = $1, updated_at = NOW()
		WHERE restaurant_id = $2 AND user_id = $3
	`

	res, err := q.ExecContext(ctx, query, role, restID, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if count == 0 {
		return fmt.Errorf("%s: %w", op, domain.ErrNotFound)
	}

	return nil
}

func (r *StaffRepo) Remove(ctx context.Context, restID, userID uuid.UUID) error {
	const op = "postgres.StaffRepo.Remove"

	q := r.Queryer(ctx)

	query := `
		DELETE FROM restaurant_staff 
		WHERE restaurant_id = $1 AND user_id = $2
	`

	res, err := q.ExecContext(ctx, query, restID, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if count == 0 {
		return fmt.Errorf("%s: %w", op, domain.ErrNotFound)
	}

	return nil
}

func (r *StaffRepo) scanRow(s scanner) (*domain.RestaurantStaff, error) {
	var staff domain.RestaurantStaff
	err := s.Scan(
		&staff.ID,
		&staff.RestaurantID,
		&staff.UserID,
		&staff.Role,
		&staff.CreatedAt,
		&staff.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &staff, nil
}
