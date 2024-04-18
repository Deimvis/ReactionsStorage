package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_makeQueryString(t *testing.T) {
	key := 42
	force := true
	testCases := []struct {
		query    interface{}
		expected string
	}{
		{
			q1{key: 42},
			"key=42",
		},
		{
			q2{},
			"",
		},
		{
			q2{key: &key},
			"key=42",
		},
		{
			q3{force: &force},
			"force=true",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			actual := makeQueryString(tc.query)
			require.Equal(t, tc.expected, actual)
		})
	}
}

type q1 struct {
	key int `query:"key"`
}

type q2 struct {
	key *int `query:"key"`
}

type q3 struct {
	force *bool `query:"force"`
}
