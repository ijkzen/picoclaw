package models

import (
	"github.com/spf13/cobra"
)

func NewModelsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "models",
		Short: "Manage model_list configuration (list/add/delete/edit)",
	}

	cmd.AddCommand(
		NewListCommand(),
		NewAddCommand(),
		NewDeleteCommand(),
		NewEditCommand(),
	)

	return cmd
}
