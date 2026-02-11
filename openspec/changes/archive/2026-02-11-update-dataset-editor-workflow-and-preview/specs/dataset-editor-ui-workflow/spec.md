## ADDED Requirements

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
