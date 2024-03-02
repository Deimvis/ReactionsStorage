package utils

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

func NewPostgresConnectionPool(lc fx.Lifecycle) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL database: %s\n", err.Error())
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
		log.Fatalf("Failed to connect to PostgreSQL database: %s\n", err.Error())
	}
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return con.Close(ctx)
		},
	})
	return con
}
