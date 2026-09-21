package provider

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"gitlab.com/2foood/backend/internal/pkg/generate"
)

type Mock struct {
	log *slog.Logger
}

func NewMock(log *slog.Logger) *Mock {
	return &Mock{log: log}
}

func (m *Mock) CreatePaymentIntent(amount float64, currency string, metadata map[string]string) (string, string, error) {
	paymentID := "pi_mock_" + generate.Code(16)
	clientSecret := "secret_" + generate.Code(24)

	m.log.Info("mock payment intent created",
		slog.String("payment_id", paymentID),
		slog.Float64("amount", amount),
		slog.String("currency", currency),
	)

	return paymentID, clientSecret, nil
}

func (m *Mock) ConfirmPayment(paymentID string) error {
	m.log.Info("mock payment confirmed", slog.String("payment_id", paymentID))
	return nil
}

func (m *Mock) CancelPaymentIntent(paymentID string) error {
	m.log.Info("mock payment intent cancelled", slog.String("payment_id", paymentID))
	return nil
}

func (m *Mock) RefundPayment(paymentID string) error {
	m.log.Info("mock payment refunded", slog.String("payment_id", paymentID))
	return nil
}

func (m *Mock) VerifyWebhookSignature(payload []byte, signature string) (string, string, error) {
	var req struct {
		OrderID   string `json:"order_id"`
		PaymentID string `json:"payment_id"`
	}

	if err := json.Unmarshal(payload, &req); err != nil {
		return "", "", fmt.Errorf("invalid webhook payload: %w", err)
	}

	m.log.Info("mock webhook verified",
		slog.String("order_id", req.OrderID),
		slog.String("payment_id", req.PaymentID),
	)

	return req.PaymentID, req.OrderID, nil
}
