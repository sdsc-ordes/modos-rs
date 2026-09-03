package runner

import (
	mdConfig "modos-rs/tools/quitsh/pkg/runner/config"
	mdRustRunner "modos-rs/tools/quitsh/pkg/runner/rust"

	"github.com/sdsc-ordes/quitsh/pkg/errors"
	"github.com/sdsc-ordes/quitsh/pkg/log"
	"github.com/sdsc-ordes/quitsh/pkg/runner/factory"

	containerfilerunner "github.com/sdsc-ordes/quitsh/pkg/runner/containerfile"
	coveragerunner "github.com/sdsc-ordes/quitsh/pkg/runner/coverage"
	execRunner "github.com/sdsc-ordes/quitsh/pkg/runner/exec"
	"github.com/sdsc-ordes/quitsh/pkg/runner/gitdiffrunner"
	gorunner "github.com/sdsc-ordes/quitsh/pkg/runner/go"
	nixrunner "github.com/sdsc-ordes/quitsh/pkg/runner/nix"
	symlinkrunner "github.com/sdsc-ordes/quitsh/pkg/runner/symlinks"
	trivyrunner "github.com/sdsc-ordes/quitsh/pkg/runner/trivy"
)

func RegisterAll(
	lintSettings *mdConfig.LintSettings,
	buildSettings *mdConfig.BuildSettings,
	testSettings *mdConfig.TestSettings,
	imageSettings *mdConfig.ImageSettings,
	nixSettings *mdConfig.NixSettings,
	factory factory.IFactory,
) {
	log.Trace("Register all runners.")
	var err error

	e := execRunner.Register(buildSettings.WrapToIBuildSettings(), factory, true)
	err = errors.Combine(err, e)

	e = gorunner.RegisterBuild(buildSettings.WrapToIBuildSettings(), factory, true)
	err = errors.Combine(err, e)

	e = gorunner.RegisterTest(testSettings.WrapToITestSettings(), factory, true)
	err = errors.Combine(err, e)

	e = gorunner.RegisterLint(&lintSettings.LintSettings, factory, true)
	err = errors.Combine(err, e)

	e = coveragerunner.Register(testSettings.WrapToITestSettings(), factory)
	err = errors.Combine(err, e)

	e = trivyrunner.Register(&lintSettings.LintSettings, factory)
	err = errors.Combine(err, e)

	e = symlinkrunner.Register(&lintSettings.LintSettings, factory)
	err = errors.Combine(err, e)

	e = nixrunner.Register(imageSettings, nixSettings, factory)
	err = errors.Combine(err, e)
	e = containerfilerunner.Register(imageSettings, factory)
	err = errors.Combine(err, e)

	e = gitdiffrunner.Register(factory, nixSettings)
	err = errors.Combine(err, e)

	e = mdRustRunner.Register(
		&lintSettings.LintSettings,
		buildSettings.WrapToIBuildSettings(),
		testSettings.WrapToITestSettings(),
		factory, true)
	err = errors.Combine(err, e)

	if err != nil {
		log.PanicE(err, "Could not register runners")
	}
}
