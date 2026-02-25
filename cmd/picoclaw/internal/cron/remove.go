package cron

import "github.com/spf13/cobra"

func newRemoveCommand(storePath func() string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove",
		Short:   "Remove a job by ID",
		Args:    cobra.ExactArgs(1),
		Example: `picoclaw cron remove 1`,
		Long: `Remove a scheduled job by its identifier.

Provide the job ID obtained from 'picoclaw cron list'. Removing a job
deletes it from the store so it will no longer run. This operation is
destructive; use 'cron list' to confirm the id before removing.
`,
		RunE: func(_ *cobra.Command, args []string) error {
			cronRemoveCmd(storePath(), args[0])
			return nil
		},
	}

	return cmd
}
