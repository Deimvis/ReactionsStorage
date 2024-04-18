package utils

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Deimvis/reactionsstorage/src/metrics"
)

type QueryStartCtxKey struct{}

type SQLTracer struct {
	log *zap.SugaredLogger
}

func (tracer *SQLTracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	tracer.log.Debugw("SQL", "query", data.SQL, "args", data.Args)
	ctx = context.WithValue(ctx, QueryStartCtxKey{}, time.Now())
	return ctx
}

func (tracer *SQLTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	start := ctx.Value(QueryStartCtxKey{}).(time.Time)
	elapsed := float64(time.Since(start)) / float64(time.Second)
	metrics.SQLReqCnt.Inc()
	metrics.SQLReqDur.Observe(elapsed)
}

func NewPostgresConnectionPool(lc fx.Lifecycle, logger *zap.SugaredLogger) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	tracer := &SQLTracer{log: logger}
	config.ConnConfig.Tracer = tracer

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		panic(fmt.Errorf("failed to connect to PostgreSQL database: %w", err))
	}
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			pool.Close()
			return nil
		},
	})
	return pool
}

func NewPostgresConnection(lc fx.Lifecycle) *pgx.Conn {
	con, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(fmt.Errorf("failed to connect to PostgreSQL database: %w", err))
	}
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return con.Close(ctx)
		},
	})
	return con
}
