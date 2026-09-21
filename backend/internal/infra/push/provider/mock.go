package provider

import (
	"context"
	"log/slog"
)

type Mock struct {
	log *slog.Logger
}

func NewMock(log *slog.Logger) *Mock {
	return &Mock{log: log}
}

func (p *Mock) Send(ctx context.Context, token string, title, body string, data map[string]any) error {
	p.log.Info("PUSH SEND",
		slog.String("token", token),
		slog.String("title", title),
		slog.String("body", body),
		slog.Any("data", data),
	)
	return nil
}

func (p *Mock) SendMulti(ctx context.Context, tokens []string, title, body string, data map[string]any) ([]string, error) {
	p.log.Info("PUSH SEND MULTI",
		slog.Any("tokens", tokens),
		slog.String("title", title),
		slog.String("body", body),
		slog.Any("data", data),
	)
	return nil, nil
}
