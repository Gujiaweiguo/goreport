# dataset-integration Specification

## Purpose
TBD - created by archiving change add-dataset-feature. Update Purpose after archive.
## Requirements
### Requirement: Dataset Metadata API for Visualization Components
The backend SHALL provide APIs for visualization components (report designer, dashboard, chart editor) to query dataset metadata including dimensions and measures.

#### Scenario: Get dataset dimensions API
- **WHEN** visualization component requests `/api/v1/datasets/{id}/dimensions`
- **THEN** system returns list of dimension fields
- **THEN** each dimension includes name, displayName, dataType, isGroupable
- **THEN** system includes computed fields marked as dimensions
- **THEN** status code 200 is returned

#### Scenario: Get dataset measures API
- **WHEN** visualization component requests `/api/v1/datasets/{id}/measures`
- **THEN** system returns list of measure fields
- **THEN** each measure includes name, displayName, dataType, isAggregateable
- **THEN** system includes computed fields marked as measures
- **THEN** status code 200 is returned

#### Scenario: Get dataset schema API
- **WHEN** visualization component requests `/api/v1/datasets/{id}/schema`
- **THEN** system returns full dataset schema
- **THEN** response includes datasource info, dimensions list, measures list
- **THEN** response includes computed fields with their expressions
- **THEN** status code 200 is returned

### Requirement: Report Designer Dataset Integration
The report designer SHALL support binding report cells to dataset dimensions and measures instead of direct datasource/SQL.

#### Scenario: Select dataset in report designer
- **WHEN** user opens report designer
- **THEN** user can select a dataset from dropdown instead of datasource
- **THEN** dropdown shows all datasets in user's tenant
- **THEN** dataset selection populates dimensions and measures lists
- **THEN** dimensions and measures appear in data binding panel

#### Scenario: Bind cell to dataset dimension
- **WHEN** user selects a cell in report designer
- **THEN** user can select a dataset dimension from dropdown
- **THEN** user can optionally configure grouping (GROUP BY)
- **THEN** cell binding includes dataset ID and dimension field name
- **THEN** dimension field is displayed with display name
- **THEN** cell type changes to `dataset-dimension`

#### Scenario: Bind cell to dataset measure
- **WHEN** user selects a cell in report designer
- **THEN** user can select a dataset measure from dropdown
- **THEN** user can optionally select aggregation function (SUM, AVG, COUNT, MAX, MIN, none)
- **THEN** cell binding includes dataset ID and measure field name
- **THEN** measure field is displayed with display name
- **THEN** cell type changes to `dataset-measure`

#### Scenario: Bind computed field
- **WHEN** user selects a cell in report designer
- **THEN** user can select a computed field (dimension or measure) from dropdown
- **THEN** computed field is displayed with its display name
- **THEN** cell binding includes computed field name
- **THEN** computed field values are calculated when report is rendered

#### Scenario: Save report with dataset bindings
- **WHEN** user saves report with dataset bindings
- **THEN** system stores report configuration with dataset IDs and field names
- **THEN** system does NOT store SQL queries
- **THEN** system stores aggregation and grouping configuration
- **THEN** report configuration references datasets by ID

#### Scenario: Render report with dataset data
- **WHEN** user previews or renders report with dataset bindings
- **THEN** backend loads dataset configurations
- **THEN** backend constructs SQL query from dataset SQL and field bindings
- **THEN** backend applies computed field expressions
- **THEN** backend applies grouping and aggregation
- **THEN** backend executes query and fills cells with data
- **THEN** report is rendered with actual data

#### Scenario: Migrate existing report to dataset
- **WHEN** user opens existing report with direct SQL bindings
- **THEN** system offers migration option to convert to dataset
- **THEN** system extracts datasource, table, and field info
- **THEN** system prompts user to create dataset or select existing
- **THEN** user can migrate report to use dataset

### Requirement: Dashboard Dataset Integration
The dashboard designer SHALL support binding dashboard widgets to dataset dimensions and measures.

#### Scenario: Select dataset in dashboard designer
- **WHEN** user adds a widget to dashboard
- **THEN** user can select a dataset from dropdown
- **THEN** dropdown shows all datasets in user's tenant
- **THEN** dataset selection populates dimensions and measures lists

#### Scenario: Bind chart to dataset dimensions
- **WHEN** user configures chart X-axis or category dimension
- **THEN** user can select dataset dimension from dropdown
- **THEN** user can optionally configure grouping
- **THEN** dimension is displayed with display name
- **THEN** chart configuration references dataset ID and dimension name

#### Scenario: Bind chart to dataset measures
- **WHEN** user configures chart Y-axis or value measure
- **THEN** user can select dataset measure from dropdown
- **THEN** user can select aggregation function
- **THEN** user can select multiple measures for multi-series charts
- **THEN** measures are displayed with display names
- **THEN** chart configuration references dataset ID and measure names

#### Scenario: Save dashboard with dataset bindings
- **WHEN** user saves dashboard with dataset bindings
- **THEN** system stores widget configuration with dataset IDs and field names
- **THEN** widget configuration references datasets by ID
- **THEN** system stores aggregation and grouping settings

#### Scenario: Render dashboard with dataset data
- **WHEN** dashboard is rendered
- **THEN** backend loads dataset configurations for each widget
- **THEN** backend queries dataset data for each widget
- **THEN** backend applies computed fields, grouping, aggregation
- **THEN** backend returns data to frontend
- **THEN** widgets are rendered with actual data

### Requirement: Chart Editor Dataset Integration
The chart editor SHALL support using datasets as data source for charts.

#### Scenario: Select dataset in chart editor
- **WHEN** user opens chart editor
- **THEN** user can select dataset as data source type
- **THEN** user can select specific dataset from dropdown
- **THEN** dataset selection shows available dimensions and measures
- **THEN** dimensions and measures are organized in separate lists

#### Scenario: Map chart dimensions to dataset dimensions
- **WHEN** user configures chart X-axis or category
- **THEN** user selects dataset dimension from list
- **THEN** dimension display name is shown in configuration
- **THEN** dimension type (string, date, number) is displayed
- **THEN** user can enable grouping if dimension is groupable

#### Scenario: Map chart measures to dataset measures
- **WHEN** user configures chart Y-axis or value
- **THEN** user selects dataset measure from list
- **THEN** measure display name is shown in configuration
- **THEN** user can select aggregation function
- **THEN** user can add multiple measures for series
- **THEN** each series references a dataset measure

#### Scenario: Preview chart with dataset data
- **WHEN** user previews chart in chart editor
- **THEN** backend queries dataset data
- **THEN** backend applies computed fields, grouping, aggregation
- **THEN** backend returns data in chart-expected format
- **THEN** frontend renders chart with actual data
- **THEN** chart updates in real-time as configuration changes

#### Scenario: Save chart configuration with dataset reference
- **WHEN** user saves chart configuration
- **THEN** system stores dataset ID in chart config
- **THEN** system stores dimension and measure references
- **THEN** system stores aggregation and grouping settings
- **THEN** chart configuration can be reused with updated dataset data

#### Scenario: Handle dataset schema changes
- **WHEN** dataset schema is updated (field added/removed/renamed)
- **THEN** chart configurations using dataset are validated
- **THEN** charts using removed fields show error in editor
- **THEN** charts using renamed fields are updated to new field names
- **THEN** system logs schema change notifications for chart owners

