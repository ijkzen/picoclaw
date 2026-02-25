package auth

import "github.com/spf13/cobra"

func NewAuthCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage authentication (login, logout, status)",
		Long: `Manage authentication for external providers and tokens.

The auth command groups subcommands for logging into provider services,
clearing stored credentials, and inspecting current authentication
status. Use 'picoclaw auth login' to add credentials (supports device
code flow for headless environments), 'picoclaw auth logout' to remove
credentials for a provider, and 'picoclaw auth status' to view what
credentials are configured.
`,
		Example: `  picoclaw auth login --provider openai
	  picoclaw auth status
	  picoclaw auth logout --provider anthropic`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(
		newLoginCommand(),
		newLogoutCommand(),
		newStatusCommand(),
		newModelsCommand(),
	)

	return cmd
}
