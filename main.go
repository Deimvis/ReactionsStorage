package main

import (
	"net/http"

	"go.uber.org/fx"

	"github.com/Deimvis/reactionsstorage/src/servers"
	"github.com/Deimvis/reactionsstorage/src/services"
	"github.com/Deimvis/reactionsstorage/src/storages"
	"github.com/Deimvis/reactionsstorage/src/utils"
)

func CreateOptions() fx.Option {
	return fx.Options(
		fx.Provide(
			utils.NewPostgresConnectionPool,
			storages.NewConfigurationStorage,
			storages.NewReactionsStorage,
			services.NewReactionsService,
			servers.NewHTTPServer,
		),
		fx.Invoke(func(s *http.Server) {}),
	)
}

func main() {
	fx.New(CreateOptions()).Run()
}
