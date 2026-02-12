## 1. Contract Alignment
- [x] 1.1 Confirm canonical dataset query contract for `/api/v1/datasets/{id}/data` (HTTP method, request payload, response envelope).
- [x] 1.2 Update backend router and handler to enforce the canonical contract consistently.
- [x] 1.3 Ensure backend returns explicit 4xx diagnostics for invalid method/payload usage.

## 2. Frontend Integration
- [x] 2.1 Align `frontend/src/api/dataset.ts` query method and typing with canonical backend contract.
- [x] 2.2 Update report-designer preview data flow to call the aligned API path only.
- [x] 2.3 Verify frontend error presentation surfaces backend diagnostics instead of opaque transport messages.

## 3. Regression Coverage
- [x] 3.1 Add backend tests for dataset query route/handler contract (success, invalid method, invalid payload).
- [x] 3.2 Add frontend tests for report-designer dataset preview query flow.
- [x] 3.3 Run `go test ./...` and frontend targeted tests to verify no regressions.

## 4. Validation
- [x] 4.1 Run `openspec validate update-dataset-query-contract-alignment --strict --no-interactive`.
- [x] 4.2 Prepare rollout note for any legacy caller compatibility behavior.
