# Maschine Plugin Manifest Updater

A tool to update plugin manifests with build artifacts from GoReleaser.

## Installation

```bash
go install maschine.io/plugin-sdk/cmd/manifest-updater@latest
```

## Usage

The manifest updater reads a plugin manifest file and updates it with:
- Version information
- Binary checksums from GoReleaser builds
- Executable paths

### Basic Usage

```bash
manifest-updater \
  -manifest plugin-manifest-v2.json \
  -checksums dist/checksums.txt \
  -version v1.0.0 \
  -project mail-plugin
```

### Options

| Option | Description | Default |
|--------|-------------|---------|
| `-manifest` | Path to the plugin manifest file | `plugin-manifest-v2.json` |
| `-checksums` | Path to GoReleaser checksums.txt file | `checksums.txt` |
| `-dist` | Distribution directory (alternative to checksums file) | `dist` |
| `-version` | Plugin version (e.g., v1.0.0) | Required |
| `-project` | Project name (binary name) | Required |
| `-output` | Output path for updated manifest | Updates in place |

## Integration with GoReleaser

### 1. Create a Script

Create `scripts/update-manifest.sh`:

```bash
#!/bin/bash
set -e

VERSION="${GORELEASER_CURRENT_TAG:-dev}"
PROJECT_NAME="your-plugin-name"
DIST_DIR="${DIST_DIR:-dist}"

manifest-updater \
    -manifest plugin-manifest-v2.json \
    -checksums "$DIST_DIR/checksums.txt" \
    -version "$VERSION" \
    -project "$PROJECT_NAME" \
    -output "$DIST_DIR/plugin-manifest-v2.json"
```

### 2. Configure GoReleaser

In `.goreleaser.yaml`:

```yaml
# Run after all artifacts are built
after:
  hooks:
    - ./scripts/update-manifest.sh

release:
  # Include the manifest in the release
  extra_files:
    - glob: dist/plugin-manifest-v2.json
```

### 3. Prepare Manifest Template

Your `plugin-manifest-v2.json` should have placeholders for dynamic values:

```json
{
  "$schema": "https://maschine.io/schemas/plugin-manifest-v2.json",
  "manifestVersion": "2.0",
  "plugin": {
    "version": "0.0.0",  // Will be updated
    ...
  },
  "runtime": {
    "executable": {
      "path": "./plugin",  // Will be updated
      "checksums": {}      // Will be populated
    },
    ...
  }
}
```

## How It Works

1. **Reads** the existing plugin manifest
2. **Updates** the version from the provided `-version` flag
3. **Parses** the GoReleaser checksums.txt file
4. **Maps** artifacts to platform identifiers (e.g., `darwin-amd64`)
5. **Updates** the manifest with SHA256 checksums
6. **Saves** the updated manifest

### Platform Mapping

The tool maps GoReleaser artifact names to platform identifiers:

| GoReleaser Format | Manifest Format |
|-------------------|-----------------|
| `plugin_1.0.0_linux_amd64.tar.gz` | `linux-amd64` |
| `plugin_1.0.0_darwin_arm64.tar.gz` | `darwin-arm64` |
| `plugin_1.0.0_windows_x86_64.zip` | `windows-amd64` |

## Example Output

After running the updater:

```json
{
  "plugin": {
    "version": "1.0.0"
  },
  "runtime": {
    "executable": {
      "path": "./mail-plugin",
      "checksums": {
        "darwin-amd64": "sha256:abc123...",
        "darwin-arm64": "sha256:def456...",
        "linux-amd64": "sha256:ghi789...",
        "windows-amd64": "sha256:jkl012..."
      }
    }
  }
}
```

## Alternative: Direct Checksum Calculation

If you don't have a checksums.txt file, the tool can calculate checksums directly:

```bash
manifest-updater \
  -manifest plugin-manifest-v2.json \
  -dist dist \
  -version v1.0.0 \
  -project mail-plugin
```

This will look for files matching `dist/mail-plugin_*` and calculate their checksums.