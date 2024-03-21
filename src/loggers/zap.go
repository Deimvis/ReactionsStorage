package loggers

import (
	"context"
	"errors"
	"syscall"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Deimvis/reactionsstorage/tests/simulation/utils"
)

func NewLogger(lc fx.Lifecycle) *zap.SugaredLogger {
	config := zap.NewDevelopmentConfig()
	config.Sampling = nil
	level := zap.InfoLevel
	if utils.IsDebugEnv() {
		level = zap.DebugLevel
	}
	config.Level.SetLevel(level)
	logger := zap.Must(config.Build()).Sugar()
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
