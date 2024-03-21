package utils

import (
	"fmt"
	"math/rand"
	"reflect"
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

// https://github.com/alexanderbez/godash/blob/703c92476f3a6a947b9f2792114ecf40d7ba2c6a/godash.go#L86-L127
func SliceEqual(slice1, slice2 interface{}) bool {
	equal, err := sliceEqual(slice1, slice2)
	if err != nil {
		panic(err)
	}
	return equal
}

func sliceEqual(slice1, slice2 interface{}) (bool, error) {
	if !IsSlice(slice1) {
		return false, fmt.Errorf("argument type '%T' is not a slice", slice1)
	} else if !IsSlice(slice2) {
		return false, fmt.Errorf("argument type '%T' is not a slice", slice2)
	}

	slice1Val := reflect.ValueOf(slice1)
	slice2Val := reflect.ValueOf(slice2)

	if slice1Val.Type().Elem() != slice2Val.Type().Elem() {
		return false, fmt.Errorf("type of '%v' does not match type of '%v'", slice1Val.Type().Elem(), slice2Val.Type().Elem())
	}

	if slice1Val.Len() != slice2Val.Len() {
		return false, nil
	}

	result := true
	i, n := 0, slice1Val.Len()

	for i < n {
		j := 0
		e := false
		for j < n && !e {
			if slice1Val.Index(i).Interface() == slice2Val.Index(j).Interface() {
				e = true
			}
			j++
		}
		if !e {
			result = false
		}
		i++
	}

	return result, nil
}

func IsSlice(value interface{}) bool {
	kind := reflect.ValueOf(value).Kind()
	return kind == reflect.Slice || kind == reflect.Array
}
