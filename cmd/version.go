// Package cmd provides the version subcommand and build metadata.
package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	// Version is populated during compile time via ldflags.
	Version = "v1.0.0"
	// GitCommit is populated during compile time via ldflags.
	GitCommit = "dev"
	// BuildDate is populated during compile time via ldflags.
	BuildDate = "unknown"
)

// versionCmd represents the version command.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number, git commit hash, and build information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("kube-cost-lens %s\n", Version)
		fmt.Printf("  Commit:     %s\n", GitCommit)
		fmt.Printf("  Build Date: %s\n", BuildDate)
		fmt.Printf("  Go Version: %s\n", runtime.Version())
		fmt.Printf("  OS/Arch:    %s/%s\n", runtime.GOOS, runtime.GOARCH)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
