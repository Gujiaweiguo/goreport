# bi-dashboard-ui Specification

## Purpose
TBD - created by archiving change build-custom-frontend. Update Purpose after archive.
## Requirements
### Requirement: Dashboard Canvas Designer

The system SHALL provide a canvas-based dashboard designer that allows users to create interactive dashboards by dragging and dropping components.

#### Scenario: User creates new dashboard
- **WHEN** user clicks "Create Dashboard" button
- **THEN** a blank canvas designer is displayed
- **AND** canvas size is configurable (default 1920x1080)
- **AND** grid is visible for alignment
- **AND** component panel is visible on left side
- **AND** property panel is visible on right side

#### Scenario: User changes canvas size
- **WHEN** user enters new width and height in canvas settings
- **THEN** canvas size is updated immediately
- **AND** all components are repositioned proportionally if needed
- **AND** canvas center is maintained

#### Scenario: User zooms canvas
- **WHEN** user uses mouse wheel or zoom controls to zoom
- **THEN** canvas zooms in or out (10% increments, 50%-200% range)
- **AND** components scale proportionally
- **AND** zoom percentage is displayed in status bar

#### Scenario: User pans canvas
- **WHEN** user holds spacebar and drags mouse
- **THEN** canvas pans in the drag direction
- **AND** smooth scrolling is applied
- **AND** canvas boundaries are respected

#### Scenario: User saves dashboard
- **WHEN** user clicks "Save" button
- **THEN** dashboard design is sent to backend via API
- **AND** success message is displayed
- **AND** "Saved" indicator appears in status bar
- **AND** auto-save timer resets

### Requirement: Component Drag and Drop

The system SHALL allow users to drag components from the component panel and drop them onto the canvas.

#### Scenario: User drags component from panel to canvas
- **WHEN** user starts dragging a component (e.g., "Text Card")
- **THEN** component follows mouse cursor with shadow
- **AND** drop target indicator appears on canvas
- **AND** grid snap points are visible

#### Scenario: User drops component on canvas
- **WHEN** user releases mouse over canvas
- **THEN** component is placed at drop location
- **AND** component is snapped to nearest grid point
- **AND** component is selected and property panel shows its properties
- **AND** drop is recorded in undo history

#### Scenario: User drags component on canvas to move
- **WHEN** user drags an existing component
- **THEN** component follows mouse cursor
- **AND** original position is maintained (ghost image)
- **AND** drop target indicator appears
- **AND** grid snap points are visible

#### Scenario: User drops moved component
- **WHEN** user releases mouse to finish move
- **THEN** component is moved to new position
- **AND** component is snapped to grid
- **AND** move is recorded in undo history
- **AND** layer order is maintained

#### Scenario: User drags component outside canvas
- **WHEN** user drags component beyond canvas boundaries
- **THEN** drop target indicator shows "invalid" (red)
- **AND** component is highlighted in red
- **AND** drop is prevented (component returns to original position)

### Requirement: Component Property Configuration

The system SHALL provide a property panel for configuring component properties.

#### Scenario: User selects component
- **WHEN** user clicks on a component on canvas
- **THEN** component is highlighted with selection border
- **AND** property panel displays component properties
- **AND** properties are grouped (Common, Style, Data, Interaction)
- **AND** layer panel highlights the component's layer

#### Scenario: User changes component common property
- **WHEN** user edits a common property (e.g., Title, ID)
- **THEN** property change is reflected immediately on canvas
- **AND** change is recorded in undo history
- **AND** property value is validated (if required)

#### Scenario: User changes component style property
- **WHEN** user edits a style property (e.g., Font Size, Color)
- **THEN** style change is reflected immediately on canvas
- **AND** change is recorded in undo history
- **AND** color picker or font selector is used for complex properties

#### Scenario: User changes component data property
- **WHEN** user configures a data property (e.g., Data Source, Field)
- **THEN** data property options are loaded from backend
- **AND** user can test data connection
- **AND** data preview is available
- **AND** change is recorded in undo history

#### Scenario: User changes component interaction property
- **WHEN** user configures an interaction (e.g., Link, Drill-down)
- **THEN** interaction configuration is saved
- **AND** interaction is tested in preview mode
- **AND** change is recorded in undo history

