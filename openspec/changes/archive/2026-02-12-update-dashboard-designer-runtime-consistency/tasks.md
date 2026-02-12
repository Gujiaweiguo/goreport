## 1. Persistence Contract
- [x] 1.1 Define canonical dashboard payload and response fields for config + components.
- [x] 1.2 Align backend dashboard model/service/handler serialization with canonical payload.
- [x] 1.3 Add compatibility handling for historical records with partial/missing fields.

## 2. Frontend Runtime Alignment
- [x] 2.1 Align `frontend/src/api/dashboard.ts` types and calls with backend canonical contract.
- [x] 2.2 Update `DashboardDesigner.vue` save/load flow to use deterministic persisted component data.
- [x] 2.3 Remove mock-dependent assumptions that bypass persisted component payload.

## 3. Preview Baseline
- [x] 3.1 Define baseline preview rendering criteria for chart/table/text components.
- [x] 3.2 Update `DashboardPreview.vue` to render baseline runtime output from persisted component data.
- [x] 3.3 Ensure preview refresh reflects latest persisted state.

## 4. Regression and Validation
- [x] 4.1 Add backend tests for dashboard create/update/get/list contract compatibility.
- [x] 4.2 Add frontend tests for save-load-preview baseline path.
- [x] 4.3 Run `openspec validate update-dashboard-designer-runtime-consistency --strict --no-interactive`.
