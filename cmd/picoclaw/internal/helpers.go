package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sipeed/picoclaw/pkg/config"
)

const Logo = "🦞"

var (
	version   = "dev"
	gitCommit string
	buildTime string
	goVersion string
)

func GetConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".picoclaw", "config.json")
}

func LoadConfig() (*config.Config, error) {
	return config.LoadConfig(GetConfigPath())
}

func SaveConfigAndRestart(cfg *config.Config) error {
	return SaveConfigPathAndRestart(GetConfigPath(), cfg)
}

func SaveConfigPathAndRestart(path string, cfg *config.Config) error {
	if err := config.SaveConfig(path, cfg); err != nil {
		return err
	}
	return restartGatewayAfterConfigChange()
}

func restartGatewayAfterConfigChange() error {
	if shouldSkipAutoGatewayRestart() {
		return nil
	}

	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to resolve executable for gateway restart: %w", err)
	}

	cmd := exec.Command(exe, "gateway", "restart")
	cmd.Stdin = nil
	output, err := cmd.CombinedOutput()
	if err != nil {
		if out := strings.TrimSpace(string(output)); out != "" {
			return fmt.Errorf("gateway restart failed: %w: %s", err, out)
		}
		return fmt.Errorf("gateway restart failed: %w", err)
	}
	return nil
}

func shouldSkipAutoGatewayRestart() bool {
	v := strings.TrimSpace(strings.ToLower(os.Getenv("PICOCLAW_DISABLE_AUTO_RESTART")))
	if v == "1" || v == "true" || v == "yes" {
		return true
	}

	// Tests run from *.test binaries; skip auto-restart there.
	return strings.HasSuffix(filepath.Base(os.Args[0]), ".test")
}

// FormatVersion returns the version string with optional git commit
func FormatVersion() string {
	v := version
	if gitCommit != "" {
		v += fmt.Sprintf(" (git: %s)", gitCommit)
	}
	return v
}

// FormatBuildInfo returns build time and go version info
func FormatBuildInfo() (string, string) {
	build := buildTime
	goVer := goVersion
	if goVer == "" {
		goVer = runtime.Version()
	}
	return build, goVer
}

// GetVersion returns the version string
func GetVersion() string {
	return version
}
