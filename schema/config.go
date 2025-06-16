package schema

import (
	"github.com/hashicorp/hcl/v2/hclsimple"
)

// Config represents the configuration structure for the Maschine plugin system.
type Config struct {
	Maschine MaschineBlock `hcl:"maschine,block"`
}

// MaschineBlock represents the top-level block in the configuration file.
type MaschineBlock struct {
	Plugins []PluginBlock `hcl:"plugin,block"`
}

// PluginBlock represents a plugin configuration block within the Maschine block.
type PluginBlock struct {
	Name    string `hcl:",label"`
	Source  string `hcl:"source"`
	Version string `hcl:"version"`
}

// LoadConfig loads a configuration from the specified HCL file.
func LoadConfig(filename string) (*Config, error) {
	var cfg Config
	err := hclsimple.DecodeFile(filename, nil, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
