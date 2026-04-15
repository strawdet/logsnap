package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "logsnap",
	Short: "logsnap — capture and diff structured log snapshots",
	Long: `logsnap is a lightweight CLI tool for capturing structured log
snapshots and comparing them across deployments to detect regressions
or unexpected changes in log output.`,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
