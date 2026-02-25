package auth

import "github.com/spf13/cobra"

func newLogoutCommand() *cobra.Command {
	var provider string

	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Remove stored credentials",
		Long: `Remove stored credentials for a given provider or clear all saved
credentials.

When a provider is specified with --provider the command removes only
that provider's stored token. If --provider is omitted, all stored
credentials managed by picoclaw will be removed. This is useful for
rotating credentials or cleaning up local state.
`,
		Example: `  picoclaw auth logout --provider openai
	  picoclaw auth logout`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return authLogoutCmd(provider)
		},
	}

	cmd.Flags().StringVarP(&provider, "provider", "p", "", "Provider to logout from (openai, anthropic); empty = all")

	return cmd
}
