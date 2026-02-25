package auth

import "github.com/spf13/cobra"

func newStatusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show current auth status",
		Long: `Display current authentication status and available provider
credentials.

This command prints which providers have credentials configured and any
relevant metadata (e.g. expiration). Use it to confirm that a provider
is ready before running commands that require API access.
`,
		Example: `  picoclaw auth status`,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return authStatusCmd()
		},
	}

	return cmd
}
