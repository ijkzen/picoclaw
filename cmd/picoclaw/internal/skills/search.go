package skills

import (
	"github.com/spf13/cobra"

	"github.com/sipeed/picoclaw/pkg/skills"
)

func newSearchCommand(installerFn func() (*skills.SkillInstaller, error)) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search available skills",
		Long: `Search remote registries or indexes for available skills.

This command assists discovering third-party skills you can install.
Depending on configuration it may query public registries or a
centralized index. After finding a skill, use 'picoclaw skills
install' to add it to your workspace.
`,
		RunE: func(_ *cobra.Command, _ []string) error {
			installer, err := installerFn()
			if err != nil {
				return err
			}
			skillsSearchCmd(installer)
			return nil
		},
	}

	return cmd
}
