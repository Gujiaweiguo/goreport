## ADDED Requirements

### Requirement: Batch Field Management
The dataset editor SHALL support batch updates for field configurations, including field type, display metadata, and ordering attributes.

#### Scenario: Batch update multiple fields
- **WHEN** user selects multiple fields in Batch Management and submits updates
- **THEN** system validates each update against field rules
- **THEN** system applies valid changes in one batch operation
- **THEN** system returns updated field states for dimensions and measures

#### Scenario: Batch update with invalid field rule
- **WHEN** batch request contains invalid field configuration
- **THEN** system rejects invalid updates with field-level error details
- **THEN** system does not silently overwrite unaffected field attributes
- **THEN** user can correct and retry without reloading page

### Requirement: Grouping Field Creation and Behavior
The dataset editor SHALL provide a dedicated entry to create grouping fields and apply grouping behavior through dimension-compatible configuration.

#### Scenario: Create grouping field from editor
- **WHEN** user clicks New Grouping Field and submits valid configuration
- **THEN** system creates a grouping-capable field configuration
- **THEN** field appears in dimension-side field list with grouping metadata
- **THEN** grouping field is available for downstream query and visualization bindings

#### Scenario: Use grouping field in query path
- **WHEN** dataset query is executed with grouping-enabled field
- **THEN** system applies grouping semantics before measure aggregation
- **THEN** grouped result set is returned with consistent field aliases

### Requirement: Grouping and Batch Validation Contract
The backend SHALL provide validation and response contracts for batch and grouping operations that are safe for non-breaking frontend evolution.

#### Scenario: Validate grouping configuration
- **WHEN** user submits grouping rule with unsupported expression or type
- **THEN** system returns validation errors with actionable messages
- **THEN** system does not persist invalid grouping configuration

#### Scenario: Backward compatibility for existing datasets
- **WHEN** existing dataset without batch/grouping metadata is loaded
- **THEN** system returns default-compatible field metadata
- **THEN** existing save and preview behavior remains functional without migration blockers
