package models

import (
	"os"
	"testing"

	setup "github.com/Deimvis/reactionsstorage/tests/setup"
)

func TestMain(m *testing.M) {
	setup.Start()
	setup.SetFakeConfiguration()

	code := m.Run()

	setup.Stop()
	os.Exit(code)
}
