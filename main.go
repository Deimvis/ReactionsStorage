package main

import (
	"go.uber.org/fx"

	"github.com/Deimvis/reactionsstorage/src/servers"
	"github.com/Deimvis/reactionsstorage/src/services"
	"github.com/Deimvis/reactionsstorage/src/storages"
	"github.com/Deimvis/reactionsstorage/src/utils"
)

func main() {
	fx.New(
		fx.Provide(
			utils.NewPostgresConnectionPool,
			storages.NewReactionsStorage,
			services.NewReactionsService,
		),
		fx.Invoke(servers.NewHTTPServer),
	).Run()
}
