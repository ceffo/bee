package testutils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SaveOrAssertEqual asserts that the given content is equal to the contents of fileName.
// If update is true, the contents of fileName are updated.
func SaveOrAssertEqual(t *testing.T, content, fileName string, update bool) {
	if update {
		save(t, content, fileName)
	}
	AssertEqual(t, content, fileName)
}

func save(t *testing.T, got, fileName string) {
	const dirPerm = 0755
	const filePerm = 0600
	// make sure path exists
	require.NoError(t, os.MkdirAll(filepath.Dir(fileName), dirPerm))
	err := os.WriteFile(fileName, []byte(got), filePerm)
	require.NoError(t, err)
}

// AssertEqual asserts that the given content is equal to the contents of fileName.
func AssertEqual(t *testing.T, content, fileName string) {
	want, err := os.ReadFile(fileName)
	require.NoError(t, err)
	assert.Equal(t, string(want), content)
}
