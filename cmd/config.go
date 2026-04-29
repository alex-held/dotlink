package cmd
import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)
// configCmd implements `dotlink config`.
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Display the active dotlink configuration",
	Long: `Print the effective dotlink runtime configuration: DOTLINK_ROOT,
HOME, and the list of dotlink.yaml manifests that would be processed.
Examples:
  dotlink config`,
	RunE: func(cmd *cobra.Command, args []string) error {
		root := DotlinkRoot()
		home := os.Getenv("HOME")
		fmt.Println("dotlink configuration")
		fmt.Println("=====================")
		fmt.Printf("DOTLINK_ROOT  : %s\n", root)
		fmt.Printf("HOME          : %s\n", home)
		fmt.Printf("dry-run       : %v\n\n", DryRun)
		files, err := findDotlinkFiles(root)
		if err != nil {
			return err
		}
		if len(files) == 0 {
			fmt.Println("No dotlink.yaml manifests found.")
			return nil
		}
		fmt.Printf("Manifests (%d):\n", len(files))
		for _, f := range files {
			cfg, err := LoadConfig(f)
			if err != nil {
				fmt.Printf("  [ERROR] %s: %v\n", f, err)
				continue
			}
			fmt.Printf("  %s  (link.to=%q, files=%d)\n", f, cfg.Link.To, len(cfg.Files))
		}
		return nil
	},
}
func init() {
	rootCmd.AddCommand(configCmd)
}
