package utils

import (
	"math/rand"
)

func Contains[T comparable](s []T, v T) bool {
	return ContainsFunc(s, func(el T) bool { return el == v })
}

func ContainsFunc[T any](s []T, fn func(v T) bool) bool {
	for _, val := range s {
		if fn(val) {
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

// Filter filters slice.
// Keeps only those elements for which given filFn returns true.
func Filter[T any](s []T, pred func(T) bool) []T {
	scopy := make([]T, len(s))
	copy(scopy, s)
	return FilterIn(&scopy, pred)
}

// FilterIn filters slice in-place.
// Keeps only those elements for which given filFn returns true.
func FilterIn[T any](s *[]T, pred func(T) bool) []T {
	newSz := 0
	for i := range *s {
		if pred((*s)[i]) {
			(*s)[newSz] = (*s)[i]
			newSz++
		}
	}
	*s = (*s)[:newSz]
	return *s
}

// FilterMapIn filters map in-place.
// Keeps only those elements for which given filFn returns true.
func FilterMapIn[K comparable, V any](m map[K]V, filFn func(k K, v V) bool) {
	for k, v := range m {
		if !filFn(k, v) {
			delete(m, k)
		}
	}
}
