package sdk

import (
	"context"
)

// grpcServer implements the gRPC server for the plugin
type grpcServer struct {
	UnimplementedPluginServer
	Impl Resource
}

func (s *grpcServer) GetMetadata(ctx context.Context, req *GetMetadataRequest) (*GetMetadataResponse, error) {
	metadata, err := s.Impl.GetMetadata(ctx)
	if err != nil {
		return nil, err
	}
	
	return &GetMetadataResponse{
		Name:               metadata.Name,
		Version:            metadata.Version,
		SupportedResources: metadata.SupportedResources,
		Capabilities:       metadata.Capabilities,
	}, nil
}

func (s *grpcServer) Execute(ctx context.Context, req *ExecuteRequest) (*ExecuteResponse, error) {
	executeReq := &ExecuteRequest{
		Resource:    req.Resource,
		Input:       req.Input,
		Parameters:  req.Parameters,
		Credentials: req.Credentials,
		Context:     req.Context,
	}
	
	resp, err := s.Impl.Execute(ctx, executeReq)
	if err != nil {
		return nil, err
	}
	
	return &ExecuteResponse{
		Output:   resp.Output,
		Error:    resp.Error,
		Metadata: resp.Metadata,
	}, nil
}

func (s *grpcServer) HealthCheck(ctx context.Context, req *HealthCheckRequest) (*HealthCheckResponse, error) {
	status, err := s.Impl.HealthCheck(ctx)
	if err != nil {
		return nil, err
	}
	
	return &HealthCheckResponse{
		Healthy: status.Healthy,
		Message: status.Message,
	}, nil
}