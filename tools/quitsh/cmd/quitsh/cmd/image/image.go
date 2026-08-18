package build

import (
	"modos-rs/tools/quitsh/pkg/runner/config"

	"github.com/sdsc-ordes/quitsh/pkg/cli"
	execstage "github.com/sdsc-ordes/quitsh/pkg/cli/cmd/exec-stage"
	imageCmd "github.com/sdsc-ordes/quitsh/pkg/cli/cmd/image"
	"github.com/sdsc-ordes/quitsh/pkg/component/stage"
	"github.com/sdsc-ordes/quitsh/pkg/dag"

	"github.com/spf13/cobra"
)

const longDescBuild = `
Build a the image of a component matching them by name patterns (glob).
`

func AddCmd(
	cl cli.ICLI,
	imageSettings *config.ImageSettings,
	execArgs *dag.ExecArgs,
) {
	execstage.AddCmdAlias(cl, cl.RootCmd(), stage.Image, execArgs, execstage.WithModifications(
		func(cmd *cobra.Command) {
			imageCmd.SetImageFlags(cmd, imageSettings)
		},
	))
}
