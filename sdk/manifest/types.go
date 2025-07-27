package manifest

import (
	"time"
)

// PluginManifest represents the complete plugin manifest structure
type PluginManifest struct {
	Schema          string           `json:"$schema,omitempty"`
	ManifestVersion string           `json:"manifestVersion"`
	Plugin          PluginInfo       `json:"plugin"`
	Runtime         RuntimeInfo      `json:"runtime"`
	Requirements    Requirements     `json:"requirements"`
	Configuration   Configuration    `json:"configuration"`
	Resources       []ResourceDef    `json:"resources"`
	Capabilities    Capabilities     `json:"capabilities"`
	Limits          Limits           `json:"limits"`
	Documentation   Documentation    `json:"documentation"`
}

// PluginInfo contains basic plugin information
type PluginInfo struct {
	ID          string     `json:"id"`          // Unique plugin identifier (e.g., "io.maschine.plugins.mail")
	Name        string     `json:"name"`        // Technical name
	DisplayName string     `json:"displayName"` // Human-readable name
	Version     string     `json:"version"`
	Description string     `json:"description"`
	Category    string     `json:"category"` // e.g., "communication", "cloud", "monitoring"
	Tags        []string   `json:"tags"`
	Author      AuthorInfo `json:"author"`
	License     string     `json:"license"`
	Homepage    string     `json:"homepage"`
	Repository  Repository `json:"repository"`
	Bugs        BugTracker `json:"bugs"`
}

// AuthorInfo contains plugin author information
type AuthorInfo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	URL   string `json:"url,omitempty"`
}

// Repository contains source repository information
type Repository struct {
	Type string `json:"type"` // "git", "svn", etc.
	URL  string `json:"url"`
}

// BugTracker contains issue tracker information
type BugTracker struct {
	URL string `json:"url"`
}

// RuntimeInfo describes how to run the plugin
type RuntimeInfo struct {
	Type            string          `json:"type"` // "native", "container", "wasm"
	Executable      ExecutableInfo  `json:"executable"`
	Protocol        string          `json:"protocol"`        // "grpc", "http", "stdio"
	ProtocolVersion string          `json:"protocolVersion"` // "1.0"
	HandshakeConfig HandshakeConfig `json:"handshakeConfig,omitempty"`
}

// ExecutableInfo contains executable details
type ExecutableInfo struct {
	Path      string            `json:"path"` // Relative path to executable
	Checksums map[string]string `json:"checksums,omitempty"` // platform -> sha256
}

// HandshakeConfig for plugin negotiation
type HandshakeConfig struct {
	ProtocolVersion  int    `json:"protocolVersion"`
	MagicCookieKey   string `json:"magicCookieKey"`
	MagicCookieValue string `json:"magicCookieValue"`
}

// Requirements specifies plugin requirements
type Requirements struct {
	MaschineVersion string   `json:"maschineVersion"` // Semantic version constraint (e.g., ">=1.0.0")
	OS              []string `json:"os"`              // ["darwin", "linux", "windows"]
	Arch            []string `json:"arch"`            // ["amd64", "arm64"]
}

// Configuration describes plugin configuration options
type Configuration struct {
	Environment []EnvVar         `json:"environment"`
	Credentials []CredentialSet  `json:"credentials"`
}

// EnvVar describes an environment variable
type EnvVar struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Required    bool     `json:"required"`
	Default     string   `json:"default,omitempty"`
	Enum        []string `json:"enum,omitempty"` // Allowed values
}

// CredentialSet describes a set of credentials
type CredentialSet struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Fields      []CredentialField `json:"fields"`
}

// CredentialField describes a single credential field
type CredentialField struct {
	Name        string `json:"name"`
	Type        string `json:"type"` // "string", "number", "boolean"
	Required    bool   `json:"required"`
	Secret      bool   `json:"secret,omitempty"` // Should be masked in UI/logs
	Pattern     string `json:"pattern,omitempty"` // Regex pattern for validation
	Description string `json:"description,omitempty"`
}

// ResourceDef defines a plugin resource/function
type ResourceDef struct {
	Type                string       `json:"type"`        // MRN type (e.g., "mrn:mail:smtp:send")
	Name                string       `json:"name"`        // Human-readable name
	Description         string       `json:"description"`
	Category            string       `json:"category"` // "action", "query", "check"
	Parameters          []Parameter  `json:"parameters"`
	RequiredCredentials []string     `json:"requiredCredentials,omitempty"`
	Output              *OutputDef   `json:"output,omitempty"`
	Examples            []Example    `json:"examples,omitempty"`
}

// Parameter describes a resource parameter
type Parameter struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"` // "string", "integer", "boolean", "array", "object"
	Required    bool        `json:"required"`
	Description string      `json:"description"`
	Default     interface{} `json:"default,omitempty"`
	Pattern     string      `json:"pattern,omitempty"`  // Regex for strings
	MinLength   int         `json:"minLength,omitempty"`
	MaxLength   int         `json:"maxLength,omitempty"`
	Minimum     float64     `json:"minimum,omitempty"`  // For numbers
	Maximum     float64     `json:"maximum,omitempty"`  // For numbers
	Enum        []string    `json:"enum,omitempty"`     // Allowed values
	Items       *Items      `json:"items,omitempty"`    // For arrays
	Properties  interface{} `json:"properties,omitempty"` // For objects
}

// Items describes array items
type Items struct {
	Type       string      `json:"type"`
	Properties interface{} `json:"properties,omitempty"`
}

// OutputDef describes resource output
type OutputDef struct {
	Type   string      `json:"type"` // "object", "array", "string", etc.
	Schema interface{} `json:"schema"`
}

// Example shows how to use a resource
type Example struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Parameters  map[string]interface{} `json:"parameters"`
	Output      interface{}            `json:"output,omitempty"`
}

// Capabilities describes plugin capabilities
type Capabilities struct {
	HealthCheck          bool `json:"health_check"`
	ConcurrentExecution  bool `json:"concurrent_execution"`
	Stateless           bool `json:"stateless"`
	SupportsCredentials bool `json:"supports_credentials,omitempty"`
}

// Limits defines resource limits
type Limits struct {
	StartupTimeout           Duration `json:"startup_timeout"`
	ExecuteTimeout           Duration `json:"execute_timeout"`
	HealthCheckInterval      Duration `json:"health_check_interval"`
	HealthCheckTimeout       Duration `json:"health_check_timeout"`
	MaxConcurrentExecutions  int      `json:"max_concurrent_executions"`
}

// Duration is a wrapper for time.Duration with JSON marshaling
type Duration struct {
	time.Duration
}

// MarshalJSON converts Duration to JSON string
func (d Duration) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.String() + `"`), nil
}

// UnmarshalJSON converts JSON string to Duration
func (d *Duration) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = s[1 : len(s)-1] // Remove quotes
	dur, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	d.Duration = dur
	return nil
}

// Documentation contains links to documentation
type Documentation struct {
	Installation  string `json:"installation"`
	Configuration string `json:"configuration"`
	Usage         string `json:"usage"`
	API           string `json:"api"`
}