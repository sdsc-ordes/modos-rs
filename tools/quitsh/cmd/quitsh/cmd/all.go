package cmd

import (
	"modos-rs/tools/quitsh/cmd/quitsh/cmd/setup"

	cnBuildCmd "gitlab.com/data-custodian/custodian/tools/quitsh/cmd/quitsh/cmd/build"
	cnFormatCmd "gitlab.com/data-custodian/custodian/tools/quitsh/cmd/quitsh/cmd/format"
	cnImageCmd "gitlab.com/data-custodian/custodian/tools/quitsh/cmd/quitsh/cmd/image"
	cnLintCmd "gitlab.com/data-custodian/custodian/tools/quitsh/cmd/quitsh/cmd/lint"
	cnManifestCmd "gitlab.com/data-custodian/custodian/tools/quitsh/cmd/quitsh/cmd/manifest"
	cnNix "gitlab.com/data-custodian/custodian/tools/quitsh/cmd/quitsh/cmd/nix"
	cnTestCmd "gitlab.com/data-custodian/custodian/tools/quitsh/cmd/quitsh/cmd/test"
	"gitlab.com/data-custodian/custodian/tools/quitsh/pkg/config"

	"github.com/sdsc-ordes/quitsh/pkg/cli"
	cleanCmd "github.com/sdsc-ordes/quitsh/pkg/cli/cmd/clean"
	configCmd "github.com/sdsc-ordes/quitsh/pkg/cli/cmd/config"
	execRunnerCmd "github.com/sdsc-ordes/quitsh/pkg/cli/cmd/exec-runner"
	execTargetCmd "github.com/sdsc-ordes/quitsh/pkg/cli/cmd/exec-target"
	listCmd "github.com/sdsc-ordes/quitsh/pkg/cli/cmd/list"
	proccompCmd "github.com/sdsc-ordes/quitsh/pkg/cli/cmd/process-compose"
	versionupcmd "github.com/sdsc-ordes/quitsh/pkg/cli/cmd/version-up"
)

func AddCommands(cl cli.ICLI, conf *config.Config) {
	// Quitsh commands.
	configCmd.AddCmd(cl.RootCmd(), conf)
	execRunnerCmd.AddCmd(cl, cl.RootCmd(), &conf.Commands.DispatchArgs)
	execTargetCmd.AddCmd(cl, cl.RootCmd(), &conf.Commands.ExecArgs)
	listCmd.AddCmd(cl, cl.RootCmd())
	cleanCmd.AddCmd(cl)
	proccompCmd.AddCmd(cl, cl.RootCmd(), conf.Nix.FlakeDirRel)
	versionupcmd.AddCmd(cl, cl.RootCmd())

	// modos commands.
	cnBuildCmd.AddCmd(cl, &conf.Build, &conf.Commands.ExecArgs)
	cnLintCmd.AddCmd(cl, &conf.Lint, &conf.Commands.ExecArgs)
	cnTestCmd.AddCmd(cl, &conf.Test, &conf.Commands.ExecArgs)
	cnImageCmd.AddCmd(cl, &conf.Image, &conf.Commands.ExecArgs)
	cnManifestCmd.AddCmd(cl, &conf.Manifest, &conf.Image, &conf.Commands.ExecArgs)
	cnNix.AddCmd(cl, &conf.Nix)
	cnFormatCmd.AddCmd(cl.RootCmd(), &conf.Nix)

	// Own commands.
	setup.AddCmd(cl.RootCmd(), &conf.Nix)
}
