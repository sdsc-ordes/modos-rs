//go:build test && integration

package test

import (
	"os"
	"path"
	"testing"

	"github.com/sdsc-ordes/quitsh/pkg/exec"
	"github.com/sdsc-ordes/quitsh/pkg/exec/git"
	fs "github.com/sdsc-ordes/quitsh/pkg/filesystem"

	"github.com/stretchr/testify/require"
)

func setup(t *testing.T) (ciTool *exec.CmdContext) {

	outputDir := os.Getenv("QUITSH_BIN_DIR")
	require.True(t, fs.Exists(outputDir))

	covDir := os.Getenv("QUITSH_COVERAGE_DIR")
	require.True(t, fs.Exists(covDir))

	ciToolExe := path.Join(outputDir, "quitsh")
	require.FileExists(t, ciToolExe)

	ciTool = exec.NewCmdCtxBuilder().
		BaseCmd(ciToolExe).
		Cwd(".").
		EnableCaptureError().
		Env(os.Environ()...).
		Env("GOCOVERDIR=" + covDir).
		Build()

	return
}

// Demo test since no useful test currently.
func TestCLIVersion(t *testing.T) {
	ciTool := setup(t)

	_, repoRoot, err := git.NewCtxAtRoot(".")
	require.NoError(t, err)

	_, err = ciTool.GetCombined(
		"-C", repoRoot,
		"list",
	)
	require.NoError(t, err)
}
