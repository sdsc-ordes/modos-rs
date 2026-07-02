package main

import (
	"modos-rs/tools/quitsh/cmd/quitsh/cmd"
	"modos-rs/tools/quitsh/pkg/build"
	dpConfig "modos-rs/tools/quitsh/pkg/config"
	dpRunner "modos-rs/tools/quitsh/pkg/runner"
	"os"

	cnConfig "gitlab.com/data-custodian/custodian/tools/quitsh/pkg/config"
	"gitlab.com/data-custodian/custodian/tools/quitsh/pkg/exec/nix"
	cnRunner "gitlab.com/data-custodian/custodian/tools/quitsh/pkg/runner"
	"gitlab.com/data-custodian/custodian/tools/quitsh/pkg/stage"

	"github.com/sdsc-ordes/quitsh/pkg/cli"
	"github.com/sdsc-ordes/quitsh/pkg/common"
	"github.com/sdsc-ordes/quitsh/pkg/config"
	"github.com/sdsc-ordes/quitsh/pkg/log"
	qRunnerExec "github.com/sdsc-ordes/quitsh/pkg/runner/exec"
	"github.com/sdsc-ordes/quitsh/pkg/toolchain"
)

func main() {
	err := log.Setup("info") // Level will be set at startup.
	if err != nil {
		log.PanicE(err, "Could not setup logger.")
	}

	args := dpConfig.New()

	cli, err := cli.New(
		&args.Commands.Root,
		&args,
		cli.WithName("quitsh"),
		cli.WithVersion(build.Version()),
		cli.WithStages(stage.AllForQuitsh()...),
		cli.WithTargetToStageMapperDefault(),
		cli.WithSignalContext(true),
		cli.WithToolchainDispatcherNix(
			nix.DefaultFlakeDirRel,
			func(c config.IConfig) *toolchain.DispatchArgs {
				cc := common.Cast[*cnConfig.Config](c)

				return &cc.Commands.DispatchArgs
			},
		),
	)
	log.PanicE(err, "Could not initialize CLI app.")

	defer func() {
		e := cli.Shutdown()
		log.WarnE(e, "Could not shutdown CLI app.")
		if err != nil {
			os.Exit(1)
		}
	}()

	// Enhance the CLI with our commands and runners.
	cmd.AddCommands(cli, &args.Config)

	cnRunner.RegisterAll(
		&args.Build,
		&args.Lint,
		&args.Test,
		&args.Image,
		&args.Manifest,
		&args.Nix,
		cli.RunnerFactory())

	dpRunner.RegisterAll(
		&args.Lint,
		&args.Build,
		&args.Test,
		cli.RunnerFactory(),
	)

	err = qRunnerExec.Register(args.Build.WrapToIBuildSettings(), cli.RunnerFactory(), true)
	log.PanicE(err, "Could not register exec runner.")

	// Run the app.
	err = cli.Run()
}
