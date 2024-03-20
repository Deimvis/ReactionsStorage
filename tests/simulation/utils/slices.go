package utils

import (
	"math/rand"
)

func Contains[T comparable](s []T, v T) bool {
	for _, val := range s {
		if val == v {
			return true
		}
	}
	return false
}

// Map slice
func Map[T, U any](s []T, mapFn func(T) U) []U {
	res := make([]U, len(s))
	for i := range s {
		res[i] = mapFn(s[i])
	}
	return res
}

// Map slice in-place
func MapIn[T any](s []T, mapFn func(T) T) []T {
	for i := range s {
		s[i] = mapFn(s[i])
	}
	return s
}

// Shuffle slice in-place
func ShuffleIn[T any](sp *[]T) []T {
	s := *sp
	rand.Shuffle(len(s), func(i, j int) { s[i], s[j] = s[j], s[i] })
	return s
}

// Filters slice.
// Removes all elements for which given filFn returns true.
func Filter[T any](s []T, filFn func(T) bool) []T {
	scopy := make([]T, len(s))
	copy(scopy, s)
	return FilterIn(&scopy, filFn)
}

// Filters slice in-place.
// Removes all elements for which given filFn returns true.
func FilterIn[T any](s *[]T, filFn func(T) bool) []T {
	newSz := 0
	for i := range *s {
		if !filFn((*s)[i]) {
			(*s)[newSz] = (*s)[i]
			newSz++
		}
	}
	*s = (*s)[:newSz]
	return *s
}
