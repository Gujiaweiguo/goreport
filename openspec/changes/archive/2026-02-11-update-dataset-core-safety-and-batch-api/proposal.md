## Why

Dataset flows are currently incomplete and fragile in production paths: frontend batch field updates call an API that is not implemented, and SQL-based query/preview paths do not have sufficient guardrails for risky statements, complexity, and timeout behavior. This change is needed now to make dataset operations reliable and safe before further editor feature expansion.

## What Changes

- Add backend batch field update API contract at `PATCH /api/v1/datasets/{id}/fields` with field-level validation and partial-failure reporting.
- Align frontend dataset batch-update call and backend response shape for deterministic error handling.
- Add SQL safety controls for dataset preview/query execution, including statement validation, complexity limits, timeout controls, and bounded pagination.
- Standardize error semantics for dataset query/preview failures to improve recoverability in UI and integration paths.
- Add focused backend and frontend tests for batch API and SQL safety behaviors.

## Capabilities

### New Capabilities

- `dataset-query-safety-controls`: SQL validation and execution guardrails for dataset query/preview endpoints.

### Modified Capabilities

- `dataset-api`: Add batch field update endpoint behavior and tighten query/preview execution constraints.
- `dataset-field-batch-and-grouping`: Define backend contract details for batch updates and field-level error reporting.
- `dataset-management`: Clarify preview/query guardrails and failure handling expectations.

## Impact

- Affected code: `backend/internal/dataset/*`, `backend/internal/httpserver/server.go`, `backend/internal/repository/*dataset*`, `frontend/src/api/dataset.ts`, `frontend/src/views/dataset/DatasetEdit.vue`.
- Affected API: `PATCH /api/v1/datasets/{id}/fields` (new) and stricter validation on `/api/v1/datasets/{id}/data` and `/api/v1/datasets/{id}/preview`.
- Dependencies: no new external dependency required.
- Risk: medium; includes API and query execution behavior changes with tenant-scoped validation paths.
