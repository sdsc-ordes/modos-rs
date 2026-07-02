package rustrunner

import (
	"github.com/sdsc-ordes/quitsh/pkg/common"
	"github.com/sdsc-ordes/quitsh/pkg/debug"
	"github.com/sdsc-ordes/quitsh/pkg/exec"
	fs "github.com/sdsc-ordes/quitsh/pkg/filesystem"
	"github.com/sdsc-ordes/quitsh/pkg/runner"
	"github.com/sdsc-ordes/quitsh/pkg/runner/config"
)

const RustBuildGenRunnerID = "quitsh::build-rust-gen"

type RustBuildGenRunner struct {
	config   *RunnerConfigBuildGen
	settings config.IBuildSettings
}

// NewRustBuildGenRunner constructs a new Rust build runner to generate files with its own config.
func NewRustBuildGenRunner(config any, settings config.IBuildSettings) (runner.IRunner, error) {
	debug.Assert(config != nil, "config is nil")

	return &RustBuildGenRunner{
		config:   common.Cast[*RunnerConfigBuildGen](config),
		settings: settings,
	}, nil
}

func (*RustBuildGenRunner) ID() runner.RegisterID {
	return RustBuildGenRunnerID
}

// Run implements [runner.IRunner].
func (r *RustBuildGenRunner) Run(ctx runner.IContext) error {
	log := ctx.Log()
	comp := ctx.Component()

	log.Info("Starting Rust build-gen for component.", "component", comp.Name())

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
		nil,
		false,
		false,
		nil,
		r.config.Features,
		false,
		false,
		comp.Version(),
		RunCommand,
	)

	cargoCtx := exec.NewCmdCtxBuilder().
		BaseCmd("cargo").
		Cwd(comp.Root()).
		Env(envs...).
		Build()

	for _, bin := range r.config.Binaries {
		log.Info("Run cargo run.")

		cmd := append([]string{"run"}, args...)
		cmd = append(cmd, r.settings.Args()...)
		cmd = append(cmd, "--bin", bin)
		err := cargoCtx.Check(cmd...)

		if err != nil {
			log.ErrorEf(err, "Cargo run failed for binary '%s'.", bin)

			return err
		}
	}

	return nil
}
