package gateway

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

func NewRestartCommand() *cobra.Command {
	var debug bool
	var delay int

	cmd := &cobra.Command{
		Use:     "restart",
		Short:   "Restart the background picoclaw gateway",
		Long:    `Restart the gateway by stopping it first and then starting it again.`,
		Example: "  picoclaw gateway restart\n  picoclaw gateway restart --debug",
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if delay > 0 {
				time.Sleep(time.Duration(delay) * time.Second)
			}
			if err := stopGateway(); err != nil {
				return err
			}
			fmt.Println("Starting gateway...")
			return startGateway(debug)
		},
	}

	cmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable debug logging")
	cmd.Flags().IntVar(&delay, "delay", 0, "Delay restart by N seconds")

	return cmd
}
