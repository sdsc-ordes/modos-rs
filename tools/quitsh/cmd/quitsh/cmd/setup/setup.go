package setup

import (
	"modos-rs/tools/quitsh/pkg/setup"

	"gitlab.com/data-custodian/custodian/tools/quitsh/pkg/runner/config"

	"github.com/spf13/cobra"
)

func AddCmd(root *cobra.Command, nixSett *config.NixSettings) {
	setupCmd := &cobra.Command{
		Use:     "setup-development",
		Aliases: []string{"setup"},
		Short:   "Setup local development.",
		Long:    "Setup the repository for local development.",
		PreRunE: func(_cmd *cobra.Command, _args []string) error {
			return nil
		},
		RunE: func(_cmd *cobra.Command, _args []string) error {
			return setup.Setup(nixSett.FlakeDirRel)
		},
	}

	root.AddCommand(setupCmd)
}
