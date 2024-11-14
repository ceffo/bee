package slices

// FoldLeft applies a function to each element of a slice, threading an accumulator argument through the computation.
func FoldLeft[T any, U any](s []T, init U, f func(U, T) U) U {
	result := init
	for _, v := range s {
		result = f(result, v)
	}
	return result
}
