package schema

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

// Variable to store the buildDownloadURL function
var buildDownloadURLFunc = buildDownloadURL

func TestDownloadPlugins(t *testing.T) {
	// Original buildDownloadURL sichern
	originalBuildDownloadURL := BuildDownloadURLFunc
	// Nach dem Test wiederherstellen
	defer func() { BuildDownloadURLFunc = originalBuildDownloadURL }()

	// Mock Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("mock plugin binary"))
	}))
	defer ts.Close()

	// buildDownloadURL für den Test überschreiben
	BuildDownloadURLFunc = func(source, version, filename string) (string, error) {
		return ts.URL, nil
	}

	tests := []struct {
		name        string
		config      *Config
		dir         string
		os          string
		arch        string
		setupDir    bool
		wantErr     bool
		errContains string
	}{
		{
			name: "erfolgreicher Download",
			config: &Config{
				Maschine: MaschineBlock{
					Plugins: []PluginBlock{
						{
							Name:    "test-plugin",
							Source:  "github.com/owner/repo",
							Version: "1.0.0",
						},
					},
				},
			},
			dir:      t.TempDir(),
			os:       "linux",
			arch:     "amd64",
			setupDir: true,
		},
		{
			name: "fehlerhaftes Source-Format",
			config: &Config{
				Maschine: MaschineBlock{
					Plugins: []PluginBlock{
						{
							Name:    "test-plugin",
							Source:  "invalid/source",
							Version: "1.0.0",
						},
					},
				},
			},
			dir:         t.TempDir(),
			os:          "linux",
			arch:        "amd64",
			setupDir:    true,
			wantErr:     true,
			errContains: "invalid source format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupDir {
				err := os.MkdirAll(tt.dir, 0755)
				require.NoError(t, err)
			}

			err := tt.config.DownloadPlugins(tt.dir, tt.os, tt.arch)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				return
			}

			assert.NoError(t, err)
			expectedFile := filepath.Join(tt.dir,
				fmt.Sprintf("%s_%s_%s_%s",
					tt.config.Maschine.Plugins[0].Name,
					tt.config.Maschine.Plugins[0].Version,
					tt.os,
					tt.arch))

			content, err := os.ReadFile(expectedFile)
			assert.NoError(t, err)
			assert.Equal(t, "mock plugin binary", string(content))
		})
	}
}
