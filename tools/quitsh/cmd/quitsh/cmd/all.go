package cmd

import (
	modosConfig "modos-rs/tools/quitsh/pkg/config"

	buildCmd "modos-rs/tools/quitsh/cmd/quitsh/cmd/build"
	ciCmd "modos-rs/tools/quitsh/cmd/quitsh/cmd/ci"
	imageCmd "modos-rs/tools/quitsh/cmd/quitsh/cmd/image"
	lintCmd "modos-rs/tools/quitsh/cmd/quitsh/cmd/lint"
	setupCmd "modos-rs/tools/quitsh/cmd/quitsh/cmd/setup"
	testCmd "modos-rs/tools/quitsh/cmd/quitsh/cmd/test"

	"github.com/sdsc-ordes/quitsh/pkg/cli"
	cleanCmd "github.com/sdsc-ordes/quitsh/pkg/cli/cmd/clean"
	configCmd "github.com/sdsc-ordes/quitsh/pkg/cli/cmd/config"
	execRunnerCmd "github.com/sdsc-ordes/quitsh/pkg/cli/cmd/exec-runner"
	execTargetCmd "github.com/sdsc-ordes/quitsh/pkg/cli/cmd/exec-target"
	formatCmd "github.com/sdsc-ordes/quitsh/pkg/cli/cmd/format"
	listCmd "github.com/sdsc-ordes/quitsh/pkg/cli/cmd/list"
	nixCmd "github.com/sdsc-ordes/quitsh/pkg/cli/cmd/nix"
	proccompCmd "github.com/sdsc-ordes/quitsh/pkg/cli/cmd/process-compose"
	versionupCmd "github.com/sdsc-ordes/quitsh/pkg/cli/cmd/version-up"
)

func AddCommands(cl cli.ICLI, conf *modosConfig.Config) {
	// Quitsh commands.
	configCmd.AddCmd(cl.RootCmd(), conf)
	execRunnerCmd.AddCmd(cl, cl.RootCmd(), &conf.Commands.DispatchArgs)
	execTargetCmd.AddCmd(cl, cl.RootCmd(), &conf.Commands.ExecArgs)
	listCmd.AddCmd(cl, cl.RootCmd())
	cleanCmd.AddCmd(cl)
	proccompCmd.AddCmd(cl, cl.RootCmd(), conf.Nix.FlakeDirRel)
	versionupCmd.AddCmd(cl, cl.RootCmd())
	nixCmd.AddCmd(cl, cl.RootCmd(), &conf.Nix)
	formatCmd.AddCmd(cl.RootCmd(), &conf.Nix)

	// modos commands.
	buildCmd.AddCmd(cl, &conf.Build, &conf.Commands.ExecArgs)
	lintCmd.AddCmd(cl, &conf.Lint, &conf.Commands.ExecArgs)
	testCmd.AddCmd(cl, &conf.Test, &conf.Commands.ExecArgs)
	imageCmd.AddCmd(cl, &conf.Image, &conf.Commands.ExecArgs)

	// Own commands.
	setupCmd.AddCmd(cl.RootCmd(), &conf.Nix)
	ciCmd.AddCmd(cl, cl.RootCmd())
}
