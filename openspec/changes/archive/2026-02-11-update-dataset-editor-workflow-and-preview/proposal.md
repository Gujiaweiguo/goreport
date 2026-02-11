## Why

The dataset editor UI has core workflows but remains operationally incomplete: key interaction paths are weakly tested, preview rendering is hardcoded to specific field names, and editing behavior is concentrated in a single large component that is difficult to evolve safely. This change is needed to improve usability, reduce regression risk, and make dataset editing behavior consistent.

## What Changes

- Improve dataset editor workflow robustness for create/edit/save/preview/refresh/batch interactions.
- Replace hardcoded preview field assumptions with schema-driven mapping so arbitrary datasets can be previewed safely.
- Refine integration between dataset edit view and field-management interactions to ensure predictable state transitions.
- Add targeted frontend tests for critical editor paths and error recovery behavior.
- Keep compatibility with existing dataset contracts while improving UI behavior and feedback quality.

## Capabilities

### New Capabilities

- None.

### Modified Capabilities

- `dataset-editor-ui-workflow`: Strengthen workflow semantics and state transitions for save/refresh/tab switching paths.
- `dataset-management`: Improve preview behavior for heterogeneous schemas and error fallback consistency.
- `dataset-field-batch-and-grouping`: Clarify editor-side batch update and grouping operation UX behavior.

## Impact

- Affected code: `frontend/src/views/dataset/DatasetEdit.vue`, `frontend/src/views/dataset/DatasetList.vue`, `frontend/src/components/dataset/DatasetPreview.vue`, `frontend/src/components/dataset/*Field*.vue`, `frontend/src/api/dataset.ts`.
- Affected API usage: no new endpoint required in this change, but stronger client handling for existing responses.
- Dependencies: no new external dependency required.
- Risk: medium; UI behavior and test coverage changes across dataset editing and preview flows.
