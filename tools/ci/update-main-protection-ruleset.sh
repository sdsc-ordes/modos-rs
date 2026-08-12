#!/usr/bin/env nix-shell
#!nix-shell -i bash
#!nix-shell -p gh jq
# shellcheck disable=SC1008

set -euo pipefail

# Sets all status checks on a branch protection rule `main-protection`.
export GH_PAGER=cat
REPO="sdsc-ordes/modos-rs"
PR=3 # any recent PR that ran all 20 checks
RULESET_NAME="main-protection"

# Keep only successful checks, dedup, shape into [{context: ...}].
checks_json=$(gh pr checks "$PR" --repo "$REPO" --json name,state \
    -q '[.[] | .name] | unique | map({context: .})')

count=$(jq 'length' <<<"$checks_json")
if [[ $count -eq 0 ]]; then
    echo "No successful checks found on PR #$PR — aborting." >&2
    exit 1
fi
echo "Requiring $count checks:" >&2
jq -r '.[].context' <<<"$checks_json" >&2

# Build the payload with the live checks spliced in.
payload=$(jq -n --argjson checks "$checks_json" '{
  name: "main-protection",
  target: "branch",
  enforcement: "active",
  conditions: { ref_name: { include: ["refs/heads/main"], exclude: [] } },
  rules: [
    { type: "pull_request",
      parameters: {
        required_approving_review_count: 1,
        dismiss_stale_reviews_on_push: true,
        require_code_owner_review: false,
        require_last_push_approval: false,
        required_review_thread_resolution: true } },
    { type: "required_status_checks",
      parameters: {
        strict_required_status_checks_policy: true,
        required_status_checks: $checks } }
  ],
  bypass_actors: []
}')

# Create-or-update by name.
existing_id=$(gh api "repos/$REPO/rulesets" --paginate \
    -q ".[] | select(.name == \"$RULESET_NAME\") | .id" | head -n1)

if [[ -n $existing_id ]]; then
    echo "Updating existing ruleset #$existing_id ($RULESET_NAME)…" >&2
    echo "$payload" | gh api -X PUT "repos/$REPO/rulesets/$existing_id" --input -
else
    echo "Creating new ruleset ($RULESET_NAME)…" >&2
    echo "$payload" | gh api -X POST "repos/$REPO/rulesets" --input -
fi
