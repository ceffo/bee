package logging

import (
	"os"
	"time"

	"github.com/charmbracelet/log"
)

type cleanupFunc func()

// SetupFileLogging sets up file logging
func SetupFileLogging(fileName string) (cleanupFunc, error) {
	const filePermissions = 0o644
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, filePermissions)
	if err != nil {
		return nil, err
	}
	log.SetOutput(f)
	log.SetFormatter(log.JSONFormatter)
	log.SetTimeFormat(time.RFC3339Nano)
	return func() {
		log.Info("Closing log file")
		f.Close()
	}, nil
}

// SetupStdoutLogging sets up stdout logging
func SetupStdoutLogging() (cleanupFunc, error) {
	log.SetOutput(os.Stdout)
	log.SetFormatter(log.TextFormatter)
	return func() {}, nil
}
