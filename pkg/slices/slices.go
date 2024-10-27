package slices

// Map
func Map(s []string, f func(string) string) []string {
	result := make([]string, len(s))
	for i, v := range s {
		result[i] = f(v)
	}
	return result
}
