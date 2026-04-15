package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/logsnap/logsnap/internal/snapshot"
	"github.com/spf13/cobra"
)

var (
	captureLabel  string
	captureOutput string
)

var captureCmd = &cobra.Command{
	Use:   "capture",
	Short: "Capture a structured log snapshot from stdin or a file",
	Long: `Read structured log lines (one JSON object per line) and save
them as a named snapshot for later diffing.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if captureLabel == "" {
			return fmt.Errorf("--label is required")
		}

		entries, err := readLogEntries(cmd, args)
		if err != nil {
			return err
		}

		meta := map[string]string{"source": "stdin"}
		if len(args) > 0 {
			meta["source"] = args[0]
		}

		snap := snapshot.New(captureLabel, entries, meta)

		out := captureOutput
		if out == "" {
			out = fmt.Sprintf("%s.snap.json", strings.ReplaceAll(captureLabel, "/", "_"))
		}

		if err := snap.Save(out); err != nil {
			return fmt.Errorf("saving snapshot: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Snapshot saved: %s (id=%s, entries=%d)\n", out, snap.ID, len(snap.Entries))
		return nil
	},
}

func readLogEntries(cmd *cobra.Command, args []string) ([]snapshot.LogEntry, error) {
	var reader *os.File
	if len(args) > 0 {
		f, err := os.Open(args[0])
		if err != nil {
			return nil, fmt.Errorf("open log file: %w", err)
		}
		defer f.Close()
		reader = f
	} else {
		reader = os.Stdin
	}

	var entries []snapshot.LogEntry
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		entries = append(entries, snapshot.LogEntry{Message: line})
	}
	return entries, scanner.Err()
}

func init() {
	captureCmd.Flags().StringVarP(&captureLabel, "label", "l", "", "Label for the snapshot (e.g. v1.2.0)")
	captureCmd.Flags().StringVarP(&captureOutput, "output", "o", "", "Output file path (default: <label>.snap.json)")
	rootCmd.AddCommand(captureCmd)
}
