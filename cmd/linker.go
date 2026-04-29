package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
)

// Linker handles the creation of symbolic links
type Linker struct {
	logger *log.Logger
	dryRun bool
}

// NewLinker creates a new Linker instance
func NewLinker(dryRun bool) *Linker {
	return &Linker{
		logger: log.New(os.Stderr),
		dryRun: dryRun,
	}
}

// LinkDirectory creates a symbolic link from the config directory to the target
func (l *Linker) LinkDirectory(configDir, targetPath string) error {
	targetPath, err := ExpandPath(targetPath)
	if err != nil {
		l.logger.Errorf("❌ Failed to expand target path: %v", err)
		return err
	}

	// Create parent directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		l.logger.Errorf("❌ Failed to create parent directory: %v", err)
		return err
	}

	// Remove existing symlink/directory if it exists
	if _, err := os.Lstat(targetPath); err == nil {
		if l.dryRun {
			l.logger.Debugf("🔍 [DRY RUN] Would remove: %s", targetPath)
		} else {
			if err := os.RemoveAll(targetPath); err != nil {
				l.logger.Errorf("❌ Failed to remove existing target: %v", err)
				return err
			}
		}
	}

	if l.dryRun {
		l.logger.Infof("🔗 [DRY RUN] Would link directory: %s -> %s", configDir, targetPath)
	} else {
		if err := os.Symlink(configDir, targetPath); err != nil {
			l.logger.Errorf("❌ Failed to create symlink: %v", err)
			return err
		}
		l.logger.Infof("🔗 Linked directory: %s -> %s", configDir, targetPath)
	}

	return nil
}

// LinkFile creates a symbolic link for a single file
func (l *Linker) LinkFile(configDir, sourceFile, targetPath string) error {
	sourcePath := filepath.Join(configDir, sourceFile)
	targetPath, err := ExpandPath(targetPath)
	if err != nil {
		l.logger.Errorf("❌ Failed to expand target path: %v", err)
		return err
	}

	// Check if source file exists
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		l.logger.Warnf("⚠️  Source file not found: %s", sourcePath)
		return fmt.Errorf("source file not found: %s", sourcePath)
	}

	// Create parent directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		l.logger.Errorf("❌ Failed to create parent directory: %v", err)
		return err
	}

	// Remove existing symlink/file if it exists
	if _, err := os.Lstat(targetPath); err == nil {
		if l.dryRun {
			l.logger.Debugf("🔍 [DRY RUN] Would remove: %s", targetPath)
		} else {
			if err := os.Remove(targetPath); err != nil {
				l.logger.Errorf("❌ Failed to remove existing target: %v", err)
				return err
			}
		}
	}

	if l.dryRun {
		l.logger.Infof("🔗 [DRY RUN] Would link file: %s -> %s", sourcePath, targetPath)
	} else {
		if err := os.Symlink(sourcePath, targetPath); err != nil {
			l.logger.Errorf("❌ Failed to create symlink: %v", err)
			return err
		}
		l.logger.Infof("🔗 Linked file: %s -> %s", sourcePath, targetPath)
	}

	return nil
}

// ProcessConfig processes a dotlink configuration
func (l *Linker) ProcessConfig(configPath string) error {
	l.logger.Infof("📖 Processing config: %s", configPath)

	config, err := LoadConfig(configPath)
	if err != nil {
		l.logger.Errorf("❌ Failed to load config: %v", err)
		return err
	}

	configDir := GetConfigDir(configPath)
	l.logger.Debugf("Config directory: %s", configDir)

	// Process directory link
	if config.Link.To != "" {
		l.logger.Infof("📁 Linking directory to: %s", config.Link.To)
		if err := l.LinkDirectory(configDir, config.Link.To); err != nil {
			l.logger.Errorf("❌ Failed to link directory: %v", err)
			return err
		}
	}

	// Process file links
	if len(config.Files) > 0 {
		l.logger.Infof("📄 Linking %d files", len(config.Files))
		for _, fileLink := range config.Files {
			if err := l.LinkFile(configDir, fileLink.File, fileLink.To); err != nil {
				l.logger.Errorf("❌ Failed to link file: %v", err)
				return err
			}
		}
	}

	l.logger.Infof("✅ Successfully processed config: %s", configPath)
	return nil
}

