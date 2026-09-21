package worker

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"gitlab.com/2foood/backend/internal/domain"
	"gitlab.com/2foood/backend/internal/infra/queue"
)

type PushSender interface {
	SendToUser(ctx context.Context, userID uuid.UUID, title, body string, data map[string]any) (int, error)
}

type UserGetter interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
}

type Translator interface {
	T(lang, messageID string, data map[string]any) string
}

type PushHandler struct {
	log        *slog.Logger
	sender     PushSender
	tr         Translator
	userGetter UserGetter
}

func NewPushHandler(log *slog.Logger, sender PushSender, tr Translator, userGetter UserGetter) *PushHandler {
	return &PushHandler{log: log, sender: sender, tr: tr, userGetter: userGetter}
}

func (h *PushHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {
	const op = "worker.PushHandler.ProcessTask"
	log := h.log.With(slog.String("op", op))

	var payload queue.PushPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		log.Error("failed to unmarshal payload", slog.String("error", err.Error()))
		return nil // don't retry if payload is invalid
	}

	lang := string(domain.LanguageDefault)
	user, err := h.userGetter.GetByID(ctx, payload.UserID)
	if err == nil && user.Language != "" {
		lang = string(user.Language)
	}

	title := h.tr.T(lang, payload.MessageID+"_Title", payload.Data)
	body := h.tr.T(lang, payload.MessageID+"_Body", payload.Data)

	count, err := h.sender.SendToUser(ctx, payload.UserID, title, body, payload.Data)
	if err != nil {
		log.Error("failed to send push", slog.String("error", err.Error()), slog.String("user_id", payload.UserID.String()))
		return err
	}

	if count > 0 {
		log.Info("push sent successfully", slog.String("user_id", payload.UserID.String()), slog.Int("count", count))
	} else {
		log.Info("no push sent", slog.String("user_id", payload.UserID.String()))
	}

	return nil
}
