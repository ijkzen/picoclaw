package gateway

import (
	"github.com/spf13/cobra"
)

func NewRunCommand() *cobra.Command {
	var debug bool

	cmd := &cobra.Command{
		Use:    "run",
		Short:  "Run picoclaw gateway in foreground (internal use)",
		Long:   `Run the HTTP gateway in foreground. This command is typically used internally by 'gateway start'.`,
		Args:   cobra.NoArgs,
		Hidden: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			return gatewayCmd(debug)
		},
	}

	cmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable debug logging")

	return cmd
}
