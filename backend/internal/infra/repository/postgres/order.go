package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gitlab.com/2foood/backend/internal/domain"
)

type OrderRepo struct {
	BaseRepo
}

func NewOrderRepo(db *sql.DB) *OrderRepo {
	return &OrderRepo{BaseRepo: NewBaseRepo(db)}
}

func (r *OrderRepo) Create(ctx context.Context, order *domain.Order) error {
	const op = "repository.OrderRepo.Create"

	q := r.Queryer(ctx)

	query := `
		INSERT INTO orders (id, user_id, restaurant_id, status, total_amount, display_number, pickup_code, payment_id, pickup_window_start, pickup_window_end, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := q.ExecContext(ctx, query,
		order.ID,
		order.UserID,
		order.RestaurantID,
		order.Status,
		order.TotalAmount,
		order.DisplayNumber,
		order.PickupCode,
		order.PaymentID,
		order.PickupWindowStart,
		order.PickupWindowEnd,
		order.CreatedAt,
		order.UpdatedAt,
	)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == pgUniqueErrCode {
			return fmt.Errorf("%s: %w", op, domain.ErrAlreadyExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	for _, item := range order.Items {
		if err := r.createItem(ctx, q, &item); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

func (r *OrderRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	const op = "repository.OrderRepo.GetByID"

	q := r.Queryer(ctx)

	query := `
		SELECT id, user_id, restaurant_id, status, total_amount, display_number, pickup_code, payment_id, pickup_window_start, pickup_window_end, created_at, updated_at, reviewed_at
		FROM orders
		WHERE id = $1
	`

	order, err := r.scanRow(q.QueryRowContext(ctx, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	items, err := r.getItems(ctx, q, order.ID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	order.Items = items

	return order, nil
}

func (r *OrderRepo) GetByUser(ctx context.Context, userID uuid.UUID) ([]*domain.Order, error) {
	const op = "repository.OrderRepo.GetByUser"

	q := r.Queryer(ctx)

	query := `
		SELECT id, user_id, restaurant_id, status, total_amount, display_number, pickup_code, payment_id, pickup_window_start, pickup_window_end, created_at, updated_at, reviewed_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	orders, err := r.scanList(ctx, q, query, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := r.loadItems(ctx, q, orders); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return orders, nil
}

func (r *OrderRepo) GetByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]*domain.Order, error) {
	const op = "repository.OrderRepo.GetByRestaurant"

	q := r.Queryer(ctx)

	query := `
		SELECT id, user_id, restaurant_id, status, total_amount, display_number, pickup_code, payment_id, pickup_window_start, pickup_window_end, created_at, updated_at, reviewed_at
		FROM orders
		WHERE restaurant_id = $1
		ORDER BY created_at DESC
	`

	orders, err := r.scanList(ctx, q, query, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := r.loadItems(ctx, q, orders); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return orders, nil
}

func (r *OrderRepo) GetActiveByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]*domain.Order, error) {
	const op = "repository.OrderRepo.GetActiveByRestaurant"

	q := r.Queryer(ctx)

	query := `
		SELECT id, user_id, restaurant_id, status, total_amount, display_number, pickup_code, payment_id, pickup_window_start, pickup_window_end, created_at, updated_at, reviewed_at
		FROM orders
		WHERE restaurant_id = $1 AND status IN ('PAID', 'READY')
		ORDER BY pickup_window_end ASC
	`

	orders, err := r.scanList(ctx, q, query, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := r.loadItems(ctx, q, orders); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return orders, nil
}

func (r *OrderRepo) GetByPickupCode(ctx context.Context, restaurantID uuid.UUID, pickupCode string) (*domain.Order, error) {
	const op = "repository.OrderRepo.GetByPickupCode"

	q := r.Queryer(ctx)

	query := `
		SELECT id, user_id, restaurant_id, status, total_amount, display_number, pickup_code, payment_id, pickup_window_start, pickup_window_end, created_at, updated_at, reviewed_at
		FROM orders
		WHERE restaurant_id = $1 AND pickup_code = $2
	`

	order, err := r.scanRow(q.QueryRowContext(ctx, query, restaurantID, pickupCode))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	items, err := r.getItems(ctx, q, order.ID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	order.Items = items

	return order, nil
}

func (r *OrderRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.OrderStatus) error {
	const op = "repository.OrderRepo.UpdateStatus"

	q := r.Queryer(ctx)

	query := `
		UPDATE orders
		SET status = $2, updated_at = NOW()
		WHERE id = $1
	`

	res, err := q.ExecContext(ctx, query, id, status)
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

func (r *OrderRepo) SetReviewed(ctx context.Context, orderID uuid.UUID) error {
	const op = "repository.OrderRepo.SetReviewed"

	q := r.Queryer(ctx)

	query := `
		UPDATE orders
		SET reviewed_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`

	_, err := q.ExecContext(ctx, query, orderID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *OrderRepo) AddStatusHistory(ctx context.Context, entry *domain.OrderStatusHistoryEntry) error {
	const op = "repository.OrderRepo.AddStatusHistory"

	q := r.Queryer(ctx)

	query := `
		INSERT INTO order_status_history (id, order_id, actor_id, old_status, new_status, reason, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := q.ExecContext(ctx, query,
		entry.ID,
		entry.OrderID,
		entry.ActorID,
		entry.OldStatus,
		entry.NewStatus,
		entry.Reason,
		entry.CreatedAt,
	)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == pgUniqueErrCode {
			return fmt.Errorf("%s: %w", op, domain.ErrAlreadyExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *OrderRepo) GetStatusHistory(ctx context.Context, orderID uuid.UUID) ([]*domain.OrderStatusHistoryEntry, error) {
	const op = "repository.OrderRepo.GetStatusHistory"

	q := r.Queryer(ctx)

	query := `
		SELECT id, order_id, actor_id, old_status, new_status, reason, created_at
		FROM order_status_history
		WHERE order_id = $1
		ORDER BY created_at DESC
	`

	rows, err := q.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var list []*domain.OrderStatusHistoryEntry
	for rows.Next() {
		e, err := r.scanStatusHistoryRow(rows)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		list = append(list, e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return list, nil
}

func (r *OrderRepo) AddPaymentHistory(ctx context.Context, entry *domain.PaymentHistoryEntry) error {
	const op = "repository.OrderRepo.AddPaymentHistory"

	q := r.Queryer(ctx)

	query := `
		INSERT INTO payment_history (id, order_id, user_id, payment_id, event, amount, currency, error_message, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := q.ExecContext(ctx, query,
		entry.ID,
		entry.OrderID,
		entry.UserID,
		entry.PaymentID,
		entry.Event,
		entry.Amount,
		entry.Currency,
		entry.ErrorMessage,
		entry.CreatedAt,
	)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == pgUniqueErrCode {
			return fmt.Errorf("%s: %w", op, domain.ErrAlreadyExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *OrderRepo) GetPaymentHistory(ctx context.Context, orderID uuid.UUID) ([]*domain.PaymentHistoryEntry, error) {
	const op = "repository.OrderRepo.GetPaymentHistory"

	q := r.Queryer(ctx)

	query := `
		SELECT id, order_id, user_id, payment_id, event, amount, currency, error_message, created_at
		FROM payment_history
		WHERE order_id = $1
		ORDER BY created_at DESC
	`

	rows, err := q.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var list []*domain.PaymentHistoryEntry
	for rows.Next() {
		e, err := r.scanPaymentHistoryRow(rows)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		list = append(list, e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return list, nil
}

func (r *OrderRepo) GetPaymentHistoryByUser(ctx context.Context, userID uuid.UUID) ([]*domain.PaymentHistoryEntry, error) {
	const op = "repository.OrderRepo.GetPaymentHistoryByUser"

	q := r.Queryer(ctx)

	query := `
		SELECT id, order_id, user_id, payment_id, event, amount, currency, error_message, created_at
		FROM payment_history
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := q.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var list []*domain.PaymentHistoryEntry
	for rows.Next() {
		e, err := r.scanPaymentHistoryRow(rows)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		list = append(list, e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return list, nil
}

func (r *OrderRepo) GetExpiredPendingOrders(ctx context.Context, threshold time.Time) ([]*domain.Order, error) {
	const op = "repository.OrderRepo.GetExpiredPendingOrders"

	q := r.Queryer(ctx)

	query := `
		SELECT id, user_id, restaurant_id, status, total_amount, display_number, pickup_code, payment_id, pickup_window_start, pickup_window_end, created_at, updated_at, reviewed_at
		FROM orders
		WHERE status = 'PENDING' AND created_at < $1
	`

	orders, err := r.scanList(ctx, q, query, threshold)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return orders, nil
}

func (r *OrderRepo) createItem(ctx context.Context, q Queryer, item *domain.OrderItem) error {
	const op = "repository.OrderRepo.createItem"

	query := `
		INSERT INTO order_items (id, order_id, product_id, product_name, price, quantity)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := q.ExecContext(ctx, query,
		item.ID,
		item.OrderID,
		item.ProductID,
		item.ProductName,
		item.Price,
		item.Quantity,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *OrderRepo) getItems(ctx context.Context, q Queryer, orderID uuid.UUID) ([]domain.OrderItem, error) {
	const op = "repository.OrderRepo.getItems"

	query := `
		SELECT id, order_id, product_id, product_name, price, quantity
		FROM order_items
		WHERE order_id = $1
	`

	rows, err := q.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var items []domain.OrderItem
	for rows.Next() {
		item, err := r.scanItemRow(rows)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		items = append(items, *item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return items, nil
}

func (r *OrderRepo) scanList(ctx context.Context, q Queryer, query string, args ...any) ([]*domain.Order, error) {
	const op = "repository.OrderRepo.scanList"

	rows, err := q.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	list := []*domain.Order{}
	for rows.Next() {
		order, err := r.scanRow(rows)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		list = append(list, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return list, nil
}

func (r *OrderRepo) loadItems(ctx context.Context, q Queryer, orders []*domain.Order) error {
	for _, order := range orders {
		items, err := r.getItems(ctx, q, order.ID)
		if err != nil {
			return err
		}
		order.Items = items
	}
	return nil
}

func (r *OrderRepo) scanRow(s scanner) (*domain.Order, error) {
	var order domain.Order
	err := s.Scan(
		&order.ID,
		&order.UserID,
		&order.RestaurantID,
		&order.Status,
		&order.TotalAmount,
		&order.DisplayNumber,
		&order.PickupCode,
		&order.PaymentID,
		&order.PickupWindowStart,
		&order.PickupWindowEnd,
		&order.CreatedAt,
		&order.UpdatedAt,
		&order.ReviewedAt,
	)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepo) scanItemRow(s scanner) (*domain.OrderItem, error) {
	var item domain.OrderItem
	err := s.Scan(
		&item.ID,
		&item.OrderID,
		&item.ProductID,
		&item.ProductName,
		&item.Price,
		&item.Quantity,
	)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *OrderRepo) scanStatusHistoryRow(s scanner) (*domain.OrderStatusHistoryEntry, error) {
	var entry domain.OrderStatusHistoryEntry
	err := s.Scan(
		&entry.ID,
		&entry.OrderID,
		&entry.ActorID,
		&entry.OldStatus,
		&entry.NewStatus,
		&entry.Reason,
		&entry.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (r *OrderRepo) scanPaymentHistoryRow(s scanner) (*domain.PaymentHistoryEntry, error) {
	var entry domain.PaymentHistoryEntry
	err := s.Scan(
		&entry.ID,
		&entry.OrderID,
		&entry.UserID,
		&entry.PaymentID,
		&entry.Event,
		&entry.Amount,
		&entry.Currency,
		&entry.ErrorMessage,
		&entry.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}
