## ADDED Requirements

### Requirement: Core Feature Entry Visibility
The frontend SHALL expose visible navigation entries and routable pages for the core modules: report designer, report preview, dashboard designer, and chart editor.

#### Scenario: Authenticated user sees core entries after login
- **WHEN** an authenticated user enters the application
- **THEN** the main navigation displays entries for report designer, report preview, dashboard designer, and chart editor
- **AND** clicking each entry navigates to a valid route without blank screen

#### Scenario: User opens core route directly
- **WHEN** an authenticated user opens a core module route directly by URL
- **THEN** the corresponding page is rendered
- **AND** unauthorized users are redirected to login or no-permission page

### Requirement: Non-Blank Feature Fallback
For any core page that is not fully implemented, the frontend SHALL render a non-blank fallback state with clear status and next-step guidance.

#### Scenario: User opens partially implemented module
- **WHEN** a user opens a module that has incomplete backend or component wiring
- **THEN** the page renders a visible fallback state instead of empty content
- **AND** the fallback state shows module name, current availability status, and a retry action

#### Scenario: Core API request fails during page initialization
- **WHEN** a core page fails to load required data during initialization
- **THEN** the page renders an error state with actionable message
- **AND** the user can trigger retry without full page refresh

### Requirement: Minimum Usable Baseline for UI Modules
The frontend SHALL provide a minimum usable baseline for dashboard and chart workflows so users can complete a visible end-to-end trial path.

#### Scenario: User completes dashboard baseline path
- **WHEN** a user creates a new dashboard
- **THEN** the user can add at least one component, configure basic properties, and open preview
- **AND** the dashboard can be saved and reloaded

#### Scenario: User completes chart baseline path
- **WHEN** a user opens chart editor
- **THEN** the user can select a chart type, bind sample data, edit basic properties, and see real-time preview
- **AND** the chart configuration can be persisted and restored
