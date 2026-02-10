## ADDED Requirements

### Requirement: Dataset CRUD API
The backend SHALL provide REST APIs for dataset CRUD operations at `/api/v1/datasets`.

#### Scenario: Create dataset API
- **WHEN** authorized user POSTs to `/api/v1/datasets` with dataset definition
- **THEN** system validates datasource ID belongs to user's tenant
- **THEN** system validates SQL query or API configuration
- **THEN** system creates dataset record
- **THEN** system extracts fields from query/API response
- **THEN** system returns created dataset with field list
- **THEN** status code 201 is returned

#### Scenario: List datasets API
- **WHEN** authorized user GETs `/api/v1/datasets`
- **THEN** system returns all datasets for user's tenant
- **THEN** each dataset includes id, name, type, datasourceName, fieldCount, createdAt, updatedAt
- **THEN** results are paginated
- **THEN** status code 200 is returned

#### Scenario: Get dataset API
- **WHEN** authorized user GETs `/api/v1/datasets/{id}`
- **THEN** system returns dataset with full details
- **THEN** response includes datasource, fields, configuration
- **THEN** status code 200 is returned
- **WHEN** dataset does not exist or belongs to other tenant
- **THEN** status code 404 is returned

#### Scenario: Update dataset API
- **WHEN** authorized user PUTs `/api/v1/datasets/{id}` with updated dataset
- **THEN** system validates dataset belongs to user's tenant
- **THEN** system updates dataset record
- **THEN** system re-extracts fields if query changed
- **THEN** system returns updated dataset
- **THEN** status code 200 is returned

#### Scenario: Delete dataset API
- **WHEN** authorized user DELETEs `/api/v1/datasets/{id}`
- **THEN** system validates dataset belongs to user's tenant
- **THEN** system soft deletes dataset record
- **THEN** system returns success message
- **THEN** status code 200 is returned

#### Scenario: Delete dataset with references
- **WHEN** dataset is used by reports, dashboards, or charts
- **THEN** system returns 409 Conflict with list of references
- **THEN** system provides option to cascade delete or abort

### Requirement: Dataset Query API
The backend SHALL provide API to query dataset data at `/api/v1/datasets/{id}/data`.

#### Scenario: Query dataset data API
- **WHEN** authorized user POSTs to `/api/v1/datasets/{id}/data` with query parameters
- **THEN** system validates dataset belongs to user's tenant
- **THEN** system executes SQL query or fetches API data
- **THEN** system applies computed field expressions
- **THEN** system applies filters, sorting, and pagination
- **THEN** system returns data array and metadata (total count, execution time)
- **THEN** status code 200 is returned

#### Scenario: Query with filters
- **WHEN** request includes filters array
- **THEN** system applies each filter to query
- **THEN** filters support operators (eq, neq, gt, gte, lt, lte, like, in)
- **THEN** filters work on both original fields and computed fields
- **THEN** filtered results are returned

#### Scenario: Query with sorting
- **WHEN** request includes sortBy and sortOrder parameters
- **THEN** system applies sorting to query
- **THEN** sorting can be applied to multiple fields
- **THEN** sorting works on both original fields and computed fields
- **THEN** sorted results are returned

#### Scenario: Query with pagination
- **WHEN** request includes page and pageSize parameters
- **THEN** system applies LIMIT and OFFSET to query
- **THEN** system returns data array for requested page
- **THEN** system returns total count for pagination controls
- **THEN** paginated results are returned

#### Scenario: Query with field selection
- **WHEN** request includes fields parameter
- **THEN** system selects only specified fields in query
- **THEN** results include only selected fields
- **THEN** performance is improved for wide tables

### Requirement: Dataset Preview API
The backend SHALL provide API to preview dataset sample data at `/api/v1/datasets/{id}/preview`.

#### Scenario: Preview dataset API
- **WHEN** authorized user GETs `/api/v1/datasets/{id}/preview`
- **THEN** system validates dataset belongs to user's tenant
- **THEN** system executes query with LIMIT 100
- **THEN** system returns first 100 rows
- **THEN** system returns column names and data types
- **THEN** system returns query execution time
- **THEN** status code 200 is returned

#### Scenario: Preview with filters
- **WHEN** request includes filters parameter
- **THEN** system applies filters before LIMIT
- **THEN** system returns filtered sample data
- **THEN** user can preview filtered data before saving

### Requirement: Dataset Field Management API
The backend SHALL provide APIs to manage dataset fields at `/api/v1/datasets/{id}/fields`.

#### Scenario: List dataset fields API
- **WHEN** authorized user GETs `/api/v1/datasets/{id}/fields`
- **THEN** system returns all fields for dataset
- **THEN** each field includes name, displayName, type (dimension/measure), dataType, isComputed, isGroupable, isSortable
- **THEN** fields are grouped by dimension and measure
- **THEN** status code 200 is returned

#### Scenario: Create computed field API
- **WHEN** authorized user POSTs to `/api/v1/datasets/{id}/fields` with field definition
- **THEN** system validates field expression syntax
- **THEN** system creates computed field record
- **THEN** system marks field as computed field
- **THEN** system returns created field
- **THEN** status code 201 is returned

#### Scenario: Update field configuration API
- **WHEN** authorized user PUTs `/api/v1/datasets/{id}/fields/{fieldId}` with updated field
- **THEN** system updates field configuration
- **THEN** system allows changing displayName, type, dataType, sort order, grouping
- **THEN** system does not allow changing original field name
- **THEN** system returns updated field
- **THEN** status code 200 is returned

#### Scenario: Delete computed field API
- **WHEN** authorized user DELETEs `/api/v1/datasets/{id}/fields/{fieldId}`
- **THEN** system validates field is computed field
- **THEN** system does not allow deleting original fields
- **THEN** system deletes computed field record
- **THEN** system returns success message
- **THEN** status code 200 is returned

### Requirement: Dataset Validation API
The backend SHALL provide API to validate dataset configuration at `/api/v1/datasets/validate`.

#### Scenario: Validate SQL query API
- **WHEN** authorized user POSTs to `/api/v1/datasets/validate` with SQL query and datasource ID
- **THEN** system validates datasource belongs to user's tenant
- **THEN** system executes EXPLAIN query
- **THEN** system returns validation result (valid/invalid)
- **THEN** system returns estimated execution time
- **THEN** system returns error message if invalid
- **THEN** status code 200 is returned

#### Scenario: Validate API endpoint API
- **WHEN** authorized user POSTs to `/api/v1/datasets/validate` with API configuration
- **THEN** system sends test request to API endpoint
- **THEN** system validates response structure
- **THEN** system returns validation result
- **THEN** system returns sample response data
- **THEN** status code 200 is returned

#### Scenario: Validate computed field expression API
- **WHEN** authorized user POSTs to `/api/v1/datasets/validate` with expression and dataset ID
- **THEN** system validates expression syntax
- **THEN** system validates field references
- **THEN** system validates function usage
- **THEN** system returns validation result with errors if any
- **THEN** status code 200 is returned
