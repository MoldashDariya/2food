package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gitlab.com/2foood/backend/internal/domain"
)

type ReviewRepo struct {
	BaseRepo
}

func NewReviewRepo(db *sql.DB) *ReviewRepo {
	return &ReviewRepo{BaseRepo: NewBaseRepo(db)}
}

func (r *ReviewRepo) Create(ctx context.Context, review *domain.Review) error {
	const op = "repository.ReviewRepo.Create"

	q := r.Queryer(ctx)

	photoJSON, err := json.Marshal(review.PhotoURLs)
	if err != nil {
		return fmt.Errorf("%s: marshal photo URLs: %w", op, err)
	}

	query := `
		INSERT INTO reviews (
			id, user_id, restaurant_id, order_id,
			rating, comment, photo_urls, is_hidden,
			created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err = q.ExecContext(ctx, query,
		review.ID,
		review.UserID,
		review.RestaurantID,
		review.OrderID,
		review.Rating,
		review.Comment,
		photoJSON,
		review.IsHidden,
		review.CreatedAt,
		review.UpdatedAt,
	)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == pgUniqueErrCode {
			return fmt.Errorf("%s: %w", op, domain.ErrAlreadyExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *ReviewRepo) GetAll(ctx context.Context) ([]*domain.Review, error) {
	const op = "repository.ReviewRepo.GetAll"

	q := r.Queryer(ctx)

	query := `
		SELECT
			rv.id, rv.user_id, rv.restaurant_id, rv.order_id,
			rv.rating, rv.comment, rv.photo_urls, 
			rv.reply, rv.replied_at, rv.is_hidden,
			rv.created_at, rv.updated_at,
			u.name, u.avatar_url, r.name
		FROM reviews rv
		LEFT JOIN users u ON rv.user_id = u.id
		LEFT JOIN restaurants r ON rv.restaurant_id = r.id
		ORDER BY rv.created_at DESC
	`

	rows, err := q.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	list := []*domain.Review{}
	for rows.Next() {
		var authorName, authorAvatar, restaurantName *string
		review, err := r.scanRow(rows, &authorName, &authorAvatar, &restaurantName)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		review.AuthorName = authorName
		review.AuthorAvatarURL = authorAvatar
		review.RestaurantName = restaurantName
		list = append(list, review)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return list, nil
}

func (r *ReviewRepo) GetByRestaurant(ctx context.Context, restID uuid.UUID) ([]*domain.Review, error) {
	const op = "repository.ReviewRepo.GetByRestaurant"

	q := r.Queryer(ctx)

	query := `
		SELECT
			rv.id, rv.user_id, rv.restaurant_id, rv.order_id,
			rv.rating, rv.comment, rv.photo_urls, 
			rv.reply, rv.replied_at, rv.is_hidden,
			rv.created_at, rv.updated_at,
			u.name, u.avatar_url
		FROM reviews rv
		LEFT JOIN users u ON rv.user_id = u.id
		WHERE rv.restaurant_id = $1 AND rv.is_hidden = false
		ORDER BY rv.created_at DESC
	`

	rows, err := q.QueryContext(ctx, query, restID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	list := []*domain.Review{}
	for rows.Next() {
		var authorName, authorAvatar *string
		review, err := r.scanRow(rows, &authorName, &authorAvatar)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		review.AuthorName = authorName
		review.AuthorAvatarURL = authorAvatar
		list = append(list, review)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return list, nil
}

func (r *ReviewRepo) GetByRestaurantAll(ctx context.Context, restID uuid.UUID) ([]*domain.Review, error) {
	const op = "repository.ReviewRepo.GetByRestaurantAll"

	q := r.Queryer(ctx)

	query := `
		SELECT
			rv.id, rv.user_id, rv.restaurant_id, rv.order_id,
			rv.rating, rv.comment, rv.photo_urls, 
			rv.reply, rv.replied_at, rv.is_hidden,
			rv.created_at, rv.updated_at,
			u.name, u.avatar_url
		FROM reviews rv
		LEFT JOIN users u ON rv.user_id = u.id
		WHERE rv.restaurant_id = $1
		ORDER BY rv.created_at DESC
	`

	rows, err := q.QueryContext(ctx, query, restID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	list := []*domain.Review{}
	for rows.Next() {
		var authorName, authorAvatar *string
		review, err := r.scanRow(rows, &authorName, &authorAvatar)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		review.AuthorName = authorName
		review.AuthorAvatarURL = authorAvatar
		list = append(list, review)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return list, nil
}

func (r *ReviewRepo) GetByUser(ctx context.Context, userID uuid.UUID) ([]*domain.Review, error) {
	const op = "repository.ReviewRepo.GetByUser"

	q := r.Queryer(ctx)

	query := `
		SELECT
			rv.id, rv.user_id, rv.restaurant_id, rv.order_id,
			rv.rating, rv.comment, rv.photo_urls, 
			rv.reply, rv.replied_at, rv.is_hidden,
			rv.created_at, rv.updated_at,
			r.name
		FROM reviews rv
		LEFT JOIN restaurants r ON rv.restaurant_id = r.id
		WHERE rv.user_id = $1
		ORDER BY rv.created_at DESC
	`

	rows, err := q.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	list := []*domain.Review{}
	for rows.Next() {
		var restaurantName *string
		review, err := r.scanRow(rows, &restaurantName)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		review.RestaurantName = restaurantName
		list = append(list, review)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return list, nil
}

func (r *ReviewRepo) ExistsByOrder(ctx context.Context, orderID uuid.UUID) (bool, error) {
	const op = "repository.ReviewRepo.ExistsByOrder"

	q := r.Queryer(ctx)

	var exists bool

	query := `
		SELECT EXISTS (
			SELECT 1 FROM reviews WHERE order_id = $1
		)
	`

	err := q.QueryRowContext(ctx, query, orderID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return exists, nil
}

func (r *ReviewRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Review, error) {
	const op = "repository.ReviewRepo.GetByID"

	q := r.Queryer(ctx)

	query := `
		SELECT
			rv.id, rv.user_id, rv.restaurant_id, rv.order_id,
			rv.rating, rv.comment, rv.photo_urls, 
			rv.reply, rv.replied_at, rv.is_hidden,
			rv.created_at, rv.updated_at,
			u.name, u.avatar_url, r.name
		FROM reviews rv
		LEFT JOIN users u ON rv.user_id = u.id
		LEFT JOIN restaurants r ON rv.restaurant_id = r.id
		WHERE rv.id = $1
	`

	var authorName, authorAvatar, restaurantName *string
	review, err := r.scanRow(q.QueryRowContext(ctx, query, id), &authorName, &authorAvatar, &restaurantName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	review.AuthorName = authorName
	review.AuthorAvatarURL = authorAvatar
	review.RestaurantName = restaurantName

	return review, nil
}

func (r *ReviewRepo) GetByOrder(ctx context.Context, orderID uuid.UUID) (*domain.Review, error) {
	const op = "repository.ReviewRepo.GetByOrder"

	q := r.Queryer(ctx)

	query := `
		SELECT
			rv.id, rv.user_id, rv.restaurant_id, rv.order_id,
			rv.rating, rv.comment, rv.photo_urls, 
			rv.reply, rv.replied_at, rv.is_hidden,
			rv.created_at, rv.updated_at,
			u.name, u.avatar_url, r.name
		FROM reviews rv
		LEFT JOIN users u ON rv.user_id = u.id
		LEFT JOIN restaurants r ON rv.restaurant_id = r.id
		WHERE rv.order_id = $1
	`

	var authorName, authorAvatar, restaurantName *string
	review, err := r.scanRow(q.QueryRowContext(ctx, query, orderID), &authorName, &authorAvatar, &restaurantName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	review.AuthorName = authorName
	review.AuthorAvatarURL = authorAvatar
	review.RestaurantName = restaurantName

	return review, nil
}

func (r *ReviewRepo) GetPhotosByRestaurant(ctx context.Context, restID uuid.UUID) ([]*domain.ReviewPhoto, error) {
	const op = "repository.ReviewRepo.GetPhotosByRestaurant"

	q := r.Queryer(ctx)

	query := `
		SELECT
			jsonb_array_elements_text(rv.photo_urls) as photo_url,
			rv.id,
			u.name, u.avatar_url,
			rv.rating, rv.comment,
			rv.reply, rv.replied_at,
			rv.created_at
		FROM reviews rv
		LEFT JOIN users u ON rv.user_id = u.id
		WHERE rv.restaurant_id = $1
			AND rv.is_hidden = false
			AND jsonb_typeof(rv.photo_urls) = 'array'
			AND jsonb_array_length(rv.photo_urls) > 0
		ORDER BY rv.created_at DESC		
	`

	rows, err := q.QueryContext(ctx, query, restID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	list := []*domain.ReviewPhoto{}
	for rows.Next() {
		p, err := r.scanPhotoRow(rows)
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

func (r *ReviewRepo) SetReply(ctx context.Context, review *domain.Review) error {
	const op = "repository.ReviewRepo.SetReply"

	q := r.Queryer(ctx)

	query := `
		UPDATE reviews
		SET reply = $2, replied_at = NOW(), updated_at = NOW()
		WHERE id = $1
		RETURNING replied_at, updated_at
	`

	err := q.QueryRowContext(ctx, query, review.ID, review.Reply).
		Scan(&review.RepliedAt, &review.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s: %w", op, domain.ErrNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *ReviewRepo) SetVisibility(ctx context.Context, id uuid.UUID, isHidden bool) error {
	const op = "repository.ReviewRepo.SetVisibility"

	q := r.Queryer(ctx)

	query := `
		UPDATE reviews
		SET is_hidden = $2, updated_at = NOW()
		WHERE id = $1
	`

	res, err := q.ExecContext(ctx, query, id, isHidden)
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

func (r *ReviewRepo) scanRow(s scanner, extra ...any) (*domain.Review, error) {
	var (
		review    domain.Review
		photoJSON []byte
	)

	dest := []any{
		&review.ID, &review.UserID, &review.RestaurantID, &review.OrderID,
		&review.Rating, &review.Comment, &photoJSON,
		&review.Reply, &review.RepliedAt,
		&review.IsHidden,
		&review.CreatedAt, &review.UpdatedAt,
	}
	dest = append(dest, extra...)

	if err := s.Scan(dest...); err != nil {
		return nil, err
	}

	if err := unmarshalJSONBValue(photoJSON, &review.PhotoURLs); err != nil {
		return nil, err
	}

	return &review, nil
}

func (r *ReviewRepo) scanPhotoRow(s scanner) (*domain.ReviewPhoto, error) {
	var p domain.ReviewPhoto

	err := s.Scan(
		&p.URL, &p.ReviewID,
		&p.AuthorName, &p.AuthorAvatarURL,
		&p.Rating, &p.Comment,
		&p.Reply, &p.RepliedAt,
		&p.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &p, nil
}
