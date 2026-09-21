package http

import (
	"io"
	"log/slog"
	"net/http"

	"gitlab.com/2foood/backend/internal/domain"
	"gitlab.com/2foood/backend/internal/pkg/respond"
)

const (
	StripeSignatureHeader = "Stripe-Signature"
	MaxBodyBytes          = int64(65536) // 64KB, plenty for a webhook body
)

// StripeConfirmPayment godoc
// @Summary      Stripe webhook handler
// @Description  Receives and verifies Stripe webhook events for payment confirmation
// @Tags         webhooks, orders
// @Accept       json
// @Produce      json
// @Param        Stripe-Signature header string true "Stripe signature for payload verification"
// @Success      200 {object} respond.MessageResponse "Webhook successfully processed"
// @Failure      400 {object} respond.ErrorResponse "Missing signature or invalid payload"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /webhooks/stripe/payment/confirm [post]
func (h *Handler) StripeConfirmPayment(w http.ResponseWriter, r *http.Request) {
	const op = "order.Handler.StripeConfirmPayment"
	log := h.log.With(slog.String("op", op))

	// no signature = not from stripe
	if r.Header.Get(StripeSignatureHeader) == "" {
		respond.Error(w, r, log, domain.ErrMissingParam)
		return
	}

	// we need the raw bytes as-is, stripe signs the exact payload
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}
	if len(payload) == 0 {
		respond.Error(w, r, log, domain.ErrInvalidJSON)
		return
	}

	signature := r.Header.Get(StripeSignatureHeader)
	if err := h.service.ConfirmPaymentWebhook(r.Context(), payload, signature); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	respond.Message(w, r, "webhook processed")
}
