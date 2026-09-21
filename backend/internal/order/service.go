package order

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"gitlab.com/2foood/backend/internal/domain"
	"gitlab.com/2foood/backend/internal/pkg/generate"
)

// everything the service needs is defined as an interface here (the consumer side),
// real implementations live in infra/, fakes for tests are in ./mocks
type TransactionManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error)
	GetByUser(ctx context.Context, userID uuid.UUID) ([]*domain.Order, error)
	GetByRestaurant(ctx context.Context, restID uuid.UUID) ([]*domain.Order, error)
	GetActiveByRestaurant(ctx context.Context, restID uuid.UUID) ([]*domain.Order, error)
	GetByPickupCode(ctx context.Context, restID uuid.UUID, pickupCode string) (*domain.Order, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.OrderStatus) error
	GetExpiredPendingOrders(ctx context.Context, threshold time.Time) ([]*domain.Order, error)
	AddStatusHistory(ctx context.Context, entry *domain.OrderStatusHistoryEntry) error
	GetStatusHistory(ctx context.Context, orderID uuid.UUID) ([]*domain.OrderStatusHistoryEntry, error)
	AddPaymentHistory(ctx context.Context, entry *domain.PaymentHistoryEntry) error
	GetPaymentHistory(ctx context.Context, orderID uuid.UUID) ([]*domain.PaymentHistoryEntry, error)
	GetPaymentHistoryByUser(ctx context.Context, userID uuid.UUID) ([]*domain.PaymentHistoryEntry, error)
}

type ProductRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	UpdateStock(ctx context.Context, id uuid.UUID, quantity int) error
}

// stripe sits behind this so we can swap it or mock it.
// we use manual capture: the app only authorizes the card, the money is taken later (see ConfirmPayment)
type PaymentGateway interface {
	CreatePaymentIntent(amount float64, metadata map[string]string) (paymentID string, clientSecret string, err error)
	VerifyWebhookSignature(payload []byte, signature string) (paymentID string, orderID string, err error)
	IsPaymentFailed(paymentID string) (bool, error)
}

// these push a job to the queue, the actual stripe call happens later in a worker (see worker/payment.go)
type PaymentDispatcher interface {
	ConfirmAsync(ctx context.Context, paymentID string, orderID uuid.UUID) error
	CancelAsync(ctx context.Context, paymentID string, orderID uuid.UUID) error
	RefundAsync(ctx context.Context, paymentID string, orderID uuid.UUID) error
	AuditAsync(ctx context.Context, paymentID string, orderID uuid.UUID, delay time.Duration) error
}

type PushDispatcher interface {
	SendAsync(ctx context.Context, userID uuid.UUID, message *domain.PushMessage) error
}

type Service struct {
	tm                  TransactionManager
	log                 *slog.Logger
	orderRepo           OrderRepository
	productRepo         ProductRepository
	paymentGateway      PaymentGateway
	paymentDispatcher   PaymentDispatcher
	pushDispatcher      PushDispatcher
	displayNumberLength int
	pickupCodeLength    int
	auditDelay          time.Duration
	currency            string
}

func NewService(
	tm TransactionManager, log *slog.Logger,
	orderRepo OrderRepository, productRepo ProductRepository,
	paymentGateway PaymentGateway, paymentDispatcher PaymentDispatcher, pushDispatcher PushDispatcher,
	displayNumberLength int, pickupCodeLength int, auditDelay time.Duration, currency string,
) *Service {
	return &Service{
		tm:                  tm,
		log:                 log,
		orderRepo:           orderRepo,
		productRepo:         productRepo,
		paymentGateway:      paymentGateway,
		paymentDispatcher:   paymentDispatcher,
		pushDispatcher:      pushDispatcher,
		displayNumberLength: displayNumberLength,
		pickupCodeLength:    pickupCodeLength,
		auditDelay:          auditDelay,
		currency:            currency,
	}
}

