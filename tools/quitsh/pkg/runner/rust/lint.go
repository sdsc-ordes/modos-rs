package rustrunner

import (
	"github.com/sdsc-ordes/quitsh/pkg/common"
	"github.com/sdsc-ordes/quitsh/pkg/debug"
	"github.com/sdsc-ordes/quitsh/pkg/exec"
	fs "github.com/sdsc-ordes/quitsh/pkg/filesystem"
	"github.com/sdsc-ordes/quitsh/pkg/runner"
	cnConfig "gitlab.com/data-custodian/custodian/tools/quitsh/pkg/runner/config"
)

const RustLintRunnerID = "quitsh::lint-rust"

type RustLintRunner struct {
	config   *RunnerConfigLint
	settings *cnConfig.LintSettings
}

// NewRustLintRunner constructs a new GoBuildRunner with its own config.
func NewRustLintRunner(config any, settings *cnConfig.LintSettings) (runner.IRunner, error) {
	debug.Assert(config != nil, "config is nil")

	return &RustLintRunner{
		config:   common.Cast[*RunnerConfigLint](config),
		settings: settings,
	}, nil
}

func (*RustLintRunner) ID() runner.RegisterID {
	return RustLintRunnerID
}

// Run implements [runner.IRunner].
func (r *RustLintRunner) Run(ctx runner.IContext) error {
	log := ctx.Log()
	comp := ctx.Component()

	log.Info("Starting Rust lint for component.", "component", comp.Name())

	fs.AssertDirs(comp.OutBuildBinDir())

	if len(r.config.PerEnvs) == 0 {
		r.config.PerEnvs = []common.EnvironmentType{common.EnvironmentDev}
	}

	// Set the output path and disable GOWORK:
	// We build in each component without looking
	// at `go.work` fs.
	targetDir := comp.OutBuildDir("target")
	binDir := comp.OutBuildBinDir()

	for _, env := range r.config.PerEnvs {
		args, envs := GetCargoArgs(
			log,
			targetDir,
			binDir,
			common.BuildRelease,
			env,
			r.config.Binaries,
			r.config.Libraries,
			r.config.Examples,
			r.config.Tests,
			r.config.Features,
			false,
			false,
			comp.Version(),
			LintCommand,
		)

		cargoCtx := exec.NewCmdCtxBuilder().
			BaseCmd("cargo").
			BaseArgs("clippy").
			BaseArgs(args...).
			Cwd(comp.Root()).
			Env(envs...).
			Build()

		log.Info("Run cargo clippy.")

		cmd := []string{}
		if r.settings.Fix {
			cmd = append(cmd, "--fix", "--allow-dirty")
		}
		err := cargoCtx.Check(cmd...)
		if err != nil {
			log.ErrorEf(err, "Cargo clippy failed for '%v'", env.String())

			return err
		}

		log.Infof("Cargo clippy successful for '%v'.", env.String())
	}

	return nil
}
