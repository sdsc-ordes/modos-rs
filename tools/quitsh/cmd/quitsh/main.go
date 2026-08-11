package main

import (
	"modos-rs/tools/quitsh/cmd/quitsh/cmd"
	"modos-rs/tools/quitsh/pkg/build"
	mdConfig "modos-rs/tools/quitsh/pkg/config"
	"modos-rs/tools/quitsh/pkg/nix"
	mdRunner "modos-rs/tools/quitsh/pkg/runner"
	"os"

	"github.com/sdsc-ordes/quitsh/pkg/cli"
	"github.com/sdsc-ordes/quitsh/pkg/common"
	"github.com/sdsc-ordes/quitsh/pkg/component/stage"
	"github.com/sdsc-ordes/quitsh/pkg/config"
	"github.com/sdsc-ordes/quitsh/pkg/log"
	"github.com/sdsc-ordes/quitsh/pkg/toolchain"
)

func main() {
	err := log.Setup("info") // Level will be set at startup.
	if err != nil {
		log.PanicE(err, "Could not setup logger.")
	}

	conf := mdConfig.New()

	cli, err := cli.New(
		&conf.Commands.Root,
		&conf,
		cli.WithName("quitsh"),
		cli.WithVersion(build.Version()),
		cli.WithStages(stage.AllStages()...),
		cli.WithTargetToStageMapperDefault(),
		cli.WithSignalContext(true),
		cli.WithToolchainDispatcherNix(
			nix.DefaultFlakeDirRel,
			func(c config.IConfig) *toolchain.DispatchArgs {
				cc := common.Cast[*mdConfig.Config](c)

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
	cmd.AddCommands(cli, &conf)

	mdRunner.RegisterAll(
		&conf.Lint,
		&conf.Build,
		&conf.Test,
		&conf.Image,
		&conf.Nix,
		cli.RunnerFactory(),
	)

	// Run the app.
	err = cli.Run()
}
