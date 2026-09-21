package domain

import (
	"time"

	"github.com/google/uuid"
)

// normal flow:  PENDING -> PAID -> READY -> COLLECTED
// (ACCEPTED is an optional step between PAID and READY)
// CANCELLED can happen any time before COLLECTED.
// UNPAID is only set after COLLECTED, when the audit finds out the payment actually failed
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "PENDING"
	OrderStatusPaid      OrderStatus = "PAID"
	OrderStatusAccepted  OrderStatus = "ACCEPTED"
	OrderStatusReady     OrderStatus = "READY"
	OrderStatusCollected OrderStatus = "COLLECTED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
	OrderStatusUnpaid    OrderStatus = "UNPAID"
)

// why an order was cancelled / marked unpaid, saved in the status history
type OrderReason string

const (
	OrderReasonOutOfStock          OrderReason = "out_of_stock"
	OrderReasonTimeout             OrderReason = "payment_timeout"
	OrderReasonInternalSystemError OrderReason = "internal_system_error"
	OrderReasonFailedAudit         OrderReason = "failed_audit"
	OrderReasonPaymentFailed       OrderReason = "payment_failed"
)

// rows in payment history, one per thing that happened with the payment
type PaymentEvent string

const (
	PaymentEventCreated       PaymentEvent = "CREATED"
	PaymentEventSucceeded     PaymentEvent = "SUCCEEDED"
	PaymentEventCancelled     PaymentEvent = "CANCELLED"
	PaymentEventRefunded      PaymentEvent = "REFUNDED"
	PaymentEventConfirmFailed PaymentEvent = "CONFIRM_FAILED"
	PaymentEventCancelFailed  PaymentEvent = "CANCEL_FAILED"
	PaymentEventRefundFailed  PaymentEvent = "REFUND_FAILED"
	PaymentEventPaymentFailed PaymentEvent = "PAYMENT_FAILED"
)

const (
	PaymentMetaKeyOrderID      = "order_id"
	PaymentMetaKeyUserID       = "user_id"
	PaymentMetaKeyRestaurantID = "restaurant_id"
)

type Order struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	RestaurantID uuid.UUID `json:"restaurant_id"`

	Status        OrderStatus `json:"status"`
	DisplayNumber string      `json:"display_number"`
	PickupCode    string      `json:"pickup_code"`
	TotalAmount   float64     `json:"total_amount"`
	PaymentID     *string     `json:"payment_id,omitempty"`

	// time frame the customer picked for coming to get the food
	PickupWindowStart time.Time `json:"pickup_window_start"`
	PickupWindowEnd   time.Time `json:"pickup_window_end"`

	Items []OrderItem `json:"items,omitempty"`

	ReviewedAt *time.Time `json:"reviewed_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type OrderItem struct {
	ID          uuid.UUID `json:"id"`
	OrderID     uuid.UUID `json:"order_id"`
	ProductID   uuid.UUID `json:"product_id"`
	ProductName string    `json:"product_name"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
}

type OrderStatusHistoryEntry struct {
	ID        uuid.UUID    `json:"id"`
	OrderID   uuid.UUID    `json:"order_id"`
	ActorID   *uuid.UUID   `json:"actor_id,omitempty"`
	OldStatus OrderStatus  `json:"old_status"`
	NewStatus OrderStatus  `json:"new_status"`
	Reason    *OrderReason `json:"reason,omitempty"`
	CreatedAt time.Time    `json:"created_at"`
}

type PaymentHistoryEntry struct {
	ID           uuid.UUID    `json:"id"`
	OrderID      uuid.UUID    `json:"order_id"`
	UserID       uuid.UUID    `json:"user_id"`
	PaymentID    string       `json:"payment_id"`
	Event        PaymentEvent `json:"event"`
	Amount       float64      `json:"amount"`
	Currency     string       `json:"currency"`
	ErrorMessage *string      `json:"error_message,omitempty"`
	CreatedAt    time.Time    `json:"created_at"`
}

type CreateOrderInput struct {
	UserID            uuid.UUID
	RestaurantID      uuid.UUID
	PickupWindowStart time.Time
	PickupWindowEnd   time.Time
	Items             []CreateOrderItemInput
}

type CreateOrderItemInput struct {
	ProductID uuid.UUID
	Quantity  int
}

type AddOrderStatusHistoryEntryInput struct {
	OrderID   uuid.UUID
	ActorID   *uuid.UUID
	OldStatus OrderStatus
	NewStatus OrderStatus
	Reason    *OrderReason
}

type AddPaymentHistoryInput struct {
	OrderID      uuid.UUID
	UserID       uuid.UUID
	PaymentID    string
	Event        PaymentEvent
	Amount       float64
	Currency     string
	ErrorMessage *string
}

func NewOrder(input CreateOrderInput, displayNumber string, pickupCode string) *Order {
	now := time.Now().UTC()
	return &Order{
		ID:                uuid.New(),
		UserID:            input.UserID,
		RestaurantID:      input.RestaurantID,
		Status:            OrderStatusPending,
		DisplayNumber:     displayNumber,
		PickupCode:        pickupCode,
		PickupWindowStart: input.PickupWindowStart,
		PickupWindowEnd:   input.PickupWindowEnd,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}

func NewOrderStatusHistoryEntry(input AddOrderStatusHistoryEntryInput) *OrderStatusHistoryEntry {
	now := time.Now().UTC()
	return &OrderStatusHistoryEntry{
		ID:        uuid.New(),
		OrderID:   input.OrderID,
		ActorID:   input.ActorID,
		OldStatus: input.OldStatus,
		NewStatus: input.NewStatus,
		Reason:    input.Reason,
		CreatedAt: now,
	}
}

func NewPaymentHistoryEntry(input AddPaymentHistoryInput) *PaymentHistoryEntry {
	now := time.Now().UTC()
	return &PaymentHistoryEntry{
		ID:           uuid.New(),
		OrderID:      input.OrderID,
		UserID:       input.UserID,
		PaymentID:    input.PaymentID,
		Event:        input.Event,
		Amount:       input.Amount,
		Currency:     input.Currency,
		ErrorMessage: input.ErrorMessage,
		CreatedAt:    now,
	}
}

func (o *Order) AddItem(productID uuid.UUID, productName string, price float64, quantity int) {
	o.Items = append(o.Items, OrderItem{
		ID:          uuid.New(),
		OrderID:     o.ID,
		ProductID:   productID,
		ProductName: productName,
		Price:       price,
		Quantity:    quantity,
	})
	o.TotalAmount += price * float64(quantity)
}

func (o *Order) PaymentMetadata() map[string]string {
	return map[string]string{
		PaymentMetaKeyOrderID:      o.ID.String(),
		PaymentMetaKeyUserID:       o.UserID.String(),
		PaymentMetaKeyRestaurantID: o.RestaurantID.String(),
	}
}

// the single source of truth for allowed status changes.
// service.updateStatusWithHistory checks this before every update
func (s OrderStatus) CanTransitionTo(to OrderStatus) bool {
	switch s {
	case OrderStatusPending:
		return to == OrderStatusPaid || to == OrderStatusCancelled
	case OrderStatusPaid:
		return to == OrderStatusAccepted || to == OrderStatusCancelled || to == OrderStatusReady // temporary allow transition to ready from paid
	case OrderStatusAccepted:
		return to == OrderStatusReady || to == OrderStatusCancelled
	case OrderStatusReady:
		return to == OrderStatusCollected || to == OrderStatusCancelled
	case OrderStatusCollected:
		return to == OrderStatusUnpaid
	default:
		return false
	}
}
