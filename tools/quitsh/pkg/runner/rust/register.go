package rustrunner

import (
	"github.com/sdsc-ordes/quitsh/pkg/component/step"
	"github.com/sdsc-ordes/quitsh/pkg/errors"
	"github.com/sdsc-ordes/quitsh/pkg/log"
	"github.com/sdsc-ordes/quitsh/pkg/runner"
	"github.com/sdsc-ordes/quitsh/pkg/runner/config"
	"github.com/sdsc-ordes/quitsh/pkg/runner/factory"

	cnConfig "gitlab.com/data-custodian/custodian/tools/quitsh/pkg/runner/config"
)

// Register registers the build/build-gen/test runner in the factory.
func Register(
	lintSettings *cnConfig.LintSettings,
	buildSettings config.IBuildSettings,
	testSettings config.ITestSettings,
	factory factory.IFactory,
	registerKey bool,
) (err error) {
	const defaultToolchain = "build-rust"

	log.Trace("Register runner.", "id", RustBuildRunnerID)
	e := factory.Register(
		RustBuildRunnerID,
		runner.RunnerData{
			Creator: func(config step.AuxConfig) (runner.IRunner, error) {
				return NewRustBuildRunner(config, buildSettings)
			},
			RunnerConfigUnmarshal: UnmarshalBuildConfig,
			DefaultToolchain:      defaultToolchain,
		})
	err = errors.Combine(err, e)

	if registerKey {
		e = factory.RegisterToKey(runner.NewRegisterKey("build", "rust"), RustBuildRunnerID)
		err = errors.Combine(err, e)
	}

	log.Trace("Register runner.", "id", RustBuildGenRunnerID)
	e = factory.Register(
		RustBuildGenRunnerID,
		runner.RunnerData{
			Creator: func(config step.AuxConfig) (runner.IRunner, error) {
				return NewRustBuildGenRunner(config, buildSettings)
			},
			RunnerConfigUnmarshal: UnmarshalBuildGenConfig,
			DefaultToolchain:      defaultToolchain,
		})
	err = errors.Combine(err, e)

	if registerKey {
		e = factory.RegisterToKey(runner.NewRegisterKey("build", "rust-gen"), RustBuildGenRunnerID)
		err = errors.Combine(err, e)
	}

	log.Trace("Register runner.", "id", RustTestRunnerID)
	e = factory.Register(
		RustTestRunnerID,
		runner.RunnerData{
			Creator: func(config step.AuxConfig) (runner.IRunner, error) {
				return NewRustTestRunner(config, testSettings)
			},
			RunnerConfigUnmarshal: UnmarshalTestConfig,
			DefaultToolchain:      defaultToolchain,
		})
	err = errors.Combine(err, e)

	if registerKey {
		e = factory.RegisterToKey(runner.NewRegisterKey("test", "rust"), RustTestRunnerID)
		err = errors.Combine(err, e)
	}

	log.Trace("Register runner.", "id", RustLintRunnerID)
	e = factory.Register(
		RustLintRunnerID,
		runner.RunnerData{
			Creator: func(config step.AuxConfig) (runner.IRunner, error) {
				return NewRustLintRunner(config, lintSettings)
			},
			RunnerConfigUnmarshal: UnmarshalLintConfig,
			DefaultToolchain:      defaultToolchain,
		})
	err = errors.Combine(err, e)

	if registerKey {
		e = factory.RegisterToKey(runner.NewRegisterKey("lint", "rust"), RustLintRunnerID)
		err = errors.Combine(err, e)
	}

	return err
}
