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
      base_url = "gitlab.company.com"  # Optional
  }
  plugin "hello-world" {
      source  = "myteam/hello-plugin"  # shorter, because host is known
      version = "1.0.0"
      # type = "github"  # Optional overwrite scm type
  }
  plugin "other-plugin" {
      type    = "github"               # overwrites default-SCM
      source  = "otherteam/plugin"
      version = "2.1.0"
  }
}
```
