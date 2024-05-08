package setup

import (
	"context"

	"github.com/Deimvis/reactionsstorage/src/models"
	"github.com/Deimvis/reactionsstorage/src/utils"
	"github.com/Deimvis/reactionsstorage/tests/fake"
)

func SetConfiguration(conf models.Configuration) {
	ctx := context.Background()
	must := utils.Must0
	must(CS.Clear(ctx))
	for _, r := range conf.Reactions {
		must(CS.AddReaction(ctx, &r))
	}
	for _, rs := range conf.ReactionSets {
		must(CS.AddReactionSet(ctx, &rs))
	}
	for _, n := range conf.Namespaces {
		must(CS.AddNamespace(ctx, &n))
	}
}

func SetFakeConfiguration() {
	SetConfiguration(models.Configuration{
		Reactions: []models.Reaction{
			fake.Reaction1,
			fake.Reaction2,
			fake.Reaction3,
			fake.FreeReaction1,
			fake.FreeReaction2,
			fake.FreeReaction3,
			fake.FreeReaction4,
			fake.FreeReaction5,
			fake.FreeReaction6,
		},
		ReactionSets: []models.ReactionSet{
			fake.ReactionSet,
		},
		Namespaces: []models.Namespace{
			fake.Namespace,
		},
	})
}

func SetUserReactions(urs []models.UserReaction) {
	ClearUserReactions()
	ctx := context.Background()
	for _, ur := range urs {
		n := utils.Must(CS.GetNamespace(ctx, ur.NamespaceId))
		utils.Must0(RS.AddUserReaction(ctx, ur, n.MaxUniqReactions, n.MutuallyExclusiveReactions, false))
	}
	utils.Must0(RS.RefreshEntityReactions(ctx))
}

func GetUserReactions() []models.UserReaction {
	return utils.Must(RS.GetUserReactions(context.Background()))
}

func ClearUserReactions() {
	utils.Must0(RS.Clear(context.Background()))
}
