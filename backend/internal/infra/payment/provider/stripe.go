package provider

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/paymentintent"
	"github.com/stripe/stripe-go/v81/refund"
	"github.com/stripe/stripe-go/v81/webhook"
	"gitlab.com/2foood/backend/internal/domain"
	"gitlab.com/2foood/backend/internal/infra/payment"
)

const (
	EventPaymentIntentAmountCapturableUpdated = "payment_intent.amount_capturable_updated"
)

type StripeConfig struct {
	SecretKey     string `envconfig:"STRIPE_SECRET_KEY" required:"true"`
	WebhookSecret string `envconfig:"STRIPE_WEBHOOK_SECRET" required:"true"`
}

type Stripe struct {
	webhookSecret string
}

func NewStripe(cfg StripeConfig) *Stripe {
	stripe.Key = cfg.SecretKey
	return &Stripe{webhookSecret: cfg.WebhookSecret}
}

func (s *Stripe) CreatePaymentIntent(amount float64, multiplier int, currency string, metadata map[string]string) (string, string, error) {
	totalAmount := int64(amount * float64(multiplier)) // stripe expects amount in smallest currency unit

	params := &stripe.PaymentIntentParams{
		Amount:        stripe.Int64(totalAmount),
		Currency:      stripe.String(currency),
		CaptureMethod: stripe.String("manual"),
	}

	for k, v := range metadata {
		params.AddMetadata(k, v)
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return "", "", fmt.Errorf("stripe create payment intent: %w", err)
	}

	return pi.ID, pi.ClientSecret, nil
}

func (s *Stripe) ConfirmPayment(paymentID string) error {
	_, err := paymentintent.Capture(paymentID, nil)
	if err != nil {
		return convertStripeError("stripe capture", err)
	}

	return nil
}

func (s *Stripe) CancelPaymentIntent(paymentID string) error {
	_, err := paymentintent.Cancel(paymentID, nil)
	if err != nil {
		return convertStripeError("stripe cancel", err)
	}

	return nil
}

func (s *Stripe) RefundPayment(paymentID string) error {
	params := &stripe.RefundParams{
		PaymentIntent: stripe.String(paymentID),
	}

	_, err := refund.New(params)
	if err != nil {
		return convertStripeError("stripe refund", err)
	}

	return nil
}

func (s *Stripe) VerifyWebhookSignature(payload []byte, signature string) (string, string, error) {
	event, err := webhook.ConstructEventWithOptions(
		payload,
		signature,
		s.webhookSecret,
		webhook.ConstructEventOptions{
			IgnoreAPIVersionMismatch: true,
		})
	if err != nil {
		return "", "", fmt.Errorf("stripe webhook verification failed: %w", err)
	}

	if event.Type != EventPaymentIntentAmountCapturableUpdated {
		return "", "", domain.ErrIgnoredWebhook
	}

	var pi stripe.PaymentIntent
	if err := json.Unmarshal(event.Data.Raw, &pi); err != nil {
		return "", "", fmt.Errorf("failed to parse payment intent: %w", err)
	}

	orderID, ok := pi.Metadata[domain.PaymentMetaKeyOrderID]
	if !ok {
		return "", "", fmt.Errorf("order_id not found in metadata")
	}

	return pi.ID, orderID, nil
}

func (s *Stripe) IsPaymentFailed(paymentID string) (bool, error) {
	pi, err := paymentintent.Get(paymentID, nil)
	if err != nil {
		return false, fmt.Errorf("stripe get payment intent: %w", err)
	}

	isFailed := pi.Status == stripe.PaymentIntentStatusCanceled ||
		pi.Status == stripe.PaymentIntentStatusRequiresPaymentMethod

	return isFailed, nil
}

func convertStripeError(op string, err error) error {
	var stripeErr *stripe.Error
	if errors.As(err, &stripeErr) {
		return &payment.Error{
			Code:    string(stripeErr.Code),
			Message: stripeErr.Msg,
		}
	}
	return fmt.Errorf("%s: %w", op, err)
}
