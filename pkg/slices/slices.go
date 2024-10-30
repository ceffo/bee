package slices

// Map applies a function to each element of a slice and returns a new slice with the results.
func Map[T any, U any](s []T, f func(T) U) []U {
	result := make([]U, len(s))
	for i, v := range s {
		result[i] = f(v)
	}
	return result
}

// FoldLeft applies a function to each element of a slice, threading an accumulator argument through the computation.
func FoldLeft[T any, U any](s []T, f func(U, T) U, init U) U {
	result := init
	for _, v := range s {
		result = f(result, v)
	}
	return result
}
