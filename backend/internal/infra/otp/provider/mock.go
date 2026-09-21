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

func (p *Mock) Send(ctx context.Context, to string, code string) error {
	p.log.Info("OTP CODE", slog.String("to", to), slog.String("code", code))
	return nil
}
