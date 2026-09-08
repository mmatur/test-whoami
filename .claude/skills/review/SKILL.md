---
name: code-review
description: Repository-specific guidance for the AI code review agent reviewing this repo
user-invocable: true
disable-model-invocation: true
---

# Code review guidance

This is `whoami` (Go): a tiny HTTP/gRPC/WebSocket server that echoes request
and host information, used as a debug/test target elsewhere in the Traefik
ecosystem. The reviewer already fixes the role, the tools, the severity
definitions, the finding mechanics and the output format; this guidance is
**additive** and must not restate them.

Review in this priority order, and keep the bar high at every level.

## Security

- This binary is deployed as a live network service. Any handler that makes
  an outbound request on behalf of a caller (SSRF), reads a caller-controlled
  path, or executes a caller-controlled command is high severity by default.
- No hardcoded certificates, keys, or tokens. Flag any relaxed TLS setting
  (`InsecureSkipVerify`, a lowered `MinVersion`, disabled client-cert
  verification) even in a flag/env-gated code path.
- The `env=true` query parameter on `/`, `/api`, and similar handlers
  intentionally dumps the process environment — this is existing, accepted
  behavior of a debug tool, not a new finding, unless a change measurably
  widens what it exposes.

## Correctness

- Resource leaks: an `http.Response.Body`, a file, or a websocket connection
  that isn't closed on every path (including error returns).
- Missing timeouts/bounds on anything that reads from the network or blocks
  on I/O (a handler that can hang or copy an unbounded amount of data is a
  real defect here, not a nit).
- Concurrent access to shared state (see `mutexHealthState`) must stay
  guarded by the existing lock; flag any new shared mutable state that
  isn't.

## Maintainability

- Idiomatic Go: early returns, the standard library over new dependencies,
  small single-purpose handlers matching the existing `xHandler` pattern in
  `app.go`.
- This repo's established convention is to discard errors from response
  writes with `_, _ = fmt.Fprintln(w, ...)` / `_, _ = w.Write(...)` — this is
  intentional (a write failure to the response is not actionable) and is
  **not** a finding. Only flag a discarded error when it hides a real
  problem (e.g. closing an upstream connection, a parse error that leaves
  state inconsistent).
- Tests use the standard library's `testing` package only (no `testify`),
  table-driven with `t.Run` and `t.Parallel()`, matching `content_test.go`.
  Don't ask for a different test framework.

## Do not flag

- Generated code (`grpc/*.pb.go`, `grpc/*_grpc.pb.go`).
- Patterns already used consistently across the codebase (see the error
  handling note above).

Before claiming a change is unsafe, check how the changed function is used
elsewhere in the repository (or in the diff itself) — the call sites decide.
