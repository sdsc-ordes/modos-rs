package oauth

import (
	devicecodeflow "modos-rs/tools/quitsh/cmd/quitsh/cmd/oauth/device-code-flow"

	"github.com/sdsc-ordes/quitsh/pkg/cli"

	"github.com/spf13/cobra"
)

func AddCmd(cl cli.ICLI, parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "oauth",
		Short: "OAuth functionality.",
	}

	devicecodeflow.AddCmd(cl, cmd)
	parent.AddCommand(cmd)
}
