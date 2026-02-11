# dataset-editor-ui-workflow Specification

## Purpose
TBD - created by archiving change update-dataset-editor-workflow. Update Purpose after archive.
## Requirements
### Requirement: Dataset Editor Workflow Layout
The dataset editor SHALL provide a workflow-oriented layout with a data source area, a field management area, and a top action area.

#### Scenario: Open dataset editor
- **WHEN** user opens dataset create or edit page
- **THEN** the page shows data source and table selection area
- **THEN** the page shows field management area organized by dimensions and measures
- **THEN** the page shows top actions including Save, Save and Return, and Refresh Data

### Requirement: Save Action Semantics
The dataset editor SHALL implement distinct semantics for Save and Save and Return actions.

#### Scenario: Save current dataset
- **WHEN** user clicks Save and validation passes
- **THEN** system persists dataset and field changes
- **THEN** user remains on the current editor page
- **THEN** system shows success feedback and updates local page state

#### Scenario: Save and return to list
- **WHEN** user clicks Save and Return and validation passes
- **THEN** system persists dataset and field changes
- **THEN** system navigates user back to dataset list page
- **THEN** system preserves successful save result for list refresh

#### Scenario: Save failure handling
- **WHEN** user triggers Save or Save and Return and backend returns error
- **THEN** user stays on current page
- **THEN** unsaved form context is preserved
- **THEN** system shows readable error feedback

### Requirement: Tab and Refresh State Flow
The dataset editor SHALL provide predictable state transitions for Data Preview, Batch Management, and Refresh Data actions.

#### Scenario: Switch from preview to batch management
- **WHEN** user switches from Data Preview tab to Batch Management tab
- **THEN** system keeps current dataset context
- **THEN** system loads or reuses field metadata for batch operations
- **THEN** loading and error states are shown consistently

#### Scenario: Refresh dataset data
- **WHEN** user clicks Refresh Data
- **THEN** system reloads preview data and related metadata for current dataset
- **THEN** system uses loading state while request is in progress
- **THEN** system shows success or error feedback after completion

### Requirement: Dataset Editor State Consistency
The dataset editor SHALL preserve valid editing context across save, preview, refresh, and tab transitions.

#### Scenario: Save failure keeps editing context
- **WHEN** user triggers Save or Save and Return and backend returns failure
- **THEN** system keeps current form state and selected dataset context
- **THEN** system displays actionable error feedback without route reset

#### Scenario: Tab transition reuses or refreshes state predictably
- **WHEN** user switches between Data Preview and Batch Management tabs
- **THEN** system reuses already loaded state when valid
- **THEN** system refreshes missing state deterministically before interaction

### Requirement: Dataset Editor Critical Path Test Coverage
The frontend SHALL provide automated tests for critical dataset editor workflows.

#### Scenario: Save and Save-and-Return behavior test
- **WHEN** automated tests execute editor save actions
- **THEN** tests verify route/state behavior for both save modes and failure paths

#### Scenario: Refresh and preview behavior test
- **WHEN** automated tests execute refresh and preview actions
- **THEN** tests verify loading, success, and failure state handling

