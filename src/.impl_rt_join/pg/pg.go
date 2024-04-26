package pg

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Deimvis/reactionsstorage/src/configs"
	"github.com/Deimvis/reactionsstorage/src/metrics"
)

type SQLTracer struct {
	log *zap.SugaredLogger
}

func (tracer *SQLTracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	queryId := uuid.New()
	tracer.log.Debugw("SQL", "query_id", queryId.String(), "query", data.SQL, "args", data.Args)
	ctx = context.WithValue(ctx, queryIdCtxKey{}, queryId)
	ctx = context.WithValue(ctx, queryStartCtxKey{}, time.Now())
	return ctx
}

func (tracer *SQLTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	queryId := ctx.Value(queryIdCtxKey{}).(uuid.UUID)
	start := ctx.Value(queryStartCtxKey{}).(time.Time)
	elapsed := float64(time.Since(start)) / float64(time.Second)
	tracer.log.Debugw("SQL done", "query_id", queryId.String(), "start", start.String(), "now", time.Now().String(), "elapsed", elapsed)
	if metrics.SQLReqCnt != nil {
		metrics.SQLReqCnt.Inc()
	}
	if metrics.SQLReqDur != nil {
		metrics.SQLReqDur.Observe(elapsed)
	}
}

func NewPostgresConnectionPool(lc fx.Lifecycle, cfg *configs.ServerConfig, logger *zap.SugaredLogger) *pgxpool.Pool {
	pgxpoolCfg := InitPgxpoolConfig(&cfg.PG.Pool, logger)
	logger.Infow("Pgxpool config was initizlied", "config", fmt.Sprintf("%v", pgxpoolCfg))
	pool, err := pgxpool.NewWithConfig(context.Background(), pgxpoolCfg)
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

func InitPgxpoolConfig(cfg *configs.PGPool, logger *zap.SugaredLogger) *pgxpool.Config {
	config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	if cfg.MinConns != nil {
		config.MinConns = *cfg.MinConns
	}
	if cfg.MaxConns != nil {
		config.MaxConns = *cfg.MaxConns
	}
	if cfg.MaxConnLifetimeJitterS != nil {
		config.MaxConnLifetimeJitter = time.Duration(*cfg.MaxConnLifetimeJitterS) * time.Second
	}
	tracer := &SQLTracer{log: logger}
	config.ConnConfig.Tracer = tracer
	return config
}

type queryIdCtxKey struct{}
type queryStartCtxKey struct{}
