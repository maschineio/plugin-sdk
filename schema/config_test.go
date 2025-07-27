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

func TestLoadConfig(t *testing.T) {
	t.Run("valid configuration", func(t *testing.T) {
		tmpDir := t.TempDir()
		validConfig := `
maschine {
  scm {
	type = "github"
  }
  plugin "testPlugin" {
	source = "owner/repo"
	version = "1.0.0"
  }
}`
		configPath := filepath.Join(tmpDir, "config.hcl")
		require.NoError(t, os.WriteFile(configPath, []byte(validConfig), 0644))

		cfg, err := LoadConfig(configPath)
		assert.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Equal(t, GitHub, cfg.Maschine.DefaultSCM.Type)
		assert.Equal(t, "testPlugin", cfg.Maschine.Plugins[0].Name)
	})

	t.Run("invalid configuration", func(t *testing.T) {
		tmpDir := t.TempDir()
		invalidConfig := `maschine { invalid }`
		configPath := filepath.Join(tmpDir, "invalid.hcl")
		require.NoError(t, os.WriteFile(configPath, []byte(invalidConfig), 0644))

		cfg, err := LoadConfig(configPath)
		assert.Error(t, err)
		assert.Nil(t, cfg)
	})

	t.Run("non existing file", func(t *testing.T) {
		cfg, err := LoadConfig("non-existent.hcl")
		assert.Error(t, err)
		assert.Nil(t, cfg)
	})
}

func TestBuildDownloadURL(t *testing.T) {
	tests := []struct {
		name        string
		scmType     SCMType
		baseURL     string
		source      string
		version     string
		filename    string
		want        string
		wantErr     bool
		errContains string
	}{
		{
			name:     "GitHub standard URL",
			scmType:  GitHub,
			source:   "owner/repo",
			version:  "1.0.0",
			filename: "plugin.bin.zip",
			want:     "https://github.com/owner/repo/releases/download/v1.0.0/plugin.bin.zip",
		},
		{
			name:     "GitLab standard URL",
			scmType:  GitLab,
			source:   "owner/repo",
			version:  "1.0.0",
			filename: "plugin.bin.zip",
			want:     "https://gitlab.com/owner/repo/-/releases/v1.0.0/downloads/plugin.bin.zip",
		},
		{
			name:     "BitBucket standard URL",
			scmType:  BitBucket,
			source:   "owner/repo",
			version:  "1.0.0",
			filename: "plugin.bin.zip",
			want:     "https://bitbucket.org/owner/repo/downloads/plugin.bin.zip",
		},
		{
			name:     "Custom base URL",
			scmType:  GitHub,
			baseURL:  "git.example.com",
			source:   "owner/repo",
			version:  "1.0.0",
			filename: "plugin.bin.zip",
			want:     "https://git.example.com/owner/repo/releases/download/v1.0.0/plugin.bin.zip",
		},
		{
			name:        "Invalid source format",
			scmType:     GitHub,
			source:      "invalid",
			version:     "1.0.0",
			filename:    "plugin.bin.zip",
			wantErr:     true,
			errContains: "invalid source format",
		},
		{
			name:        "Unsupported SCM type",
			scmType:     "unsupported",
			source:      "owner/repo",
			version:     "1.0.0",
			filename:    "plugin.bin.zip",
			wantErr:     true,
			errContains: "unsupported SCM type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildDownloadURL(tt.scmType, tt.baseURL, tt.source, tt.version, tt.filename)
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

func TestDownloadPlugins(t *testing.T) {
	origBuildDownloadURL := BuildDownloadURL
	defer func() { BuildDownloadURL = origBuildDownloadURL }()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("mock plugin binary"))
	}))
	defer ts.Close()

	tests := []struct {
		name        string
		setupMock   func()
		config      *Config
		dir         string
		os          string
		arch        string
		wantErr     bool
		errContains string
	}{
		{
			name: "successful download",
			setupMock: func() {
				BuildDownloadURL = func(scmType SCMType, baseURL, source, version, filename string) (string, error) {
					return ts.URL, nil
				}
			},
			config: &Config{
				Maschine: MaschineBlock{
					DefaultSCM: SCMConfig{Type: GitHub},
					Plugins: []PluginBlock{
						{Name: "test", Source: "owner/repo", Version: "1.0.0"},
					},
				},
			},
			dir:  t.TempDir(),
			os:   "linux",
			arch: "amd64",
		},
		{
			name: "URL construction failed",
			setupMock: func() {
				BuildDownloadURL = func(scmType SCMType, baseURL, source, version, filename string) (string, error) {
					return "", fmt.Errorf("URL error")
				}
			},
			config: &Config{
				Maschine: MaschineBlock{
					DefaultSCM: SCMConfig{Type: GitHub},
					Plugins: []PluginBlock{
						{Name: "test", Source: "invalid", Version: "1.0.0"},
					},
				},
			},
			dir:         t.TempDir(),
			os:          "linux",
			arch:        "amd64",
			wantErr:     true,
			errContains: "URL construction failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			err := tt.config.DownloadPlugins(tt.dir, tt.os, tt.arch)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				return
			}
			assert.NoError(t, err)

			pluginFile := filepath.Join(tt.dir, fmt.Sprintf("%s_%s_%s_%s.zip",
				tt.config.Maschine.Plugins[0].Name,
				tt.config.Maschine.Plugins[0].Version,
				tt.os,
				tt.arch))

			content, err := os.ReadFile(pluginFile)
			assert.NoError(t, err)
			assert.Equal(t, "mock plugin binary", string(content))
		})
	}
}
