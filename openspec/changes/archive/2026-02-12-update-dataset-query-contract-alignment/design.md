## Context
Report-designer preview currently depends on dataset query execution. The implementation has drifted: route registration and handler expectations are not guaranteed to match in all runtime paths.

## Goals / Non-Goals
- Goals:
  - Define one canonical dataset query contract and enforce it in backend route + handler.
  - Ensure frontend preview path uses only this contract.
  - Provide actionable error messages for contract violations.
- Non-Goals:
  - Redesign query engine internals.
  - Introduce new dataset query capabilities beyond contract alignment.

## Decisions
- Decision: Keep `POST /api/v1/datasets/{id}/data` as the canonical query endpoint.
  - Alternatives considered:
    - `GET` with query-string body serialization (rejected: poor expressiveness for complex filters/aggregations).
    - Dual-support `GET` + `POST` long-term (rejected: increases maintenance and ambiguity).
- Decision: Keep response envelope shape stable (`success/result/message`) and ensure query result payload is deterministic.
- Decision: Add explicit 4xx diagnostics for method/payload mismatch.

## Risks / Trade-offs
- Risk: Existing hidden callers may rely on incorrect method.
  - Mitigation: Add regression tests and temporary compatibility note in rollout docs.
- Risk: Frontend fallback paths may mask backend mismatch.
  - Mitigation: Add integration-focused tests in report designer preview flow.

## Migration Plan
1. Update backend route/handler contract enforcement.
2. Update frontend dataset API client and report preview caller.
3. Add regression tests.
4. Validate end-to-end in docker compose development environment.

## Open Questions
- Should temporary compatibility handling for legacy callers be retained for one release window or removed immediately?
