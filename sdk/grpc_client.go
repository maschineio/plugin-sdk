package sdk

import (
	"context"
)

// grpcClient implements the client side of the gRPC plugin
type grpcClient struct {
	client PluginClient
}

func (c *grpcClient) GetMetadata(ctx context.Context) (*Metadata, error) {
	resp, err := c.client.GetMetadata(ctx, &GetMetadataRequest{})
	if err != nil {
		return nil, err
	}
	
	return &Metadata{
		Name:               resp.Name,
		Version:            resp.Version,
		SupportedResources: resp.SupportedResources,
		Capabilities:       resp.Capabilities,
	}, nil
}

func (c *grpcClient) Execute(ctx context.Context, req *ExecuteRequest) (*ExecuteResponse, error) {
	grpcReq := &ExecuteRequest{
		Resource:    req.Resource,
		Input:       req.Input,
		Parameters:  req.Parameters,
		Credentials: req.Credentials,
		Context:     req.Context,
	}
	
	resp, err := c.client.Execute(ctx, grpcReq)
	if err != nil {
		return nil, err
	}
	
	return &ExecuteResponse{
		Output:   resp.Output,
		Error:    resp.Error,
		Metadata: resp.Metadata,
	}, nil
}

func (c *grpcClient) HealthCheck(ctx context.Context) (*HealthStatus, error) {
	resp, err := c.client.HealthCheck(ctx, &HealthCheckRequest{})
	if err != nil {
		return nil, err
	}
	
	return &HealthStatus{
		Healthy: resp.Healthy,
		Message: resp.Message,
	}, nil
}