package mocks

import (
	"context"
)

type TransactionManagerFake struct{}

func (f *TransactionManagerFake) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}
