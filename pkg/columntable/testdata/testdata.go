package testdata

import (
	"path/filepath"
	"runtime"
	"strings"
)

// CurrentPath returns the current path
func CurrentPath() string {
	_, file, _, _ := runtime.Caller(0)
	// get folder
	return filepath.Dir(file)
}

// Absolute returns the full path of a file
func Absolute(relativePath string) string {
	return filepath.Join(CurrentPath(), relativePath)
}

// Golden returns the full path of a golden file
func Golden(relativePath string) string {
	// transform spaces to _
	relativePath = strings.Replace(relativePath, " ", "_", -1)
	return Absolute(relativePath + ".golden")
}

func ChangeExtension(fileName string, newExtension string) string {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))] + newExtension
}
