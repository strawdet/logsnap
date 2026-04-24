package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"logsnap/internal/snapshot"
)

func init() {
	schemaCmd := &cobra.Command{
		Use:   "schema",
		Short: "Manage validation schemas for snapshot entries",
	}

	saveCmd := &cobra.Command{
		Use:   "save <name>",
		Short: "Save a new validation schema",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, _ := cmd.Flags().GetString("dir")
			required, _ := cmd.Flags().GetStringSlice("required")
			allowed, _ := cmd.Flags().GetStringSlice("allowed-levels")
			desc, _ := cmd.Flags().GetString("description")
			s := &snapshot.Schema{
				ID:             args[0],
				Name:           args[0],
				Description:    desc,
				RequiredFields: required,
				AllowedLevels:  allowed,
			}
			if err := snapshot.SaveSchema(dir, s); err != nil {
				return fmt.Errorf("save schema: %w", err)
			}
			fmt.Printf("Schema %q saved.\n", args[0])
			return nil
		},
	}
	saveCmd.Flags().String("dir", defaultSnapshotDir(), "snapshot directory")
	saveCmd.Flags().StringSlice("required", []string{"message", "level"}, "required fields")
	saveCmd.Flags().StringSlice("allowed-levels", nil, "allowed log levels (empty = any)")
	saveCmd.Flags().String("description", "", "schema description")

	validateCmd := &cobra.Command{
		Use:   "validate <snapshot-id> <schema-name>",
		Short: "Validate a snapshot against a schema",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, _ := cmd.Flags().GetString("dir")
			snap, err := snapshot.Load(dir, args[0])
			if err != nil {
				return fmt.Errorf("load snapshot: %w", err)
			}
			schema, err := snapshot.LoadSchema(dir, args[1])
			if err != nil {
				return fmt.Errorf("load schema: %w", err)
			}
			violations := snapshot.ValidateSnapshot(snap, schema)
			if len(violations) == 0 {
				fmt.Println("Validation passed: no violations found.")
				return nil
			}
			fmt.Fprintf(os.Stderr, "Validation failed: %d violation(s)\n", len(violations))
			for _, v := range violations {
				fmt.Fprintf(os.Stderr, "  [entry %d] field=%q reason=%s\n", v.EntryIndex, v.Field, v.Reason)
			}
			return fmt.Errorf("schema validation failed")
		},
	}
	validateCmd.Flags().String("dir", defaultSnapshotDir(), "snapshot directory")

	showCmd := &cobra.Command{
		Use:   "show <name>",
		Short: "Show a schema as JSON",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, _ := cmd.Flags().GetString("dir")
			s, err := snapshot.LoadSchema(dir, args[0])
			if err != nil {
				return err
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(s)
		},
	}
	showCmd.Flags().String("dir", defaultSnapshotDir(), "snapshot directory")

	deleteCmd := &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a schema",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, _ := cmd.Flags().GetString("dir")
			if err := snapshot.DeleteSchema(dir, args[0]); err != nil {
				return err
			}
			fmt.Printf("Schema %q deleted.\n", args[0])
			return nil
		},
	}
	deleteCmd.Flags().String("dir", defaultSnapshotDir(), "snapshot directory")

	_ = strings.TrimSpace("")
	schemaCmd.AddCommand(saveCmd, validateCmd, showCmd, deleteCmd)
	rootCmd.AddCommand(schemaCmd)
}
