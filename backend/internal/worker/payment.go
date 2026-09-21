package worker

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"gitlab.com/2foood/backend/internal/domain"
	"gitlab.com/2foood/backend/internal/infra/payment"
	"gitlab.com/2foood/backend/internal/infra/queue"
)

const paymentErrCodeInternalError = "internal_error"

type PaymentManager interface {
	ConfirmPayment(paymentID string) error
	CancelPaymentIntent(paymentID string) error
	RefundPayment(paymentID string) error
}

type PaymentEventRecorder interface {
	RecordPaymentEvent(ctx context.Context, orderID uuid.UUID, paymentID string, event domain.PaymentEvent, errMsg *string) error
}

type OrderAuditor interface {
	AuditPayment(ctx context.Context, orderID uuid.UUID, paymentID string) error
}

// handles all payment jobs: confirm / cancel / refund / audit.
// returning an error = asynq retries the job (up to MaxRetries), returning nil = done
type PaymentHandler struct {
	log      *slog.Logger
	manager  PaymentManager
	auditor  OrderAuditor
	recorder PaymentEventRecorder
}

func NewPaymentHandler(log *slog.Logger, manager PaymentManager, auditor OrderAuditor, recorder PaymentEventRecorder) *PaymentHandler {
	return &PaymentHandler{
		log:     log,
		manager: manager,
		auditor: auditor,
		recorder: recorder,
	}
}

func (h *PaymentHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {
	const op = "worker.PaymentHandler.ProcessTask"
	log := h.log.With(slog.String("op", op))

	// a payload we can't parse gets `return nil` (no retry), retrying wouldn't fix it
	var payload queue.PaymentPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		log.Error("failed to unmarshal payload", slog.String("error", err.Error()))
		return nil // don't retry if payload is invalid
	}

	log = log.With(
		slog.String("task_type", task.Type()),
		slog.String("payment_id", payload.PaymentID),
		slog.String("order_id", payload.OrderID.String()),
	)

	switch task.Type() {
	case queue.TypePaymentConfirm:
		err := h.manager.ConfirmPayment(payload.PaymentID)
		if err != nil {
			h.record(ctx, log, payload, domain.PaymentEventConfirmFailed, err)
			return err
		}
		h.record(ctx, log, payload, domain.PaymentEventSucceeded, nil)

	case queue.TypePaymentCancel:
		err := h.manager.CancelPaymentIntent(payload.PaymentID)
		if err != nil {
			h.record(ctx, log, payload, domain.PaymentEventCancelFailed, err)
			return err
		}
		h.record(ctx, log, payload, domain.PaymentEventCancelled, nil)

	case queue.TypePaymentRefund:
		err := h.manager.RefundPayment(payload.PaymentID)
		if err != nil {
			h.record(ctx, log, payload, domain.PaymentEventRefundFailed, err)
			return err
		}
		h.record(ctx, log, payload, domain.PaymentEventRefunded, nil)

	case queue.TypePaymentAudit:
		err := h.auditor.AuditPayment(ctx, payload.OrderID, payload.PaymentID)
		if err != nil {
			log.Error("audit task failed, will retry", slog.String("error", err.Error()))
			return err
		}

	default:
		log.Warn("unknown task type received")
		return nil // don't retry if task type is unknown
	}

	log.Info("payment task processed successfully")
	return nil
}

// writes the result (success or failure code) into payment history.
// failing to record it is only logged, it shouldn't fail the job itself
func (h *PaymentHandler) record(ctx context.Context, log *slog.Logger, payload queue.PaymentPayload, event domain.PaymentEvent, err error) {
	var errMsg *string
	if err != nil {
		var paymentErr *payment.Error
		if errors.As(err, &paymentErr) {
			code := paymentErr.Code
			errMsg = &code
		} else {
			code := paymentErrCodeInternalError
			errMsg = &code
		}
	}

	if recErr := h.recorder.RecordPaymentEvent(ctx, payload.OrderID, payload.PaymentID, event, errMsg); recErr != nil {
		log.Error("failed to record payment history", slog.String("error", recErr.Error()))
	}
}
