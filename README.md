# CLIProxyAPI Gateway Plugin

A plugin-first gateway and control plane for CLIProxyAPI / CPAMC.

This plugin is built around CPAMC top-level `Access & Authentication` API keys. It does not create a second downstream user key model. Instead, it lets you attach policy, routing, rewrite, governance, preview, and audit behavior directly to the API keys that already exist in the host.

## What It Does

- Binds policy per top-level API key
- Filters and rewrites incoming requests
- Routes traffic across models and providers
- Enforces provider/model allow and deny rules
- Applies quotas, rate limits, concurrency limits, validity windows, and schedules
- Provides preview/apply workflows for route-pool operations
- Records operator audit with before/after snapshots and structured diffs
- Exposes an embedded management UI and management API

## Current Capabilities

- Model rewrite, provider prefix forcing, and endpoint compatibility rewrite
- Weighted routing, route pools, provider affinity, fallback, failover, and mirror targets
- Pool operator actions such as drain, resume, canary split, restore weights, provider traffic shift, and rebalance by health
- Dry-run with stage trace and routing hints
- Safe apply with preview token flow
- Signed preview token verification on the backend
- Operator audit with structured member diffs

## Directory Layout

- `go/main.go`
  Plugin entrypoint, management routes, embedded UI, state handling, audit, preview/apply flow, and tests-facing helpers.
- `go/rules.go`
  Match, rewrite, routing, failover, and shared gateway execution logic.
- `go/main_test.go`
  Policy, routing, management API, governance, preview token, and operator audit coverage.

## Management API

The plugin exposes management routes under `/v0/management/plugins/gateway/...`.

Core route groups:

- `GET /keys`
- `GET/PUT/PATCH/DELETE /policies`
- `POST /policies/add`
- `POST /policies/clone`
- `POST /rules/add`
- `PATCH/DELETE /rules`
- `POST /route-members/op`
- `POST /route-members/preview`
- `GET /usage`
- `POST /usage/reset`
- `GET /audit`
- `GET /audit/detail`
- `GET /audit/summary`
- `POST /dry-run`

## Host Contract

This plugin depends on a small host-side metadata pass-through so interceptors can see the top-level access context.

In the current CLIProxyAPI workspace, that support lives in:

- `sdk/api/handlers/handlers.go`
- `sdk/api/handlers/handlers_interceptors_test.go`
- `sdk/cliproxy/executor/types.go`
- `sdk/pluginapi/types.go`

If this plugin is moved to its own repository, those host-side changes still need to remain in the CLIProxyAPI host repository unless the plugin interface is expanded upstream.

## Local Development

From the plugin module directory:

```powershell
Set-Location go
go test ./...
```

In restricted or sandboxed environments, local cache overrides are useful:

```powershell
$env:GOCACHE='.tmp\gocache'
$env:GOMODCACHE='.tmp\gomodcache'
go test ./...
```

## Split Recommendation

This plugin is large enough to live as its own repository.

Suggested split:

- Move `examples/plugin/gateway` into a standalone repository
- Keep the host-side metadata plumbing in the CLIProxyAPI repository
- Replace the example implementation in the host repository with a thin pointer document that links to the standalone plugin repository

That split keeps plugin releases, issues, and documentation independent while preserving a minimal host integration contract.

## Status

This is a strong plugin-level `v1` foundation:

- The gateway behavior is already substantial
- The management UI is already usable
- The backend execution path is largely converged across preview and apply

Good next steps after the split:

- Finish converging preview/apply into a single execution engine
- Add batch provider-family and model-family traffic actions
- Add approval and rollback-oriented control-plane flows
