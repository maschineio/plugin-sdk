package sdk

import (
	"context"
	
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	pluginv1 "maschine.io/plugin-sdk/proto/plugin/v1"
)

// Handshake is a common handshake that is shared by plugin and host
var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "MASCHINE_PLUGIN",
	MagicCookieValue: "93f6bc9f-f0bc-4b65-a0e5-9c8c7e5d3f4b",
}

// PluginMap is the map of plugins we can dispense
var PluginMap = map[string]plugin.Plugin{
	"maschine": &MaschinePlugin{},
}

// MaschineResource is the interface that we're exposing as a plugin
type MaschineResource interface {
	GetMetadata(context.Context, *GetMetadataRequest) (*GetMetadataResponse, error)
	Execute(context.Context, *ExecuteRequest) (*ExecuteResponse, error)
	HealthCheck(context.Context, *HealthCheckRequest) (*HealthCheckResponse, error)
}

// GetMetadataRequest is the request for metadata
type GetMetadataRequest struct{}

// GetMetadataResponse contains plugin metadata
type GetMetadataResponse struct {
	Name               string
	Version            string
	SupportedResources []string
	Capabilities       map[string]string
}

// ExecuteRequest contains execution parameters
type ExecuteRequest struct {
	Resource    string
	Input       []byte
	Parameters  map[string][]byte
	Credentials map[string]string
	Context     map[string]string
}

// ExecuteResponse contains execution results
type ExecuteResponse struct {
	Output   []byte
	Error    string
	Metadata map[string]string
}

// HealthCheckRequest is the request for health check
type HealthCheckRequest struct{}

// HealthCheckResponse contains health status
type HealthCheckResponse struct {
	Healthy bool
	Message string
}

// MaschinePlugin is the implementation of plugin.Plugin so we can serve/consume this
type MaschinePlugin struct {
	plugin.Plugin
	// Impl Injection
	Impl MaschineResource
}

func (p *MaschinePlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	pluginv1.RegisterPluginServer(s, &grpcServer{Impl: p.Impl})
	return nil
}

func (p *MaschinePlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &grpcClient{client: pluginv1.NewPluginClient(c)}, nil
}