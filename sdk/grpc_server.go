package sdk

import (
	"context"
	
	pluginv1 "maschine.io/plugin-sdk/proto/plugin/v1"
)

// grpcServer is the server that GRPCServer will construct
type grpcServer struct {
	pluginv1.UnimplementedPluginServer
	// This is our real implementation
	Impl MaschineResource
}

func (s *grpcServer) GetMetadata(ctx context.Context, req *pluginv1.GetMetadataRequest) (*pluginv1.GetMetadataResponse, error) {
	resp, err := s.Impl.GetMetadata(ctx, &GetMetadataRequest{})
	if err != nil {
		return nil, err
	}
	
	return &pluginv1.GetMetadataResponse{
		Name:               resp.Name,
		Version:            resp.Version,
		SupportedResources: resp.SupportedResources,
		Capabilities:       resp.Capabilities,
	}, nil
}

func (s *grpcServer) Execute(ctx context.Context, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	resp, err := s.Impl.Execute(ctx, &ExecuteRequest{
		Resource:    req.Resource,
		Input:       req.Input,
		Parameters:  req.Parameters,
		Credentials: req.Credentials,
		Context:     req.Context,
	})
	if err != nil {
		// If Execute returns an error, wrap it in the response
		return &pluginv1.ExecuteResponse{
			Error: err.Error(),
		}, nil
	}
	
	return &pluginv1.ExecuteResponse{
		Output:   resp.Output,
		Error:    resp.Error,
		Metadata: resp.Metadata,
	}, nil
}

func (s *grpcServer) HealthCheck(ctx context.Context, req *pluginv1.HealthCheckRequest) (*pluginv1.HealthCheckResponse, error) {
	resp, err := s.Impl.HealthCheck(ctx, &HealthCheckRequest{})
	if err != nil {
		return nil, err
	}
	
	return &pluginv1.HealthCheckResponse{
		Healthy: resp.Healthy,
		Message: resp.Message,
	}, nil
}