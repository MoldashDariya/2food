package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"gitlab.com/2foood/backend/internal/domain"
)

type CityRepo struct {
	BaseRepo
}

func NewCityRepo(db *sql.DB) *CityRepo {
	return &CityRepo{BaseRepo: NewBaseRepo(db)}
}

func (r *CityRepo) Create(ctx context.Context, city *domain.City) error {
	const op = "repository.CityRepo.Create"

	q := r.Queryer(ctx)

	query := `
		INSERT INTO cities (name, label)
		VALUES ($1, $2)
		RETURNING id
	`

	err := q.QueryRowContext(ctx, query, city.Name, city.Label).Scan(&city.ID)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == pgUniqueErrCode {
			return fmt.Errorf("%s: %w", op, domain.ErrAlreadyExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *CityRepo) GetAll(ctx context.Context) ([]*domain.City, error) {
	const op = "repository.CityRepo.GetAll"

	q := r.Queryer(ctx)

	query := `
		SELECT id, name, label
		FROM cities
		ORDER BY name
	`

	rows, err := q.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	list := []*domain.City{}
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

func (r *CityRepo) GetByID(ctx context.Context, id int) (*domain.City, error) {
	const op = "repository.CityRepo.GetByID"

	q := r.Queryer(ctx)

	query := `
		SELECT id, name, label
		FROM cities
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

func (r *CityRepo) GetByName(ctx context.Context, name string) (*domain.City, error) {
	const op = "repository.CityRepo.GetByName"

	q := r.Queryer(ctx)

	query := `
		SELECT id, name, label
		FROM cities
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

func (r *CityRepo) Delete(ctx context.Context, id int) error {
	const op = "repository.CityRepo.Delete"

	q := r.Queryer(ctx)

	query := `
		DELETE FROM cities 
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

func (r *CityRepo) scanRow(s scanner) (*domain.City, error) {
	var city domain.City
	err := s.Scan(
		&city.ID,
		&city.Name,
		&city.Label,
	)
	if err != nil {
		return nil, err
	}
	return &city, nil
}
