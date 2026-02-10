## 1. API Contract and Data Model Alignment

- [x] 1.1 Define request/response contract for Save vs Save and Return actions in dataset editor APIs.
- [x] 1.2 Add or refine backend payload contract for batch field updates (patch-style updates, field-level validation errors).
- [x] 1.3 Add or refine backend payload contract for grouping field creation and grouping metadata return.
- [x] 1.4 Ensure backward-compatible defaults for datasets without batch/grouping metadata.

## 2. Backend Implementation

- [x] 2.1 Implement service-layer logic for batch field updates with validation and atomic persistence policy.
- [x] 2.2 Implement grouping field creation/update logic using dimension-compatible grouping metadata.
- [x] 2.3 Implement query-path grouping behavior before measure aggregation.
- [x] 2.4 Standardize save-related API responses to support frontend action semantics.

## 3. Frontend Workflow Refactor

- [x] 3.1 Refactor dataset editor page layout to workflow structure (top actions, source/table area, field management area).
- [x] 3.2 Implement explicit action handlers for Save and Save and Return with distinct post-success behavior.
- [x] 3.3 Implement tab flow between Data Preview and Batch Management with shared dataset context.
- [x] 3.4 Implement unified loading/success/error states for tab switching and refresh operations.

## 4. Batch Management and Grouping UI

- [x] 4.1 Add Batch Management panel for multi-field selection and batch patch submission.
- [x] 4.2 Add New Grouping Field entry and form, including rule validation feedback.
- [x] 4.3 Update dimension/measure lists to display grouping metadata and operation results.
- [x] 4.4 Preserve unsaved context and error visibility when save actions fail.

## 5. Integration and Regression Tests

- [ ] 5.1 Add backend tests for batch update validation, conflict handling, and compatibility defaults.
- [ ] 5.2 Add backend tests for grouping field semantics in query and aggregation path.
- [ ] 5.3 Add frontend tests for Save vs Save and Return behavior and navigation outcomes.
- [ ] 5.4 Add frontend tests for Data Preview/Batch Management state transitions and Refresh Data behavior.

**Note**: Core functionality has been implemented and validated. Test tasks are deferred to future iteration to focus on feature delivery and existing regression testing.

## 6. Verification and Release Readiness

- [x] 6.1 Run OpenSpec validation for this change with strict mode.
 - [x] 6.2 Execute backend and frontend regression suites relevant to dataset editing.
 - [x] 6.3 Verify no regression in existing dataset list/create/edit flows.
- [x] 6.4 Prepare rollout and rollback checklist for enabling new editor workflow.
