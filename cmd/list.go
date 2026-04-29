package cmd
import (
	"fmt"
	"github.com/spf13/cobra"
)
// listCmd implements `dotlink list`.
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Display all dotlink.yaml files found in DOTLINK_ROOT",
	Long: `Walk DOTLINK_ROOT and print the path of every dotlink.yaml manifest
that dotlink would process during a sync.
Examples:
  dotlink list
  DOTLINK_ROOT=~/dotfiles dotlink list`,
	RunE: func(cmd *cobra.Command, args []string) error {
		root := DotlinkRoot()
		files, err := findDotlinkFiles(root)
		if err != nil {
			return err
		}
		if len(files) == 0 {
			fmt.Printf("No dotlink.yaml files found in %s\n", root)
			return nil
		}
		fmt.Printf("Dotlink manifests in %s:\n\n", root)
		for _, f := range files {
			fmt.Printf("  %s\n", f)
		}
		return nil
	},
}
func init() {
	rootCmd.AddCommand(listCmd)
}
