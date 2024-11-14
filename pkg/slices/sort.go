package slices

import (
	"slices"
)

// Less is a weak ordering comparison function.
type Less[T any] func(a, b T) int

// compose returns a comparison function that composes a list of comparison functions.
func compose[T any](cmps ...Less[T]) Less[T] {
	return func(a, b T) int {
		for _, cmp := range cmps {
			if r := cmp(a, b); r != 0 {
				return r
			}
		}
		return 0
	}
}

// SortBy sorts a slice composing multiple comparison functions.
func SortBy[T any](s []T, cmps ...Less[T]) {
	slices.SortFunc(s, compose(cmps...))
}
