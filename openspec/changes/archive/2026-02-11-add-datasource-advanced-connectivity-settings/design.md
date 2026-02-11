## Context

Datasource connections may need to traverse isolated networks and must be bounded for operational safety. This change introduces SSH tunnel support and execution controls in datasource configuration while preserving existing tenant isolation and connection test flows.

## Goals / Non-Goals

**Goals:**
- Add explicit SSH tunnel settings for datasource connectivity.
- Support both password and private-key authentication for SSH.
- Add configurable connection pool size and query timeout settings.
- Keep existing datasource CRUD and test API compatibility.

**Non-Goals:**
- No replacement of existing datasource connector architecture.
- No introduction of new datasource types beyond profile declaration.
- No UI redesign outside datasource configuration forms.

## Decisions

### Decision 1: Model advanced connectivity as optional datasource config block
- Keep top-level datasource model stable.
- Add optional `advanced` config section for SSH and runtime controls.
- Rationale: backward compatible and easier migration.

### Decision 2: SSH tunnel settings include mode-specific validation
- `authMode=password` requires ssh username/password.
- `authMode=key` requires ssh username/private key and optional passphrase.
- Rationale: fail fast at validation and connection-test time.

### Decision 3: Runtime control defaults and bounds
- Add `maxConnections` and `queryTimeoutSeconds` with defaults and max bounds.
- Enforce bounds at datasource save time and connection use time.
- Rationale: reduce runaway resource use.

## Risks / Trade-offs

- [Risk] Misconfigured SSH settings may block connections -> Mitigation: explicit validation and diagnostic messages.
- [Risk] Aggressive timeout values can break slow queries -> Mitigation: documented defaults and configurable limits.
- [Risk] Connector profile drift -> Mitigation: central profile registry and capability checks.

## Migration Plan

1. Extend datasource config schema with optional advanced fields.
2. Add validation and connection test handling for SSH and runtime controls.
3. Update connection builders to apply timeout/pool settings.
4. Add tests for both SSH auth modes and bounds behavior.

Rollback strategy:
- Ignore advanced config block and fall back to existing direct connection behavior.
