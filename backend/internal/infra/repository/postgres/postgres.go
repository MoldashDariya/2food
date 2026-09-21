package postgres

import (
	"context"
	"database/sql"
	"encoding/json"

	"gitlab.com/2foood/backend/internal/infra/transactor"
)

const (
	pgUniqueErrCode = "23505"
)

type scanner interface {
	Scan(dest ...any) error
}

type Queryer interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type BaseRepo struct {
	db *sql.DB
}

func NewBaseRepo(db *sql.DB) BaseRepo {
	return BaseRepo{db: db}
}

func (r *BaseRepo) Queryer(ctx context.Context) Queryer {
	if tx := transactor.ExtractTx(ctx); tx != nil {
		return tx
	}
	return r.db
}

func unmarshalJSONB[T any](data []byte, dest **T) error {
	if data == nil {
		return nil
	}

	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*dest = &v
	return nil
}

func unmarshalJSONBValue[T any](data []byte, dest *T) error {
	if data == nil {
		return nil
	}
	return json.Unmarshal(data, dest)
}
