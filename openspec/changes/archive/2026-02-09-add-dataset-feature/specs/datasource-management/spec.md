## ADDED Requirements

### Requirement: Dataset Metadata Query
The backend SHALL provide metadata query capabilities for datasets used in visualizations.

#### Scenario: Query dataset tables and fields
- **WHEN** a user requests metadata for a dataset
- **THEN** system returns dataset ID, name, type
- **THEN** system returns list of dimensions (name, displayName, dataType)
- **THEN** system returns list of measures (name, displayName, dataType)
- **THEN** system indicates if dimension is groupable
- **THEN** system indicates if measure is aggregateable

#### Scenario: Query dataset computed fields
- **WHEN** a user requests metadata for a dataset with computed fields
- **THEN** system returns computed fields separately
- **THEN** each computed field includes expression
- **THEN** each computed field indicates if it's dimension or measure

#### Scenario: Cross-validate datasource and dataset access
- **WHEN** a user queries metadata for a dataset
- **THEN** system validates user has access to both dataset and underlying datasource
- **THEN** request is rejected if user lacks datasource access
