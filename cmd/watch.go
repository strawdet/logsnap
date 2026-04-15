package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"logsnap/internal/snapshot"
)

var (
	watchInterval int
	watchLabels   map[string]string
)

func init() {
	watchCmd := &cobra.Command{
		Use:   "watch <logfile>",
		Short: "Continuously watch a log file and capture snapshots on change",
		Args:  cobra.ExactArgs(1),
		RunE:  runWatch,
	}

	watchCmd.Flags().StringVar(&snapshotDir, "dir", ".logsnap", "Directory to store snapshots")
	watchCmd.Flags().IntVar(&watchInterval, "interval", 30, "Poll interval in seconds")
	watchCmd.Flags().StringToStringVar(&watchLabels, "label", nil, "Labels to attach (key=value)")

	rootCmd.AddCommand(watchCmd)
}

func runWatch(cmd *cobra.Command, args []string) error {
	logFile := args[0]

	if _, err := os.Stat(logFile); err != nil {
		return fmt.Errorf("log file not accessible: %w", err)
	}

	if err := os.MkdirAll(snapshotDir, 0o755); err != nil {
		return fmt.Errorf("create snapshot dir: %w", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	fmt.Fprintf(cmd.OutOrStdout(), "Watching %s every %ds — press Ctrl+C to stop\n", logFile, watchInterval)

	opts := snapshot.WatchOptions{
		LogFile:  logFile,
		Dir:      snapshotDir,
		Interval: time.Duration(watchInterval) * time.Second,
		Labels:   watchLabels,
		OnSnap: func(id string, err error) {
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "error capturing snapshot: %v\n", err)
				return
			}
			fmt.Fprintf(cmd.OutOrStdout(), "snapshot captured: %s\n", id)
		},
	}

	if err := snapshot.Watch(ctx, opts); err != nil && err != context.Canceled {
		return err
	}
	fmt.Fprintln(cmd.OutOrStdout(), "watch stopped")
	return nil
}
