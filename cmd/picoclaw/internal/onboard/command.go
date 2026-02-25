package onboard

import (
	"embed"

	"github.com/spf13/cobra"
)

//go:generate cp -r ../../../../workspace .
//go:embed workspace
var embeddedFiles embed.FS

func NewOnboardCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "onboard",
		Aliases: []string{"o"},
		Short:   "Initialize picoclaw configuration and workspace",
		Long: `Initialize a default picoclaw configuration and workspace in the
current user's home directory.

This command copies embedded workspace templates into ~/.picoclaw and
creates default files needed to run picoclaw. Run this once when
setting up picoclaw for the first time.
`,
		Example: `  picoclaw onboard`,
		Run: func(cmd *cobra.Command, args []string) {
			onboard()
		},
	}

	return cmd
}
