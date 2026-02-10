## ADDED Requirements

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
