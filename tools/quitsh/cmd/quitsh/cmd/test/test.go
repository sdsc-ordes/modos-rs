package build

import (
	"fmt"
	"modos-rs/tools/quitsh/pkg/runner/config"

	"github.com/sdsc-ordes/quitsh/pkg/cli"
	execstage "github.com/sdsc-ordes/quitsh/pkg/cli/cmd/exec-stage"
	"github.com/sdsc-ordes/quitsh/pkg/common"
	"github.com/sdsc-ordes/quitsh/pkg/component/stage"
	"github.com/sdsc-ordes/quitsh/pkg/dag"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const longDescBuild = `
Test a component matching them by name patterns (glob).
`

func AddCmd(
	cl cli.ICLI,
	testSettings *config.TestSettings,
	execArgs *dag.ExecArgs,
) {
	execstage.AddCmdAlias(cl, cl.RootCmd(), stage.Test, execArgs, execstage.WithModifications(
		func(cmd *cobra.Command) {
			cmd.PreRun = func(cmd *cobra.Command, _ []string) {
				adjustDefaultValues(cmd, testSettings)
			}

			cmd.Flags().
				VarP(&testSettings.BuildType,
					"build-type", "b",
					fmt.Sprintf("The build type (set by env. type if not set) (%v).", common.GetAllBuildTypes()),
				)

			cmd.Flags().
				VarP(&testSettings.EnvironmentType,
					"env-type", "e", fmt.Sprintf("The environment type. (%s)", common.GetEnvTypesHelp()))

			cmd.Flags().
				BoolVar(&testSettings.ShowTestLog,
					"show-test-log", testSettings.ShowTestLog,
					"Show the log of the tests.")
		},
	))

}

func adjustDefaultValues(cmd *cobra.Command, setts *config.TestSettings) {
	// Adjust the `BuildType` if its not set.
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		if flag.Name == "build-type" && !flag.Changed {
			setts.BuildType = common.NewBuildTypeFromEnv(setts.EnvironmentType)
		}
	})
}
