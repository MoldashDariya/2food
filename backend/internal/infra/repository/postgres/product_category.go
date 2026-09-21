package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"gitlab.com/2foood/backend/internal/domain"
)

type ProductCategoryRepo struct {
	BaseRepo
}

func NewProductCategoryRepo(db *sql.DB) *ProductCategoryRepo {
	return &ProductCategoryRepo{BaseRepo: NewBaseRepo(db)}
}
func (r *ProductCategoryRepo) Create(ctx context.Context, category *domain.ProductCategory) error {
	const op = "repository.ProductCategoryRepo.Create"

	query := `
		INSERT INTO product_categories (name, label) 
		VALUES ($1, $2) 
		RETURNING id`

	err := r.Queryer(ctx).QueryRowContext(ctx, query, category.Name, category.Label).Scan(&category.ID)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == pgUniqueErrCode {
			return fmt.Errorf("%s: %w", op, domain.ErrAlreadyExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *ProductCategoryRepo) GetAll(ctx context.Context) ([]*domain.ProductCategory, error) {
	const op = "repository.ProductCategoryRepo.GetAll"

	q := r.Queryer(ctx)

	query := `
		SELECT id, name, label
		FROM product_categories
		ORDER BY name
	`

	rows, err := q.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	list := []*domain.ProductCategory{}
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

func (r *ProductCategoryRepo) GetByID(ctx context.Context, id int) (*domain.ProductCategory, error) {
	const op = "repository.ProductCategoryRepo.GetByID"

	q := r.Queryer(ctx)

	query := `
		SELECT id, name, label
		FROM product_categories
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

func (r *ProductCategoryRepo) GetByName(ctx context.Context, name string) (*domain.ProductCategory, error) {
	const op = "repository.ProductCategoryRepo.GetByName"

	q := r.Queryer(ctx)

	query := `
		SELECT id, name, label
		FROM product_categories
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

func (r *ProductCategoryRepo) Delete(ctx context.Context, id int) error {
	const op = "repository.ProductCategoryRepo.Delete"

	q := r.Queryer(ctx)

	query := `
		DELETE FROM product_categories 
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

func (r *ProductCategoryRepo) scanRow(s scanner) (*domain.ProductCategory, error) {
	var category domain.ProductCategory
	err := s.Scan(
		&category.ID,
		&category.Name,
		&category.Label,
	)
	if err != nil {
		return nil, err
	}
	return &category, nil
}
