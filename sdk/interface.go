package sdk

import (
	"context"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

// Handshake config for plugin negotiation
var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "MASCHINE_PLUGIN",
	MagicCookieValue: "maschine-io-plugin-sdk-v1",
}

// PluginMap is the map of plugins we can serve
var PluginMap = map[string]plugin.Plugin{
	"resource": &ResourcePlugin{},
}

// Resource is the interface that plugins must implement
type Resource interface {
	// GetMetadata returns plugin metadata
	GetMetadata(ctx context.Context) (*Metadata, error)
	
	// GetVersion returns detailed version information
	GetVersion(ctx context.Context) (*VersionInfo, error)
	
	// Execute runs the plugin function
	Execute(ctx context.Context, req *ExecuteRequest) (*ExecuteResponse, error)
	
	// HealthCheck checks plugin health
	HealthCheck(ctx context.Context) (*HealthStatus, error)
}

// Metadata contains plugin information (converted from protobuf)
type Metadata struct {
	Name               string
	Version            string
	SupportedResources []string
	Capabilities       map[string]string
}

// HealthStatus contains health check results (converted from protobuf)
type HealthStatus struct {
	Healthy bool
	Message string
}

// VersionInfo contains detailed version information
type VersionInfo struct {
	Version      string            `json:"version"`
	GitCommit    string            `json:"git_commit"`
	BuildDate    string            `json:"build_date"`
	GoVersion    string            `json:"go_version"`
	Platform     string            `json:"platform"`
	Dependencies map[string]string `json:"dependencies,omitempty"`
}

// ResourcePlugin is the plugin implementation
type ResourcePlugin struct {
	plugin.Plugin
	Impl Resource
}

// GRPCServer is required by the plugin.Plugin interface
func (p *ResourcePlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	RegisterPluginServer(s, &grpcServer{Impl: p.Impl})
	return nil
}

// GRPCClient is required by the plugin.Plugin interface
func (p *ResourcePlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &grpcClient{client: NewPluginClient(c)}, nil
}