package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/hibiken/asynq"
)

type ExpiredPendingOrdersCanceler interface {
	CancelExpiredPendingOrders(ctx context.Context, expirationTime time.Duration) (int, error)
}

// runs on a cron (see app.go): cancels orders that stayed PENDING (unpaid) for too long
type OrderCleanerHandler struct {
	log            *slog.Logger
	canceler       ExpiredPendingOrdersCanceler
	expirationTime time.Duration
}

func NewOrderCleanerHandler(log *slog.Logger, canceler ExpiredPendingOrdersCanceler, expirationTime time.Duration) *OrderCleanerHandler {
	return &OrderCleanerHandler{
		log:            log,
		canceler:       canceler,
		expirationTime: expirationTime,
	}
}

func (h *OrderCleanerHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {
	const op = "worker.OrderCleanerHandler.ProcessTask"
	log := h.log.With(slog.String("op", op))

	// don't let one run hang forever
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	count, err := h.canceler.CancelExpiredPendingOrders(ctx, h.expirationTime)
	if err != nil {
		log.Error("failed to cleanup expired orders", slog.String("error", err.Error()))
		return err
	}

	if count > 0 {
		log.Info("cleaned up expired orders", slog.Int("count", count))
	} else {
		log.Info("no expired orders to cleanup")
	}

	return nil
}
