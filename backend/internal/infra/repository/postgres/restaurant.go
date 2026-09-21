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

type RestaurantRepo struct {
	BaseRepo
}

func NewRestaurantRepo(db *sql.DB) *RestaurantRepo {
	return &RestaurantRepo{BaseRepo: NewBaseRepo(db)}
}

func (r *RestaurantRepo) Create(ctx context.Context, rest *domain.Restaurant) error {
	const op = "repository.RestaurantRepo.Create"

	q := r.Queryer(ctx)

	whJSON, err := rest.MarshalWorkingHours()
	if err != nil {
		return fmt.Errorf("%s: marshal working hours: %w", op, err)
	}
	slJSON, err := rest.MarshalSocialLinks()
	if err != nil {
		return fmt.Errorf("%s: marshal social links: %w", op, err)
	}

	query := `
		INSERT INTO restaurants (
			id, name, description, address, phone, website,
			image_url, latitude, longitude, working_hours, social_links,
			city_id, is_active, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`

	_, err = q.ExecContext(ctx, query,
		rest.ID, rest.Name, rest.Description, rest.Address, rest.Phone, rest.Website,
		rest.ImageURL, rest.Latitude, rest.Longitude, whJSON, slJSON,
		rest.CityID, rest.IsActive, rest.CreatedAt, rest.UpdatedAt,
	)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == pgUniqueErrCode {
			return fmt.Errorf("%s: %w", op, domain.ErrAlreadyExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *RestaurantRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Restaurant, error) {
	const op = "repository.RestaurantRepo.GetByID"

	q := r.Queryer(ctx)

	query := `
		SELECT
			r.id, r.name, r.description, r.address, r.phone, r.website,
			r.image_url, r.latitude, r.longitude, r.working_hours, r.social_links,
			r.city_id, r.is_active, r.created_at, r.updated_at,
			r.rating_avg, r.rating_count,
			c.id, c.name, c.label
		FROM restaurants r
		LEFT JOIN cities c ON r.city_id = c.id
		WHERE r.id = $1
	`

	rest, err := r.scanRow(q.QueryRowContext(ctx, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	cuisines, err := r.getCuisines(ctx, q, rest.ID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	rest.CuisineTypes = cuisines

	return rest, nil
}

func (r *RestaurantRepo) GetAll(ctx context.Context) ([]*domain.Restaurant, error) {
	const op = "repository.RestaurantRepo.GetAll"

	q := r.Queryer(ctx)

	query := `
		SELECT
			r.id, r.name, r.description, r.address, r.phone, r.website,
			r.image_url, r.latitude, r.longitude, r.working_hours, r.social_links,
			r.city_id, r.is_active, r.created_at, r.updated_at,
			r.rating_avg, r.rating_count,
			c.id, c.name, c.label
		FROM restaurants r
		LEFT JOIN cities c ON r.city_id = c.id
	`

	rests, err := r.scanList(ctx, q, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := r.loadCuisines(ctx, q, rests); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return rests, nil
}

func (r *RestaurantRepo) GetByPartner(ctx context.Context, userID uuid.UUID) ([]*domain.Restaurant, error) {
	const op = "repository.RestaurantRepo.GetByPartner"

	q := r.Queryer(ctx)

	query := `
		SELECT
			r.id, r.name, r.description, r.address, r.phone, r.website,
			r.image_url, r.latitude, r.longitude, r.working_hours, r.social_links,
			r.city_id, r.is_active, r.created_at, r.updated_at,
			r.rating_avg, r.rating_count,
			c.id, c.name, c.label,
			rs.role AS role
		FROM restaurants r
		LEFT JOIN cities c ON r.city_id = c.id
		LEFT JOIN restaurant_staff rs ON r.id = rs.restaurant_id
		WHERE rs.user_id = $1
		ORDER BY r.created_at DESC
	`

	rows, err := q.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var list []*domain.Restaurant
	for rows.Next() {
		var role domain.StaffRole
		rest, err := r.scanRow(rows, &role)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		rest.Role = role
		list = append(list, rest)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := r.loadCuisines(ctx, q, list); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return list, nil
}

func (r *RestaurantRepo) UpdateInfo(ctx context.Context, rest *domain.Restaurant) error {
	const op = "postgres.RestaurantRepo.UpdateInfo"

	q := r.Queryer(ctx)

	whJSON, err := rest.MarshalWorkingHours()
	if err != nil {
		return fmt.Errorf("%s: marshal working hours: %w", op, err)
	}
	slJSON, err := rest.MarshalSocialLinks()
	if err != nil {
		return fmt.Errorf("%s: marshal social links: %w", op, err)
	}

	query := `
		UPDATE restaurants
		SET
			name = $2, description = $3, address = $4, phone = $5, website = $6,
			image_url = $7, latitude = $8, longitude = $9,
			working_hours = $10, social_links = $11,
			city_id = $12, is_active = $13, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`

	err = q.QueryRowContext(ctx, query,
		rest.ID, rest.Name, rest.Description, rest.Address, rest.Phone, rest.Website,
		rest.ImageURL, rest.Latitude, rest.Longitude, whJSON, slJSON,
		rest.CityID, rest.IsActive,
	).Scan(&rest.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *RestaurantRepo) UpdateCuisines(ctx context.Context, restID uuid.UUID, cuisineIDs []int) error {
	const op = "repository.RestaurantRepo.UpdateCuisines"

	q := r.Queryer(ctx)

	// Delete existing cuisines
	query := `
		DELETE FROM restaurant_cuisine_types
		WHERE restaurant_id = $1
	`

	_, err := q.ExecContext(ctx, query, restID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// Insert new cuisines
	query = `
		INSERT INTO restaurant_cuisine_types (restaurant_id, cuisine_type_id)
		VALUES ($1, $2)
	`

	for _, cid := range cuisineIDs {
		_, err := q.ExecContext(ctx, query, restID, cid)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

func (r *RestaurantRepo) IncrementRating(ctx context.Context, restID uuid.UUID, rating int8) error {
	const op = "repository.RestaurantRepo.IncrementRating"

	q := r.Queryer(ctx)

	query := `
		UPDATE restaurants
		SET 
			rating_avg = (rating_avg * rating_count + $2) / (rating_count + 1), 
			rating_count = rating_count + 1, 
			updated_at = NOW()
		WHERE id = $1
	`

	_, err := q.ExecContext(ctx, query, restID, rating)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *RestaurantRepo) Delete(ctx context.Context, id uuid.UUID) error {
	const op = "repository.RestaurantRepo.Delete"

	q := r.Queryer(ctx)

	query := `
		DELETE FROM restaurants 
		WHERE id = $1
	`

	res, err := q.ExecContext(ctx, query, id)
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

func (r *RestaurantRepo) getCuisines(ctx context.Context, q Queryer, restID uuid.UUID) ([]domain.CuisineType, error) {
	const op = "repository.RestaurantRepo.getCuisines"

	query := `
		SELECT ct.id, ct.name, ct.label
		FROM restaurant_cuisine_types rct
		JOIN cuisine_types ct ON rct.cuisine_type_id = ct.id
		WHERE rct.restaurant_id = $1
		ORDER BY ct.name
	`

	rows, err := q.QueryContext(ctx, query, restID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.CuisineType
	for rows.Next() {
		cuisine, err := r.scanCuisineRow(rows)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		list = append(list, *cuisine)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return list, nil
}

func (r *RestaurantRepo) loadCuisines(ctx context.Context, q Queryer, list []*domain.Restaurant) error {
	for _, rest := range list {
		cuisines, err := r.getCuisines(ctx, q, rest.ID)
		if err != nil {
			return err
		}
		rest.CuisineTypes = cuisines
	}
	return nil
}

func (r *RestaurantRepo) scanList(ctx context.Context, q Queryer, query string, args ...any) ([]*domain.Restaurant, error) {
	const op = "repository.RestaurantRepo.scanList"

	rows, err := q.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	list := []*domain.Restaurant{}
	for rows.Next() {
		rest, err := r.scanRow(rows)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		list = append(list, rest)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return list, nil
}

func (r *RestaurantRepo) scanRow(s scanner, extra ...any) (*domain.Restaurant, error) {
	var (
		rest      domain.Restaurant
		whJSON    []byte
		slJSON    []byte
		cityID    sql.NullInt64
		cityName  sql.NullString
		cityLabel sql.NullString
	)

	dest := []any{
		&rest.ID, &rest.Name, &rest.Description, &rest.Address, &rest.Phone, &rest.Website,
		&rest.ImageURL, &rest.Latitude, &rest.Longitude, &whJSON, &slJSON,
		&rest.CityID, &rest.IsActive, &rest.CreatedAt, &rest.UpdatedAt,
		&rest.RatingAvg, &rest.RatingCount,
		&cityID, &cityName, &cityLabel,
	}
	dest = append(dest, extra...)

	if err := s.Scan(dest...); err != nil {
		return nil, err
	}

	if err := unmarshalJSONB(whJSON, &rest.WorkingHours); err != nil {
		return nil, err
	}
	if err := unmarshalJSONB(slJSON, &rest.SocialLinks); err != nil {
		return nil, err
	}
	if cityID.Valid {
		rest.City = &domain.City{ID: int(cityID.Int64), Name: cityName.String, Label: cityLabel.String}
	}

	return &rest, nil
}

func (r *RestaurantRepo) scanCuisineRow(s scanner) (*domain.CuisineType, error) {
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
