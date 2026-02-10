## ADDED Requirements

### Requirement: Dataset CRUD
The backend SHALL allow authorized users to create, update, list, and delete dataset definitions.

#### Scenario: Create SQL-based dataset
- **WHEN** a user submits a dataset definition with SQL query type
- **THEN** the dataset is stored with SQL query and datasource ID
- **THEN** the dataset becomes selectable for reports, dashboards, and charts
- **THEN** the dataset's fields are automatically extracted from SQL query results

#### Scenario: Create API-based dataset
- **WHEN** a user submits a dataset definition with API data source type
- **THEN** the dataset is stored with API endpoint URL and HTTP method
- **THEN** the dataset's fields are automatically extracted from API response
- **THEN** the dataset becomes selectable for reports, dashboards, and charts

#### Scenario: Create file-import dataset
- **WHEN** a user submits a dataset definition with file import (Excel, CSV)
- **THEN** the dataset is stored with file reference
- **THEN** the dataset's fields are automatically extracted from file headers
- **THEN** the dataset becomes selectable for reports, dashboards, and charts

#### Scenario: List datasets
- **WHEN** a user requests dataset list
- **THEN** system returns all datasets belonging to user's tenant
- **THEN** each dataset includes name, type, datasource name, field count, and timestamps

#### Scenario: Update dataset
- **WHEN** a user updates their dataset
- **THEN** changes are persisted
- **THEN** dataset fields are re-extracted if query or source changed
- **THEN** updated dataset is returned

#### Scenario: Delete dataset
- **WHEN** a user deletes their dataset
- **THEN** dataset is soft deleted
- **THEN** dataset no longer appears in list
- **THEN** reports, dashboards, and charts using this dataset show reference error

#### Scenario: Dataset preview
- **WHEN** a user requests dataset preview
- **THEN** system executes the dataset's query or fetches data from source
- **THEN** system returns first 100 rows of data
- **THEN** system returns column names and data types

### Requirement: Dataset Field Configuration
The frontend and backend SHALL support configuring dataset fields including name, type (dimension/measure), data type, sort order, and grouping.

#### Scenario: Configure field as dimension
- **WHEN** user selects a dataset field
- **THEN** user can set field type to "dimension"
- **THEN** dimension is used for categorical data (grouping, filtering)
- **THEN** dimension appears in dimension list for visualization components

#### Scenario: Configure field as measure
- **WHEN** user selects a dataset field
- **THEN** user can set field type to "measure"
- **THEN** measure is used for numerical data (aggregation, calculation)
- **THEN** measure appears in measure list for visualization components

#### Scenario: Set field data type
- **WHEN** user selects a dataset field
- **THEN** user can set field data type (string, number, date, boolean)
- **THEN** field data type is used for type validation and formatting
- **THEN** visualization components use correct type for rendering

#### Scenario: Configure field sort order
- **WHEN** user selects a dataset field
- **THEN** user can set default sort order (asc, desc, none)
- **THEN** sort order is applied when dataset is queried
- **THEN** user can override sort order in visualization components

#### Scenario: Configure field grouping
- **WHEN** user selects a dimension field
- **THEN** user can enable grouping
- **THEN** data is grouped by this field when used in charts
- **THEN** measures are aggregated per group

#### Scenario: Configure field display name
- **WHEN** user selects a dataset field
- **THEN** user can set display name (alias)
- **THEN** display name is shown in visualization component UI instead of original field name
- **THEN** original field name is preserved in queries

### Requirement: Dataset Tenant Isolation
The system SHALL ensure dataset access is scoped by tenant.

#### Scenario: Cross-tenant dataset protection
- **WHEN** a user attempts to access a dataset from another tenant
- **THEN** request is rejected with 403 Forbidden

#### Scenario: Query by tenant
- **WHEN** a user queries dataset list
- **THEN** only datasets belonging to user's tenant are returned

#### Scenario: Tenant isolation in data queries
- **WHEN** a user queries data from a dataset
- **THEN** system enforces tenant isolation on the underlying datasource
- **THEN** user cannot access data from other tenants
