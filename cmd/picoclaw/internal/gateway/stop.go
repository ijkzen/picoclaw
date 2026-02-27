package gateway

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/sipeed/picoclaw/pkg/logger"
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

func stopGateway() error {
	logger.InfoC("gateway", "Gateway stop requested")

	running, pid, err := isRunning()
	if err != nil {
		logger.ErrorCF("gateway", "Gateway stop status check failed", map[string]any{"error": err.Error()})
		return fmt.Errorf("failed to check gateway status: %w", err)
	}
	if !running {
		logger.InfoC("gateway", "Gateway stop skipped: not running")
		// Gateway is not running - do nothing as requested
		return nil
	}

	// Send SIGTERM to the process
	if err := signalProcess(pid, syscall.SIGTERM); err != nil {
		logger.ErrorCF("gateway", "Gateway stop signal failed", map[string]any{
			"pid":   pid,
			"error": err.Error(),
		})
		return fmt.Errorf("failed to stop gateway (PID: %d): %w", pid, err)
	}

	// Wait for process to stop (with timeout)
	pidFile := getPIDFile()
	for i := 0; i < 30; i++ {
		if !isProcessRunning(pid) {
			// Process no longer exists
			_ = os.Remove(pidFile)
			logger.InfoCF("gateway", "Gateway stopped", map[string]any{"pid": pid})
			fmt.Println("✓ Gateway stopped")
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}

	// If still running, force kill
	if err := signalProcess(pid, syscall.SIGKILL); err == nil {
		_ = os.Remove(pidFile)
		logger.WarnCF("gateway", "Gateway stopped with force kill", map[string]any{"pid": pid})
		fmt.Println("✓ Gateway stopped (forced)")
		return nil
	}

	logger.ErrorCF("gateway", "Gateway stop failed", map[string]any{"pid": pid})
	return fmt.Errorf("failed to stop gateway (PID: %d)", pid)
}

func NewStopCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "stop",
		Short:   "Stop the background picoclaw gateway",
		Long:    `Stop the gateway that is running in the background.`,
		Example: `  picoclaw gateway stop`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return stopGateway()
		},
	}

	return cmd
}
