package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"maschine.io/plugin-sdk/sdk/manifest"
)

func main() {
	var (
		manifestPath  = flag.String("manifest", "plugin-manifest-v2.json", "Path to manifest file")
		checksumsPath = flag.String("checksums", "checksums.txt", "Path to goreleaser checksums file")
		distDir       = flag.String("dist", "dist", "Distribution directory")
		version       = flag.String("version", "", "Plugin version")
		projectName   = flag.String("project", "", "Project name")
		output        = flag.String("output", "", "Output manifest path (default: update in place)")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Maschine Plugin Manifest Updater\n\n")
		fmt.Fprintf(os.Stderr, "Updates a plugin manifest with build artifacts from goreleaser\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample:\n")
		fmt.Fprintf(os.Stderr, "  %s -manifest plugin-manifest-v2.json -checksums dist/checksums.txt -version v1.0.0 -project mail-plugin\n\n", os.Args[0])
	}

	flag.Parse()

	if *version == "" || *projectName == "" {
		fmt.Fprintf(os.Stderr, "Error: -version and -project are required\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// Load manifest
	m, err := manifest.Load(*manifestPath)
	if err != nil {
		log.Fatalf("Failed to load manifest: %v", err)
	}

	// Update version
	m.Plugin.Version = strings.TrimPrefix(*version, "v")

	// Update executable path
	m.Runtime.Executable.Path = fmt.Sprintf("./%s", *projectName)

	// Parse checksums if provided
	if *checksumsPath != "" {
		checksums, err := parseChecksumsFile(*checksumsPath)
		if err != nil {
			log.Fatalf("Failed to parse checksums: %v", err)
		}

		// Update checksums in manifest
		if m.Runtime.Executable.Checksums == nil {
			m.Runtime.Executable.Checksums = make(map[string]string)
		}

		// Map goreleaser artifacts to manifest format
		for filename, hash := range checksums {
			if platform := extractPlatform(filename, *projectName); platform != "" {
				m.Runtime.Executable.Checksums[platform] = fmt.Sprintf("sha256:%s", hash)
			}
		}
	}

	// Alternatively, calculate checksums from dist directory
	if *checksumsPath == "" && *distDir != "" {
		checksums, err := calculateChecksumsFromDist(*distDir, *projectName)
		if err != nil {
			log.Printf("Warning: Failed to calculate checksums from dist: %v", err)
		} else {
			m.Runtime.Executable.Checksums = checksums
		}
	}

	// Determine output path
	outputPath := *manifestPath
	if *output != "" {
		outputPath = *output
	}

	// Save updated manifest
	if err := m.Save(outputPath); err != nil {
		log.Fatalf("Failed to save manifest: %v", err)
	}

	fmt.Printf("✓ Updated manifest '%s'\n", outputPath)
	fmt.Printf("  Version: %s\n", m.Plugin.Version)
	fmt.Printf("  Executable: %s\n", m.Runtime.Executable.Path)
	fmt.Printf("  Checksums: %d platforms\n", len(m.Runtime.Executable.Checksums))
}

// parseChecksumsFile parses a goreleaser checksums.txt file
func parseChecksumsFile(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	checksums := make(map[string]string)
	scanner := bufio.NewScanner(file)
	
	// Pattern: <hash>  <filename>
	re := regexp.MustCompile(`^([a-f0-9]{64})\s+(.+)$`)
	
	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)
		if len(matches) == 3 {
			checksums[matches[2]] = matches[1]
		}
	}

	return checksums, scanner.Err()
}

// extractPlatform extracts OS-arch from goreleaser filename
func extractPlatform(filename, projectName string) string {
	// Pattern: projectname_version_os_arch.ext
	pattern := fmt.Sprintf(`%s_.*_([^_]+)_([^.]+)\.(tar\.gz|zip)`, projectName)
	re := regexp.MustCompile(pattern)
	
	matches := re.FindStringSubmatch(filename)
	if len(matches) >= 3 {
		os := matches[1]
		arch := matches[2]
		
		// Normalize architecture names
		if arch == "x86_64" {
			arch = "amd64"
		} else if arch == "i386" {
			arch = "386"
		}
		
		return fmt.Sprintf("%s-%s", os, arch)
	}
	
	return ""
}

// calculateChecksumsFromDist calculates checksums directly from dist directory
func calculateChecksumsFromDist(distDir, projectName string) (map[string]string, error) {
	checksums := make(map[string]string)
	
	// Find all binary files
	pattern := filepath.Join(distDir, fmt.Sprintf("%s_*", projectName))
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		// Skip checksums.txt and other non-archive files
		if strings.HasSuffix(file, ".txt") || strings.HasSuffix(file, ".json") {
			continue
		}

		// Calculate SHA256
		hash, err := calculateSHA256(file)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate checksum for %s: %w", file, err)
		}

		// Extract platform from filename
		basename := filepath.Base(file)
		if platform := extractPlatform(basename, projectName); platform != "" {
			checksums[platform] = fmt.Sprintf("sha256:%s", hash)
		}
	}

	return checksums, nil
}

// calculateSHA256 calculates SHA256 hash of a file
func calculateSHA256(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}