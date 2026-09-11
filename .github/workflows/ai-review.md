---
# AI code review for traefik/test-whoami, the staging repository of this
# package: roll-outs land here first. The review itself (firewall, tools, safe
# outputs, prompt) is shared/base.ai-review.md; this file holds only what gh-aw
# does not accept from an import: trigger, permissions, engine, budgets.
#
# Installed copy: .github/workflows/ai-review.md in the target repository, next
# to .github/workflows/shared/base.ai-review.md. Never edit those copies. To
# roll out a change to either file, run from the target repository:
#   gh aw add traefik/ai/test-whoami.ai-review -n ai-review --force
#   gh aw compile --strict --actionlint
# (zizmor and poutine run in traefik/ai CI on the same lock content; in a
# target repository, without the zizmor configuration kept there, they fail.)
# Not `gh aw update`: it refreshes this file but leaves shared/ as it was.
name: "AI code review"
description: "Reviews a pull request when a maintainer applies the ai/review label. Read-only agent; writes only through safe outputs."
emoji: "🔍"

# pull_request_target, not pull_request: under a plain pull_request trigger a
# fork PR gets neither secrets nor an OIDC id-token, so the Anthropic OIDC
# exchange fails. pull_request_target runs in the base repository context, so
# it works for forks too. The compiler refuses this trigger with a checkout
# (see `checkout: false` below): the PR head is never cloned, which is what
# makes the trigger safe.
#
# The `ai/review` label is not removed automatically. To re-review the same
# commit, remove and re-apply it.
on:
  pull_request_target:
    types: [labeled]
  roles: [admin, maintainer, write]   # exact-match allowlist of who may apply the label
  reaction: eyes
  status-comment: false               # one less write to the PR timeline

# pull_request_target has no label-name filter; gate on the exact label.
if: github.event.label.name == 'ai/review'

checkout: false

# Cache the pre-fetched diff/metadata/comments on the PR head SHA so a
# re-review of the same commit skips the GitHub API calls.
cache:
  key: pr-prefetch-${{ github.event.pull_request.head.sha }}
  path: /tmp/gh-aw/agent
  restore-keys:
    - pr-prefetch-${{ github.event.pull_request.number }}-

# The agent job is read-only. All writes happen in separate safe-output jobs.
permissions:
  contents: read
  pull-requests: read
  id-token: write

# Keyless Anthropic auth (Workload Identity Federation). No ANTHROPIC_API_KEY
# secret exists in any repository; the federation rule decides which GitHub
# repositories may exchange their OIDC token. These IDs are identifiers, not
# credentials. Do not add PATs, custom github-token overrides, or MCP servers.
engine:
  id: claude
  auth:
    type: github-oidc
    provider: anthropic
    federation-rule-id: fdrl_013Yw3g9LVvJzRJgQoP5zFnR
    organization-id: 797e3cc2-9e62-4091-a6e2-9bd04249babc
    service-account-id: svac_015rSKKoTmnBF2WbzqpemYW5
    workspace-id: wrkspc_01EZUP6bV8tj87UdRfaC3zKR

# Cost guardrails. A review that needs more than this is a review that should
# be split; fail closed rather than run away.
timeout-minutes: 20
max-turns: 100

# Smaller than the shared default (6): the staging repository is where the
# per-repository override of a safe-output type is exercised end to end.
safe-outputs:
  create-pull-request-review-comment:
    max: 4
    target: triggering

imports:
  - shared/base.ai-review.md
source: traefik/ai/workflows/test-whoami.ai-review.md@main
---

# Repository

`traefik/test-whoami` is the public staging repository of the AI review; its
pull requests exist to exercise the reviewer.
