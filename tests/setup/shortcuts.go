package setup

import (
	"context"

	"github.com/Deimvis/reactionsstorage/src/models"
	"github.com/Deimvis/reactionsstorage/tests/fake"
	"github.com/Deimvis/reactionsstorage/tests/simulation/utils"
)

func SetConfiguration(conf models.Configuration) {
	ctx := context.Background()
	CS.ClearStrict(ctx)
	for _, r := range conf.Reactions {
		CS.AddReaction(ctx, &r)
	}
	for _, rs := range conf.ReactionSets {
		CS.AddReactionSet(ctx, &rs)
	}
	for _, n := range conf.Namespaces {
		CS.AddNamespace(ctx, &n)
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
}

func GetUserReactions() []models.UserReaction {
	return utils.Must(RS.GetUserReactions(context.Background()))
}

func ClearUserReactions() {
	RS.ClearStrict(context.Background())
}
