package must

// NoError panics if err is not nil. It returns v.
func NoError[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
