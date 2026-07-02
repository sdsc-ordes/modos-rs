package runner

import (
	rustrunner "modos-rs/tools/quitsh/pkg/runner/rust"

	"github.com/sdsc-ordes/quitsh/pkg/errors"
	"github.com/sdsc-ordes/quitsh/pkg/log"
	"github.com/sdsc-ordes/quitsh/pkg/runner/factory"
	"gitlab.com/data-custodian/custodian/tools/quitsh/pkg/runner/config"
)

func RegisterAll(
	lintSettings *config.LintSettings,
	buildSettings *config.BuildSettings,
	testSettings *config.TestSettings,
	factory factory.IFactory,
) {
	log.Trace("Register all dac-portal runners.")
	var err error

	e := rustrunner.Register(
		lintSettings,
		buildSettings.WrapToIBuildSettings(),
		testSettings.WrapToITestSettings(),
		factory, true)
	err = errors.Combine(err, e)

	if err != nil {
		log.PanicE(err, "Could not register runners.")
	}
}
