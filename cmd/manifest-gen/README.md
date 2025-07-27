# Maschine Plugin Manifest Generator

A command-line tool for generating, validating, and updating Maschine plugin manifest files.

## Installation

```bash
go install maschine.io/plugin-sdk/cmd/manifest-gen@latest
```

## Usage

### Generate a New Manifest

Create a new plugin manifest with basic configuration:

```bash
manifest-gen \
  -id io.maschine.plugins.example \
  -name example-plugin \
  -display-name "Example Plugin" \
  -description "An example plugin for Maschine" \
  -author "John Doe" \
  -email john@example.com
```

This generates a `plugin-manifest-v2.json` file with:
- Basic plugin metadata
- Default runtime configuration for native gRPC plugins
- Example resource definition
- Standard capabilities and limits

### Validate an Existing Manifest

Check if a manifest file is valid:

```bash
manifest-gen -validate -output ./plugin-manifest-v2.json
```

### Update Manifest Version

Update the version in an existing manifest:

```bash
manifest-gen -update -output ./plugin-manifest-v2.json -version 1.0.0
```

## Command-Line Options

| Option | Description | Default |
|--------|-------------|---------|
| `-id` | Plugin ID in reverse domain notation (e.g., io.maschine.plugins.example) | Required |
| `-name` | Plugin technical name (e.g., example-plugin) | Required |
| `-display-name` | Plugin display name (e.g., "Example Plugin") | Generated from name |
| `-description` | Plugin description | Generated |
| `-category` | Plugin category (communication, cloud, monitoring, security, data, general) | general |
| `-version` | Plugin version | 0.1.0 |
| `-author` | Author name | Unknown |
| `-email` | Author email | unknown@example.com |
| `-license` | License (SPDX identifier) | Apache-2.0 |
| `-output` | Output file path | plugin-manifest-v2.json |
| `-validate` | Validate existing manifest file | false |
| `-update` | Update existing manifest file | false |

## Examples

### Generate Manifest for Mail Plugin

```bash
manifest-gen \
  -id io.maschine.plugins.mail \
  -name mail-plugin \
  -display-name "Mail Plugin" \
  -description "Send and receive emails via SMTP/IMAP protocols" \
  -category communication \
  -author "Maschine.io Team" \
  -email plugins@maschine.io
```

### Generate Manifest for Azure Plugin

```bash
manifest-gen \
  -id io.maschine.plugins.azure \
  -name azure-plugin \
  -display-name "Azure Plugin" \
  -description "Azure cloud platform integration for Maschine" \
  -category cloud \
  -author "Maschine.io Team" \
  -email plugins@maschine.io
```

## Generated Manifest Structure

The tool generates a manifest with:

1. **Plugin Information**: ID, name, version, description, author details
2. **Runtime Configuration**: Native gRPC plugin settings
3. **Requirements**: OS and architecture support
4. **Example Resource**: A sample resource definition to get started
5. **Standard Capabilities**: Health check, concurrent execution, stateless operation
6. **Default Limits**: Timeouts and concurrency limits

## Next Steps

After generating a manifest:

1. Edit the manifest to add your plugin's specific resources
2. Update credential requirements if needed
3. Modify capabilities and limits based on your plugin's behavior
4. Add proper documentation URLs
5. Validate the manifest using `-validate` flag

## Manifest Schema

The generated manifest follows the Maschine Plugin Manifest v2.0 schema.
See the schema documentation for detailed field descriptions.