package schema

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_ValidFile(t *testing.T) {
	tmpDir := t.TempDir()
	hclContent := `
maschine {
  plugin "testPlugin" {
	source  = "github.com/test/plugin"
	version = "v1.2.3"
  }
}
`
	filePath := filepath.Join(tmpDir, "valid_config.hcl")
	err := os.WriteFile(filePath, []byte(hclContent), 0644)
	assert.NoError(t, err)

	cfg, loadErr := LoadConfig(filePath)
	assert.NoError(t, loadErr)
	assert.NotNil(t, cfg)
	assert.Equal(t, 1, len(cfg.Maschine.Plugins))
	assert.Equal(t, "testPlugin", cfg.Maschine.Plugins[0].Name)
	assert.Equal(t, "github.com/test/plugin", cfg.Maschine.Plugins[0].Source)
	assert.Equal(t, "v1.2.3", cfg.Maschine.Plugins[0].Version)
}

func TestLoadConfig_InvalidFile(t *testing.T) {
	tmpDir := t.TempDir()
	invalidHCL := `maschine { plugin { source = "github.com/test/plugin" version } }`
	filePath := filepath.Join(tmpDir, "invalid_config.hcl")
	err := os.WriteFile(filePath, []byte(invalidHCL), 0644)
	assert.NoError(t, err)

	cfg, loadErr := LoadConfig(filePath)
	assert.Error(t, loadErr)
	assert.Nil(t, cfg)
}

func TestLoadConfig_NonExistentFile(t *testing.T) {
	cfg, err := LoadConfig("nonexistent_config.hcl")
	assert.Error(t, err)
	assert.Nil(t, cfg)
}
