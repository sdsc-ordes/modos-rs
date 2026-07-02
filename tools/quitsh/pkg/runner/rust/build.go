package rustrunner

import (
	"github.com/sdsc-ordes/quitsh/pkg/common"
	"github.com/sdsc-ordes/quitsh/pkg/debug"
	"github.com/sdsc-ordes/quitsh/pkg/exec"
	fs "github.com/sdsc-ordes/quitsh/pkg/filesystem"
	"github.com/sdsc-ordes/quitsh/pkg/runner"
	"github.com/sdsc-ordes/quitsh/pkg/runner/config"
)

const RustBuildRunnerID = "quitsh::build-rust"

type RustBuildRunner struct {
	config   *RunnerConfigBuild
	settings config.IBuildSettings
}

// NewRustBuildRunner constructs a new GoBuildRunner with its own config.
func NewRustBuildRunner(config any, settings config.IBuildSettings) (runner.IRunner, error) {
	debug.Assert(config != nil, "config is nil")

	return &RustBuildRunner{
		config:   common.Cast[*RunnerConfigBuild](config),
		settings: settings,
	}, nil
}

func (*RustBuildRunner) ID() runner.RegisterID {
	return RustBuildRunnerID
}

// Run implements [runner.IRunner].
func (r *RustBuildRunner) Run(ctx runner.IContext) error {
	log := ctx.Log()
	comp := ctx.Component()

	log.Info("Starting Rust build for component.", "component", comp.Name())

	fs.AssertDirs(comp.OutBuildBinDir())

	// Set the output path and disable GOWORK:
	// We build in each component without looking
	// at `go.work` fs.
	targetDir := comp.OutBuildDir("target")
	var binDir string
	if r.settings.Coverage() {
		binDir = comp.OutCoverageBinDir()
	} else {
		binDir = comp.OutBuildBinDir()
	}

	// Build everything into `outputDir`.
	args, envs := GetCargoArgs(
		log,
		targetDir,
		binDir,
		r.settings.BuildType(),
		r.settings.EnvironmentType(),
		r.config.Binaries,
		r.config.Libraries,
		r.config.Examples,
		r.config.Tests,
		r.config.Features,
		false,
		false,
		comp.Version(),
		BuildCommand,
	)

	cargoCtx := exec.NewCmdCtxBuilder().
		BaseCmd("cargo").
		Cwd(comp.Root()).
		Env(envs...).
		Build()

	log.Info("Run cargo build.")

	cmd := append([]string{"build"}, args...)
	cmd = append(cmd, r.settings.Args()...)
	err := cargoCtx.Check(cmd...)

	if err != nil {
		log.ErrorE(err, "Cargo build failed.")

		return err
	}

	return nil
}
