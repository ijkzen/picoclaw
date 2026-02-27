// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/sipeed/picoclaw/cmd/picoclaw/internal"
	"github.com/sipeed/picoclaw/cmd/picoclaw/internal/agent"
	"github.com/sipeed/picoclaw/cmd/picoclaw/internal/auth"
	"github.com/sipeed/picoclaw/cmd/picoclaw/internal/channel"
	"github.com/sipeed/picoclaw/cmd/picoclaw/internal/cron"
	"github.com/sipeed/picoclaw/cmd/picoclaw/internal/gateway"
	"github.com/sipeed/picoclaw/cmd/picoclaw/internal/migrate"
	"github.com/sipeed/picoclaw/cmd/picoclaw/internal/models"
	"github.com/sipeed/picoclaw/cmd/picoclaw/internal/onboard"
	"github.com/sipeed/picoclaw/cmd/picoclaw/internal/skills"
	"github.com/sipeed/picoclaw/cmd/picoclaw/internal/status"
	"github.com/sipeed/picoclaw/cmd/picoclaw/internal/version"
	"github.com/sipeed/picoclaw/pkg/logger"
)

func NewPicoclawCommand() *cobra.Command {
	short := fmt.Sprintf("%s picoclaw - Personal AI Assistant v%s\n\n", internal.Logo, internal.GetVersion())

	cmd := &cobra.Command{
		Use:     "picoclaw",
		Short:   short,
		Example: "picoclaw list",
		PersistentPreRun: func(command *cobra.Command, _ []string) {
			if command == nil {
				return
			}

			logger.InfoCF("cli", "Command operation",
				map[string]any{
					"command": command.CommandPath(),
				})
		},
	}

	cmd.AddCommand(
		onboard.NewOnboardCommand(),
		agent.NewAgentCommand(),
		auth.NewAuthCommand(),
		channel.NewChannelCommand(),
		gateway.NewGatewayCommand(),
		models.NewModelsCommand(),
		status.NewStatusCommand(),
		cron.NewCronCommand(),
		migrate.NewMigrateCommand(),
		skills.NewSkillsCommand(),
		version.NewVersionCommand(),
	)

	return cmd
}

func main() {
	initFileLogging()

	cmd := NewPicoclawCommand()
	if err := cmd.Execute(); err != nil {
		logger.ErrorCF("cli", "Command execution failed", map[string]any{
			"error": err.Error(),
		})
		os.Exit(1)
	}
}

func initFileLogging() {
	if strings.TrimSpace(os.Getenv("PICOCLAW_DISABLE_FILE_LOGGING")) == "1" {
		return
	}

	if strings.HasSuffix(filepath.Base(os.Args[0]), ".test") {
		return
	}

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: failed to resolve home directory for logging: %v\n", err)
		return
	}

	logDir := filepath.Join(home, ".picoclaw", "logs")
	if err := logger.EnableDailyFileLogging(logDir, 7); err != nil {
		fmt.Fprintf(os.Stderr, "warning: failed to initialize file logging: %v\n", err)
	}
}
