package channel

import (
	"github.com/spf13/cobra"
)

// NewChannelCommand creates a new channel command
func NewChannelCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "channel",
		Aliases: []string{"ch"},
		Short:   "Configure chat channel integrations",
		Long: `Configure and manage chat channel integrations for picoclaw.

This command opens an interactive TUI for configuring various chat platforms
including Telegram, Discord, Slack, WeCom, Feishu, and more.

Each channel can be enabled/disabled and configured with the necessary
credentials and settings. The configuration is tested before saving.

Supported channels:
  - Telegram, Discord, Slack, QQ
  - DingTalk, WeCom (Bot & App), Feishu
  - LINE, OneBot, MaixCam, WhatsApp`,
		Example: `  picoclaw channel`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return runChannelTUI()
		},
	}

	return cmd
}
