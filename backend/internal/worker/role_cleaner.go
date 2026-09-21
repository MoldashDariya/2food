package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/hibiken/asynq"
)

type RoleCleaner interface {
	CleanupRoles(ctx context.Context) (int64, error)
}

type RoleCleanerHandler struct {
	log     *slog.Logger
	cleaner RoleCleaner
}

func NewRoleCleanerHandler(log *slog.Logger, cleaner RoleCleaner) *RoleCleanerHandler {
	return &RoleCleanerHandler{
		log:     log,
		cleaner: cleaner,
	}
}

func (h *RoleCleanerHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {
	const op = "worker.RoleCleanerHandler.ProcessTask"
	log := h.log.With(slog.String("op", op))

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	count, err := h.cleaner.CleanupRoles(ctx)
	if err != nil {
		log.Error("failed to cleanup roles", slog.String("error", err.Error()))
		return err
	}

	if count > 0 {
		log.Info("cleaned up roles", slog.Int("count", int(count)))
	} else {
		log.Info("no roles to cleanup")
	}

	return nil
}
