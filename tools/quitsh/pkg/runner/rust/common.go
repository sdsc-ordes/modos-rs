package rustrunner

import (
	"strings"

	cm "github.com/sdsc-ordes/quitsh/pkg/common"
	"github.com/sdsc-ordes/quitsh/pkg/log"

	"github.com/hashicorp/go-version"
)

const RunCommand CommandType = 0
const TestCommand CommandType = 1
const BuildCommand CommandType = 2
const LintCommand CommandType = 3

type CommandType int

//nolint:gocognit,funlen // Its fine.
func GetCargoArgs(
	log log.ILog,
	targetDir string,
	binDir string,
	buildType cm.BuildType,
	envType cm.EnvironmentType,
	binaries []string,
	libraries bool,
	examples bool,
	tests []string,
	features []string,
	coverage bool,
	verbose bool,
	version *version.Version,
	cmdType CommandType,
) (args []string, envs []string) {
	args = []string{
		"--locked",
		"--color=auto",
		"--features", envType.String(),
	}

	if verbose {
		args = append(args, "--verbose")
	}

	// Output dir.
	switch cmdType { //nolint:exhaustive // Ok.
	case TestCommand:
		break
	case BuildCommand:
		args = append(args, "-Z", "unstable-options", "--artifact-dir", binDir)

		fallthrough
	default:
		args = append(args, "--target-dir", targetDir)
	}

	// Features.
	if len(features) != 0 {
		args = append(args, "--features", strings.Join(features, ","))
	}

	// Profile.
	switch buildType { //nolint:exhaustive // This is correct.
	case cm.BuildRelease:
		args = append(args, "--release")
	default:
		log.Warn("Building/testing in debug mode.")
		args = append(args, "--profile", "dev")
	}

	// Adding tests targets.
	switch {
	case cmdType == RunCommand || cmdType == BuildCommand:
		// No test for run or build.
		break
	case tests == nil:
		fallthrough
	case cmdType == TestCommand:
		args = append(args, "--tests")
	default:
		for _, t := range tests {
			args = append(args, "--test", t)
		}
	}

	// Adding binary targets.
	switch {
	case cmdType == RunCommand:
		// Only run the first binary.
		if len(binaries) > 1 {
			log.Warn("Cargo can run only one binary. Specified '%v'.", binaries)
			binaries = binaries[:1]
		}

		fallthrough
	case binaries != nil:
		for _, b := range binaries {
			args = append(args, "--bin", b)
		}
	default:
		args = append(args, "--bins")
	}

	// Add library targets.
	switch cmdType { //nolint:exhaustive // Ok.
	case RunCommand:
		break
	default:
		if libraries {
			args = append(args, "--lib")
		}
	}

	// Add example targets.
	switch cmdType { //nolint:exhaustive // Ok.
	case RunCommand:
		break
	default:
		if examples {
			args = append(args, "--examples")
		}
	}

	envs = []string{
		"CARGO_TARGET_DIR=" + targetDir, // Just in case.
		"QUITSH_COMPONENT_VERSION=" + version.String(),
	}

	if coverage && cmdType == TestCommand {
		envs = append(envs, "RUSTFLAGS=-Cinstrument-coverage")
	}

	log.Info("Cargo args:", "args", args)

	return args, envs
}
