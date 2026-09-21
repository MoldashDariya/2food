package queue

import (
	"context"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"gitlab.com/2foood/backend/internal/domain"
)

type PushPayload struct {
	UserID    uuid.UUID      `json:"user_id"`
	MessageID string         `json:"message_id"`
	Data      map[string]any `json:"data,omitempty"`
}

type PushAdapter struct {
	client     *Client
	maxRetries int
}

func NewPushAdapter(client *Client, maxRetries int) *PushAdapter {
	return &PushAdapter{client: client, maxRetries: maxRetries}
}

func (a *PushAdapter) SendAsync(ctx context.Context, userID uuid.UUID, message *domain.PushMessage) error {
	payload := PushPayload{
		UserID:    userID,
		MessageID: message.MessageID,
		Data:      message.Data,
	}
	return a.client.EnqueueTask(ctx, TypePushSend, payload, asynq.MaxRetry(a.maxRetries))
}
