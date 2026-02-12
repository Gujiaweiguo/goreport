# bi-dashboard Specification

## Purpose
TBD - created by archiving change migrate-go-backend. Update Purpose after archive.
## Requirements
### Requirement: Dashboard Designer Access
The backend SHALL provide the dashboard designer entry route `/drag/list` for authorized users.

#### Scenario: Access dashboard designer
- **WHEN** an authenticated user visits `/drag/list`
- **THEN** the dashboard designer UI is available

### Requirement: Dashboard CRUD
The backend SHALL provide APIs to create, update, list, and delete dashboard definitions.

#### Scenario: Update dashboard
- **WHEN** a user updates an existing dashboard
- **THEN** the updated dashboard is persisted and available for preview

### Requirement: Dashboard Dataset Binding
The dashboard designer SHALL support binding dashboard widgets to dataset dimensions and measures.

#### Scenario: Bind widget to dataset
- **WHEN** user adds or edits a widget in dashboard designer
- **THEN** user can select dataset as data source
- **THEN** dataset dropdown shows all available datasets
- **THEN** selecting dataset loads dimensions and measures

#### Scenario: Configure widget with dataset dimensions
- **WHEN** user configures widget dimension axis (X-axis, category)
- **THEN** user can select dimension from dataset's dimensions list
- **THEN** dimension display name is shown
- **THEN** user can enable grouping
- **THEN** widget configuration references dataset ID and dimension name

#### Scenario: Configure widget with dataset measures
- **WHEN** user configures widget measure axis (Y-axis, value)
- **THEN** user can select measure from dataset's measures list
- **THEN** user can select aggregation function
- **THEN** user can select multiple measures for multi-series widgets
- **THEN** widget configuration references dataset ID and measure names

#### Scenario: Save dashboard with dataset bindings
- **WHEN** user saves dashboard with dataset-bound widgets
- **THEN** system stores widget configurations with dataset IDs
- **THEN** system stores dimension and measure references
- **THEN** system stores aggregation and grouping settings

#### Scenario: Render dashboard with dataset data
- **WHEN** dashboard is loaded or previewed
- **THEN** backend loads dataset configurations for all widgets
- **THEN** backend queries dataset data for each widget
- **THEN** backend applies computed fields, grouping, aggregation
- **THEN** backend returns widget-specific data
- **THEN** widgets render with actual data

### Requirement: Dashboard Persistence Contract Compatibility
The system SHALL persist and retrieve dashboard definitions using a deterministic payload schema that includes dashboard config and component definitions.

#### Scenario: Save and load dashboard definition deterministically
- **WHEN** user saves a dashboard from designer
- **THEN** backend persists dashboard configuration and component payload using canonical schema
- **WHEN** user reloads the dashboard
- **THEN** backend returns equivalent structure that reproduces the original layout and component settings

#### Scenario: Backward compatibility for historical records
- **WHEN** a stored dashboard record has partial or legacy configuration fields
- **THEN** backend applies compatibility defaults and returns valid response envelope
- **THEN** frontend can render and edit without runtime crash

