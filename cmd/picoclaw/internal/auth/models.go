package auth

import "github.com/spf13/cobra"

func newModelsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "models",
		Short: "Show available models",
		Long: `List configured or available LLM models known to picoclaw.

This command queries the configured providers and prints a list of
models that can be used with the --model flag on agent commands. Use
this to discover model identifiers and ensure your configuration maps
model names to providers correctly.
`,
		Example: `  picoclaw auth models`,
		RunE: func(_ *cobra.Command, _ []string) error {
			return authModelsCmd()
		},
	}

	return cmd
}
