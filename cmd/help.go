package cmd
import (
	"fmt"
	"github.com/spf13/cobra"
)
// helpManualCmd implements `dotlink help` with a comprehensive manual.
var helpManualCmd = &cobra.Command{
	Use:   "manual",
	Short: "Display a comprehensive manual with examples",
	Long:  "Print the full dotlink manual.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(manual)
	},
}
const manual = `
dotlink — declarative dotfile symlink manager
=============================================
SYNOPSIS
  dotlink <command> [flags]
DESCRIPTION
  dotlink manages dotfiles by scanning a root directory for dotlink.yaml
  manifest files and creating symbolic links on your filesystem.
  Each directory in DOTLINK_ROOT may contain a dotlink.yaml that declares:
    - A directory-level link  (the whole directory is linked somewhere)
    - Individual file links   (specific files are linked to specific paths)
ENVIRONMENT VARIABLES
  DOTLINK_ROOT    Root directory scanned recursively for dotlink.yaml files.
                  Default: $HOME/.custom
GLOBAL FLAGS
  -n, --dry-run   Preview changes without writing anything to disk.
COMMANDS
  sync            Scan DOTLINK_ROOT and create / update all symlinks.
  list            Print every dotlink.yaml path that would be processed.
  status          Show the link status (linked/broken/missing) per target.
  config          Print the active runtime configuration.
  new [dir]       Create a skeleton dotlink.yaml in dir (default ~/.dotlink).
  manual          Show this manual.
DOTLINK.YAML FORMAT
  link:
    to: $HOME/.config/nvim      # symlink this directory to the target path
  files:
    - file: init.lua            # path relative to the dotlink.yaml directory
      to:   $HOME/.config/nvim/init.lua
  Environment variables are expanded in all path fields.
EXAMPLES
  # Run a full sync
  dotlink sync
  # Preview what would happen
  dotlink sync --dry-run
  # Use a custom dotfile root
  DOTLINK_ROOT=~/dotfiles dotlink sync
  # List all manifests
  dotlink list
  # Check link health
  dotlink status
  # Initialise a manifest for your neovim config
  dotlink new ~/dotfiles/nvim
  # See the active config
  dotlink config
SEE ALSO
  https://pkg.go.dev/dotlink
`
func init() {
	rootCmd.AddCommand(helpManualCmd)
}