// creates the order + stripe payment intent.
// returns the client secret so the app can show the payment sheet.
// stock is NOT decreased here, that only happens once the payment is confirmed (ConfirmPayment)
func (s *Service) Create(ctx context.Context, input domain.CreateOrderInput) (*domain.Order, string, error) {
	const op = "order.Service.Create"

	var order *domain.Order
	var clientSecret string

	err := s.tm.Do(ctx, func(txCtx context.Context) error {
		// all in one db transaction, if something fails nothing gets saved
		displayNumber := generate.Code(s.displayNumberLength)
		pickupCode := generate.Code(s.pickupCodeLength)
		order = domain.NewOrder(input, displayNumber, pickupCode)

		for _, itemInput := range input.Items {
			product, err := s.productRepo.GetByID(txCtx, itemInput.ProductID)
			if err != nil {
				return fmt.Errorf("failed to get product: %w", err)
			}

			if !product.IsActive {
				return fmt.Errorf("product %s: %w", product.Name, domain.ErrProductUnavailable)
			}

			if product.RestaurantID != input.RestaurantID {
				return fmt.Errorf("product %s: %w", product.Name, domain.ErrPermissionDenied)
			}

			if product.QuantityAvailable < itemInput.Quantity {
				return fmt.Errorf("product %s: %w", product.Name, domain.ErrInsufficientStock)
			}

			if err := validatePickupWindow(input.PickupWindowStart, input.PickupWindowEnd, product.PickupTimeStart, product.PickupTimeEnd); err != nil {
				return fmt.Errorf("pickup window mismatch: %w", err)
			}

			// price is always taken from the db, never from the client
			order.AddItem(itemInput.ProductID, product.Name, product.DiscountedPrice, itemInput.Quantity)
		}

		// metadata lets us find the order again when the stripe webhook comes in
		paymentID, secret, err := s.paymentGateway.CreatePaymentIntent(order.TotalAmount, order.PaymentMetadata())
		if err != nil {
			return fmt.Errorf("failed to create payment intent: %w", err)
		}

		order.PaymentID = &paymentID
		clientSecret = secret

		// recording the fact of creating a payment intent
		entry := domain.NewPaymentHistoryEntry(domain.AddPaymentHistoryInput{
			OrderID:   order.ID,
			UserID:    order.UserID,
			PaymentID: paymentID,
			Event:     domain.PaymentEventCreated,
			Amount:    order.TotalAmount,
			Currency:  s.currency,
		})

		if err := s.orderRepo.Create(txCtx, order); err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		if err := s.orderRepo.AddPaymentHistory(txCtx, entry); err != nil {
			return fmt.Errorf("failed to add payment history: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, "", fmt.Errorf("%s: %w", op, err)
	}

	return order, clientSecret, nil
}

// called from the stripe webhook once the card is authorized (money is not taken yet).
// three steps: 1) db transaction: check stock, take it, mark PAID
// 2) queue the stripe side: capture the money, or cancel the hold if something ran out
// 3) push notification to the user
func (s *Service) ConfirmPayment(ctx context.Context, orderID uuid.UUID, paymentID string) error {
	const op = "order.Service.ConfirmPayment"

	var order *domain.Order

	// transcation step
	err := s.tm.Do(ctx, func(txCtx context.Context) error {
		var err error
		order, err = s.orderRepo.GetByID(txCtx, orderID)
		if err != nil {
			return fmt.Errorf("failed to get order: %w", err)
		}

		if order.PaymentID == nil || *order.PaymentID != paymentID {
			return fmt.Errorf("invalid payment ID: %w", domain.ErrPermissionDenied)
		}

		type stockUpdate struct {
			productID uuid.UUID
			newStock  int
		}
		var updates []stockUpdate

		// check if all products are in stock
		for _, item := range order.Items {
			// FOR UPDATE = row lock, so two people paying for the last item at the same time can't both get it
			product, err := s.productRepo.GetByIDForUpdate(txCtx, item.ProductID)
			if err != nil {
				return fmt.Errorf("failed to get product: %w", err)
			}

			newStock := product.QuantityAvailable - item.Quantity

			// someone else took the last one while this user was paying.
			// the `return nil` below is on purpose: we want the CANCELLED status committed, not rolled back
			// if product is out of stock, cancel order and exit transaction
			if newStock < 0 {
				reason := domain.OrderReasonOutOfStock
				if err := s.updateStatusWithHistory(txCtx, order, domain.OrderStatusCancelled, nil, &reason); err != nil {
					return fmt.Errorf("failed to update order history: %w", err)
				}
				return nil
			}

			updates = append(updates, stockUpdate{
				productID: item.ProductID,
				newStock:  newStock,
			})
		}

		// all products in stock,  update stock
		for _, u := range updates {
			if err := s.productRepo.UpdateStock(txCtx, u.productID, u.newStock); err != nil {
				return fmt.Errorf("failed to update product stock: %w", err)
			}
		}

		// success, update order status to paid
		if err := s.updateStatusWithHistory(txCtx, order, domain.OrderStatusPaid, nil, nil); err != nil {
			return fmt.Errorf("failed to update order status: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// after the tx is committed. the payment call goes through the queue so it can be retried
	// async tasks step
	switch order.Status {
	case domain.OrderStatusCancelled:
		// product is out of stock, cancel payment
		if err := s.paymentDispatcher.CancelAsync(ctx, paymentID, orderID); err != nil {
			s.log.Error("failed to cancel payment async",
				slog.String("payment_id", paymentID),
				slog.String("error", err.Error()),
			)
		}

	case domain.OrderStatusPaid:
		// everything is fine, confirm payment
		if err := s.paymentDispatcher.ConfirmAsync(ctx, paymentID, orderID); err != nil {
			s.log.Error("failed to confirm payment async. Starting db rollback!", slog.String("err", err.Error()))

			// couldn't even queue the confirmation -> undo the order in our db so we don't keep a PAID order with no confirmed payment
			reason := domain.OrderReasonInternalSystemError
			safeCtx := context.WithoutCancel(ctx) // don't cancel the context if the payment confirmation fails

			// db rollback
			if rollbackErr := s.cancelOrder(safeCtx, order, nil, &reason); rollbackErr != nil {
				s.log.Error("failed to rollback order after internal system error!",
					slog.Bool("is_critical", true),
					slog.String("order_id", orderID.String()),
					slog.String("rollback_err", rollbackErr.Error()),
				)
			}
			return fmt.Errorf("%s: internal system error: %w", op, err)
		}

		// confirmation sent successfully, schedule audit
		// double check with stripe a bit later that the payment really went through (AuditPayment)
		_ = s.paymentDispatcher.AuditAsync(ctx, paymentID, orderID, s.auditDelay)
	}

	// notifications step
	if order != nil {
		s.notifyOrderStatusChange(ctx, order)
	}

	return nil
}

// entry point for the webhook: verify the signature first, only then trust the payload
func (s *Service) ConfirmPaymentWebhook(ctx context.Context, payload []byte, signature string) error {
	const op = "order.Service.ConfirmPaymentWebhook"

	paymentID, orderIDStr, err := s.paymentGateway.VerifyWebhookSignature(payload, signature)
	if err != nil {
		// stripe sends lots of event types, the provider only handles one of them (amount_capturable_updated).
		// the rest get a plain 200 so stripe doesn't keep retrying
		if errors.Is(err, domain.ErrIgnoredWebhook) {
			return nil // ignored webhook
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := s.ConfirmPayment(ctx, orderID, paymentID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// safety net that runs some time after payment confirmation.
// if our db says "paid" but stripe says the payment failed, we fix the order here
func (s *Service) AuditPayment(ctx context.Context, orderID uuid.UUID, paymentID string) error {
	const op = "order.Service.AuditPayment"
	log := s.log.With(
		slog.String("op", op),
		slog.String("order_id", orderID.String()),
		slog.String("payment_id", paymentID),
	)

	// check if payment is failed
	isFailed, err := s.paymentGateway.IsPaymentFailed(paymentID)
	if err != nil {
		return fmt.Errorf("%s: failed to check if payment is failed: %w", op, err)
	}

	if isFailed {
		order, err := s.orderRepo.GetByID(ctx, orderID)
		if err != nil {
			return fmt.Errorf("%s: failed to get order: %w", op, err)
		}

		switch order.Status {
		case domain.OrderStatusPaid, domain.OrderStatusReady:
			log.Warn("payment audit failed! DB indicates success, but bank says failed. Starting late rollback", slog.String("db_status", string(order.Status)))

			reason := domain.OrderReasonFailedAudit
			safeCtx := context.WithoutCancel(ctx) // don't cancel the context if the rollback fails

			// trying to rollback the order
			if rollbackErr := s.cancelOrder(safeCtx, order, nil, &reason); rollbackErr != nil {
				s.log.Error("failed to rollback order during audit", slog.Bool("is_critical", true), slog.String("error", rollbackErr.Error()))
				return fmt.Errorf("%s: failed to cancel order: %w", op, rollbackErr)
			}

			// notify the user that the order has been cancelled
			s.notifyOrderStatusChange(safeCtx, order)

		case domain.OrderStatusCollected:
			// worst case: the food is already gone but the payment failed.
			// can't cancel anymore, so mark it UNPAID and log it as critical
			s.log.Error("payment failed but food was already collected", slog.Bool("is_critical", true), slog.String("user_id", order.UserID.String()))

			reason := domain.OrderReasonPaymentFailed
			safeCtx := context.WithoutCancel(ctx) // don't cancel the context if the rollback fails

			err := s.tm.Do(safeCtx, func(txCtx context.Context) error {
				if err := s.updateStatusWithHistory(txCtx, order, domain.OrderStatusUnpaid, nil, &reason); err != nil {
					return err
				}

				// recording the fact of payment failure
				entry := domain.NewPaymentHistoryEntry(domain.AddPaymentHistoryInput{
					OrderID:   orderID,
					UserID:    order.UserID,
					PaymentID: paymentID,
					Event:     domain.PaymentEventPaymentFailed,
					Amount:    order.TotalAmount,
					Currency:  s.currency,
				})

				return s.orderRepo.AddPaymentHistory(txCtx, entry)
			})

			if err != nil {
				s.log.Error("failed to mark order as unpaid during audit", slog.Bool("is_critical", true), slog.String("error", err.Error()))
				return fmt.Errorf("%s: failed to update status to unpaid: %w", op, err)
			}

			// notify the user that the order has been marked as unpaid
			s.notifyOrderStatusChange(safeCtx, order)

		default:
			// if order is in expected safe state, skip audit
			s.log.Info("audit skipped: order is in expected safe state", slog.String("order_id", orderID.String()), slog.String("db_status", string(order.Status)))
		}
	}

	return nil
}

func (s *Service) GetByID(ctx context.Context, orderID uuid.UUID) (*domain.Order, error) {
	const op = "order.Service.GetByID"

	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return order, nil
}

func (s *Service) GetByUser(ctx context.Context, userID uuid.UUID) ([]*domain.Order, error) {
	const op = "order.Service.GetByUser"

	orders, err := s.orderRepo.GetByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return orders, nil
}

func (s *Service) GetByRestaurant(ctx context.Context, restID uuid.UUID) ([]*domain.Order, error) {
	const op = "order.Service.GetByRestaurant"

	orders, err := s.orderRepo.GetByRestaurant(ctx, restID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return orders, nil
}

func (s *Service) GetActiveByRestaurant(ctx context.Context, restID uuid.UUID) ([]*domain.Order, error) {
	const op = "order.Service.GetActiveByRestaurant"

	orders, err := s.orderRepo.GetActiveByRestaurant(ctx, restID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return orders, nil
}

func (s *Service) GetByPickupCode(ctx context.Context, restID uuid.UUID, pickupCode string) (*domain.Order, error) {
	const op = "order.Service.GetByPickupCode"

	order, err := s.orderRepo.GetByPickupCode(ctx, restID, pickupCode)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return order, nil
}

// note: there's no route for this in router.go right now, restaurants go PAID -> READY directly
func (s *Service) MarkAsAccepted(ctx context.Context, actor *domain.RestaurantStaff, orderID uuid.UUID) error {
	const op = "order.Service.MarkAsAccepted"

	order, err := s.changeStatusByRestaurant(ctx, actor, orderID, domain.OrderStatusAccepted)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	s.notifyOrderStatusChange(ctx, order)

	return nil
}

func (s *Service) MarkAsReady(ctx context.Context, actor *domain.RestaurantStaff, orderID uuid.UUID) error {
	const op = "order.Service.MarkAsReady"

	order, err := s.changeStatusByRestaurant(ctx, actor, orderID, domain.OrderStatusReady)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	s.notifyOrderStatusChange(ctx, order)

	return nil
}

func (s *Service) MarkAsCollected(ctx context.Context, actor *domain.RestaurantStaff, orderID uuid.UUID, pickupCode string, force bool) error {
	const op = "order.Service.MarkAsCollected"

	// staff types in the code the customer shows. force = skip the check (logged as a warning)
	validatePickupCode := func(order *domain.Order) error {
		if !force {
			if pickupCode != order.PickupCode {
				return domain.ErrInvalidPickupCode
			}
		} else {
			s.log.Warn("force collect trigerred", slog.String("order_id", order.ID.String()), slog.String("actor_id", actor.UserID.String()))
		}
		return nil
	}

	order, err := s.changeStatusByRestaurant(ctx, actor, orderID, domain.OrderStatusCollected, validatePickupCode)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	s.notifyOrderStatusChange(ctx, order)

	return nil
}

func (s *Service) RecordPaymentEvent(ctx context.Context, orderID uuid.UUID, paymentID string, event domain.PaymentEvent, errMsg *string) error {
	const op = "order.Service.RecordPaymentEvent"

	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	entry := domain.NewPaymentHistoryEntry(domain.AddPaymentHistoryInput{
		OrderID:      orderID,
		UserID:       order.UserID,
		PaymentID:    paymentID,
		Event:        event,
		Amount:       order.TotalAmount,
		Currency:     s.currency,
		ErrorMessage: errMsg,
	})

	if err := s.orderRepo.AddPaymentHistory(ctx, entry); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) GetStatusHistoryByUser(ctx context.Context, userID uuid.UUID, orderID uuid.UUID) ([]*domain.OrderStatusHistoryEntry, error) {
	const op = "order.Service.GetStatusHistoryByUser"

	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if userID != order.UserID {
		return nil, fmt.Errorf("%s: order does not belong to user: %w", op, domain.ErrPermissionDenied)
	}

	history, err := s.orderRepo.GetStatusHistory(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return history, nil
}

func (s *Service) GetStatusHistoryByRestaurant(ctx context.Context, actor *domain.RestaurantStaff, orderID uuid.UUID) ([]*domain.OrderStatusHistoryEntry, error) {
	const op = "order.Service.GetStatusHistoryByRestaurant"

	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if actor.RestaurantID != order.RestaurantID {
		return nil, fmt.Errorf("%s: order does not belong to restaurant: %w", op, domain.ErrPermissionDenied)
	}

	history, err := s.orderRepo.GetStatusHistory(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return history, nil
}

func (s *Service) GetPaymentHistoryByOrder(ctx context.Context, userID uuid.UUID, orderID uuid.UUID) ([]*domain.PaymentHistoryEntry, error) {
	const op = "order.Service.GetPaymentHistoryByOrder"

	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if userID != order.UserID {
		return nil, fmt.Errorf("%s: order does not belong to user: %w", op, domain.ErrPermissionDenied)
	}

	history, err := s.orderRepo.GetPaymentHistory(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return history, nil
}

func (s *Service) GetPaymentHistoryByUser(ctx context.Context, userID uuid.UUID) ([]*domain.PaymentHistoryEntry, error) {
	const op = "order.Service.GetPaymentHistoryByUser"

	history, err := s.orderRepo.GetPaymentHistoryByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return history, nil
}

func (s *Service) CancelByUser(ctx context.Context, userID uuid.UUID, orderID uuid.UUID, reason *domain.OrderReason) error {
	const op = "order.Service.CancelByUser"

	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if userID != order.UserID {
		return fmt.Errorf("%s: order does not belong to user: %w", op, domain.ErrPermissionDenied)
	}

	if err := s.cancelOrder(ctx, order, &userID, reason); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}

func (s *Service) CancelByRestaurant(ctx context.Context, actor *domain.RestaurantStaff, orderID uuid.UUID, reason *domain.OrderReason) error {
	const op = "order.Service.CancelByRestaurant"

	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if actor.RestaurantID != order.RestaurantID {
		return fmt.Errorf("%s: order does not belong to restaurant: %w", op, domain.ErrPermissionDenied)
	}

	if err := s.cancelOrder(ctx, order, &actor.UserID, reason); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	s.notifyOrderStatusChange(ctx, order)

	return nil
}

// called by the cron worker. cancels orders that were never paid
func (s *Service) CancelExpiredPendingOrders(ctx context.Context, expirationTime time.Duration) (int, error) {
	const op = "order.Service.CancelExpiredPendingOrders"

	threshold := time.Now().UTC().Add(-expirationTime)

	orders, err := s.orderRepo.GetExpiredPendingOrders(ctx, threshold)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	if len(orders) == 0 {
		return 0, nil
	}

	reason := domain.OrderReasonTimeout
	var count int

	for _, order := range orders {
		// one bad order shouldn't stop the rest, the next cron run will pick it up again
		if err := s.cancelOrder(ctx, order, nil, &reason); err != nil {
			continue // ignore errors, we will retry later
		}
		count++
	}

	return count, nil
}

// shared by every cancel path (user, restaurant, timeout, audit).
// db part is in a transaction, the stripe part (cancel / refund) goes to the queue afterwards
func (s *Service) cancelOrder(ctx context.Context, order *domain.Order, actorID *uuid.UUID, reason *domain.OrderReason) error {
	const op = "order.Service.cancelOrder"

	// save original status to restore payment if needed
	originalStatus := order.Status

	err := s.tm.Do(ctx, func(txCtx context.Context) error {
		// stock is only taken after payment, so it's only given back for paid orders.
		// todo: ACCEPTED isn't handled here (or in processPaymentCancellation), fine while accept isn't used
		if order.Status == domain.OrderStatusPaid || order.Status == domain.OrderStatusReady {
			if err := s.restoreStock(txCtx, order); err != nil {
				return err
			}
		}

		if err := s.updateStatusWithHistory(txCtx, order, domain.OrderStatusCancelled, actorID, reason); err != nil {
			return fmt.Errorf("failed to update order status: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if order.PaymentID != nil {
		safeCtx := context.WithoutCancel(ctx)
		// detached ctx above: the order is already cancelled in the db, so the
		// refund/cancel must still be queued even if the request gets cancelled
		if err := s.processPaymentCancellation(safeCtx, order.ID, *order.PaymentID, originalStatus); err != nil {
			s.log.Error("order cancelled in DB, but failed to process payment cancellation",
				slog.Bool("is_critical", true),
				slog.String("order_id", order.ID.String()),
				slog.String("error", err.Error()),
			)
		}
	}

	return nil
}

func (s *Service) restoreStock(ctx context.Context, order *domain.Order) error {
	const op = "order.Service.restoreStock"

	for _, item := range order.Items {
		product, err := s.productRepo.GetByIDForUpdate(ctx, item.ProductID)
		if err != nil {
			return fmt.Errorf("%s: failed to get product %s: %w", op, item.ProductID, err)
		}

		newStock := product.QuantityAvailable + item.Quantity
		if err := s.productRepo.UpdateStock(ctx, item.ProductID, newStock); err != nil {
			return fmt.Errorf("%s: failed to update stock for product %s: %w", op, item.ProductID, err)
		}
	}

	return nil
}

func (s *Service) changeStatusByRestaurant(
	ctx context.Context,
	actor *domain.RestaurantStaff,
	orderID uuid.UUID,
	newStatus domain.OrderStatus,
	validators ...func(*domain.Order) error,
) (*domain.Order, error) {
	const op = "order.Service.changeStatusByRestaurant"

	var updatedOrder *domain.Order
	err := s.tm.Do(ctx, func(txCtx context.Context) error {
		order, err := s.orderRepo.GetByID(txCtx, orderID)
		if err != nil {
			return err
		}

		if actor.RestaurantID != order.RestaurantID {
			return domain.ErrPermissionDenied
		}

		for _, validate := range validators {
			if err := validate(order); err != nil {
				return err
			}
		}

		if err := s.updateStatusWithHistory(txCtx, order, newStatus, &actor.UserID, nil); err != nil {
			return fmt.Errorf("failed to update order status: %w", err)
		}

		updatedOrder = order
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return updatedOrder, nil
}

// every status change goes through here so the history table always stays in sync
func (s *Service) updateStatusWithHistory(ctx context.Context, order *domain.Order, newStatus domain.OrderStatus, actorID *uuid.UUID, reason *domain.OrderReason) error {
	const op = "order.Service.updateStatusWithHistory"

	if !order.Status.CanTransitionTo(newStatus) {
		return fmt.Errorf("%s: %w", op, domain.ErrInvalidStatusTransition)
	}

	if err := s.orderRepo.UpdateStatus(ctx, order.ID, newStatus); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	entry := domain.NewOrderStatusHistoryEntry(domain.AddOrderStatusHistoryEntryInput{
		OrderID:   order.ID,
		ActorID:   actorID,
		OldStatus: order.Status,
		NewStatus: newStatus,
		Reason:    reason,
	})

	if err := s.orderRepo.AddStatusHistory(ctx, entry); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// update order status in memory
	order.Status = newStatus

	return nil
}

// not paid yet -> just cancel the payment intent. already paid -> refund the money
func (s *Service) processPaymentCancellation(ctx context.Context, orderID uuid.UUID, paymentID string, status domain.OrderStatus) error {
	switch status {
	case domain.OrderStatusPending:
		return s.paymentDispatcher.CancelAsync(ctx, paymentID, orderID)
	case domain.OrderStatusPaid, domain.OrderStatusReady:
		return s.paymentDispatcher.RefundAsync(ctx, paymentID, orderID)
	default:
		return nil
	}
}

// MessageID is a key for the translator (internal/pkg/translator/locales), the text is
// rendered per user language in the push worker
func (s *Service) notifyOrderStatusChange(ctx context.Context, order *domain.Order) {
	var messageID string
	data := map[string]any{
		"order_id": order.ID.String(),
		"status":   string(order.Status),
	}

	switch order.Status {
	case domain.OrderStatusPaid:
		messageID = "OrderStatusPaid"
	case domain.OrderStatusAccepted:
		messageID = "OrderStatusAccepted"
	case domain.OrderStatusReady:
		messageID = "OrderStatusReady"
		data["pickup_code"] = order.PickupCode // Passed as a variable for the translator
	case domain.OrderStatusCollected:
		messageID = "OrderStatusCollected"
	case domain.OrderStatusCancelled:
		messageID = "OrderStatusCancelled"
	case domain.OrderStatusUnpaid:
		messageID = "OrderStatusUnpaid"
	default:
		return // no notification for other statuses
	}

	msg := &domain.PushMessage{
		MessageID: messageID,
		Data:      data,
	}

	// ignore errors because we don't want to block the flow
	err := s.pushDispatcher.SendAsync(ctx, order.UserID, msg)
	if err != nil {
		slog.Error("failed to send push notification async", slog.String("error", err.Error()), slog.String("order_id", order.ID.String()), slog.String("user_id", order.UserID.String()))
	}
}

// the chosen pickup window has to fit inside the time the product is sold.
// note: it's compared against today's date in UTC
func validatePickupWindow(orderStart, orderEnd time.Time, productStart, productEnd string) error {
	prodStart, err := parseTimeOfDay(productStart)
	if err != nil {
		return fmt.Errorf("invalid product start time: %w", err)
	}
	prodEnd, err := parseTimeOfDay(productEnd)
	if err != nil {
		return fmt.Errorf("invalid product end time: %w", err)
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)
	ps := today.Add(prodStart)
	pe := today.Add(prodEnd)

	if orderStart.Before(ps) || orderEnd.After(pe) {
		return domain.ErrPickupWindowMismatch
	}

	return nil
}

// product times can be "15:04" or a full RFC3339 timestamp, only the time of day is used
func parseTimeOfDay(s string) (time.Duration, error) {
	if d, err := time.Parse("15:04", s); err == nil {
		return time.Duration(d.Hour())*time.Hour + time.Duration(d.Minute())*time.Minute, nil
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return time.Duration(t.Hour())*time.Hour + time.Duration(t.Minute())*time.Minute, nil
	}
	return 0, fmt.Errorf("expected 15:04 or RFC3339, got %q", s)
}
