package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gitlab.com/2foood/backend/internal/domain"
)

type UserRepo struct {
	BaseRepo
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{BaseRepo: NewBaseRepo(db)}
}

func (r *UserRepo) Create(ctx context.Context, user *domain.User) error {
	const op = "repository.UserRepo.Create"

	q := r.Queryer(ctx)

	notifJSON, err := json.Marshal(user.NotificationSettings)
	if err != nil {
		return fmt.Errorf("%s: marshal notification settings: %w", op, err)
	}
	dietJSON, err := json.Marshal(user.DietaryPreferences)
	if err != nil {
		return fmt.Errorf("%s: marshal dietary preferences: %w", op, err)
	}

	query := `
		INSERT INTO users (id, role, phone, email, google_id, name, avatar_url, city_id, language, notification_settings, dietary_preferences, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err = q.ExecContext(ctx, query,
		user.ID,
		user.Role,
		user.Phone,
		user.Email,
		user.GoogleID,
		user.Name,
		user.AvatarURL,
		user.CityID,
		user.Language,
		notifJSON,
		dietJSON,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == pgUniqueErrCode {
			return fmt.Errorf("%s: %w", op, domain.ErrAlreadyExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *UserRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	const op = "repository.UserRepo.GetByID"

	q := r.Queryer(ctx)

	query := `
		SELECT id, role, phone, email, google_id, name, avatar_url, city_id, language, notification_settings, dietary_preferences, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	u, err := r.scanRow(q.QueryRowContext(ctx, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return u, nil
}

func (r *UserRepo) GetByPhone(ctx context.Context, phone string) (*domain.User, error) {
	const op = "repository.UserRepo.GetByPhone"

	q := r.Queryer(ctx)

	query := `
		SELECT id, role, phone, email, google_id, name, avatar_url, city_id, language, notification_settings, dietary_preferences, created_at, updated_at
		FROM users
		WHERE phone = $1
	`

	u, err := r.scanRow(q.QueryRowContext(ctx, query, phone))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return u, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	const op = "repository.UserRepo.GetByEmail"

	q := r.Queryer(ctx)

	query := `
		SELECT id, role, phone, email, google_id, name, avatar_url, city_id, language, notification_settings, dietary_preferences, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	u, err := r.scanRow(q.QueryRowContext(ctx, query, email))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return u, nil
}

func (r *UserRepo) GetByGoogleID(ctx context.Context, googleID string) (*domain.User, error) {
	const op = "repository.UserRepo.GetByGoogleID"

	q := r.Queryer(ctx)

	query := `
		SELECT id, role, phone, email, google_id, name, avatar_url, city_id, language, notification_settings, dietary_preferences, created_at, updated_at
		FROM users
		WHERE google_id = $1
	`

	u, err := r.scanRow(q.QueryRowContext(ctx, query, googleID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return u, nil
}

func (r *UserRepo) Update(ctx context.Context, user *domain.User) error {
	const op = "repository.UserRepo.Update"

	q := r.Queryer(ctx)

	notifJSON, err := json.Marshal(user.NotificationSettings)
	if err != nil {
		return fmt.Errorf("%s: marshal notification settings: %w", op, err)
	}
	dietJSON, err := json.Marshal(user.DietaryPreferences)
	if err != nil {
		return fmt.Errorf("%s: marshal dietary preferences: %w", op, err)
	}

	query := `
		UPDATE users
		SET name = $2, avatar_url = $3, city_id = $4, language = $5, notification_settings = $6, dietary_preferences = $7, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`

	err = q.QueryRowContext(ctx, query,
		user.ID,
		user.Name,
		user.AvatarURL,
		user.CityID,
		user.Language,
		notifJSON,
		dietJSON,
	).Scan(&user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *UserRepo) SetPhone(ctx context.Context, id uuid.UUID, phone string) error {
	const op = "repository.UserRepo.SetPhone"

	q := r.Queryer(ctx)

	query := `
		UPDATE users
		SET phone = $2, updated_at = NOW()
		WHERE id = $1
	`

	res, err := q.ExecContext(ctx, query, id, phone)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == pgUniqueErrCode {
			return fmt.Errorf("%s: %w", op, domain.ErrAlreadyExists)
		}
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

func (r *UserRepo) SetEmail(ctx context.Context, id uuid.UUID, email string) error {
	const op = "repository.UserRepo.SetEmail"

	q := r.Queryer(ctx)

	query := `
		UPDATE users
		SET email = $2, updated_at = NOW()
		WHERE id = $1
	`

	res, err := q.ExecContext(ctx, query, id, email)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == pgUniqueErrCode {
			return fmt.Errorf("%s: %w", op, domain.ErrAlreadyExists)
		}
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

func (r *UserRepo) SetRole(ctx context.Context, id uuid.UUID, role domain.UserRole) error {
	const op = "repository.UserRepo.SetRole"

	q := r.Queryer(ctx)

	query := `
		UPDATE users
		SET role = $2, updated_at = NOW()
		WHERE id = $1 AND role != 'admin'
	`

	res, err := q.ExecContext(ctx, query, id, role)
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

func (r *UserRepo) LinkGoogle(ctx context.Context, userID uuid.UUID, info *domain.OAuthInfo) error {
	const op = "repository.UserRepo.LinkGoogle"

	q := r.Queryer(ctx)

	query := `
		UPDATE users
		SET 
			google_id = $2,
			email = COALESCE(NULLIF($3, ''), email),
			name = COALESCE(NULLIF($4, ''), name),
			avatar_url = COALESCE(NULLIF($5, ''), avatar_url),
			updated_at = NOW()
		WHERE id = $1
	`

	// TODO: add city_id, language, notification_settings, dietary_preferences
	res, err := q.ExecContext(ctx, query,
		userID,
		info.OAuthID,
		info.Email,
		info.Name,
		info.AvatarURL,
	)

	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") || strings.Contains(err.Error(), "duplicate key") {
			return fmt.Errorf("%s: email already in use by another user: %w", op, domain.ErrAlreadyExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *UserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	const op = "repository.UserRepo.Delete"

	q := r.Queryer(ctx)

	query := `
		DELETE FROM users
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

func (r *UserRepo) CleanupRoles(ctx context.Context) (int64, error) {
	const op = "repository.UserRepo.CleanupRoles"

	q := r.Queryer(ctx)

	query := `
		UPDATE users
		SET role = 'customer', updated_at = NOW()
		WHERE role = 'partner'
			AND id NOT IN (SELECT DISTINCT user_id FROM restaurant_staff)
	`

	res, err := q.ExecContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return count, nil
}

func (r *UserRepo) scanRow(s scanner) (*domain.User, error) {
	var u domain.User
	var notifJSON, dietJSON []byte

	err := s.Scan(
		&u.ID,
		&u.Role,
		&u.Phone,
		&u.Email,
		&u.GoogleID,
		&u.Name,
		&u.AvatarURL,
		&u.CityID,
		&u.Language,
		&notifJSON,
		&dietJSON,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if err := unmarshalJSONBValue(notifJSON, &u.NotificationSettings); err != nil {
		return nil, err
	}
	if err := unmarshalJSONBValue(dietJSON, &u.DietaryPreferences); err != nil {
		return nil, err
	}

	return &u, nil
}
