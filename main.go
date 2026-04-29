// dotlink is a declarative dotfile symlink manager.
//
// It scans DOTLINK_ROOT (default: $HOME/.custom) for dotlink.yaml manifests
// and creates symbolic links on your filesystem so that configuration files
// live in a single version-controlled location and are linked wherever the
// tools expect them.
//
// # Usage
//
//	dotlink <command> [flags]
//
// # Commands
//
//	sync            Create / update all symlinks declared in dotlink.yaml files.
//	list            Print every dotlink.yaml path found in DOTLINK_ROOT.
//	status          Show the link status (linked/broken/missing) per target.
//	config          Print the active runtime configuration.
//	new [directory] Create a skeleton dotlink.yaml (default: ~/.dotlink).
//	manual          Print a comprehensive manual with examples.
//
// # Environment Variables
//
//	DOTLINK_ROOT   Root directory scanned for dotlink.yaml files.
//	               Defaults to $HOME/.custom.
//
// # Examples
//
//	dotlink sync
//	dotlink sync --dry-run
//	DOTLINK_ROOT=~/dotfiles dotlink sync
//	dotlink list
//	dotlink status
//	dotlink new ~/dotfiles/nvim
package main

import "dotlink/cmd"

func main() {
	cmd.Execute()
}
