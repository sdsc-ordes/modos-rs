package lint

import (
	"modos-rs/tools/quitsh/pkg/runner/config"

	"github.com/sdsc-ordes/quitsh/pkg/cli"
	execstage "github.com/sdsc-ordes/quitsh/pkg/cli/cmd/exec-stage"
	"github.com/sdsc-ordes/quitsh/pkg/component/stage"
	"github.com/sdsc-ordes/quitsh/pkg/dag"

	"github.com/spf13/cobra"
)

const longDescBuild = `
Lint a component matching them by name patterns (glob).
`

func AddCmd(
	cl cli.ICLI,
	lintSettings *config.LintSettings,
	execArgs *dag.ExecArgs,
) {
	execstage.AddCmdAlias(cl, cl.RootCmd(), stage.Lint, execArgs, execstage.WithModifications(
		func(cmd *cobra.Command) {
			cmd.Flags().
				BoolVar(&lintSettings.Fix, "fix", lintSettings.Fix, "Try to fix all linting issues automatically.")
		},
	))

}
