package http

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/render"
)

const checkTimeout = 5 * time.Second

type CheckFunc func(ctx context.Context) error

type Handler struct {
	log    *slog.Logger
	checks map[string]CheckFunc
}

func NewHandler(log *slog.Logger) *Handler {
	return &Handler{
		log:    log,
		checks: make(map[string]CheckFunc),
	}
}

func (h *Handler) Add(name string, check CheckFunc) {
	h.checks[name] = check
}

func (h *Handler) Live(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Readiness(w http.ResponseWriter, r *http.Request) {
	const op = "health.Handler.Readiness"
	log := h.log.With(slog.String("op", op))

	ctx, cancel := context.WithTimeout(r.Context(), checkTimeout)
	defer cancel()

	status := make(map[string]string)
	hasError := false

	for name, check := range h.checks {
		if err := check(ctx); err != nil {
			log.Error("health check failed",
				slog.String("service", name),
				slog.String("error", err.Error()),
			)
			status[name] = "down"
			hasError = true
			continue
		}
		status[name] = "up"
	}

	httpStatus := http.StatusOK
	if hasError {
		httpStatus = http.StatusServiceUnavailable
	}

	render.Status(r, httpStatus)
	render.JSON(w, r, status)
}
