package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Basisdata interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, argumen ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, argumen ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, argumen ...any) pgx.Row
}
