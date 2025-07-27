package logger

import (
	"os"
	"sync"
	"testing"

	"github.com/hashicorp/go-hclog"
)

func TestInitializeFromEnv(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected hclog.Level
	}{
		{"default to info", "", hclog.Info},
		{"debug level", "debug", hclog.Debug},
		{"trace level", "trace", hclog.Trace},
		{"warn level", "warn", hclog.Warn},
		{"error level", "error", hclog.Error},
		{"invalid defaults to info", "invalid", hclog.NoLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset singleton
			instance = nil
			once = sync.Once{}
			
			// Set environment
			if tt.envValue != "" {
				os.Setenv("MASCHINE_PLUGIN_LOG_LEVEL", tt.envValue)
				defer os.Unsetenv("MASCHINE_PLUGIN_LOG_LEVEL")
			}

			// Initialize logger
			logger := InitializeFromEnv("test-plugin")
			
			// Get should return the same instance
			if Get() != logger {
				t.Error("Get() should return initialized logger")
			}
			
			// Named should work
			named := Named("test-component")
			if named == nil {
				t.Error("Named() should return a logger")
			}
		})
	}
}

func TestGetWithoutInitialization(t *testing.T) {
	// Reset singleton
	instance = nil
	once = sync.Once{}
	
	// Should return default logger
	logger := Get()
	if logger == nil {
		t.Error("Get() should return default logger when not initialized")
	}
}

func TestSingletonBehavior(t *testing.T) {
	// Reset singleton
	instance = nil
	once = sync.Once{}
	
	// Initialize multiple times
	logger1 := hclog.NewNullLogger()
	logger2 := hclog.NewNullLogger()
	
	Initialize(logger1)
	Initialize(logger2) // Should be ignored
	
	// Should return first logger
	if Get() != logger1 {
		t.Error("Initialize should only set logger once")
	}
}