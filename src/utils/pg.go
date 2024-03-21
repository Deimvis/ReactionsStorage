package utils

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type myQueryTracer struct {
	log *zap.SugaredLogger
}

func (tracer *myQueryTracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	// TODO: make it look ok
	tracer.log.Infow("Executing command", "sql", data.SQL, "args", data.Args)
	return ctx
}

func (tracer *myQueryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
}

func NewPostgresConnectionPool(lc fx.Lifecycle, logger *zap.SugaredLogger) *pgxpool.Pool {
	// TODO: pass tracer as a dependency
	config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL")) // Using environment variables instead of a connection string.
	if err != nil {
		panic(err)
	}
	tracer := &myQueryTracer{log: logger}
	config.ConnConfig.Tracer = tracer

	// ---

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
