package rustrunner

import (
	"github.com/sdsc-ordes/quitsh/pkg/ci"
	cm "github.com/sdsc-ordes/quitsh/pkg/common"
	"github.com/sdsc-ordes/quitsh/pkg/component"
	"github.com/sdsc-ordes/quitsh/pkg/debug"
	"github.com/sdsc-ordes/quitsh/pkg/exec"
	fs "github.com/sdsc-ordes/quitsh/pkg/filesystem"
	"github.com/sdsc-ordes/quitsh/pkg/log"
	"github.com/sdsc-ordes/quitsh/pkg/runner"
	"github.com/sdsc-ordes/quitsh/pkg/runner/config"
)

const RustTestRunnerID = "quitsh::test-rust"

type RustTestRunner struct {
	config   *RunnerTestConfig
	settings config.ITestSettings
}

func NewRustTestRunner(config any, settings config.ITestSettings) (runner.IRunner, error) {
	debug.Assert(config != nil, "config is nil")

	return &RustTestRunner{
		config:   cm.Cast[*RunnerTestConfig](config),
		settings: settings,
	}, nil
}

func (r *RustTestRunner) ID() runner.RegisterID {
	return RustTestRunnerID
}

func generateCoverageReport(
	log log.ILog,
	cargoCtx *exec.CmdContext,
	comp *component.Component) error {
	covLcov := comp.OutCoverageDataDir("coverage.lcov")
	covCodecov := comp.OutCoverageDataDir("coverage-codecov.json")

	err := cargoCtx.Check("report", "--lcov", "--output-path", covLcov)
	if err != nil {
		log.ErrorE(err, "Rust coverage 'lcov' report failed.")
	}

	err = cargoCtx.Check("report", "--codecov", "--output-path", covCodecov)
	if err != nil {
		log.ErrorE(err, "Rust coverage 'codecov' conversion failed.")
	}

	htmlCmd := []string{"report", "--html", "--output-dir", comp.OutCoverageDataDir()}
	if !ci.IsRunning() {
		htmlCmd = append(htmlCmd, "--open")
	}
	err = cargoCtx.Check(htmlCmd...)
	if err != nil {
		log.ErrorE(err, "Rust coverage html conversion failed.")
	}

	return err
}

func (r *RustTestRunner) Run(ctx runner.IContext) error {
	log := ctx.Log()
	comp := ctx.Component()

	log.Info("Starting Rust test for component.", "component", comp.Name())

	fs.AssertDirs(comp.OutBuildBinDir())

	// Set the output path and disable GOWORK:
	// We build in each component without looking
	// at `go.work` fs.
	covDataDir := comp.OutCoverageDataDir()
	fs.AssertDirs(comp.OutBuildBinDir(), covDataDir)

	// Build everything into `outputDir`.
	args, envs := GetCargoArgs(
		log,
		covDataDir,
		"",
		r.settings.BuildType(),
		cm.EnvironmentDev, // Not used.
		r.config.Binaries,
		r.config.Libraries,
		r.config.Examples,
		r.config.Tests,
		r.config.Features,
		true,
		r.settings.ShowTestLog(),
		comp.Version(),
		TestCommand,
	)

	cargoCtx := exec.NewCmdCtxBuilder().
		BaseCmd("cargo").
		BaseArgs("llvm-cov").
		Cwd(comp.Root()).
		Env(envs...).
		Build()

	log.Info("Run cargo build.")

	cmd := []string{"test", "--no-report", "--output-dir", covDataDir}
	cmd = append(cmd, args...)
	cmd = append(cmd, r.settings.Args()...)
	cmd = append(cmd, r.config.Args...)
	if len(r.config.TestArgs) != 0 {
		cmd = append(cmd, r.config.TestArgs...)
	}
	err := cargoCtx.Check(cmd...)

	if err != nil {
		log.ErrorE(err, "Cargo build failed.")

		return err
	}

	err = generateCoverageReport(log, cargoCtx, comp)
	if err != nil {
		return err
	}

	return nil
}
