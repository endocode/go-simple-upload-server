package logging

import (
	"log"
	"os"

	"github.com/hashicorp/logutils"
)

// Start starts the logging
func Start(debug bool) *log.Logger {
	// Logging preferences
	ll := "INFO"
	if debug {
		ll = "DEBUG"
	}

	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "ERROR"},
		MinLevel: logutils.LogLevel(ll),
		Writer:   os.Stderr,
	}

	logger := log.New(os.Stderr, "", 1)
	logger.SetOutput(filter)
	logger.Printf("[INFO] Starting logging.")
	logger.Printf("[INFO] Setting Log Level to %v.", ll)
	return logger
}
