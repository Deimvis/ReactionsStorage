package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRemove(t *testing.T) {
	testCases := []RemoveTC[int]{
		{
			[]int{1, 2, 3},
			func(x int) bool { return x%2 == 0 },
			[]int{1, 3},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			initialCopy := append([]int(nil), tc.initial...)
			actual := Remove(tc.initial, tc.filFn)
			require.Equal(t, tc.expected, actual)
			require.Equal(t, initialCopy, tc.initial) // not changed
		})
	}
}

func TestRemoveIn(t *testing.T) {
	testCases := []RemoveTC[int]{
		{
			[]int{1, 2, 3},
			func(x int) bool { return x%2 == 0 },
			[]int{1, 3},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			actual := RemoveIn(&tc.initial, tc.filFn)
			require.Equal(t, tc.expected, actual)
			require.Equal(t, tc.expected, tc.initial) // changed
		})
	}
}

func TestShuffleIn(t *testing.T) {
	testCases := []struct {
		initial []int
	}{
		{[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			initialCopy := append([]int(nil), tc.initial...)
			actual := ShuffleIn(&tc.initial)
			require.Equal(t, len(initialCopy), len(actual))
			require.NotEqual(t, initialCopy, actual)
			require.Equal(t, tc.initial, actual) // changed

		})
	}
}

type RemoveTC[T any] struct {
	initial  []T
	filFn    func(T) bool
	expected []T
}
