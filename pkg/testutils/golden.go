package testutils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func SaveOrAssertEqual(t *testing.T, got, fileName string, update bool) {
	if update {
		Save(t, got, fileName)
	} else {
		AssertEqual(t, got, fileName)
	}
}

func Save(t *testing.T, got, fileName string) {
	const dirPerm = 0755
	const filePerm = 0600
	// make sure path exists
	require.NoError(t, os.MkdirAll(filepath.Dir(fileName), dirPerm))
	err := os.WriteFile(fileName, []byte(got), filePerm)
	require.NoError(t, err)
}

func AssertEqual(t *testing.T, got, fileName string) {
	want, err := os.ReadFile(fileName)
	require.NoError(t, err)
	assert.Equal(t, string(want), got)
}
