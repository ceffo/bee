package common

import "strings"

// CountLines counts the number of lines in a string
func CountLines(str string) int {
	lines := strings.Split(str, "\n")
	return len(lines)
}
