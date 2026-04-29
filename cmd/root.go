// Package cmd provides the CLI commands for dotlink.
//
// dotlink is a dotfile symlink manager. It scans a root directory for
// dotlink.yaml manifests and creates symbolic links as configured.
//
// # Environment Variables
//
//	DOTLINK_ROOT  Root directory to scan for dotlink.yaml files.
//	              Defaults to $HOME/.custom
//
// # Examples
//
//	dotlink sync
//	dotlink list
//	dotlink status
//	dotlink config
//	dotlink new ~/.dotfiles/nvim
package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var (
	// DryRun controls whether sync actually writes symlinks.
	DryRun bool

	logger *log.Logger
)

// DotlinkRoot returns the effective DOTLINK_ROOT directory.
func DotlinkRoot() string {
	if r := os.Getenv("DOTLINK_ROOT"); r != "" {
		return r
	}
	return os.ExpandEnv("$HOME/.custom")
}

// rootCmd is the base command when dotlink is called without subcommands.
var rootCmd = &cobra.Command{
	Use:   "dotlink",
	Short: "Manage dotfile symlinks declaratively",
	Long: `dotlink — a declarative dotfile symlink manager.

It scans DOTLINK_ROOT (default: $HOME/.custom) for dotlink.yaml manifests
and creates symbolic links on your filesystem so that your configuration files
live in one place and are linked wherever the tools expect them.

Environment Variables:
  DOTLINK_ROOT   Root directory scanned for dotlink.yaml files.
                 Defaults to $HOME/.custom

Run 'dotlink help' for a full manual with examples.`,
}

// Execute runs the root command. Call this from main.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	logger = log.New(os.Stderr)

	rootCmd.PersistentFlags().BoolVarP(&DryRun, "dry-run", "n", false,
		"Print what would be done without making any changes")
}

