package cmd
import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)
// statusCmd implements `dotlink status`.
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show linking status for all dotlink.yaml files",
	Long: `Walk DOTLINK_ROOT and for every dotlink.yaml manifest display the
current link status (linked / broken / missing) of each configured target.
Examples:
  dotlink status
  DOTLINK_ROOT=~/dotfiles dotlink status`,
	RunE: func(cmd *cobra.Command, args []string) error {
		root := DotlinkRoot()
		fmt.Printf("DOTLINK_ROOT : %s\n", root)
		fmt.Printf("HOME         : %s\n\n", os.Getenv("HOME"))
		files, err := findDotlinkFiles(root)
		if err != nil {
			return err
		}
		if len(files) == 0 {
			fmt.Println("No dotlink.yaml files found.")
			return nil
		}
		for _, configFile := range files {
			cfg, err := LoadConfig(configFile)
			if err != nil {
				fmt.Printf("  [ERROR] %s: %v\n", configFile, err)
				continue
			}
			configDir := GetConfigDir(configFile)
			fmt.Printf("\n%s\n", configFile)
			if cfg.Link.To != "" {
				target, _ := ExpandPath(cfg.Link.To)
				status := linkStatus(target, configDir)
				fmt.Printf("  dir  %-45s  -> %-35s  [%s]\n", configDir, target, status)
			}
			for _, fl := range cfg.Files {
				src := configDir + "/" + fl.File
				dst, _ := ExpandPath(fl.To)
				status := linkStatus(dst, src)
				fmt.Printf("  file %-45s  -> %-35s  [%s]\n", src, dst, status)
			}
		}
		return nil
	},
}
// linkStatus returns a human-readable status string for a symlink target.
func linkStatus(target, expected string) string {
	info, err := os.Lstat(target)
	if os.IsNotExist(err) {
		return "missing"
	}
	if err != nil {
		return "error"
	}
	if info.Mode()&os.ModeSymlink == 0 {
		return "not a symlink"
	}
	dest, err := os.Readlink(target)
	if err != nil {
		return "error"
	}
	if dest == expected {
		return "✅ linked"
	}
	return fmt.Sprintf("⚠️  points to %s", dest)
}
func init() {
	rootCmd.AddCommand(statusCmd)
}
