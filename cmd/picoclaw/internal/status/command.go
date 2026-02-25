package status

import (
	"github.com/spf13/cobra"
)

func NewStatusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "status",
		Aliases: []string{"s"},
		Short:   "Show picoclaw status",
		Long: `Display high-level status information about the picoclaw
runtime, including configured providers, workspace path, and basic
health checks.

Useful to quickly confirm environment configuration before running
commands or starting the gateway. This command is informational and
does not modify state.
`,
		Example: `  picoclaw status`,
		Run: func(cmd *cobra.Command, args []string) {
			statusCmd()
		},
	}

	return cmd
}