### Requirement: Layer Management

The system SHALL provide a layer panel for managing component layers (order, visibility, locking).

#### Scenario: User views layer panel
- **WHEN** dashboard designer is open
- **THEN** layer panel is visible on right side
- **AND** all components on canvas are listed in layer panel
- **AND** layers are ordered from top to bottom
- **AND** visibility toggles are shown for each layer

#### Scenario: User reorders layers via drag
- **WHEN** user drags a layer in layer panel
- **THEN** layer follows mouse cursor
- **AND** drop target indicator appears between layers
- **AND** layer order changes on canvas immediately when dropped

#### Scenario: User reorders layers via buttons
- **WHEN** user clicks "Move Up" or "Move Down" button for a layer
- **THEN** layer moves up or down one position
- **AND** layer order changes on canvas immediately
- **AND** buttons are disabled for top and bottom layers

#### Scenario: User toggles layer visibility
- **WHEN** user clicks visibility toggle (eye icon) for a layer
- **THEN** component is hidden on canvas
- **AND** toggle changes to "hidden" state
- **AND** component remains in layer panel
- **AND** hidden component is not exported or previewed

#### Scenario: User locks layer
- **WHEN** user clicks lock toggle (lock icon) for a layer
- **THEN** component cannot be selected or moved on canvas
- **AND** lock icon changes to "locked" state
- **AND** locked layer cannot be reordered
- **AND** locked layer properties cannot be edited

#### Scenario: User deletes layer
- **WHEN** user clicks delete button for a layer
- **THEN** confirmation dialog appears
- **AND** if confirmed, component is removed from canvas
- **AND** layer is removed from layer panel
- **AND** deletion is recorded in undo history

### Requirement: Data Binding

The system SHALL allow users to bind components to data sources for dynamic data display.

#### Scenario: User binds component to database field
- **WHEN** user selects a component and opens data binding dialog
- **THEN** user can select a data source from dropdown
- **AND** user can select a table from data source
- **AND** user can select fields from table
- **AND** selected fields are shown in binding preview
- **AND** binding is saved to component properties

#### Scenario: User configures SQL query for component
- **WHEN** user selects "Custom SQL" data source option
- **THEN** SQL editor is displayed
- **AND** user can write or edit SQL query
- **AND** syntax highlighting is applied
- **AND** SQL can be tested before saving

#### Scenario: User tests data binding
- **WHEN** user clicks "Test" button in data binding dialog
- **THEN** system executes SQL query or fetches data
- **AND** sample data is displayed in preview table
- **AND** data row count is shown
- **AND** query execution time is displayed

#### Scenario: User configures data refresh interval
- **WHEN** user sets a refresh interval (e.g., 60 seconds)
- **THEN** component data is refreshed at specified interval
- **AND** refresh indicator is shown on component
- **AND** refresh can be disabled by setting interval to 0

#### Scenario: User configures data filtering
- **WHEN** user adds filter conditions to data binding
- **THEN** filter conditions are applied to data query
- **AND** multiple filters can be combined (AND/OR logic)
- **AND** filter operators are supported (=, >, <, LIKE, IN, etc.)
- **AND** filters can be tested before saving

### Requirement: Dashboard Preview

The system SHALL provide a preview mode to view the dashboard as it would appear to end users.

#### Scenario: User opens preview mode
- **WHEN** user clicks "Preview" button
- **THEN** dashboard is displayed in preview mode
- **AND** designer elements (grid, selection borders, resize handles) are hidden
- **AND** property panel and layer panel are hidden
- **AND** toolbar changes to preview toolbar (Refresh, Fullscreen, Close)

#### Scenario: User previews dashboard with real data
- **WHEN** dashboard is in preview mode
- **THEN** all components display real data from data sources
- **AND** data is refreshed at configured intervals
- **AND** interactive features (drill-down, hover) are enabled
- **AND** no editing is allowed

#### Scenario: User refreshes dashboard in preview
- **WHEN** user clicks "Refresh" button in preview
- **THEN** all component data is refreshed
- **AND** refresh indicator is shown during refresh
- **AND** last refresh timestamp is updated
- **AND** errors are displayed if any component fails to refresh

