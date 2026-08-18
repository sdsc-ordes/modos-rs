// Package ci groups CI-related helper commands.
package ci

import (
	updateruleset "modos-rs/tools/quitsh/cmd/quitsh/cmd/ci/update-main-protection-ruleset"

	"github.com/sdsc-ordes/quitsh/pkg/cli"

	"github.com/spf13/cobra"
)

func AddCmd(cl cli.ICLI, parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "ci",
		Short: "CI helper commands.",
	}

	updateruleset.AddCmd(cl, cmd)
	parent.AddCommand(cmd)
}
