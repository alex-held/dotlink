package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestLoadConfig tests loading a valid config file
func TestLoadConfig(t *testing.T) {
	// Create a temporary directory
	tmpDir := t.TempDir()

	// Create a test config file
	configPath := filepath.Join(tmpDir, "dotlink.yaml")
	configContent := `link:
  to: $HOME/.zsh
files:
  - file: .zshrc
    to: $HOME/.zshrc
  - file: .zsh_profile
    to: $HOME/.zsh_profile
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	// Load the config
	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Verify config
	if config.Link.To != "$HOME/.zsh" {
		t.Errorf("Expected link.to=$HOME/.zsh, got %s", config.Link.To)
	}

	if len(config.Files) != 2 {
		t.Errorf("Expected 2 files, got %d", len(config.Files))
	}

	if config.Files[0].File != ".zshrc" {
		t.Errorf("Expected first file .zshrc, got %s", config.Files[0].File)
	}

	if config.Files[0].To != "$HOME/.zshrc" {
		t.Errorf("Expected first file to=$HOME/.zshrc, got %s", config.Files[0].To)
	}
}

// TestGetConfigDir tests getting the directory from a config path
func TestGetConfigDir(t *testing.T) {
	testCases := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "simple path",
			path:     "/home/user/.config/dotlink.yaml",
			expected: "/home/user/.config",
		},
		{
			name:     "nested path",
			path:     "/a/b/c/dotlink.yaml",
			expected: "/a/b/c",
		},
		{
			name:     "root path",
			path:     "/dotlink.yaml",
			expected: "/",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GetConfigDir(tc.path)
			if result != tc.expected {
				t.Errorf("GetConfigDir(%s) = %s, want %s", tc.path, result, tc.expected)
			}
		})
	}
}

// TestExpandPath tests environment variable expansion
func TestExpandPath(t *testing.T) {
	// Set a test environment variable
	_ = os.Setenv("TEST_VAR", "/test/path")

	testCases := []struct {
		name   string
		input  string
		hasErr bool
	}{
		{
			name:   "with environment variable",
			input:  "$TEST_VAR/subdir",
			hasErr: false,
		},
		{
			name:   "with HOME variable",
			input:  "$HOME/.zshrc",
			hasErr: false,
		},
		{
			name:   "absolute path",
			input:  "/absolute/path",
			hasErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ExpandPath(tc.input)
			if (err != nil) != tc.hasErr {
				t.Errorf("ExpandPath(%s) error = %v, wantErr %v", tc.input, err, tc.hasErr)
			}
			if err == nil && result == "" {
				t.Errorf("ExpandPath(%s) returned empty string", tc.input)
			}
		})
	}
}

// TestLinkerLinkFile tests the LinkFile functionality
func TestLinkerLinkFile(t *testing.T) {
	tmpDir := t.TempDir()
	linker := NewLinker(false)

	// Create source file
	sourceDir := filepath.Join(tmpDir, "config")
	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}

	sourceFile := filepath.Join(sourceDir, ".zshrc")
	if err := os.WriteFile(sourceFile, []byte("# test"), 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Create target directory
	targetDir := filepath.Join(tmpDir, "target")
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		t.Fatalf("Failed to create target directory: %v", err)
	}

	// Set target path without environment variables
	targetPath := filepath.Join(targetDir, ".zshrc")

	// Link the file
	err := linker.LinkFile(sourceDir, ".zshrc", targetPath)
	if err != nil {
		t.Fatalf("LinkFile failed: %v", err)
	}

	// Verify symlink was created
	if _, err := os.Lstat(targetPath); os.IsNotExist(err) {
		t.Errorf("Symlink not created at %s", targetPath)
	}

	// Verify symlink points to correct file
	link, err := os.Readlink(targetPath)
	if err != nil {
		t.Fatalf("Failed to read symlink: %v", err)
	}
	if link != sourceFile {
		t.Errorf("Symlink points to %s, expected %s", link, sourceFile)
	}
}

// TestLinkerLinkDirectory tests the LinkDirectory functionality
func TestLinkerLinkDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	linker := NewLinker(false)

	// Create source directory
	sourceDir := filepath.Join(tmpDir, "config")
	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}

	// Create file in source directory
	if err := os.WriteFile(filepath.Join(sourceDir, "test.txt"), []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create target directory
	targetDir := filepath.Join(tmpDir, "target")
	targetPath := filepath.Join(targetDir, "linked")

	// Link the directory
	err := linker.LinkDirectory(sourceDir, targetPath)
	if err != nil {
		t.Fatalf("LinkDirectory failed: %v", err)
	}

	// Verify symlink was created
	if _, err := os.Lstat(targetPath); os.IsNotExist(err) {
		t.Errorf("Symlink not created at %s", targetPath)
	}

	// Verify symlink points to correct directory
	link, err := os.Readlink(targetPath)
	if err != nil {
		t.Fatalf("Failed to read symlink: %v", err)
	}
	if link != sourceDir {
		t.Errorf("Symlink points to %s, expected %s", link, sourceDir)
	}
}

// TestLinkerDryRun tests dry-run mode
func TestLinkerDryRun(t *testing.T) {
	tmpDir := t.TempDir()
	linker := NewLinker(true) // Enable dry-run

	// Create source file
	sourceDir := filepath.Join(tmpDir, "config")
	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}

	sourceFile := filepath.Join(sourceDir, ".zshrc")
	if err := os.WriteFile(sourceFile, []byte("# test"), 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	targetDir := filepath.Join(tmpDir, "target")
	targetPath := filepath.Join(targetDir, ".zshrc")

	// Link the file in dry-run mode
	err := linker.LinkFile(sourceDir, ".zshrc", targetPath)
	if err != nil {
		t.Fatalf("LinkFile failed in dry-run: %v", err)
	}

	// Verify symlink was NOT created
	if _, err := os.Lstat(targetPath); err == nil {
		t.Errorf("Symlink should not be created in dry-run mode")
	}
}

// TestProcessConfig tests the full config processing flow
func TestProcessConfig(t *testing.T) {
	tmpDir := t.TempDir()

	// Create config directory
	configDir := filepath.Join(tmpDir, "config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	// Create link target directory
	linkTargetDir := filepath.Join(tmpDir, "link_target")
	if err := os.MkdirAll(linkTargetDir, 0755); err != nil {
		t.Fatalf("Failed to create link target directory: %v", err)
	}

	// Create file target directory
	fileTargetDir := filepath.Join(tmpDir, "file_target")
	if err := os.MkdirAll(fileTargetDir, 0755); err != nil {
		t.Fatalf("Failed to create file target directory: %v", err)
	}

	// Create config file
	configPath := filepath.Join(configDir, "dotlink.yaml")
	configContent := `link:
  to: ` + filepath.Join(linkTargetDir, "linked") + `
files:
  - file: .zshrc
    to: ` + filepath.Join(fileTargetDir, ".zshrc") + `
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	// Create source file
	sourceFile := filepath.Join(configDir, ".zshrc")
	if err := os.WriteFile(sourceFile, []byte("# zshrc"), 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Process config with actual linking
	linker := NewLinker(false)
	err := linker.ProcessConfig(configPath)
	if err != nil {
		t.Fatalf("ProcessConfig failed: %v", err)
	}

	// Verify links were created
	if _, err := os.Lstat(filepath.Join(linkTargetDir, "linked")); os.IsNotExist(err) {
		t.Errorf("Directory link not created")
	}

	if _, err := os.Lstat(filepath.Join(fileTargetDir, ".zshrc")); os.IsNotExist(err) {
		t.Errorf("File link not created")
	}
}

// TestLoadConfigInvalidYAML tests error handling for invalid YAML
func TestLoadConfigInvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()

	configPath := filepath.Join(tmpDir, "dotlink.yaml")
	// Invalid YAML
	configContent := `link:
  to: $HOME/.zsh
  invalid: [
files:
  - file: .zshrc
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	_, err := LoadConfig(configPath)
	if err == nil {
		t.Error("Expected error for invalid YAML, but got none")
	}
}

