package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/go-plugin"
	"maschine.io/plugin-sdk/sdk"
	"maschine.io/plugin-sdk/sdk/manifest"
)

// Request structs with manifest tags
type SendEmailRequest struct {
	From    string   `json:"from" manifest:"description=Sender email address,pattern=^[^@]+@[^@]+\\.[^@]+$"`
	To      string   `json:"to" manifest:"description=Recipient email address,pattern=^[^@]+@[^@]+\\.[^@]+$"`
	Subject string   `json:"subject" manifest:"description=Email subject line,maxLength=255"`
	Body    string   `json:"body" manifest:"description=Email body content"`
	CC      []string `json:"cc,omitempty" manifest:"description=CC recipients"`
}

type SendEmailV2Request struct {
	From        string   `json:"from" manifest:"description=Sender email address"`
	To          string   `json:"to" manifest:"description=Recipient email address"`
	CC          []string `json:"cc,omitempty" manifest:"description=CC recipients"`
	BCC         []string `json:"bcc,omitempty" manifest:"description=BCC recipients"`
	Subject     string   `json:"subject" manifest:"description=Email subject"`
	Body        string   `json:"body" manifest:"description=Email body (plain text)"`
	HTML        string   `json:"html,omitempty" manifest:"description=HTML body content"`
	Timeout     int      `json:"timeout,omitempty" manifest:"description=Timeout in seconds,default=10"`
}

func main() {
	// Create plugin with manifest support
	pluginInfo := manifest.PluginInfo{
		ID:          "io.maschine.plugins.example",
		Name:        "example-plugin",
		DisplayName: "Example Plugin",
		Version:     "1.0.0",
		Description: "Example plugin demonstrating manifest generation",
		Category:    "general",
		Tags:        []string{"example", "demo"},
		Author: manifest.AuthorInfo{
			Name:  "Maschine.io Team",
			Email: "plugins@maschine.io",
		},
		License: "Apache-2.0",
	}

	p := sdk.NewManifestPlugin(pluginInfo)

	// Set configuration requirements
	p.SetConfiguration(manifest.Configuration{
		Environment: []manifest.EnvVar{
			{
				Name:        "PLUGIN_LOG_LEVEL",
				Description: "Log level for the plugin",
				Required:    false,
				Default:     "info",
				Enum:        []string{"debug", "info", "warn", "error"},
			},
		},
		Credentials: []manifest.CredentialSet{
			{
				Name:        "smtp",
				Description: "SMTP server credentials",
				Fields: []manifest.CredentialField{
					{
						Name:     "user",
						Type:     "string",
						Required: true,
					},
					{
						Name:     "password",
						Type:     "string",
						Required: true,
						Secret:   true,
					},
					{
						Name:     "server",
						Type:     "string",
						Required: true,
						Pattern:  "^[^:]+:[0-9]+$",
					},
				},
			},
		},
	})

	// Register resources with automatic parameter extraction
	err := p.RegisterResourceWithRequest(
		"mrn:example:email:send",
		"Send Email",
		"Send an email via SMTP",
		"action",
		SendEmailRequest{},
		sendEmail,
	)
	if err != nil {
		log.Fatal(err)
	}

	err = p.RegisterResourceWithRequest(
		"mrn:example:email:sendv2",
		"Send Email V2",
		"Send an email with advanced features",
		"action",
		SendEmailV2Request{},
		sendEmailV2,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Generate manifest (for development/documentation)
	if len(os.Args) > 1 && os.Args[1] == "generate-manifest" {
		m, err := p.GenerateManifest()
		if err != nil {
			log.Fatal(err)
		}
		
		if err := m.Save("plugin-manifest-v2.json"); err != nil {
			log.Fatal(err)
		}
		
		fmt.Println("✓ Generated plugin-manifest-v2.json")
		return
	}

	// Serve plugin
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: sdk.Handshake,
		Plugins: map[string]plugin.Plugin{
			"resource": &sdk.ResourcePlugin{
				Impl: p,
			},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}

func sendEmail(ctx context.Context, req *sdk.TypedExecuteRequest) (any, error) {
	var params SendEmailRequest
	if err := req.GetParameters(&params); err != nil {
		return nil, err
	}

	// Implementation here...
	return map[string]interface{}{
		"success":   true,
		"messageId": "12345",
	}, nil
}

func sendEmailV2(ctx context.Context, req *sdk.TypedExecuteRequest) (any, error) {
	var params SendEmailV2Request
	if err := req.GetParameters(&params); err != nil {
		return nil, err
	}

	// Implementation here...
	return map[string]interface{}{
		"success":   true,
		"messageId": "67890",
	}, nil
}