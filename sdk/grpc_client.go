package sdk

import (
	"context"
	
	pluginv1 "maschine.io/plugin-sdk/proto/plugin/v1"
)

// grpcClient is an implementation of MaschineResource that talks over RPC
type grpcClient struct {
	client pluginv1.PluginClient
}

func (c *grpcClient) GetMetadata(ctx context.Context, req *GetMetadataRequest) (*GetMetadataResponse, error) {
	resp, err := c.client.GetMetadata(ctx, &pluginv1.GetMetadataRequest{})
	if err != nil {
		return nil, err
	}
	
	return &GetMetadataResponse{
		Name:               resp.Name,
		Version:            resp.Version,
		SupportedResources: resp.SupportedResources,
		Capabilities:       resp.Capabilities,
	}, nil
}

func (c *grpcClient) Execute(ctx context.Context, req *ExecuteRequest) (*ExecuteResponse, error) {
	resp, err := c.client.Execute(ctx, &pluginv1.ExecuteRequest{
		Resource:    req.Resource,
		Input:       req.Input,
		Parameters:  req.Parameters,
		Credentials: req.Credentials,
		Context:     req.Context,
	})
	if err != nil {
		return nil, err
	}
	
	return &ExecuteResponse{
		Output:   resp.Output,
		Error:    resp.Error,
		Metadata: resp.Metadata,
	}, nil
}

func (c *grpcClient) HealthCheck(ctx context.Context, req *HealthCheckRequest) (*HealthCheckResponse, error) {
	resp, err := c.client.HealthCheck(ctx, &pluginv1.HealthCheckRequest{})
	if err != nil {
		return nil, err
	}
	
	return &HealthCheckResponse{
		Healthy: resp.Healthy,
		Message: resp.Message,
	}, nil
}