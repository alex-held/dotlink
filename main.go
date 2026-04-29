package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
)


func main() {
	root := os.Getenv("DOTLINK_ROOT")
	if root == "" {
		root = os.ExpandEnv("$HOME/.custom")
	}

	logger := log.New(os.Stderr)
	logger.SetLevel(log.DebugLevel)

	logger.Infof("🔍 Scanning root directory: %s", root)

	var files []string
	err := filepath.Walk(root, func(path string, d fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		logger.Debugf("📂 Scanning: %s (type=%s)", path, d.Mode().Type().String())

		if strings.HasSuffix(path, "dotlink.yaml") {
			files = append(files, path)
			logger.Debugf("✨ Found dotlink config: %s", path)
		}

		return nil
	})
	if err != nil {
		logger.Fatalf("❌ Error walking directory: %v", err)
	}

	if len(files) == 0 {
		logger.Warn("⚠️  No dotlink.yaml files found")
		return
	}

	logger.Infof("📋 Found %d dotlink configuration(s)", len(files))

	// Process each config file
	linker := NewLinker(false) // Set to true for dry-run mode
	for _, configFile := range files {
		if err := linker.ProcessConfig(configFile); err != nil {
			logger.Errorf("❌ Error processing config %s: %v", configFile, err)
		}
	}

	logger.Infof("✅ Dotlink completed successfully")
}
