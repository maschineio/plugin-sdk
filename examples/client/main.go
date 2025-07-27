package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	sdk "maschine.io/plugin-sdk/sdk"
)

func main() {
	// Create logger
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin-client",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})
	
	// We're a host. Start by launching the plugin process.
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: sdk.Handshake,
		Plugins:         sdk.PluginMap,
		Cmd:             exec.Command("../mail/mail-plugin"),
		Logger:          logger,
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolGRPC,
		},
	})
	defer client.Kill()
	
	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		log.Fatal(err)
	}
	
	// Request the plugin
	raw, err := rpcClient.Dispense("maschine")
	if err != nil {
		log.Fatal(err)
	}
	
	// We should have a MaschineResource now
	maschinePlugin := raw.(sdk.MaschineResource)
	
	// Test GetMetadata
	fmt.Println("=== Testing GetMetadata ===")
	metadata, err := maschinePlugin.GetMetadata(context.Background(), &sdk.GetMetadataRequest{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Plugin: %s v%s\n", metadata.Name, metadata.Version)
	fmt.Printf("Resources: %v\n", metadata.SupportedResources)
	fmt.Printf("Capabilities: %v\n", metadata.Capabilities)
	
	// Test HealthCheck
	fmt.Println("\n=== Testing HealthCheck ===")
	health, err := maschinePlugin.HealthCheck(context.Background(), &sdk.HealthCheckRequest{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Healthy: %v\n", health.Healthy)
	fmt.Printf("Message: %s\n", health.Message)
	
	// Test Execute
	fmt.Println("\n=== Testing Execute ===")
	
	// Prepare parameters
	params := map[string]interface{}{
		"to":      "user@example.com",
		"from":    "sender@example.com",
		"subject": "Test Email",
		"body":    "This is a test email",
	}
	
	paramBytes := make(map[string][]byte)
	for k, v := range params {
		data, _ := json.Marshal(v)
		paramBytes[k] = data
	}
	
	execReq := &sdk.ExecuteRequest{
		Resource:   "mrn:mail:smtp:send",
		Parameters: paramBytes,
		Credentials: map[string]string{
			"smtp_server":   "smtp.example.com:587",
			"smtp_user":     "testuser",
			"smtp_password": "testpass",
		},
	}
	
	resp, err := maschinePlugin.Execute(context.Background(), execReq)
	if err != nil {
		log.Fatal(err)
	}
	
	if resp.Error != "" {
		fmt.Printf("Error: %s\n", resp.Error)
	} else {
		var result map[string]interface{}
		json.Unmarshal(resp.Output, &result)
		fmt.Printf("Result: %+v\n", result)
		fmt.Printf("Metadata: %+v\n", resp.Metadata)
	}
}