package setup

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"

	"go.uber.org/fx"

	"github.com/Deimvis/reactionsstorage/src/loggers"
	"github.com/Deimvis/reactionsstorage/src/servers"
	"github.com/Deimvis/reactionsstorage/src/services"
	"github.com/Deimvis/reactionsstorage/src/storages"
	"github.com/Deimvis/reactionsstorage/src/utils"
	thelpers "github.com/Deimvis/reactionsstorage/tests/helpers"
	rs "github.com/Deimvis/reactionsstorage/tests/simulation/rs_client"
)

var CS *storages.ConfigurationStorage
var RS *storages.ReactionsStorage
var SRV *http.Server
var RSClient rs.Client

func Start() {
	thelpers.CheckEnv("TEST_DATABASE_URL")
	os.Setenv("DATABASE_URL", os.Getenv("TEST_DATABASE_URL"))
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
	RSClient = makeClient(SRV)
	app.Start(ctx)
}

func Stop() {
	ctx := context.Background()
	app.Stop(ctx)
}

func makeClient(s *http.Server) rs.Client {
	host, port, err := net.SplitHostPort(s.Addr)
	if err != nil {
		panic(fmt.Errorf("failed to parse server address: %w", err))
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		panic(fmt.Errorf("failed to parse server port: %w", err))
	}
	return rs.NewClientHTTP(host, portInt, false)
}

var app *fx.App
