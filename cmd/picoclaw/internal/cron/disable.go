package cron

import "github.com/spf13/cobra"

func newDisableCommand(storePath func() string) *cobra.Command {
	return &cobra.Command{
		Use:     "disable",
		Short:   "Disable a job",
		Args:    cobra.ExactArgs(1),
		Example: `picoclaw cron disable 1`,
		Long: `Disable a scheduled job so it will not be executed until
re-enabled.

Provide the job id from 'picoclaw cron list'. Use this to temporarily
pause reminders without removing their configuration.
`,
		RunE: func(_ *cobra.Command, args []string) error {
			cronSetJobEnabled(storePath(), args[0], false)
			return nil
		},
	}
}