#### Scenario: User views fullscreen preview
- **WHEN** user clicks "Fullscreen" button in preview
- **THEN** dashboard expands to full screen
- **AND** browser fullscreen API is used
- **AND** "Exit Fullscreen" button is shown in corner
- **AND** dashboard is centered and scaled to fit screen

#### Scenario: User exits preview mode
- **WHEN** user clicks "Close" button in preview
- **THEN** preview mode is closed
- **AND** designer mode is reopened
- **AND** all designer elements are restored
- **AND** unsaved changes in preview are discarded

### Requirement: Common Dashboard Components

The system SHALL provide common dashboard components for building dashboards.

#### Scenario: User adds Text Card component
- **WHEN** user drags "Text Card" to canvas
- **THEN** text card is created with default content
- **AND** property panel shows text card properties:
  - Title
  - Content (rich text)
  - Background color
  - Border style
  - Padding

#### Scenario: User adds Image component
- **WHEN** user drags "Image" to canvas
- **THEN** image component is created
- **AND** property panel shows image properties:
  - Image source (upload or URL)
  - Width and height
  - Border radius
  - Fit mode (cover, contain, stretch)

#### Scenario: User adds Number Card component
- **WHEN** user drags "Number Card" to canvas
- **THEN** number card is created
- **AND** property panel shows number card properties:
  - Title
  - Data binding
  - Number format
  - Trend indicator (arrow, color)
  - Unit

#### Scenario: User adds Table component
- **WHEN** user drags "Table" to canvas
- **THEN** table component is created with sample data
- **AND** property panel shows table properties:
  - Data source binding
  - Column configuration (width, alignment, format)
  - Row height
  - Header style
  - Pagination (yes/no)

#### Scenario: User adds Border Frame component
- **WHEN** user drags "Border Frame" to canvas
- **THEN** border frame is created
- **AND** property panel shows border properties:
  - Border style (solid, dashed, double)
  - Border color
  - Border width
  - Corner radius
  - Background color
  - Title text and style

#### Scenario: User adds Chart component
- **WHEN** user drags "Chart" to canvas
- **THEN** chart component is created with sample chart
- **AND** property panel shows chart properties:
  - Chart type (bar, line, pie, etc.)
  - Data source binding
  - Title and subtitle
  - Legend settings
  - Axis settings
  - Color scheme
  - Animation

### Requirement: Performance Requirements

The dashboard designer SHALL perform smoothly with many components and large datasets.

#### Scenario: Designer loads dashboard with 50 components
- **WHEN** user opens a dashboard with 50 components
- **THEN** canvas loads within 2 seconds
- **AND** all components are rendered
- **AND** layer panel lists all components
- **AND** user can interact with components immediately

#### Scenario: Designer handles drag-drop with many components
- **WHEN** user drags a component on canvas with 50+ components
- **THEN** drag-drop is smooth (60 FPS)
- **AND** drop indicator updates in real-time
- **AND** no lag or freeze occurs

#### Scenario: Preview loads dashboard with real data
- **WHEN** user previews a dashboard with 50 components and 10 data sources
- **THEN** preview loads within 3 seconds
- **AND** all components display data
- **AND** data refresh completes within 2 seconds

#### Scenario: Dashboard handles real-time data refresh
- **WHEN** dashboard is in preview mode with 20+ components configured to refresh every 60 seconds
- **THEN** all components refresh on schedule
- **AND** refresh completes within 1 second per component
- **AND** no UI lag or freeze occurs
- **AND** refresh errors are displayed if any component fails

### Requirement: Designer Save-Load Runtime Consistency
The dashboard designer SHALL restore persisted component state consistently after save and reload.

#### Scenario: User reloads saved dashboard in designer
- **WHEN** user saves dashboard and later loads it in designer
- **THEN** component list, layer order, style settings, and data bindings are restored from persisted payload
- **THEN** designer does not rely on hardcoded mock component defaults for restoration

### Requirement: Preview Baseline Renders Persisted Components
Dashboard preview SHALL render a minimum runtime baseline based on persisted component payload rather than placeholder-only view.

#### Scenario: Preview renders baseline component output
- **WHEN** user enters preview mode after loading a saved dashboard
- **THEN** preview renders baseline output for supported component types using persisted data
- **THEN** preview refresh action updates the rendered state without switching to placeholder-only content

