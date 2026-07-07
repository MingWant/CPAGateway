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
- Optional JSON state persistence for policies, templates, audit, and usage counters
- Optional Redis-backed shared counters for multi-instance quotas, rate limits, token quotas, and inflight limits
- Optional plugin-level read/admin/UI tokens layered on top of host management authentication
- Read-only health diagnostics for CPAMC compatibility, persistence, Redis, and security status

## Production-Oriented Configuration

The host still authenticates `/v0/management/...` requests. For stronger plugin-local controls and durable state, add these fields to the plugin config:

```yaml
persistence:
  state_path: ./gateway-state.json
  persist_runtime: true

cluster:
  backend: redis
  redis:
    addr: 127.0.0.1:6379
    key_prefix: cpa-gateway
    failure_mode: reject # reject | allow | local_fallback

security:
  require_management_token: true
  admin_tokens:
    - change-this-admin-token
  read_tokens:
    - change-this-read-token
  ui_access_tokens:
    - change-this-ui-token
```

Notes:

- Redis is not required for normal CPAMC usage. If `cluster.backend` is unset, the plugin keeps using local in-process counters and CPAMC does not need a Redis service.
- `state_path` stores policy/template state and, when `persist_runtime` is true, runtime usage and audit snapshots. Protect this file because it can contain policy API-key match secrets.
- `cluster.backend: redis` moves quota, per-minute, token, and inflight counters into Redis so multiple plugin instances share enforcement state.
- `redis.failure_mode` controls Redis outage behavior: `reject` fails closed with `503`, `allow` fails open, and `local_fallback` temporarily uses per-process counters. Keep `reject` for strict multi-instance quotas; use `local_fallback` only when availability matters more than exact cluster-wide limits during an outage.
- `admin_tokens` are required for mutating management routes when `require_management_token` is true. `read_tokens` can only call read and preview routes.
- `ui_access_tokens` protects the browser resource route, which the host exposes separately from authenticated management API routes.

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
- `GET /health`
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

`GET /health` returns non-secret diagnostics only: plugin version, counter backend, whether Redis is required, Redis ping status when enabled, persistence mode, token counts, and policy/template counters. It does not return API keys, Redis passwords, or management token values.

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
