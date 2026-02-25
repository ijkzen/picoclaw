package skills

import "github.com/spf13/cobra"

func newListBuiltinCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list-builtin",
		Short:   "List available builtin skills",
		Example: `picoclaw skills list-builtin`,
		Long: `List skills that are bundled with the picoclaw distribution.

Builtin skills are provided by the project and can be installed into
the workspace with 'picoclaw skills install-builtin'. Use this to view
what builtin capabilities are available.
`,
		Run: func(_ *cobra.Command, _ []string) {
			skillsListBuiltinCmd()
		},
	}

	return cmd
}
