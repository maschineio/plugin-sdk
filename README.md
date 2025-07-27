# pluging-sdk

[![Go Report Card](https://goreportcard.com/badge/github.com/maschineio/plugin-sdk)](https://goreportcard.com/report/github.com/maschineio/plugin-sdk) [![Go Reference](https://pkg.go.dev/badge/maschine.io/plugin-sdk.svg)](https://pkg.go.dev/maschine.io/plugin-sdk) [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=maschineio_plugin-sdk&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=maschineio_plugin-sdk) [![Coverage](https://sonarcloud.io/api/project_badges/measure?project=maschineio_plugin-sdk&metric=coverage)](https://sonarcloud.io/summary/new_code?id=maschineio_plugin-sdk)

Maschine Plugin SDK enables plugins for state machine tasks.

## Installation

```shell
go get maschine.io/plugin-sdk
```

## Schema configuration

Maschine can be configured with the help of a simple `hcl` configuration file, that is used by `maschine init` command, that loads the plugins from the configured scm repository.

This example shows the usage for a configuration.

```hcl
maschine {
  scm {
      type     = "gitlab"
      base_url = "gitlab.company.com"  # optional (use your own gitlab server)
  }

  plugin "hello-world" {
      source  = "myteam/hello-plugin"  # shorter, because host is known and gitlab will be used
      version = "1.0.0"
  }

  plugin "other-plugin" {
      type    = "github"               # overwrites default-SCM (which is in this case gitlab)
      source  = "otherteam/plugin"
      version = "2.1.0"
  }
}
```

## Rules

All plugins must be compressed with zip and must have a `filename.zip` suffix.

## Plugin Development

### Quick Start

```go
package main

import (
    "context"

    "github.com/hashicorp/go-plugin"
    "maschine.io/plugin-sdk/sdk"
    "maschine.io/plugin-sdk/sdk/logger"
)

func main() {
    // Initialize logger from environment
    logger.InitializeFromEnv("my-plugin")

    // Create your plugin
    p := sdk.NewBasePlugin("my-plugin", "1.0.0")

    // Register functions
    p.RegisterSimpleFunction("resource:action", MyFunction, "description")

    // Serve the plugin
    plugin.Serve(&plugin.ServeConfig{
        HandshakeConfig: sdk.Handshake,
        Plugins: map[string]plugin.Plugin{
            "resource": &sdk.ResourcePlugin{Impl: p},
        },
        Logger:     logger.Get(),
        GRPCServer: plugin.DefaultGRPCServer,
    })
}

func MyFunction(ctx context.Context, req *sdk.TypedExecuteRequest) (any, error) {
    // Use the logger
    log := logger.Named("my-function")
    log.Info("Processing request")

    // Your logic here
    return "result", nil
}
```

### Logging

The SDK provides built-in logging support using HashiCorp's hclog:

- Automatic initialization from `MASCHINE_PLUGIN_LOG_LEVEL` environment variable
- Structured JSON logging to STDERR
- Available log levels: trace, debug, info, warn, error

```bash
# Run plugin with debug logging
export MASCHINE_PLUGIN_LOG_LEVEL=debug
```
