## ADDED Requirements
### Requirement: Report Rendering Endpoint
The backend SHALL provide a report view endpoint that renders a report by identifier with optional parameters.

#### Scenario: Render report by id
- **WHEN** a user requests a report with a valid report id and parameters
- **THEN** the system returns the rendered report view

### Requirement: Route Compatibility for Reports
The backend SHALL preserve the existing report route prefix `/jmreport/` for report rendering and view access.

#### Scenario: Existing UI route access
- **WHEN** the UI requests a view under `/jmreport/`
- **THEN** the backend serves a compatible response without requiring UI changes
