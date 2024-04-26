package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"go.uber.org/fx"

	"github.com/Deimvis/reactionsstorage/src/configs"
	"github.com/Deimvis/reactionsstorage/src/loggers"
	"github.com/Deimvis/reactionsstorage/src/pg"
	"github.com/Deimvis/reactionsstorage/src/servers"
	"github.com/Deimvis/reactionsstorage/src/services"
	"github.com/Deimvis/reactionsstorage/src/storages"
	"github.com/Deimvis/reactionsstorage/src/utils"
)

var cfgFilePath *string

func ParseArgs() {
	cfgFilePath = flag.String("config", "configs/server.yaml", "Path to the server configuration file")
	flag.Parse()
}

func CreateOptions() fx.Option {
	return fx.Options(
		fx.Provide(
			configs.NewServerConfig(cfgFilePath),
			loggers.NewLogger,
			pg.NewPostgresConnectionPool,
			storages.NewConfigurationStorage,
			storages.NewReactionsStorage,
			services.NewConfigurationService,
			services.NewReactionsService,
			servers.NewHTTPServer,
		),
		fx.Invoke(func(s *http.Server) {}),
	)
}

func main() {
	ParseArgs()
	fx.New(CreateOptions()).Run()
}

func setupSigHandlers() {
	if utils.IsDebugEnv() {
		sigChan := make(chan os.Signal)
		go func() {
			for range sigChan {
				debug.PrintStack()
			}
		}()
		signal.Notify(sigChan, syscall.SIGQUIT)
	}
}
