package manifest

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	// DefaultManifestFile is the default manifest filename
	DefaultManifestFile = "plugin-manifest.json"
	
	// SchemaURL is the URL for the manifest schema
	SchemaURL = "https://maschine.io/schemas/plugin-manifest-v2.json"
	
	// CurrentVersion is the current manifest version
	CurrentVersion = "2.0"
)

// Load loads a plugin manifest from a file
func Load(path string) (*PluginManifest, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open manifest file: %w", err)
	}
	defer file.Close()

	return Read(file)
}

// Read reads a plugin manifest from a reader
func Read(r io.Reader) (*PluginManifest, error) {
	var manifest PluginManifest
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&manifest); err != nil {
		return nil, fmt.Errorf("failed to decode manifest: %w", err)
	}

	// Validate after loading
	if err := manifest.Validate(); err != nil {
		return nil, fmt.Errorf("manifest validation failed: %w", err)
	}

	return &manifest, nil
}

// Save saves a plugin manifest to a file
func (m *PluginManifest) Save(path string) error {
	// Ensure schema and version are set
	if m.Schema == "" {
		m.Schema = SchemaURL
	}
	if m.ManifestVersion == "" {
		m.ManifestVersion = CurrentVersion
	}

	// Validate before saving
	if err := m.Validate(); err != nil {
		return fmt.Errorf("manifest validation failed: %w", err)
	}

	// Create directory if needed
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create manifest file: %w", err)
	}
	defer file.Close()

	return m.Write(file)
}

// Write writes a plugin manifest to a writer
func (m *PluginManifest) Write(w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(m); err != nil {
		return fmt.Errorf("failed to encode manifest: %w", err)
	}
	return nil
}

// FindManifest searches for a manifest file in common locations
func FindManifest(dir string) (string, error) {
	// Check for new manifest first
	candidates := []string{
		filepath.Join(dir, DefaultManifestFile),
		filepath.Join(dir, "plugin-manifest-v2.json"),
		filepath.Join(dir, "manifest.json"),
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("no manifest file found in %s", dir)
}

// New creates a new plugin manifest with defaults
func New(pluginName, pluginID string) *PluginManifest {
	return &PluginManifest{
		Schema:          SchemaURL,
		ManifestVersion: CurrentVersion,
		Plugin: PluginInfo{
			ID:          pluginID,
			Name:        pluginName,
			DisplayName: pluginName,
			Version:     "0.1.0",
			Category:    "general",
			Tags:        []string{},
			Author: AuthorInfo{
				Name:  "Unknown",
				Email: "unknown@example.com",
			},
			License: "Apache-2.0",
		},
		Runtime: RuntimeInfo{
			Type: "native",
			Executable: ExecutableInfo{
				Path: fmt.Sprintf("./%s", pluginName),
			},
			Protocol:        "grpc",
			ProtocolVersion: "1.0",
		},
		Requirements: Requirements{
			MaschineVersion: ">=1.0.0",
			OS:              []string{"darwin", "linux", "windows"},
			Arch:            []string{"amd64", "arm64"},
		},
		Configuration: Configuration{
			Environment: []EnvVar{},
			Credentials: []CredentialSet{},
		},
		Resources: []ResourceDef{},
		Capabilities: Capabilities{
			HealthCheck:         true,
			ConcurrentExecution: true,
			Stateless:          true,
		},
		Limits: Limits{
			StartupTimeout:          Duration{30 * 1000000000}, // 30s
			ExecuteTimeout:          Duration{300 * 1000000000}, // 5m
			HealthCheckInterval:     Duration{60 * 1000000000}, // 60s
			HealthCheckTimeout:      Duration{10 * 1000000000}, // 10s
			MaxConcurrentExecutions: 100,
		},
		Documentation: Documentation{},
	}
}