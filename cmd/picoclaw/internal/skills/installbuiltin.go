package skills

import "github.com/spf13/cobra"

func newInstallBuiltinCommand(workspaceFn func() (string, error)) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "install-builtin",
		Short:   "Install all builtin skills to workspace",
		Example: `picoclaw skills install-builtin`,
		Long: `Copy all builtin skills bundled with the picoclaw binary into
the workspace so they become available to the local agent.

This is convenient after a fresh install or when restoring a
workspace. It installs the default set of builtin skills into the
workspace skills folder.
`,
		RunE: func(_ *cobra.Command, _ []string) error {
			workspace, err := workspaceFn()
			if err != nil {
				return err
			}
			skillsInstallBuiltinCmd(workspace)
			return nil
		},
	}

	return cmd
}
