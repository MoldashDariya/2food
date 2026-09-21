package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"gitlab.com/2foood/backend/internal/domain"
)

type CuisineTypeRepo struct {
	BaseRepo
}

func NewCuisineTypeRepo(db *sql.DB) *CuisineTypeRepo {
	return &CuisineTypeRepo{BaseRepo: NewBaseRepo(db)}
}

func (r *CuisineTypeRepo) Create(ctx context.Context, cuisine *domain.CuisineType) error {
	const op = "repository.CuisineTypeRepo.Create"

	q := r.Queryer(ctx)

	query := `
		INSERT INTO cuisine_types (name, label)
		VALUES ($1, $2)
		RETURNING id
	`

	err := q.QueryRowContext(ctx, query, cuisine.Name, cuisine.Label).Scan(&cuisine.ID)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == pgUniqueErrCode {
			return fmt.Errorf("%s: %w", op, domain.ErrAlreadyExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *CuisineTypeRepo) GetAll(ctx context.Context) ([]*domain.CuisineType, error) {
	const op = "repository.CuisineTypeRepo.GetAll"

	q := r.Queryer(ctx)

	query := `
		SELECT id, name, label
		FROM cuisine_types
		ORDER BY name
	`

	rows, err := q.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	list := []*domain.CuisineType{}
	for rows.Next() {
		c, err := r.scanRow(rows)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		list = append(list, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return list, nil
}

func (r *CuisineTypeRepo) GetByID(ctx context.Context, id int) (*domain.CuisineType, error) {
	const op = "repository.CuisineTypeRepo.GetByID"

	q := r.Queryer(ctx)

	query := `
		SELECT id, name, label
		FROM cuisine_types
		WHERE id = $1
	`

	c, err := r.scanRow(q.QueryRowContext(ctx, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return c, nil
}

func (r *CuisineTypeRepo) GetByName(ctx context.Context, name string) (*domain.CuisineType, error) {
	const op = "repository.CuisineTypeRepo.GetByName"

	q := r.Queryer(ctx)

	query := `
		SELECT id, name, label
		FROM cuisine_types
		WHERE name = $1
	`

	c, err := r.scanRow(q.QueryRowContext(ctx, query, name))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return c, nil
}

func (r *CuisineTypeRepo) Delete(ctx context.Context, id int) error {
	const op = "repository.CuisineTypeRepo.Delete"

	q := r.Queryer(ctx)

	query := `
		DELETE FROM cuisine_types 
		WHERE id = $1
	`

	res, err := q.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rows == 0 {
		return fmt.Errorf("%s: %w", op, domain.ErrNotFound)
	}

	return nil
}

func (r *CuisineTypeRepo) scanRow(s scanner) (*domain.CuisineType, error) {
	var cuisine domain.CuisineType
	err := s.Scan(
		&cuisine.ID,
		&cuisine.Name,
		&cuisine.Label,
	)
	if err != nil {
		return nil, err
	}
	return &cuisine, nil
}
