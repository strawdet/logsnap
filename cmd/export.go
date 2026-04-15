package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"logsnap/internal/snapshot"
)

var (
	exportFormat string
	exportOutput string
)

func init() {
	exportCmd := &cobra.Command{
		Use:   "export <snapshot-id>",
		Short: "Export a snapshot to JSON or CSV",
		Args:  cobra.ExactArgs(1),
		RunE:  runExport,
	}

	exportCmd.Flags().StringVarP(&exportFormat, "format", "f", "json", "Output format: json or csv")
	exportCmd.Flags().StringVarP(&exportOutput, "output", "o", "", "Destination file path (default: <id>.<format>)")

	rootCmd.AddCommand(exportCmd)
}

func runExport(cmd *cobra.Command, args []string) error {
	snapshotID := args[0]

	dir := os.Getenv("LOGSNAP_DIR")
	if dir == "" {
		dir = ".logsnap"
	}

	snap, err := snapshot.Load(dir, snapshotID)
	if err != nil {
		return fmt.Errorf("loading snapshot %q: %w", snapshotID, err)
	}

	fmt := snapshot.ExportFormat(exportFormat)

	dest := exportOutput
	if dest == "" {
		dest = filepath.Join(dir, fmt.Sprintf("%s.%s", snapshotID, exportFormat))
	}

	if err := snapshot.ExportSnapshot(snap, dest, fmt); err != nil {
		return fmt.Errorf("exporting snapshot: %w", err)
	}

	cmd.Printf("Snapshot %q exported to %s\n", snapshotID, dest)
	return nil
}
