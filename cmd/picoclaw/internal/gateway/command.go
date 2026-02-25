package gateway

import (
	"github.com/spf13/cobra"
)

func NewGatewayCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "gateway",
		Aliases: []string{"g"},
		Short:   "Manage picoclaw gateway",
		Long: `Manage the HTTP gateway which exposes webhook and health endpoints
for channel integrations (Telegram, Discord, etc.).

The gateway runs an HTTP server and handles incoming messages from
integrated chat platforms. Use 'gateway start' to run in background,
or 'gateway run' to run in foreground.
`,
		Args: cobra.NoArgs,
	}

	// Add subcommands
	cmd.AddCommand(NewStartCommand())
	cmd.AddCommand(NewStopCommand())
	cmd.AddCommand(NewStatusCommand())
	cmd.AddCommand(NewRunCommand())

	return cmd
}
