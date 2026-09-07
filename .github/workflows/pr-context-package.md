---
# .github/workflows/pr-context-package.md
# Compile with:  gh aw compile --strict --actionlint --zizmor --poutine
#
# Companion to ai-review-2job.md — see that file for the full two-job
# rationale. This is the unprivileged half: it must stay a plain
# `pull_request` trigger (never pull_request_target) so it keeps GitHub's
# default fork restriction — read-only GITHUB_TOKEN, no secrets, no OIDC —
# even though it checks out the fork's PR head. Nothing sensitive is
# available here to leak, and this workflow never writes back to the PR.
#
# It only packages data; it never needs the AI engine. The "Skip the AI
# agent" step below always fires a `noop`, so the harness exits before any
# inference cost is incurred (see reference/cost-management#skip-the-agent-
# from-steps-using-noop). `engine:` is still declared because gh-aw requires
# one, but it is never actually invoked.
name: PR Context Package (untrusted, no secrets)

on:
    pull_request:
        types: [ labeled ]
        names: [ ai/review-2job ]         # native label filter, not a hand-rolled check
        # gh-aw blocks pull_request triggers from forks by default (repository
        # ID check). This workflow is safe to allow broadly: it has no secrets
        # and the roles: check below still gates on who applied the label, not
        # on where the PR head lives — but scope it to the known test fork
        # rather than "*" since this is only a comparison exercise.
        forks: [ "mmatur/*" ]
    roles: [ admin, maintainer, write ] # native allowlist for who may apply the label
    reaction: eyes
    status-comment: false

permissions:
    contents: read
    id-token: write   # never actually usable here (fork PR, no OIDC) — see engine: below

# Full history so `git diff <base_sha> <head_sha>` below has both commits
# locally — the default pull_request checkout is otherwise fetch-depth: 1.
checkout:
    fetch-depth: 0

# gh-aw validates the engine's credentials in the activation job unconditionally,
# before the agent job (and its noop skip) even starts — so a bare `engine: id:
# claude` fails here expecting a static ANTHROPIC_API_KEY secret. Declaring the
# same github-oidc auth as ai-review.md/ai-review-2job.md satisfies that check.
# It's never actually exercised: this trigger is a fork PR, which never gets an
# OIDC token anyway, but the noop step always fires before the engine starts.
engine:
    id: claude
    auth:
        type: github-oidc
        provider: anthropic
        federation-rule-id: fdrl_013Yw3g9LVvJzRJgQoP5zFnR
        organization-id: 797e3cc2-9e62-4091-a6e2-9bd04249babc
        service-account-id: svac_015rSKKoTmnBF2WbzqpemYW5
        workspace-id: wrkspc_01EZUP6bV8tj87UdRfaC3zKR

steps:
    -   name: Package PR context (source tree + diff + metadata)
        run: |
            set -euo pipefail
            mkdir -p /tmp/pr-context
            printf '%s' "${{ github.event.pull_request.number }}" > /tmp/pr-context/pr_number.txt
            printf '%s' "${{ github.event.pull_request.head.sha }}" > /tmp/pr-context/head_sha.txt
            printf '%s' "${{ github.event.pull_request.base.sha }}" > /tmp/pr-context/base_sha.txt
            git diff "${{ github.event.pull_request.base.sha }}" "${{ github.event.pull_request.head.sha }}" > /tmp/pr-context/pr.diff
            tar -czf /tmp/pr-context/source.tar.gz --exclude=.git .

    -   name: Upload PR context artifact
        uses: actions/upload-artifact@043fb46d1a93c77aae656e7c1c64a875d1fc6a0a # v7.0.1
        with:
            name: pr-context-${{ github.event.pull_request.number }}
            path: /tmp/pr-context
            retention-days: 1

    -   name: Skip the AI agent — this workflow only packages data for ai-review-2job.md
        # GH_AW_SAFE_OUTPUTS is a step OUTPUT of the compiler's own
        # set-runtime-paths step, not a job-wide env var — every step that
        # writes to it (including custom steps:) must reference it explicitly.
        env:
            GH_AW_SAFE_OUTPUTS: ${{ steps.set-runtime-paths.outputs.GH_AW_SAFE_OUTPUTS }}
        run: |
            msg="Packaged PR #${{ github.event.pull_request.number }} context (source tree + diff) for the workflow_run consumer; no AI review runs in this workflow."
            printf '{"type":"noop","message":"%s"}\n' "$msg" >> "$GH_AW_SAFE_OUTPUTS"
---

# PR Context Packager

This workflow only runs deterministic steps; the AI agent is always skipped
(see the `noop` step above). This section is unused but required by the
compiler.
