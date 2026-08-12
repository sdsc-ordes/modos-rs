// Package updateruleset implements the `ci update-main-protection-ruleset`
// command which (re)creates the `main-protection` branch ruleset on the
// repository, requiring all status checks observed on a reference PR.
package updateruleset

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/sdsc-ordes/quitsh/pkg/cli"
	"github.com/sdsc-ordes/quitsh/pkg/errors"
	"github.com/sdsc-ordes/quitsh/pkg/exec"
	"github.com/sdsc-ordes/quitsh/pkg/log"

	"github.com/spf13/cobra"
)

// Settings holds the flags for the command.
type Settings struct {
	Repo        string // owner/name of the repository.
	RulesetName string // name of the ruleset to create-or-update.
}

// prCheck is one entry of `gh pr checks --json name,state`.
type prCheck struct {
	Name  string `json:"name"`
	State string `json:"state"`
}

// rulesetRef is the subset of `gh api repos/<repo>/rulesets` we need.
type rulesetRef struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// The payload types below mirror the GitHub "repository ruleset" schema.
//
//nolint:tagliatelle // GitHub REST API uses snake_case.
type (
	statusCheck struct {
		Context       string `json:"context"`
		IntegrationID int64  `json:"integration_id,omitempty"`
	}

	refNameCondition struct {
		Include []string `json:"include"`
		Exclude []string `json:"exclude"`
	}

	conditions struct {
		RefName refNameCondition `json:"ref_name"`
	}

	pullRequestParams struct {
		RequiredApprovingReviewCount   int  `json:"required_approving_review_count"`
		DismissStaleReviewsOnPush      bool `json:"dismiss_stale_reviews_on_push"`
		RequireCodeOwnerReview         bool `json:"require_code_owner_review"`
		RequireLastPushApproval        bool `json:"require_last_push_approval"`
		RequiredReviewThreadResolution bool `json:"required_review_thread_resolution"`
	}

	statusCheckParams struct {
		StrictRequiredStatusChecksPolicy bool          `json:"strict_required_status_checks_policy"`
		RequiredStatusChecks             []statusCheck `json:"required_status_checks"`
	}

	rule struct {
		Type       string `json:"type"`
		Parameters any    `json:"parameters"`
	}

	ruleset struct {
		Name         string     `json:"name"`
		Target       string     `json:"target"`
		Enforcement  string     `json:"enforcement"`
		Conditions   conditions `json:"conditions"`
		Rules        []rule     `json:"rules"`
		BypassActors []any      `json:"bypass_actors"`
	}
)

// AddCmd registers the `update-main-protection-ruleset` command onto `parent`.
//
// Note: the context is taken from `cl.Ctx()` (the CLI's signal context) rather
// than `cmd.Context()`. quitsh runs cobra via plain `Execute()`, so
// `cmd.Context()` is an uncancellable `context.Background()` and would not abort
// on CTRL+C.
func AddCmd(cl cli.ICLI, parent *cobra.Command) {
	sett := Settings{
		Repo:        "sdsc-ordes/modos-rs",
		RulesetName: "main-protection",
	}

	cmd := &cobra.Command{
		Use:     "update-main-protection-ruleset",
		Aliases: []string{"update-ruleset"},
		Short:   "Create-or-update the `main-protection` branch ruleset.",
		Long: "Set all status checks observed on a reference PR as required " +
			"status checks on the `main-protection` branch ruleset.",
		RunE: func(_ *cobra.Command, _ []string) error {
			return run(cl.Ctx(), &sett)
		},
	}

	cmd.Flags().StringVar(&sett.Repo, "repo", sett.Repo,
		"Repository as 'owner/name'.")
	cmd.Flags().StringVar(&sett.RulesetName, "ruleset-name", sett.RulesetName,
		"Name of the ruleset to create-or-update.")

	parent.AddCommand(cmd)
}

// run executes the full create-or-update flow.
func run(ctx context.Context, sett *Settings) error {
	// A single `gh` command context; `GH_PAGER=cat` disables the interactive
	// pager just like the shell script did.
	gh := exec.NewCmdCtxBuilder().
		BaseCmd("gh").
		Context(ctx).
		Env("GH_PAGER=cat").
		Build()

	checks := []statusCheck{
		{Context: "format", IntegrationID: 15368},
		{Context: "checks-all", IntegrationID: 15368},
	}

	log.Infof("Requiring %d checks:", len(checks))
	for _, c := range checks {
		log.Infof("  - '%v'", c)
	}

	payload, err := buildPayload(sett, checks)
	if err != nil {
		return err
	}

	return applyRuleset(gh, sett, payload)
}

// buildPayload marshals the ruleset payload with the live checks spliced in.
func buildPayload(sett *Settings, checks []statusCheck) ([]byte, error) {
	rs := ruleset{
		Name:        sett.RulesetName,
		Target:      "branch",
		Enforcement: "active",
		Conditions: conditions{
			RefName: refNameCondition{
				Include: []string{"refs/heads/main"},
				Exclude: []string{},
			},
		},
		Rules: []rule{
			{
				Type: "pull_request",
				Parameters: pullRequestParams{
					RequiredApprovingReviewCount:   1,
					DismissStaleReviewsOnPush:      false,
					RequireCodeOwnerReview:         false,
					RequireLastPushApproval:        false,
					RequiredReviewThreadResolution: true,
				},
			},
			{
				Type: "required_status_checks",
				Parameters: statusCheckParams{
					StrictRequiredStatusChecksPolicy: true,
					RequiredStatusChecks:             checks,
				},
			},
		},
		BypassActors: []any{},
	}

	payload, err := json.Marshal(rs)
	if err != nil {
		return nil, errors.AddContext(err, "could not marshal ruleset payload")
	}

	return payload, nil
}

// applyRuleset creates the ruleset or updates the existing one by name.
func applyRuleset(gh *exec.CmdContext, sett *Settings, payload []byte) error {
	endpoint := "repos/" + sett.Repo + "/rulesets"

	existingID, err := findRulesetID(gh, endpoint, sett.RulesetName)
	if err != nil {
		return err
	}

	if existingID != 0 {
		log.Infof("Updating existing ruleset #%d (%s)…", existingID, sett.RulesetName)

		target := fmt.Sprintf("%s/%d", endpoint, existingID)

		return gh.WithStdin(bytes.NewReader(payload)).
			Check("api", "-X", "PUT", target, "--input", "-")
	}

	log.Infof("Creating new ruleset (%s)…", sett.RulesetName)

	return gh.WithStdin(bytes.NewReader(payload)).
		Check("api", "-X", "POST", endpoint, "--input", "-")
}

// findRulesetID returns the id of the ruleset named `name`, or 0 if none exists.
func findRulesetID(gh *exec.CmdContext, endpoint, name string) (int64, error) {
	out, err := gh.Get("api", endpoint, "--paginate")
	if err != nil {
		return 0, errors.AddContext(err, "could not list rulesets")
	}

	var rulesets []rulesetRef
	if e := json.Unmarshal([]byte(out), &rulesets); e != nil {
		return 0, errors.AddContext(e, "could not decode rulesets response")
	}

	for _, r := range rulesets {
		if r.Name == name {
			return r.ID, nil
		}
	}

	return 0, nil
}
