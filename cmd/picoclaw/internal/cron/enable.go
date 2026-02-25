package cron

import "github.com/spf13/cobra"

func newEnableCommand(storePath func() string) *cobra.Command {
	return &cobra.Command{
		Use:     "enable",
		Short:   "Enable a job",
		Args:    cobra.ExactArgs(1),
		Example: `picoclaw cron enable 1`,
		Long: `Enable a previously disabled job so it will be scheduled again.

Provide the job id from 'picoclaw cron list'. Enabling a job does not
change its schedule or message, only toggles its active state.
`,
		RunE: func(_ *cobra.Command, args []string) error {
			cronSetJobEnabled(storePath(), args[0], true)
			return nil
		},
	}
}
