package schema

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigValidFile(t *testing.T) {
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

func TestLoadConfigInvalidFile(t *testing.T) {
	tmpDir := t.TempDir()
	invalidHCL := `maschine { plugin { source = "github.com/test/plugin" version } }`
	filePath := filepath.Join(tmpDir, "invalid_config.hcl")
	err := os.WriteFile(filePath, []byte(invalidHCL), 0644)
	assert.NoError(t, err)

	cfg, loadErr := LoadConfig(filePath)
	assert.Error(t, loadErr)
	assert.Nil(t, cfg)
}

func TestLoadConfigNonExistentFile(t *testing.T) {
	cfg, err := LoadConfig("nonexistent_config.hcl")
	assert.Error(t, err)
	assert.Nil(t, cfg)
}

func TestBuildDownloadURL(t *testing.T) {
	tests := []struct {
		name        string
		source      string
		version     string
		filename    string
		want        string
		wantErr     bool
		errContains string
	}{
		{
			name:     "github url",
			source:   "github.com/owner/repo",
			version:  "1.0.0",
			filename: "plugin_1.0.0",
			want:     "https://github.com/owner/repo/releases/download/v1.0.0/plugin_1.0.0",
		},
		{
			name:     "gitlab url",
			source:   "gitlab.com/owner/repo",
			version:  "1.0.0",
			filename: "plugin_1.0.0",
			want:     "https://gitlab.com/owner/repo/-/releases/v1.0.0/downloads/plugin_1.0.0",
		},
		{
			name:        "invalid source",
			source:      "invalid/source",
			version:     "1.0.0",
			filename:    "plugin_1.0.0",
			wantErr:     true,
			errContains: "invalid source format",
		},
		{
			name:        "unsupported host",
			source:      "unsupported.com/owner/repo",
			version:     "1.0.0",
			filename:    "plugin_1.0.0",
			wantErr:     true,
			errContains: "unsupported SCM system",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildDownloadURL(tt.source, tt.version, tt.filename)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
