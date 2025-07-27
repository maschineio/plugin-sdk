// Package logger provides a global logger infrastructure for Maschine.io plugins.
//
// The logger is based on HashiCorp's hclog and provides structured logging
// with configurable log levels via environment variables.
//
// Usage:
//
// In your plugin's main function:
//
//	func main() {
//		// Initialize logger from environment
//		logger.InitializeFromEnv("my-plugin")
//		
//		// Create and run your plugin...
//	}
//
// In your plugin functions:
//
//	func MyFunction(ctx context.Context, req *sdk.TypedExecuteRequest) (any, error) {
//		log := logger.Named("my-function")
//		log.Info("Processing request")
//		// ...
//	}
//
// Configuration:
//
// Set the log level using the MASCHINE_PLUGIN_LOG_LEVEL environment variable:
//
//	export MASCHINE_PLUGIN_LOG_LEVEL=debug
//
// Supported levels: trace, debug, info, warn, error
//
// The logger outputs JSON-formatted logs to STDERR by default.
package logger