package gateway

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewRestartCommand() *cobra.Command {
	var debug bool

	cmd := &cobra.Command{
		Use:     "restart",
		Short:   "Restart the background picoclaw gateway",
		Long:    `Restart the gateway by stopping it first and then starting it again.`,
		Example: "  picoclaw gateway restart\n  picoclaw gateway restart --debug",
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := stopGateway(); err != nil {
				return err
			}
			fmt.Println("Starting gateway...")
			return startGateway(debug)
		},
	}

	cmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable debug logging")

	return cmd
}
