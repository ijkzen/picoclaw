package gateway

import (
	"github.com/spf13/cobra"
)

func NewGatewayCommand() *cobra.Command {
	var debug bool

	cmd := &cobra.Command{
		Use:     "gateway",
		Aliases: []string{"g"},
		Short:   "Start picoclaw gateway",
		Long: `Start the HTTP gateway which exposes webhook and health endpoints
for channel integrations (Telegram, Discord, etc.).

The gateway runs an HTTP server and handles incoming messages from
integrated chat platforms. Use --debug to enable verbose logging. The
gateway should generally be run on a server or container and may need
ports opened for webhooks.
`,
		Example: `  picoclaw gateway
	  picoclaw gateway --debug`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return gatewayCmd(debug)
		},
	}

	cmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable debug logging")

	return cmd
}
