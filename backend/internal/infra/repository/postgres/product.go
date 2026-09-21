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

type ProductRepo struct {
	BaseRepo
}

func NewProductRepo(db *sql.DB) *ProductRepo {
	return &ProductRepo{BaseRepo: NewBaseRepo(db)}
}

func (r *ProductRepo) Create(ctx context.Context, product *domain.Product) error {
	const op = "repository.ProductRepo.Create"

	query := `
		INSERT INTO products (
			id, restaurant_id, category_id, name, description, image_url,
			original_price, discounted_price, pickup_time_start, pickup_time_end,
			quantity_available, is_active, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	_, err := r.Queryer(ctx).ExecContext(ctx, query,
		product.ID,
		product.RestaurantID,
		product.CategoryID,
		product.Name,
		product.Description,
		product.ImageURL,
		product.OriginalPrice,
		product.DiscountedPrice,
		product.PickupTimeStart,
		product.PickupTimeEnd,
		product.QuantityAvailable,
		product.IsActive,
		product.CreatedAt,
		product.UpdatedAt,
	)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == pgUniqueErrCode {
			return fmt.Errorf("%s: %w", op, domain.ErrAlreadyExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *ProductRepo) GetByRestaurant(ctx context.Context, restID uuid.UUID) ([]*domain.Product, error) {
	const op = "repository.ProductRepo.GetByRestaurant"

	query := `
		SELECT
			p.id, p.restaurant_id, p.category_id,
			p.name, p.description, p.image_url,
			p.original_price, p.discounted_price,
			p.pickup_time_start, p.pickup_time_end,
			p.quantity_available, p.is_active,
			p.created_at, p.updated_at,
			c.name, c.label
		FROM products p
		LEFT JOIN product_categories c ON p.category_id = c.id
		WHERE p.restaurant_id = $1
		ORDER BY p.created_at DESC
	`

	return r.scanList(r.Queryer(ctx), ctx, op, query, restID)
}

func (r *ProductRepo) GetActiveByRestaurant(ctx context.Context, restID uuid.UUID) ([]*domain.Product, error) {
	const op = "repository.ProductRepo.GetActiveByRestaurant"

	query := `
		SELECT
			p.id, p.restaurant_id, p.category_id,
			p.name, p.description, p.image_url,
			p.original_price, p.discounted_price,
			p.pickup_time_start, p.pickup_time_end,
			p.quantity_available, p.is_active,
			p.created_at, p.updated_at,
			c.name, c.label
		FROM products p
		LEFT JOIN product_categories c ON p.category_id = c.id
		WHERE p.restaurant_id = $1 AND p.is_active = true
		ORDER BY p.created_at DESC
	`

	return r.scanList(r.Queryer(ctx), ctx, op, query, restID)
}

func (r *ProductRepo) GetByCategory(ctx context.Context, categoryID int) ([]*domain.Product, error) {
	const op = "repository.ProductRepo.GetByCategory"

	query := `
		SELECT
			p.id, p.restaurant_id, p.category_id,
			p.name, p.description, p.image_url,
			p.original_price, p.discounted_price,
			p.pickup_time_start, p.pickup_time_end,
			p.quantity_available, p.is_active,
			p.created_at, p.updated_at,
			c.name, c.label
		FROM products p
		LEFT JOIN product_categories c ON p.category_id = c.id
		WHERE p.category_id = $1
		ORDER BY p.is_active DESC, p.name ASC
	`

	return r.scanList(r.Queryer(ctx), ctx, op, query, categoryID)
}

func (r *ProductRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	const op = "repository.ProductRepo.GetByID"

	query := `
		SELECT
			p.id, p.restaurant_id, p.category_id,
			p.name, p.description, p.image_url,
			p.original_price, p.discounted_price,
			p.pickup_time_start, p.pickup_time_end,
			p.quantity_available, p.is_active,
			p.created_at, p.updated_at,
			c.name, c.label
		FROM products p
		LEFT JOIN product_categories c ON p.category_id = c.id
		WHERE p.id = $1
	`

	p, err := r.scanRow(r.Queryer(ctx).QueryRowContext(ctx, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return p, nil
}

func (r *ProductRepo) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	const op = "repository.ProductRepo.GetByIDForUpdate"

	query := `
		SELECT
			p.id, p.restaurant_id, p.category_id,
			p.name, p.description, p.image_url,
			p.original_price, p.discounted_price,
			p.pickup_time_start, p.pickup_time_end,
			p.quantity_available, p.is_active,
			p.created_at, p.updated_at,
			c.name, c.label
		FROM products p
		LEFT JOIN product_categories c ON p.category_id = c.id
		WHERE p.id = $1
		FOR UPDATE OF p
	`

	p, err := r.scanRow(r.Queryer(ctx).QueryRowContext(ctx, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return p, nil
}

func (r *ProductRepo) GetCategoriesByRestaurant(ctx context.Context, restID uuid.UUID) ([]*domain.ProductCategory, error) {
	const op = "repository.ProductRepo.GetCategoriesByRestaurant"

	q := r.Queryer(ctx)

	query := `
		SELECT DISTINCT c.id, c.name, c.label
		FROM products p
		JOIN product_categories c ON p.category_id = c.id
		WHERE p.restaurant_id = $1 AND p.is_active = true
		ORDER BY c.name
	`

	rows, err := q.QueryContext(ctx, query, restID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	list := []*domain.ProductCategory{}
	for rows.Next() {
		var c domain.ProductCategory
		if err := rows.Scan(&c.ID, &c.Name, &c.Label); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		list = append(list, &c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return list, nil
}

func (r *ProductRepo) Update(ctx context.Context, product *domain.Product) error {
	const op = "repository.ProductRepo.Update"

	query := `
		UPDATE products
		SET category_id = $2, name = $3, description = $4, image_url = $5,
		    original_price = $6, discounted_price = $7,
		    pickup_time_start = $8, pickup_time_end = $9,
		    updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`

	err := r.Queryer(ctx).QueryRowContext(ctx, query,
		product.ID,
		product.CategoryID,
		product.Name,
		product.Description,
		product.ImageURL,
		product.OriginalPrice,
		product.DiscountedPrice,
		product.PickupTimeStart,
		product.PickupTimeEnd,
	).Scan(&product.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s: %w", op, domain.ErrNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *ProductRepo) UpdateStatus(ctx context.Context, id uuid.UUID, isActive bool) error {
	const op = "repository.ProductRepo.UpdateStatus"

	res, err := r.Queryer(ctx).ExecContext(ctx, `
		UPDATE products SET is_active = $2, updated_at = NOW() WHERE id = $1
	`, id, isActive)
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

func (r *ProductRepo) UpdateStock(ctx context.Context, id uuid.UUID, quantity int) error {
	const op = "repository.ProductRepo.UpdateStock"

	res, err := r.Queryer(ctx).ExecContext(ctx, `
		UPDATE products SET quantity_available = $2, updated_at = NOW() WHERE id = $1
	`, id, quantity)
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

func (r *ProductRepo) Delete(ctx context.Context, id uuid.UUID) error {
	const op = "repository.ProductRepo.Delete"

	res, err := r.Queryer(ctx).ExecContext(ctx, `DELETE FROM products WHERE id = $1`, id)
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

func (r *ProductRepo) scanList(q Queryer, ctx context.Context, op, query string, args ...any) ([]*domain.Product, error) {
	rows, err := q.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	list := []*domain.Product{}
	for rows.Next() {
		p, err := r.scanRow(rows)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		list = append(list, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return list, nil
}

func (r *ProductRepo) scanRow(s scanner) (*domain.Product, error) {
	var (
		product       domain.Product
		categoryName  sql.NullString
		categoryLabel sql.NullString
	)

	err := s.Scan(
		&product.ID, &product.RestaurantID, &product.CategoryID,
		&product.Name, &product.Description, &product.ImageURL,
		&product.OriginalPrice, &product.DiscountedPrice,
		&product.PickupTimeStart, &product.PickupTimeEnd,
		&product.QuantityAvailable, &product.IsActive,
		&product.CreatedAt, &product.UpdatedAt,
		&categoryName, &categoryLabel,
	)
	if err != nil {
		return nil, err
	}

	if product.CategoryID != nil {
		product.Category = &domain.ProductCategory{
			ID:    *product.CategoryID,
			Name:  categoryName.String,
			Label: categoryLabel.String,
		}
	}

	return &product, nil
}
