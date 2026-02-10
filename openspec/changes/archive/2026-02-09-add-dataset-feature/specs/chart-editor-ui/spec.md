## ADDED Requirements

### Requirement: Chart Dataset Data Source
The chart editor SHALL support using datasets as data source for charts.

#### Scenario: Select dataset as chart data source
- **WHEN** user opens chart editor data source configuration
- **THEN** user can select "Dataset" as data source type
- **THEN** user can choose specific dataset from dropdown
- **THEN** dataset selection shows available dimensions and measures
- **THEN** dimensions and measures are organized separately

#### Scenario: Map chart dimensions to dataset dimensions
- **WHEN** user configures chart dimension (X-axis, category, series)
- **THEN** user selects dimension from dataset's dimensions list
- **THEN** dimension display name is shown in UI
- **THEN** dimension data type is displayed
- **THEN** user can enable grouping if dimension is groupable
- **THEN** chart configuration references dataset ID and dimension field name

#### Scenario: Map chart measures to dataset measures
- **WHEN** user configures chart measure (Y-axis, value, size, angle)
- **THEN** user selects measure from dataset's measures list
- **THEN** measure display name is shown in UI
- **THEN** user can select aggregation function
- **THEN** user can select multiple measures for multi-series charts
- **THEN** chart configuration references dataset ID and measure field names

#### Scenario: Preview chart with dataset data
- **WHEN** user previews chart or modifies configuration
- **THEN** backend queries dataset data
- **THEN** backend applies computed field expressions
- **THEN** backend applies grouping and aggregation
- **THEN** backend returns data in chart-expected format
- **THEN** chart renders with actual data
- **THEN** preview updates in real-time

#### Scenario: Save chart configuration with dataset
- **WHEN** user saves chart configuration
- **THEN** system stores dataset ID in chart config
- **THEN** system stores dimension and measure field references
- **THEN** system stores aggregation and grouping settings
- **THEN** chart can be reused with updated dataset data

#### Scenario: Use computed fields in chart
- **WHEN** dataset contains computed fields
- **THEN** computed fields appear in dimensions or measures list
- **THEN** computed field display name is shown
- **THEN** user can select computed field for chart axis
- **THEN** computed field expression is evaluated when chart data is queried
- **THEN** chart displays computed field values

#### Scenario: Validate chart-dataset compatibility
- **WHEN** user selects dataset for chart type
- **THEN** system validates dataset has required fields for chart type
- **THEN** system warns if dataset lacks dimensions or measures
- **THEN** system suggests compatible chart types based on available fields

#### Scenario: Handle dataset schema changes
- **WHEN** dataset schema is updated (field added/removed/renamed)
- **THEN** charts using dataset are validated
- **THEN** charts using removed fields show error in editor
- **THEN** charts using renamed fields are updated to new field names
- **THEN** system logs schema change notifications for chart owners
