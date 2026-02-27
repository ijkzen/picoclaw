package gateway

import (
	"fmt"
	"time"

	"github.com/sipeed/picoclaw/pkg/logger"
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
			logger.InfoCF("gateway", "Gateway restart requested",
				map[string]any{
					"debug": debug,
					"delay": delay,
				})

			if delay > 0 {
				time.Sleep(time.Duration(delay) * time.Second)
			}
			if err := stopGateway(); err != nil {
				logger.ErrorCF("gateway", "Gateway restart stop phase failed", map[string]any{"error": err.Error()})
				return err
			}
			fmt.Println("Starting gateway...")
			if err := startGateway(debug); err != nil {
				logger.ErrorCF("gateway", "Gateway restart start phase failed", map[string]any{"error": err.Error()})
				return err
			}

			logger.InfoC("gateway", "Gateway restart completed")
			return nil
		},
	}

	cmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable debug logging")
	cmd.Flags().IntVar(&delay, "delay", 0, "Delay restart by N seconds")

	return cmd
}
