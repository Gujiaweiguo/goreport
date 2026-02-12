# Change: Make Dashboard Designer Runtime and Persistence Consistent

## Why
Dashboard designer currently mixes mock behavior and incomplete persistence contracts, causing save/load/preview flows to diverge from expected production behavior. Users can create layouts, but runtime restoration and data-backed preview are not consistently guaranteed.

## What Changes
- Align backend dashboard API contract with frontend designer persistence needs (`config`, component payloads, response shape).
- Ensure dashboard save/load cycle restores component definitions deterministically.
- Replace mock-only preview assumptions with baseline runtime rendering behavior suitable for end-to-end trial flow.
- Add regression tests for create/save/load/preview baseline path.

## Impact
- Affected specs: `bi-dashboard`, `bi-dashboard-ui`, `frontend-feature-availability`
- Affected code:
  - `backend/internal/dashboard/*`
  - `frontend/src/api/dashboard.ts`
  - `frontend/src/views/DashboardDesigner.vue`
  - `frontend/src/components/dashboard/DashboardPreview.vue`
