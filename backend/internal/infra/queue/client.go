package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

type Client struct {
	enqueuer *asynq.Client
}

func NewClient(redisOpt asynq.RedisClientOpt) *Client {
	enqueuer := asynq.NewClient(redisOpt)
	return &Client{enqueuer: enqueuer}
}

func (c *Client) Close() error {
	return c.enqueuer.Close()
}

func (c *Client) EnqueueTask(ctx context.Context, taskType string, payload any, opts ...asynq.Option) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(taskType, payloadBytes)

	// enqueue task with options
	if _, err = c.enqueuer.EnqueueContext(ctx, task, opts...); err != nil {
		return fmt.Errorf("failed to enqueue task %s: %w", taskType, err)
	}

	return nil
}
