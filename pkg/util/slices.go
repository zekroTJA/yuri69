package util

import (
	"sort"

	"golang.org/x/exp/constraints"
)

func IndexOf[T comparable](s []T, v T) int {
	if len(s) == 0 {
		return -1
	}

	for i, c := range s {
		if c == v {
			return i
		}
	}

	return -1
}

func Contains[T comparable](s []T, v T) bool {
	return IndexOf(s, v) != -1
}

func ContainsAny[T comparable](s []T, v []T) bool {
	if len(s) == 0 || len(v) == 0 {
		return false
	}

	for _, c := range v {
		if Contains(s, c) {
			return true
		}
	}

	return false
}

func ContainsAll[T comparable](s []T, v []T) bool {
	if len(s) == 0 && len(v) == 0 {
		return true
	}
	if len(s) == 0 {
		return false
	}

	for _, c := range v {
		if !Contains(s, c) {
			return false
		}
	}

	return true
}

func HasDuplicates[T constraints.Ordered](v []T) bool {
	if len(v) < 2 {
		return false
	}

	s := make([]T, len(v))
	copy(s, v)

	sort.Slice(s, func(i, j int) bool {
		return s[i] < s[j]
	})

	for i := 1; i < len(s); i++ {
		if s[i-1] == s[i] {
			return true
		}
	}

	return false
}

func ApplyToAll[T any](s []T, f func(v T) T) {
	for i, v := range s {
		s[i] = f(v)
	}
}

func AppendIfNotContains[T comparable](s []T, v T) []T {
	if !Contains(s, v) {
		s = append(s, v)
	}
	return s
}

func Remove[T comparable](s []T, v T) []T {
	i := IndexOf(s, v)
	if i != -1 {
		s = append(s[:i], s[i+1:]...)
	}
	return s
}
