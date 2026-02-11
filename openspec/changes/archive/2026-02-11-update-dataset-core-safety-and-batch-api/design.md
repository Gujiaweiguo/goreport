## Context

Dataset API and editor flows currently depend on field-level batch updates and SQL execution paths that are not consistently protected. Frontend already invokes a batch field update endpoint, but backend route coverage is incomplete for that contract. Query and preview paths execute user-provided SQL with limited runtime guardrails, which increases security and stability risk under malformed or expensive statements.

## Goals / Non-Goals

**Goals:**
- Provide a backend batch field update endpoint aligned with current frontend usage and field-level error reporting.
- Add SQL safety controls before execution for dataset query and preview operations.
- Enforce bounded execution behavior (timeout, pagination caps, complexity checks) for dataset data paths.
- Keep current tenant and authorization boundaries intact while improving error transparency.

**Non-Goals:**
- No redesign of the dataset editor UI layout.
- No migration of expression language or replacement of existing query executor architecture.
- No introduction of external SQL parser dependencies in this change.

## Decisions

### Decision 1: Add explicit batch update contract at dataset field collection endpoint
- Endpoint: `PATCH /api/v1/datasets/{id}/fields`.
- Request: list of field updates keyed by field id.
- Response: includes updated field ids and field-level errors for rejected items.
- Rationale: frontend already has batch workflow and needs deterministic partial-failure handling.
- Alternative considered: loop single-field `PUT` from frontend only. Rejected because it increases round trips and makes consistent partial-failure semantics harder to maintain.

### Decision 2: Introduce server-side SQL safety gate before query/preview execution
- Validate disallowed statements/keywords for dataset SQL operations.
- Enforce query complexity limits (for example max nested subquery depth and join count).
- Reject unsafe payloads with explicit validation errors.
- Rationale: reduces operational risk without changing datasource model or query executor ownership.
- Alternative considered: rely only on datasource permissions. Rejected because application-level constraints are still needed for predictable behavior.

### Decision 3: Standardize execution bounds for preview/query
- Apply timeout to preview and query execution contexts.
- Cap page size and enforce default pagination behavior.
- Keep existing preview row cap but make it part of explicit guardrail behavior.
- Rationale: protects service stability and gives frontend predictable failure modes.
- Alternative considered: only DB-side timeout. Rejected due to inconsistent enforcement across environments.

### Decision 4: Preserve backward compatibility for existing clients
- Existing list/get/create/update/delete/query/preview contracts remain available.
- New batch endpoint is additive.
- Query/preview safety rejects unsafe requests with clear errors, but does not change successful response shape.

## Risks / Trade-offs

- [Risk] Strict SQL safety checks may reject previously accepted edge queries -> Mitigation: start with a documented allowlist/denylist and targeted regression tests.
- [Risk] Partial-failure batch behavior may be misinterpreted by callers -> Mitigation: enforce explicit response contract with updated/failed arrays and frontend handling tests.
- [Risk] Timeout bounds may affect long-running valid queries -> Mitigation: define default values in config and tune with staged verification.

## Migration Plan

1. Add backend endpoint and validation logic behind existing dataset auth/tenant checks.
2. Add SQL safety gate and execution bounds in query/preview path.
3. Align frontend batch workflow handling to backend response details.
4. Add tests for endpoint contract and safety gate behavior.
5. Roll out with staging verification on representative datasets.

Rollback strategy:
- Remove new route registration and handler wiring for batch endpoint if needed.
- Revert safety gate checks to previous permissive behavior while preserving logging for diagnosis.

## Open Questions

- Should partial-success batch updates return HTTP 200 with error details, or a dedicated non-2xx contract when any item fails?
- Should SQL safety violations be surfaced as validation errors (`400`) or forbidden operations (`422`)?
