package schema

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

// DownloadPlugins downloads the plugins specified in the configuration to the given directory.
func (c *Config) DownloadPlugins(dir, operatingSystem, arch string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create plugin directory: %v", err)
	}

	for _, plugin := range c.Maschine.Plugins {
		// Construct download URL
		parts := strings.Split(plugin.Source, "/")
		if len(parts) < 3 {
			return fmt.Errorf("invalid source format for plugin %s", plugin.Name)
		}

		filename := fmt.Sprintf("%s_%s_%s_%s", plugin.Name, plugin.Version, operatingSystem, arch)
		url, err := buildDownloadURL(plugin.Source, plugin.Version, filename)
		if err != nil {
			return fmt.Errorf("URL construction failed for %s: %v", plugin.Name, err)
		}

		// Perform download
		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("download failed for plugin %s: %v", plugin.Name, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("download failed for plugin %s: HTTP %d", plugin.Name, resp.StatusCode)
		}

		// Save binary in the specified directory
		destPath := filepath.Join(dir, filename)
		out, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			return fmt.Errorf("failed to create plugin file %s: %v", destPath, err)
		}
		defer out.Close()

		if _, err := io.Copy(out, resp.Body); err != nil {
			return fmt.Errorf("failed to save plugin %s: %v", plugin.Name, err)
		}
	}

	return nil
}

// buildDownloadURL creates the download URL based on the SCM system
func buildDownloadURL(source, version, filename string) (string, error) {
	parts := strings.Split(source, "/")
	if len(parts) < 3 {
		return "", fmt.Errorf("invalid source format: %s", source)
	}

	host := parts[0]
	owner := parts[1]
	repo := parts[2]

	switch host {
	case "github.com":
		return fmt.Sprintf("https://%s/%s/%s/releases/download/v%s/%s",
			host, owner, repo, version, filename), nil
	case "gitlab.com":
		return fmt.Sprintf("https://%s/%s/%s/-/releases/v%s/downloads/%s",
			host, owner, repo, version, filename), nil
	case "bitbucket.org":
		return fmt.Sprintf("https://%s/%s/%s/downloads/%s",
			host, owner, repo, filename), nil
	default:
		return "", fmt.Errorf("unsupported SCM system: %s", host)
	}
}
