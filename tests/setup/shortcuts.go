package setup

import (
	"context"

	"github.com/Deimvis/reactionsstorage/tests/fake"
)

func SetFakeConfiguration() {
	ctx := context.Background()
	CS.ClearStrict(ctx)
	CS.AddReactionStrict(ctx, &fake.Reaction)
	CS.AddReactionStrict(ctx, &fake.Reaction2)
	CS.AddReactionStrict(ctx, &fake.Reaction3)
	CS.AddReactionSetStrict(ctx, &fake.ReactionSet)
	CS.AddNamespaceStrict(ctx, &fake.Namespace)
}

func ClearUserReactions() {
	RS.ClearStrict(context.Background())
}
