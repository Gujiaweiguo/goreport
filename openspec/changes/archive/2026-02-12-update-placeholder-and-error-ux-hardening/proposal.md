# Change: Harden Placeholder Features and Error UX

## Why
Several user-facing flows still expose "开发中" placeholders or generic error messages, creating dead-end interactions and increasing support/debug cost. Core modules should fail with actionable guidance instead of opaque status text.

## What Changes
- Replace production-facing placeholder actions with either minimum viable workflows or explicit disabled states with guidance.
- Standardize backend error responses for request binding and contract violations to improve diagnosability.
- Ensure frontend error surfaces prioritize backend diagnostic messages over generic transport text.
- Add targeted tests for placeholder-to-usable transitions and error rendering behavior.

## Impact
- Affected specs: `frontend-feature-availability`, `datasource-management`
- Affected code:
  - `frontend/src/views/DatasourceManage.vue`
  - `frontend/src/views/DashboardDesigner.vue`
  - `frontend/src/api/*`
  - `backend/internal/*/handler.go` (error message consistency)
