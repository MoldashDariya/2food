package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"gitlab.com/2foood/backend/internal/domain"
)

type DeviceTokenRepo struct {
	BaseRepo
}

func NewDeviceTokenRepo(db *sql.DB) *DeviceTokenRepo {
	return &DeviceTokenRepo{BaseRepo: NewBaseRepo(db)}
}

func (r *DeviceTokenRepo) Register(ctx context.Context, dt *domain.DeviceToken) error {
	const op = "repository.DeviceTokenRepo.Register"

	q := r.Queryer(ctx)

	query := `
		INSERT INTO device_tokens (id, user_id, token, platform, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (token) DO UPDATE
			SET user_id = EXCLUDED.user_id,
				platform = EXCLUDED.platform,
				updated_at = EXCLUDED.updated_at
	`

	_, err := q.ExecContext(ctx, query,
		dt.ID, dt.UserID, dt.Token, dt.Platform, dt.CreatedAt, dt.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *DeviceTokenRepo) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.DeviceToken, error) {
	const op = "repository.DeviceTokenRepo.GetByUserID"

	q := r.Queryer(ctx)

	query := `
		SELECT id, user_id, token, platform, created_at, updated_at
		FROM device_tokens
		WHERE user_id = $1
	`

	rows, err := q.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	list := []*domain.DeviceToken{}
	for rows.Next() {
		dt, err := r.scanRow(rows)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		list = append(list, dt)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return list, nil
}

func (r *DeviceTokenRepo) Delete(ctx context.Context, userID uuid.UUID, token string) error {
	const op = "repository.DeviceTokenRepo.Delete"

	q := r.Queryer(ctx)

	query := `
		DELETE FROM device_tokens
		WHERE user_id = $1 AND token = $2
	`

	res, err := q.ExecContext(ctx, query, userID, token)
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

func (r *DeviceTokenRepo) scanRow(s scanner) (*domain.DeviceToken, error) {
	var dt domain.DeviceToken
	err := s.Scan(
		&dt.ID,
		&dt.UserID,
		&dt.Token,
		&dt.Platform,
		&dt.CreatedAt,
		&dt.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &dt, nil
}
