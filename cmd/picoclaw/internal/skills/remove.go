package skills

import (
	"github.com/spf13/cobra"

	"github.com/sipeed/picoclaw/pkg/skills"
)

func newRemoveCommand(installerFn func() (*skills.SkillInstaller, error)) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove",
		Aliases: []string{"rm", "uninstall"},
		Short:   "Remove installed skill",
		Args:    cobra.ExactArgs(1),
		Example: `picoclaw skills remove weather`,
		Long: `Uninstall a skill from the workspace by name.

This removes the skill files so the agent no longer loads them. Use
the exact skill name as shown by 'picoclaw skills list'. Removal is
permanent for the workspace copy; if you need to keep the files back
them up first.
`,
		RunE: func(_ *cobra.Command, args []string) error {
			installer, err := installerFn()
			if err != nil {
				return err
			}
			skillsRemoveCmd(installer, args[0])
			return nil
		},
	}

	return cmd
}
