package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	sdk "maschine.io/plugin-sdk/sdk"
)

// MailPlugin is our implementation of the MaschineResource interface
type MailPlugin struct {
	logger hclog.Logger
}

// GetMetadata returns plugin metadata
func (p *MailPlugin) GetMetadata(ctx context.Context, req *sdk.GetMetadataRequest) (*sdk.GetMetadataResponse, error) {
	return &sdk.GetMetadataResponse{
		Name:    "mail-plugin",
		Version: "1.0.0",
		SupportedResources: []string{
			"mrn:mail:smtp:send",
			"mrn:mail:imap:fetch",
		},
		Capabilities: map[string]string{
			"smtp": "Send emails via SMTP",
			"imap": "Fetch emails via IMAP",
		},
	}, nil
}

// Execute runs a plugin command
func (p *MailPlugin) Execute(ctx context.Context, req *sdk.ExecuteRequest) (*sdk.ExecuteResponse, error) {
	p.logger.Info("executing resource", "resource", req.Resource)
	
	switch req.Resource {
	case "mrn:mail:smtp:send":
		return p.sendMail(ctx, req)
	case "mrn:mail:imap:fetch":
		return p.fetchMail(ctx, req)
	default:
		return &sdk.ExecuteResponse{
			Error: fmt.Sprintf("unknown resource: %s", req.Resource),
		}, nil
	}
}

// HealthCheck checks if the plugin is healthy
func (p *MailPlugin) HealthCheck(ctx context.Context, req *sdk.HealthCheckRequest) (*sdk.HealthCheckResponse, error) {
	return &sdk.HealthCheckResponse{
		Healthy: true,
		Message: "Mail plugin is operational",
	}, nil
}

func (p *MailPlugin) sendMail(ctx context.Context, req *sdk.ExecuteRequest) (*sdk.ExecuteResponse, error) {
	// Parse parameters
	params := make(map[string]interface{})
	for key, data := range req.Parameters {
		var value interface{}
		if err := json.Unmarshal(data, &value); err != nil {
			params[key] = string(data)
		} else {
			params[key] = value
		}
	}
	
	// Extract email parameters
	to, _ := params["to"].(string)
	from, _ := params["from"].(string) 
	subject, _ := params["subject"].(string)
	body, _ := params["body"].(string)
	
	// Validate
	if to == "" || from == "" || subject == "" {
		return &sdk.ExecuteResponse{
			Error: "missing required parameters: to, from, subject",
		}, nil
	}
	
	// Get credentials
	server := req.Credentials["smtp_server"]
	user := req.Credentials["smtp_user"]
	password := req.Credentials["smtp_password"]
	
	p.logger.Info("sending email",
		"to", to,
		"from", from,
		"subject", subject,
		"body_length", len(body),
		"server", server,
		"user", user,
		"has_password", password != "",
	)
	
	// TODO: Actual SMTP implementation here
	
	// Return result
	result := map[string]interface{}{
		"status": "success",
		"message": fmt.Sprintf("Email sent to %s", to),
		"messageId": "12345",
	}
	
	output, _ := json.Marshal(result)
	
	return &sdk.ExecuteResponse{
		Output: output,
		Metadata: map[string]string{
			"operation": "smtp_send",
		},
	}, nil
}

func (p *MailPlugin) fetchMail(ctx context.Context, req *sdk.ExecuteRequest) (*sdk.ExecuteResponse, error) {
	// TODO: Implement IMAP fetch
	return &sdk.ExecuteResponse{
		Error: "not implemented",
	}, nil
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:       "mail-plugin",
		Output:     os.Stderr,
		Level:      hclog.Info,
		JSONFormat: true,
	})
	
	// Check for debug mode
	if os.Getenv("MASCHINE_PLUGIN_DEBUG") == "1" {
		logger.SetLevel(hclog.Debug)
	}
	
	logger.Info("starting mail plugin", "version", "1.0.0")
	
	// pluginMap is the map of plugins we can dispense
	var pluginMap = map[string]plugin.Plugin{
		"maschine": &sdk.MaschinePlugin{Impl: &MailPlugin{logger: logger}},
	}
	
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: sdk.Handshake,
		Plugins:         pluginMap,
		Logger:          logger,
		GRPCServer:      plugin.DefaultGRPCServer,
	})
}