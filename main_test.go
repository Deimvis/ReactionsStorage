package main

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestDepGraphCorrectness(t *testing.T) {
	err := fx.ValidateApp(CreateOptions())
	require.NoError(t, err)
}

func TestLaunch(t *testing.T) {
	app := fxtest.New(t, CreateOptions())
	app.RequireStart()
	app.RequireStop()
}
