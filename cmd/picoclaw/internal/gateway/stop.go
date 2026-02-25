package gateway

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

func isProcessRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	// Signal 0 is used to check if process exists
	err = process.Signal(syscall.Signal(0))
	return err == nil
}

func NewStopCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "stop",
		Short:   "Stop the background picoclaw gateway",
		Long:    `Stop the gateway that is running in the background.`,
		Example: `  picoclaw gateway stop`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			running, pid, err := isRunning()
			if err != nil {
				return fmt.Errorf("failed to check gateway status: %w", err)
			}
			if !running {
				// Gateway is not running - do nothing as requested
				return nil
			}

			// Send SIGTERM to the process
			if err := signalProcess(pid, syscall.SIGTERM); err != nil {
				return fmt.Errorf("failed to stop gateway (PID: %d): %w", pid, err)
			}

			// Wait for process to stop (with timeout)
			pidFile := getPIDFile()
			for i := 0; i < 30; i++ {
				if !isProcessRunning(pid) {
					// Process no longer exists
					_ = os.Remove(pidFile)
					fmt.Println("✓ Gateway stopped")
					return nil
				}
				time.Sleep(100 * time.Millisecond)
			}

			// If still running, force kill
			if err := signalProcess(pid, syscall.SIGKILL); err == nil {
				_ = os.Remove(pidFile)
				fmt.Println("✓ Gateway stopped (forced)")
				return nil
			}

			return fmt.Errorf("failed to stop gateway (PID: %d)", pid)
		},
	}

	return cmd
}
