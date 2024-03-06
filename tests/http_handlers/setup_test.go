package http_handlers_test

import (
	"context"
	"net/http"
	"os"
	"testing"

	"go.uber.org/fx"

	"github.com/Deimvis/reactionsstorage/src/loggers"
	"github.com/Deimvis/reactionsstorage/src/servers"
	"github.com/Deimvis/reactionsstorage/src/services"
	"github.com/Deimvis/reactionsstorage/src/storages"
	"github.com/Deimvis/reactionsstorage/src/utils"
	thelpers "github.com/Deimvis/reactionsstorage/tests/helpers"
)

var cs *storages.ConfigurationStorage
var rs *storages.ReactionsStorage
var srv *http.Server

func TestMain(m *testing.M) {
	thelpers.CheckEnv("TEST_DATABASE_URL")
	os.Setenv("DATABASE_URL", os.Getenv("TEST_DATABASE_URL"))
	ctx := context.Background()
	app := fx.New(
		fx.Provide(
			loggers.NewLogger,
			utils.NewPostgresConnectionPool,
			storages.NewConfigurationStorage,
			storages.NewReactionsStorage,
			services.NewConfigurationService,
			services.NewReactionsService,
			servers.NewHTTPServer,
		),
		fx.Invoke(func(s *http.Server) {}),
		fx.Invoke(func(c *storages.ConfigurationStorage) {
			cs = c
		}),
		fx.Invoke(func(r *storages.ReactionsStorage) {
			rs = r
		}),
		fx.Invoke(func(s *http.Server) {
			srv = s
		}),
	)
	app.Start(ctx)
	setFakeConfiguration()
	clearUserReactions()

	code := m.Run()

	app.Stop(ctx)
	os.Exit(code)
}
