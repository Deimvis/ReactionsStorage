package loggers

import (
	"context"
	"errors"
	"fmt"
	"syscall"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewLogger(lc fx.Lifecycle) *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Errorf("failed to create logger: %w", err))
	}
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			err := logger.Sync()
			// https://github.com/uber-go/zap/issues/991
			if err != nil && (!errors.Is(err, syscall.EBADF) && !errors.Is(err, syscall.ENOTTY)) {
				return err
			}
			return nil
		},
	})
	return logger
}
