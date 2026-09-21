package transactor

import (
	"context"
	"database/sql"
	"fmt"
)

// the tx is stored in the context, repositories call ExtractTx() to see if they should use it
type txKey struct{}

type Transactor struct {
	db *sql.DB
}

func New(db *sql.DB) *Transactor {
	return &Transactor{db: db}
}

// runs fn in a transaction: commit if it returns nil, rollback on error or panic
func (t *Transactor) Do(ctx context.Context, fn func(ctx context.Context) error) (err error) {
	// nested Do just joins the outer transaction (no savepoints)
	// if transaction already exists, execute the function
	if _, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return fn(ctx)
	}

	tx, err := t.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p) // re-panic to maintain the original panic
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			commitErr := tx.Commit()
			if commitErr != nil {
				err = fmt.Errorf("failed to commit transaction: %w", commitErr)
			}
		}
	}()

	ctxWithTx := context.WithValue(ctx, txKey{}, tx)

	err = fn(ctxWithTx)
	return err
}

func ExtractTx(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx
	}
	return nil
}
