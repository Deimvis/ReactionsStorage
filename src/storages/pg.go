package storages

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Deimvis/reactionsstorage/src/utils"
)

type PG interface {
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	Begin(ctx context.Context) (pgx.Tx, error)
}

type PGStorage interface {
	GetPool() *pgxpool.Pool
}

// AcquirePG returns connection from context, if any, and PGStorage's pgxpool otherwise
func AcquirePG(ctx context.Context, s PGStorage) PG {
	var pg PG
	ctxConn := ctx.Value(ConnCtxKey{})
	if ctxConn != nil {
		pg = ctxConn.(*pgxpool.Conn)
	} else {
		pg = s.GetPool()
	}
	return pg
}

func CtxWithConn(ctx context.Context, s PGStorage, fn func(ctx context.Context)) {
	ctx = CtxAcquireConn(ctx, s)
	defer CtxReleaseConn(&ctx) // return conn to the pool
	fn(ctx)
}

func CtxAcquireConn(ctx context.Context, s PGStorage) context.Context {
	conn := utils.Must(s.GetPool().Acquire(ctx))
	return context.WithValue(ctx, ConnCtxKey{}, conn)
}

func CtxReleaseConn(ctx *context.Context) {
	conn := (*ctx).Value(ConnCtxKey{}).(*pgxpool.Conn)
	conn.Release()
	*ctx = context.WithValue(*ctx, ConnCtxKey{}, nil)
}

type ConnCtxKey struct{}
