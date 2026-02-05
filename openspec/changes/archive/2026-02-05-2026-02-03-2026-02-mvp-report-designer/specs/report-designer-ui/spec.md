# report-designer-ui Delta Specification

## ADDED Requirements

### Requirement: Report Preview Page
The system SHALL provide a report preview page that displays the rendered report with actual data from the database.

#### Scenario: User views report preview
- **WHEN** user opens report preview page with a valid report ID
- **THEN** system loads report configuration from backend
- **THEN** system calls report preview API to render the report
- **THEN** system displays the rendered HTML with actual data
- **AND** toolbar is visible with refresh and export buttons

#### Scenario: User refreshes report preview
- **WHEN** user clicks refresh button on preview page
- **THEN** system re-fetches data from database
- **THEN** system re-renders the report
- **THEN** preview displays updated data

#### Scenario: User exports report
- **WHEN** user clicks export button on preview page
- **THEN** system displays export options (Excel, PDF, etc.)
- **THEN** user can download the report in selected format

#### Scenario: Preview with parameters
- **WHEN** user opens preview page with URL parameters
- **THEN** system passes parameters to rendering engine
- **THEN** report renders with parameterized data

### Requirement: Frontend Report API Integration
The frontend SHALL provide API client functions to interact with backend report CRUD and preview endpoints.

#### Scenario: Create report via API
- **WHEN** frontend calls report create API
- **THEN** request includes report name and configuration JSON
- **THEN** backend creates report record and returns report ID

#### Scenario: Update report via API
- **WHEN** frontend calls report update API
- **THEN** request includes report ID and updated configuration
- **THEN** backend persists changes and returns updated report

#### Scenario: Preview report via API
- **WHEN** frontend calls report preview API
- **THEN** request includes report ID and optional parameters
- **THEN** backend renders report with data and returns HTML

### Requirement: Report Designer Routing
The frontend SHALL provide routes for accessing report designer and preview pages with proper authentication guards.

#### Scenario: User navigates to report designer
- **WHEN** authenticated user navigates to `/report/designer`
- **THEN** report designer page is displayed
- **AND** canvas is initialized with grid
- **AND** toolbar and property panel are visible

#### Scenario: User navigates to report preview
- **WHEN** authenticated user navigates to `/report/preview?id=<report-id>`
- **THEN** report preview page is displayed
- **AND** report is rendered with actual data
- **AND** unauthenticated users are redirected to login

### Requirement: Property Panel Styling Configuration
The property panel SHALL provide UI for configuring cell styling including font, alignment, borders, and colors.

#### Scenario: User configures cell font
- **WHEN** user selects a cell in designer
- **THEN** property panel shows font controls (size, family, weight, style)
- **THEN** user can change font properties
- **AND** changes are applied to canvas immediately

#### Scenario: User configures cell alignment
- **WHEN** user selects a cell in designer
- **THEN** property panel shows alignment controls (left, center, right, top, middle, bottom)
- **THEN** user can change alignment
- **AND** changes are applied to canvas immediately

#### Scenario: User configures cell borders
- **WHEN** user selects a cell in designer
- **THEN** property panel shows border controls (top, bottom, left, right, all)
- **THEN** user can change border style and width
- **AND** changes are applied to canvas immediately
