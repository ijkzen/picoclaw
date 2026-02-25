package gateway

import (
	"fmt"

	"github.com/sipeed/picoclaw/cmd/picoclaw/internal"
	"github.com/spf13/cobra"
)

func NewStatusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "status",
		Short:   "Show picoclaw gateway status",
		Long:    `Display whether the gateway is running and its PID if available.`,
		Example: `  picoclaw gateway status`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			running, pid, err := isRunning()
			if err != nil {
				return fmt.Errorf("failed to check gateway status: %w", err)
			}

			cfg, err := internal.LoadConfig()
			if err != nil {
				return fmt.Errorf("error loading config: %w", err)
			}

			if running {
				fmt.Printf("Gateway is running\n")
				fmt.Printf("  PID: %d\n", pid)
				fmt.Printf("  Address: %s:%d\n", cfg.Gateway.Host, cfg.Gateway.Port)
				fmt.Printf("  Health: http://%s:%d/health\n", cfg.Gateway.Host, cfg.Gateway.Port)
			} else {
				fmt.Println("Gateway is not running")
				fmt.Printf("  Configured address: %s:%d\n", cfg.Gateway.Host, cfg.Gateway.Port)
			}

			return nil
		},
	}

	return cmd
}
