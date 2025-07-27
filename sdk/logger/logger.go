package logger

import (
	"os"
	"sync"

	"github.com/hashicorp/go-hclog"
)

var (
	instance hclog.Logger
	once     sync.Once
)

// Initialize sets up the global logger instance
// This is called automatically when a plugin starts
func Initialize(logger hclog.Logger) {
	once.Do(func() {
		instance = logger
	})
}

// InitializeFromEnv initializes the logger from environment variables
// Uses MASCHINE_PLUGIN_LOG_LEVEL to set the log level
func InitializeFromEnv(name string) hclog.Logger {
	logLevel := os.Getenv("MASCHINE_PLUGIN_LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	level := hclog.LevelFromString(logLevel)
	
	logger := hclog.New(&hclog.LoggerOptions{
		Name:       name,
		Level:      level,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	Initialize(logger)
	
	// Log initialization
	logger.Info("Plugin logger initialized", 
		"name", name,
		"log_level", logLevel,
		"configured_level", level.String())
	
	return logger
}

// Get returns the global logger instance
// If not initialized, returns a default logger
func Get() hclog.Logger {
	if instance == nil {
		// Fallback to default logger if not initialized
		return hclog.Default()
	}
	return instance
}

// Named returns a named sub-logger
func Named(name string) hclog.Logger {
	return Get().Named(name)
}

// With returns a logger with additional fields
func With(args ...interface{}) hclog.Logger {
	return Get().With(args...)
}