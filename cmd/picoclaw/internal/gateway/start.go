package gateway

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/sipeed/picoclaw/cmd/picoclaw/internal"
	"github.com/spf13/cobra"
)

func getPIDFile() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".picoclaw", "gateway.pid")
}

func isRunning() (bool, int, error) {
	pidFile := getPIDFile()
	data, err := os.ReadFile(pidFile)
	if err != nil {
		if os.IsNotExist(err) {
			return false, 0, nil
		}
		return false, 0, err
	}

	var pid int
	if _, err := fmt.Sscanf(string(data), "%d", &pid); err != nil {
		return false, 0, nil
	}

	// Check if process exists by sending signal 0
	process, err := os.FindProcess(pid)
	if err != nil {
		_ = os.Remove(pidFile)
		return false, 0, nil
	}
	if err := process.Signal(syscall.Signal(0)); err != nil {
		// Process not running, clean up stale pid file
		_ = os.Remove(pidFile)
		return false, 0, nil
	}

	return true, pid, nil
}

func signalProcess(pid int, sig os.Signal) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return process.Signal(sig)
}

func startGateway(debug bool) error {
	// Check if already running
	running, pid, err := isRunning()
	if err != nil {
		return fmt.Errorf("failed to check gateway status: %w", err)
	}
	if running {
		fmt.Printf("Gateway is already running (PID: %d)\n", pid)
		return nil
	}

	// Get the executable path
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Load config to get port
	cfg, err := internal.LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	// Build command arguments - use the original gateway command for foreground run
	// but we need to run it in background
	args := []string{"gateway", "run"}
	if debug {
		args = append(args, "--debug")
	}

	// Create command
	cmd := exec.Command(exe, args...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = nil

	// Start the process
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start gateway: %w", err)
	}

	// Write PID file
	pidFile := getPIDFile()
	if err := os.WriteFile(pidFile, []byte(fmt.Sprintf("%d", cmd.Process.Pid)), 0o644); err != nil {
		// Try to kill the process if we can't write PID file
		_ = cmd.Process.Kill()
		return fmt.Errorf("failed to write PID file: %w", err)
	}

	fmt.Printf("✓ Gateway started in background (PID: %d)\n", cmd.Process.Pid)
	fmt.Printf("  Listening on %s:%d\n", cfg.Gateway.Host, cfg.Gateway.Port)
	fmt.Println("  Use 'picoclaw gateway stop' to stop")

	// Detach from the child process
	_ = cmd.Process.Release()

	return nil
}

func NewStartCommand() *cobra.Command {
	var debug bool

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start picoclaw gateway in the background",
		Long: `Start the HTTP gateway in the background which exposes webhook and health endpoints
for channel integrations (Telegram, Discord, etc.).

The gateway will continue running even after the command line is closed.
Use 'picoclaw gateway stop' to stop the background process.
`,
		Example: `  picoclaw gateway start
  picoclaw gateway start --debug`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return startGateway(debug)
		},
	}

	cmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable debug logging")

	return cmd
}
