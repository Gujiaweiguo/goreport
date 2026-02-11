## 1. Batch Field Update API

- [x] 1.1 Add route registration for `PATCH /api/v1/datasets/:id/fields` in dataset router wiring.
- [x] 1.2 Implement handler request binding and tenant/permission checks for batch field updates.
- [x] 1.3 Implement service-layer batch update orchestration with field-level validation and partial-failure collection.
- [x] 1.4 Define and return deterministic batch response contract (updated ids + field-level errors).

## 2. SQL Safety Controls

- [x] 2.1 Implement SQL safety validation for disallowed statements before preview/query execution.
- [x] 2.2 Implement query complexity checks before execution for guarded dataset SQL.
- [x] 2.3 Apply timeout bounds to preview/query execution context.
- [x] 2.4 Enforce default/max pagination limits consistently in query path.

## 3. Frontend Contract Alignment

- [x] 3.1 Align frontend batch update request/response handling with new backend PATCH contract.
- [x] 3.2 Surface field-level batch failures in actionable UI messages.
- [x] 3.3 Ensure query/preview guardrail errors are distinguishable from transport failures.

## 4. Verification

- [x] 4.1 Add backend tests for batch update success/partial-failure/tenant-reject scenarios.
- [x] 4.2 Add backend tests for SQL safety rejection, timeout, and bounded pagination behavior.
- [x] 4.3 Add frontend tests for batch update partial failure handling.
- [x] 4.4 Run `openspec validate update-dataset-core-safety-and-batch-api --strict --no-interactive`.
