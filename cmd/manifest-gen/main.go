package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"maschine.io/plugin-sdk/sdk/manifest"
)

func main() {
	var (
		pluginID    = flag.String("id", "", "Plugin ID in reverse domain notation (e.g., io.maschine.plugins.example)")
		name        = flag.String("name", "", "Plugin technical name (e.g., example-plugin)")
		displayName = flag.String("display-name", "", "Plugin display name (e.g., Example Plugin)")
		description = flag.String("description", "", "Plugin description")
		category    = flag.String("category", "general", "Plugin category (communication|cloud|monitoring|security|data|general)")
		version     = flag.String("version", "0.1.0", "Plugin version")
		author      = flag.String("author", "", "Author name")
		email       = flag.String("email", "", "Author email")
		license     = flag.String("license", "Apache-2.0", "License (SPDX identifier)")
		output      = flag.String("output", "plugin-manifest-v2.json", "Output file path")
		validate    = flag.Bool("validate", false, "Validate existing manifest file")
		update      = flag.Bool("update", false, "Update existing manifest file")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Maschine Plugin Manifest Generator\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  # Generate a new manifest\n")
		fmt.Fprintf(os.Stderr, "  %s -id io.maschine.plugins.example -name example-plugin -display-name \"Example Plugin\" -description \"An example plugin\" -author \"John Doe\" -email john@example.com\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  # Validate an existing manifest\n")
		fmt.Fprintf(os.Stderr, "  %s -validate -output ./plugin-manifest-v2.json\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  # Update an existing manifest\n")
		fmt.Fprintf(os.Stderr, "  %s -update -output ./plugin-manifest-v2.json -version 1.0.0\n\n", os.Args[0])
	}

	flag.Parse()

	// Validate mode
	if *validate {
		if err := validateManifest(*output); err != nil {
			log.Fatalf("Validation failed: %v", err)
		}
		fmt.Printf("✓ Manifest '%s' is valid\n", *output)
		return
	}

	// Update mode
	if *update {
		if err := updateManifest(*output, *version); err != nil {
			log.Fatalf("Update failed: %v", err)
		}
		fmt.Printf("✓ Updated manifest '%s'\n", *output)
		return
	}

	// Generate mode - validate required fields
	if *pluginID == "" || *name == "" {
		fmt.Fprintf(os.Stderr, "Error: -id and -name are required when generating a new manifest\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// Generate display name if not provided
	if *displayName == "" {
		*displayName = generateDisplayName(*name)
	}

	// Generate description if not provided
	if *description == "" {
		*description = fmt.Sprintf("%s plugin for Maschine", *displayName)
	}

	// Validate plugin ID format
	if !isValidPluginID(*pluginID) {
		log.Fatalf("Invalid plugin ID format. Must be in reverse domain notation (e.g., io.maschine.plugins.example)")
	}

	// Validate plugin name format
	if !isValidPluginName(*name) {
		log.Fatalf("Invalid plugin name format. Must contain only lowercase letters, numbers, and hyphens (e.g., example-plugin)")
	}

	// Create manifest
	m := manifest.New(*name, *pluginID)
	m.Plugin.DisplayName = *displayName
	m.Plugin.Description = *description
	m.Plugin.Category = *category
	m.Plugin.Version = *version
	
	if *author != "" {
		m.Plugin.Author.Name = *author
	}
	if *email != "" {
		m.Plugin.Author.Email = *email
	}
	if *license != "" {
		m.Plugin.License = *license
	}

	// Add example resource
	m.Resources = []manifest.ResourceDef{
		{
			Type:        fmt.Sprintf("mrn:%s:example:action", getPluginShortName(*pluginID)),
			Name:        "Example Action",
			Description: "An example action resource",
			Category:    "action",
			Parameters: []manifest.Parameter{
				{
					Name:        "input",
					Type:        "string",
					Required:    true,
					Description: "Input parameter",
				},
			},
			Output: &manifest.OutputDef{
				Type: "object",
				Schema: map[string]interface{}{
					"result": map[string]string{"type": "string"},
					"status": map[string]string{"type": "string"},
				},
			},
			Examples: []manifest.Example{
				{
					Name: "Basic example",
					Parameters: map[string]interface{}{
						"input": "Hello, World!",
					},
					Output: map[string]interface{}{
						"result": "Processed: Hello, World!",
						"status": "success",
					},
				},
			},
		},
	}

	// Save manifest
	if err := m.Save(*output); err != nil {
		log.Fatalf("Failed to save manifest: %v", err)
	}

	fmt.Printf("✓ Generated manifest '%s'\n", *output)
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("1. Edit '%s' to add your plugin's resources and configuration\n", *output)
	fmt.Printf("2. Run '%s -validate -output %s' to validate your changes\n", os.Args[0], *output)
	fmt.Printf("3. Use the manifest in your plugin's distribution\n")
}

func validateManifest(path string) error {
	m, err := manifest.Load(path)
	if err != nil {
		return fmt.Errorf("failed to load manifest: %w", err)
	}

	// Additional validation beyond basic schema
	fmt.Printf("Manifest details:\n")
	fmt.Printf("  Plugin ID: %s\n", m.Plugin.ID)
	fmt.Printf("  Name: %s\n", m.Plugin.Name)
	fmt.Printf("  Version: %s\n", m.Plugin.Version)
	fmt.Printf("  Category: %s\n", m.Plugin.Category)
	fmt.Printf("  Resources: %d\n", len(m.Resources))
	
	for i, res := range m.Resources {
		fmt.Printf("    %d. %s (%s)\n", i+1, res.Type, res.Category)
	}

	return nil
}

func updateManifest(path string, newVersion string) error {
	m, err := manifest.Load(path)
	if err != nil {
		return fmt.Errorf("failed to load manifest: %w", err)
	}

	// Update version if provided
	if newVersion != "" {
		m.Plugin.Version = newVersion
	}

	// Save updated manifest
	return m.Save(path)
}

func isValidPluginID(id string) bool {
	parts := strings.Split(id, ".")
	if len(parts) < 2 {
		return false
	}
	for _, part := range parts {
		if part == "" || !isLowerAlphaNum(part) {
			return false
		}
	}
	return true
}

func isValidPluginName(name string) bool {
	if name == "" {
		return false
	}
	for i, ch := range name {
		if i == 0 && !isLowerAlpha(ch) {
			return false
		}
		if !isLowerAlpha(ch) && ch != '-' && !isDigit(ch) {
			return false
		}
	}
	return true
}

func isLowerAlpha(ch rune) bool {
	return ch >= 'a' && ch <= 'z'
}

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func isLowerAlphaNum(s string) bool {
	for _, ch := range s {
		if !isLowerAlpha(ch) && !isDigit(ch) {
			return false
		}
	}
	return true
}

func generateDisplayName(name string) string {
	// Convert kebab-case to Title Case
	parts := strings.Split(name, "-")
	for i, part := range parts {
		if part != "" {
			parts[i] = strings.Title(part)
		}
	}
	return strings.Join(parts, " ")
}

func getPluginShortName(pluginID string) string {
	// Extract short name from plugin ID (e.g., "io.maschine.plugins.mail" -> "mail")
	parts := strings.Split(pluginID, ".")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return "plugin"
}