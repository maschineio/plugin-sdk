package manifest

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestManifestValidation(t *testing.T) {
	tests := []struct {
		name    string
		modify  func(*PluginManifest)
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid manifest",
			modify:  func(m *PluginManifest) {},
			wantErr: false,
		},
		{
			name: "missing manifest version",
			modify: func(m *PluginManifest) {
				m.ManifestVersion = ""
			},
			wantErr: true,
			errMsg:  "manifestVersion: is required",
		},
		{
			name: "invalid plugin ID",
			modify: func(m *PluginManifest) {
				m.Plugin.ID = "invalid_id"
			},
			wantErr: true,
			errMsg:  "must follow reverse domain notation",
		},
		{
			name: "missing resources",
			modify: func(m *PluginManifest) {
				m.Resources = []ResourceDef{}
			},
			wantErr: true,
			errMsg:  "at least one resource must be defined",
		},
		{
			name: "invalid MRN",
			modify: func(m *PluginManifest) {
				m.Resources[0].Type = "invalid-mrn"
			},
			wantErr: true,
			errMsg:  "must be a valid MRN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestManifest()
			tt.modify(m)
			
			err := m.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("expected error containing '%s', got '%s'", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestManifestReadWrite(t *testing.T) {
	m := createTestManifest()
	
	// Write to buffer
	var buf bytes.Buffer
	err := m.Write(&buf)
	if err != nil {
		t.Fatalf("failed to write manifest: %v", err)
	}
	
	// Read back
	m2, err := Read(&buf)
	if err != nil {
		t.Fatalf("failed to read manifest: %v", err)
	}
	
	// Compare key fields
	if m2.Plugin.ID != m.Plugin.ID {
		t.Errorf("plugin ID mismatch: got %s, want %s", m2.Plugin.ID, m.Plugin.ID)
	}
	if m2.Plugin.Name != m.Plugin.Name {
		t.Errorf("plugin name mismatch: got %s, want %s", m2.Plugin.Name, m.Plugin.Name)
	}
	if len(m2.Resources) != len(m.Resources) {
		t.Errorf("resource count mismatch: got %d, want %d", len(m2.Resources), len(m.Resources))
	}
}

func TestDurationJSON(t *testing.T) {
	d := Duration{30 * time.Second}
	
	// Marshal
	data, err := d.MarshalJSON()
	if err != nil {
		t.Fatalf("failed to marshal duration: %v", err)
	}
	
	expected := `"30s"`
	if string(data) != expected {
		t.Errorf("unexpected JSON: got %s, want %s", string(data), expected)
	}
	
	// Unmarshal
	var d2 Duration
	err = d2.UnmarshalJSON(data)
	if err != nil {
		t.Fatalf("failed to unmarshal duration: %v", err)
	}
	
	if d2.Duration != d.Duration {
		t.Errorf("duration mismatch: got %v, want %v", d2.Duration, d.Duration)
	}
}

func createTestManifest() *PluginManifest {
	m := New("test-plugin", "io.test.plugin")
	m.Plugin.Description = "Test plugin"
	m.Plugin.Author.Name = "Test Author"
	m.Plugin.Author.Email = "test@example.com"
	
	// Add a test resource
	m.Resources = append(m.Resources, ResourceDef{
		Type:        "mrn:test:resource:action",
		Name:        "Test Resource",
		Description: "A test resource",
		Category:    "action",
		Parameters: []Parameter{
			{
				Name:        "param1",
				Type:        "string",
				Required:    true,
				Description: "A test parameter",
			},
		},
	})
	
	return m
}