package setup

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Deimvis/reactionsstorage/src/loggers"
	"github.com/Deimvis/reactionsstorage/src/servers"
	"github.com/Deimvis/reactionsstorage/src/services"
	"github.com/Deimvis/reactionsstorage/src/storages"
	"github.com/Deimvis/reactionsstorage/src/utils"
	thelpers "github.com/Deimvis/reactionsstorage/tests/helpers"
	rs "github.com/Deimvis/reactionsstorage/tests/simulation/rs_client"
)

var Logger *zap.SugaredLogger
var CS *storages.ConfigurationStorage
var RS *storages.ReactionsStorage
var SRV *http.Server
var RSClient rs.Client

func Start() {
	Logger = makeLogger()
	zap.ReplaceGlobals(Logger.Desugar())
	thelpers.CheckEnv("TEST_DATABASE_URL")
	thelpers.CheckEnv("TEST_PORT")
	os.Setenv("DATABASE_URL", os.Getenv("TEST_DATABASE_URL"))
	os.Setenv("PORT", os.Getenv("TEST_PORT"))
	os.Setenv("DEBUG", "1")
	ctx := context.Background()
	app = fx.New(
		fx.Provide(
			loggers.NewLogger,
			utils.NewPostgresConnectionPool,
			storages.NewConfigurationStorage,
			storages.NewReactionsStorage,
			services.NewConfigurationService,
			services.NewReactionsService,
			servers.NewHTTPServer,
		),
		fx.Invoke(func(c *storages.ConfigurationStorage) {
			CS = c
		}),
		fx.Invoke(func(r *storages.ReactionsStorage) {
			RS = r
		}),
		fx.Invoke(func(s *http.Server) {
			SRV = s
		}),
	)
	RSClient = makeClient(SRV, Logger)
	app.Start(ctx)
}

func Stop() {
	defer Logger.Sync()
	ctx := context.Background()
	app.Stop(ctx)
}

func makeClient(s *http.Server, logger *zap.SugaredLogger) rs.Client {
	host, port, err := net.SplitHostPort(s.Addr)
	if err != nil {
		panic(fmt.Errorf("failed to parse server address: %w", err))
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		panic(fmt.Errorf("failed to parse server port: %w", err))
	}
	return rs.NewClientHTTP(host, portInt, false, logger, nil)
}

func makeLogger() *zap.SugaredLogger {
	config := zap.NewDevelopmentConfig()
	config.Sampling = nil
	config.Level.SetLevel(zap.DebugLevel)
	return zap.Must(config.Build()).Sugar()
}

var app *fx.App
