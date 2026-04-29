package cmd

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// syncCmd implements `dotlink sync`.
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Create all symlinks defined in dotlink.yaml files",
	Long: `Scan DOTLINK_ROOT recursively for dotlink.yaml manifests and create
every symlink that is defined in those files.

If a symlink target already exists it will be removed and re-created.
Use --dry-run to preview changes without modifying the filesystem.

Examples:
  dotlink sync
  dotlink sync --dry-run
  DOTLINK_ROOT=~/dotfiles dotlink sync`,
	RunE: func(cmd *cobra.Command, args []string) error {
		root := DotlinkRoot()
		logger.Infof("Scanning root directory: %s", root)

		files, err := findDotlinkFiles(root)
		if err != nil {
			return err
		}

		if len(files) == 0 {
			logger.Warn("No dotlink.yaml files found")
			return nil
		}

		logger.Infof("Found %d dotlink configuration(s)", len(files))

		linker := NewLinker(DryRun)
		for _, configFile := range files {
			if err := linker.ProcessConfig(configFile); err != nil {
				logger.Errorf("Error processing config %s: %v", configFile, err)
			}
		}

		if DryRun {
			logger.Info("Dry-run complete, no changes were made")
		} else {
			logger.Info("Dotlink sync completed successfully")
		}
		return nil
	},
}

// findDotlinkFiles walks root and returns all dotlink.yaml paths found.
func findDotlinkFiles(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, d fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, "dotlink.yaml") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func init() {
	rootCmd.AddCommand(syncCmd)
}

