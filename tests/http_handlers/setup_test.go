package http_handlers_test

import (
	"os"
	"testing"

	setup "github.com/Deimvis/reactionsstorage/tests/setup"
)

func TestMain(m *testing.M) {
	setup.Start()
	setFakeConfiguration()
	clearUserReactions()

	// TODO: remove
	// TLDR: I need a fake for small unit tests
	// I don't want to implement fake, so I will use real service
	// I need to make sure that it won't conflict with real server running for integration tests

	code := m.Run()

	setup.Stop()
	os.Exit(code)
}
