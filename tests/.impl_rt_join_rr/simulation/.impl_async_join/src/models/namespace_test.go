package models

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Deimvis/reactionsstorage/tests/fake"
	"github.com/Deimvis/reactionsstorage/tests/setup"
)

func TestNamespace_GetAvailableReactionIds(t *testing.T) {
	n := NewNamespace(fake.NamespaceId, setup.RSClient)
	require.Equal(t, fake.ReactionSet.ReactionIds, n.GetAvailableReactionIds())
}
