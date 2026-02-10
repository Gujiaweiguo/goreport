## ADDED Requirements

### Requirement: Report Dataset Binding
The frontend and backend SHALL support binding report cells to dataset dimensions and measures.

#### Scenario: Bind cell to dataset dimension
- **WHEN** user selects a cell and chooses dataset as data source
- **THEN** user can select dataset from available datasets
- **THEN** user can select dimension from dataset's dimensions list
- **THEN** user can enable grouping (GROUP BY) for dimension
- **THEN** cell is marked as bound to dataset dimension
- **THEN** cell displays dimension's display name

#### Scenario: Bind cell to dataset measure
- **WHEN** user selects a cell and chooses dataset as data source
- **THEN** user can select measure from dataset's measures list
- **THEN** user can select aggregation function (SUM, AVG, COUNT, MAX, MIN, none)
- **THEN** cell is marked as bound to dataset measure
- **THEN** cell displays measure's display name

#### Scenario: Save report with dataset bindings
- **WHEN** user saves report with dataset-bound cells
- **THEN** system stores dataset ID in report configuration
- **THEN** system stores dimension/measure field names
- **THEN** system stores aggregation and grouping configuration
- **THEN** system does NOT store SQL queries in report configuration

#### Scenario: Render report with dataset data
- **WHEN** user previews or prints report with dataset bindings
- **THEN** backend loads dataset configuration from ID
- **THEN** backend queries dataset data using dataset's SQL or API
- **THEN** backend applies computed field expressions
- **THEN** backend applies grouping and aggregation from cell bindings
- **THEN** backend fills cells with query results
- **THEN** report is rendered with actual data

#### Scenario: Migrate existing report to dataset
- **WHEN** user opens report with direct datasource bindings
- **THEN** system provides option to create dataset from bindings
- **THEN** system offers to select existing dataset as replacement
- **THEN** user can choose to migrate report to use dataset
- **THEN** system converts datasource/field bindings to dataset/field bindings
