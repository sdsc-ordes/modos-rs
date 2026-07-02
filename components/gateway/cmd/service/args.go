package main

import (
	"fmt"
	"os"

	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/log"
	"gitlab.com/data-custodian/dac-portal/components/nats-server/internal/build"

	"github.com/alexflint/go-arg"
	"github.com/creasty/defaults"
)

type Args struct {
	ConfigDir string `arg:"-c,--config-dir" help:"The config directory path." default:"config"`
	DataDir   string `arg:"-d,--data-dir"   help:"The data directory path." default:"data"`
	Version   bool   `arg:"-v,--version"    help:"The version of the service"`

	LogToFile bool `arg:"--log-to-file" help:"If we should log to files too."`
}

func parseArgs() (args Args) {
	err := defaults.Set(&args)
	log.FatalE(err, "Could not set defaults on arguments.")

	//nolint:exhaustruct,nolintlint
	p, err := arg.NewParser(arg.Config{
		Out: os.Stderr,
	}, &args)
	log.FatalE(err, "Could not setup parser.")

	err = p.Parse(os.Args[1:])
	log.FatalE(err, "Could not parse arguments.")

	if args.Version {
		printVersionAndExit()
	}

	log.Debug("Arguments", "arguments", args)

	return args
}

func printVersionAndExit() {
	fmt.Printf("version: %v\n", build.GetBuildVersion().String()) //nolint:forbidigo,nolintlint
	os.Exit(0)
}
