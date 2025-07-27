package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"sort"
)

// PluginFunction defines a plugin function with metadata
type PluginFunction struct {
	Handler  func(context.Context, *TypedExecuteRequest) (any, error)
	Metadata map[string]string
}

// BasePlugin provides a base implementation for plugins using a function registry
type BasePlugin struct {
	name         string
	version      string
	functions    map[string]PluginFunction
	versionInfo  *VersionInfo // Optional detailed version info
}

// NewBasePlugin creates a new base plugin
func NewBasePlugin(name, version string) *BasePlugin {
	return &BasePlugin{
		name:      name,
		version:   version,
		functions: make(map[string]PluginFunction),
	}
}

// RegisterFunction adds a function to the plugin
func (p *BasePlugin) RegisterFunction(resource string, fn PluginFunction) {
	p.functions[resource] = fn
}

// RegisterSimpleFunction adds a function with default metadata
func (p *BasePlugin) RegisterSimpleFunction(resource string, handler func(context.Context, *TypedExecuteRequest) (any, error), operation string) {
	p.RegisterFunction(resource, PluginFunction{
		Handler: handler,
		Metadata: map[string]string{
			"operation": operation,
		},
	})
}

// GetMetadata returns plugin metadata
func (p *BasePlugin) GetMetadata(ctx context.Context) (*Metadata, error) {
	// Collect all supported resources
	resources := make([]string, 0, len(p.functions))
	for resource := range p.functions {
		resources = append(resources, resource)
	}
	sort.Strings(resources) // For consistent output
	
	return &Metadata{
		Name:               p.name,
		Version:            p.version,
		SupportedResources: resources,
	}, nil
}

// Execute runs a plugin function
func (p *BasePlugin) Execute(ctx context.Context, req *ExecuteRequest) (*ExecuteResponse, error) {
	// Create typed request
	typedReq, err := NewTypedExecuteRequest(req)
	if err != nil {
		return &ExecuteResponse{
			Error: fmt.Sprintf("failed to create typed request: %v", err),
		}, nil
	}
	
	// Look up function
	fn, exists := p.functions[req.Resource]
	if !exists {
		return &ExecuteResponse{
			Error: fmt.Sprintf("unsupported resource: %s", req.Resource),
		}, nil
	}
	
	// Execute function
	result, err := fn.Handler(ctx, typedReq)
	if err != nil {
		return &ExecuteResponse{
			Error: err.Error(),
		}, nil
	}
	
	// Marshal result
	output, err := json.Marshal(result)
	if err != nil {
		return &ExecuteResponse{
			Error: fmt.Sprintf("failed to marshal result: %v", err),
		}, nil
	}
	
	// Return response with metadata
	return &ExecuteResponse{
		Output:   output,
		Metadata: fn.Metadata,
	}, nil
}

// HealthCheck provides a default health check implementation
func (p *BasePlugin) HealthCheck(ctx context.Context) (*HealthStatus, error) {
	return &HealthStatus{
		Healthy: true,
		Message: fmt.Sprintf("%s plugin is operational", p.name),
	}, nil
}

// GetVersion returns detailed version information
func (p *BasePlugin) GetVersion(ctx context.Context) (*VersionInfo, error) {
	if p.versionInfo != nil {
		return p.versionInfo, nil
	}
	
	// Return default version info if not set
	return &VersionInfo{
		Version:   p.version,
		GoVersion: runtime.Version(),
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}, nil
}

// SetVersionInfo allows plugins to provide detailed version information
func (p *BasePlugin) SetVersionInfo(info *VersionInfo) {
	p.versionInfo = info
}

// SetHealthCheck allows plugins to override the default health check
func (p *BasePlugin) SetHealthCheck(check func(context.Context) (*HealthStatus, error)) {
	// This would need a field to store the custom health check
	// For now, plugins can embed BasePlugin and override HealthCheck
}