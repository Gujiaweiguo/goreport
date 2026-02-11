## 1. Editor Workflow Reliability

- [x] 1.1 Refine save/save-and-return flow to preserve editor context on failure.
- [x] 1.2 Standardize tab-switch behavior to reuse loaded state and refresh only when missing.
- [x] 1.3 Standardize refresh and preview loading/error feedback behavior.

## 2. Schema-Driven Preview

- [x] 2.1 Replace fixed preview field assumptions with schema-aware mapping logic.
- [x] 2.2 Add deterministic fallback preview mode when chart mapping is unavailable.
- [x] 2.3 Ensure preview rendering remains stable for heterogeneous dataset schemas.

## 3. Batch Interaction Consistency

- [x] 3.1 Enforce selection and payload pre-checks before batch update submission.
- [x] 3.2 Present partial-failure field details without losing successful update context.
- [x] 3.3 Refresh schema state after batch operations to prevent stale field metadata.

## 4. Verification

- [x] 4.1 Add frontend tests for save and save-and-return success/failure routes.
- [x] 4.2 Add frontend tests for preview and refresh state handling.
- [x] 4.3 Add frontend tests for batch update no-selection and partial-failure paths.
- [x] 4.4 Run `openspec validate update-dataset-editor-workflow-and-preview --strict --no-interactive`.
