package schema

import (
	"github.com/hashicorp/hcl/v2/hclsimple"
)

type Config struct {
	Maschine MaschineBlock `hcl:"maschine,block"`
}

type MaschineBlock struct {
	Plugins []PluginBlock `hcl:"plugin,block"`
}

type PluginBlock struct {
	Name    string `hcl:",label"`
	Source  string `hcl:"source"`
	Version string `hcl:"version"`
}

// Beispiel, wie man ein solches HCL-File parsen könnte:
func LoadConfig(filename string) (*Config, error) {
	var cfg Config
	err := hclsimple.DecodeFile(filename, nil, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
