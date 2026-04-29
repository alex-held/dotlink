package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const dotlinkTemplate = `# dotlink.yaml
# ---------------------------------------------------------------------------
# Declarative dotfile symlink manifest.
#
# link.to    - create a symlink from this directory to the given path.
#              Supports environment variables such as $HOME.
# files      - individual file symlinks.  Each entry needs:
#                file: relative path inside this directory
#                to:   absolute (or env-expanded) destination path
# ---------------------------------------------------------------------------
link:
  to: ""          # e.g.  $HOME/.config/nvim
files: []
# - file: init.lua
#   to:   $HOME/.config/nvim/init.lua
`

// newCmd implements `dotlink new [directory]`.
var newCmd = &cobra.Command{
	Use:   "new [directory]",
	Short: "Initialise a new dotlink.yaml in the given directory",
	Long: `Create a skeleton dotlink.yaml in DIRECTORY (default: ~/.dotlink).
The file is populated with commented examples so you can get started quickly.
If the file already exists the command exits without overwriting it.
Examples:
  dotlink new
  dotlink new ~/dotfiles/nvim
  dotlink new ~/.config/zsh`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("cannot get current directory: %w", err)
		}
		if len(args) == 1 {
			expanded, err := ExpandPath(args[0])
			if err != nil {
				return fmt.Errorf("invalid directory: %w", err)
			}
			dir = expanded
		}
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("cannot create directory %s: %w", dir, err)
		}
		dest := filepath.Join(dir, "dotlink.yaml")
		if _, err := os.Stat(dest); err == nil {
			fmt.Printf("dotlink.yaml already exists at %s\n", dest)
			return nil
		}
		if err := os.WriteFile(dest, []byte(dotlinkTemplate), 0644); err != nil {
			return fmt.Errorf("cannot write %s: %w", dest, err)
		}
		fmt.Printf("✅ Created %s\n", dest)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
