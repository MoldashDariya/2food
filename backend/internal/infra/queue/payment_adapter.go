package queue

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

type PaymentPayload struct {
	PaymentID string    `json:"payment_id"`
	OrderID   uuid.UUID `json:"order_id"`
}

type PaymentAdapter struct {
	client *Client
	maxRetries int
}

func NewPaymentAdapter(client *Client, maxRetries int) *PaymentAdapter {
	return &PaymentAdapter{client: client, maxRetries: maxRetries}
}

func (a *PaymentAdapter) ConfirmAsync(ctx context.Context, paymentID string, orderID uuid.UUID) error {
	payload := PaymentPayload{PaymentID: paymentID, OrderID: orderID}
	return a.client.EnqueueTask(ctx, TypePaymentConfirm, payload, asynq.MaxRetry(a.maxRetries))
}

func (a *PaymentAdapter) CancelAsync(ctx context.Context, paymentID string, orderID uuid.UUID) error {
	payload := PaymentPayload{PaymentID: paymentID, OrderID: orderID}
	return a.client.EnqueueTask(ctx, TypePaymentCancel, payload, asynq.MaxRetry(a.maxRetries))
}

func (a *PaymentAdapter) RefundAsync(ctx context.Context, paymentID string, orderID uuid.UUID) error {
	payload := PaymentPayload{PaymentID: paymentID, OrderID: orderID}
	return a.client.EnqueueTask(ctx, TypePaymentRefund, payload, asynq.MaxRetry(a.maxRetries))
}

func (a *PaymentAdapter) AuditAsync(ctx context.Context, paymentID string, orderID uuid.UUID, delay time.Duration) error {
	payload := PaymentPayload{PaymentID: paymentID, OrderID: orderID}
	return a.client.EnqueueTask(ctx, TypePaymentAudit, payload, asynq.ProcessIn(delay), asynq.MaxRetry(a.maxRetries))
}
