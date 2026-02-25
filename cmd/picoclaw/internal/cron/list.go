package cron

import "github.com/spf13/cobra"

func newListCommand(storePath func() string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all scheduled jobs",
		Long: `Display all scheduled jobs stored in the workspace cron store.

This shows job id, name, schedule, next run time and enabled status.
Useful to inspect current reminders and confirm which jobs are active.
`,
		Example: `  picoclaw cron list`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			cronListCmd(storePath())
			return nil
		},
	}

	return cmd
}
