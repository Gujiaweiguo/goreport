## ADDED Requirements

### Requirement: Datasource Organizational Operations
The system SHALL provide datasource copy, move, and rename operations under tenant-scoped authorization.

#### Scenario: Copy datasource
- **WHEN** authorized user triggers datasource copy within the same tenant
- **THEN** system creates a new datasource entry with duplicated connection settings
- **THEN** copied datasource receives a new unique identifier

#### Scenario: Move datasource
- **WHEN** authorized user moves datasource to another folder/location namespace
- **THEN** system updates datasource location metadata
- **THEN** datasource remains accessible to the same tenant scope

#### Scenario: Rename datasource
- **WHEN** authorized user renames datasource with valid target name
- **THEN** system updates datasource display name
- **THEN** update does not break existing references by datasource id

### Requirement: Datasource Search Capability
The system SHALL support keyword search for datasource listings.

#### Scenario: Search datasource by keyword
- **WHEN** user queries datasource list with search keyword
- **THEN** system returns tenant-scoped datasource results matching name or configured searchable fields

#### Scenario: Empty search keyword
- **WHEN** user queries datasource list without keyword
- **THEN** system returns default tenant-scoped datasource list behavior

### Requirement: Datasource Metadata Description Retrieval
The system SHALL return table and field description metadata when datasource supports description/comment introspection.

#### Scenario: Retrieve table descriptions
- **WHEN** user requests table metadata from a supported datasource
- **THEN** system includes table description/comment fields in response when available

#### Scenario: Retrieve field descriptions
- **WHEN** user requests field metadata from a supported datasource table
- **THEN** system includes field description/comment fields in response when available
