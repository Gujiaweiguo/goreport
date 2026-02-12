# Change: Align Dataset Query Contract Across Frontend and Backend

## Why
The dataset query path is currently fragile because frontend and backend do not consistently enforce the same request method and payload contract. This causes runtime failures in report preview flows and blocks the baseline data-preview experience.

## What Changes
- Standardize dataset query endpoint contract for `/api/v1/datasets/{id}/data` (method, payload binding, and response envelope).
- Align frontend dataset query client usage and report-designer preview integration with the canonical backend contract.
- Add compatibility behavior and explicit error responses for invalid method/payload combinations.
- Add backend and frontend regression tests for dataset query preview workflow.

## Impact
- Affected specs: `dataset-api`
- Affected code:
  - `backend/internal/httpserver/server.go`
  - `backend/internal/dataset/handler.go`
  - `frontend/src/api/dataset.ts`
  - `frontend/src/views/ReportDesigner.vue`
