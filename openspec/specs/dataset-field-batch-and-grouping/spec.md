# dataset-field-batch-and-grouping Specification

## Purpose
TBD - created by archiving change update-dataset-editor-workflow. Update Purpose after archive.
## Requirements
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

### Requirement: Batch Update Interaction Reliability
The dataset editor SHALL provide reliable interaction feedback for batch field updates.

#### Scenario: Batch submit without selection
- **WHEN** user triggers batch update with no selected fields
- **THEN** system blocks request submission
- **THEN** system shows clear guidance to select fields first

#### Scenario: Batch update partial failure feedback
- **WHEN** batch update response contains both successes and field-level failures
- **THEN** system presents failed field details without dropping successful updates
- **THEN** system refreshes field metadata to match persisted backend state

### Requirement: Batch Field Update Response Contract
The backend SHALL provide field-level response details for dataset batch field updates.

#### Scenario: Return per-field update result
- **WHEN** batch update request is processed
- **THEN** response includes a list of successfully updated field ids
- **THEN** response includes per-field error objects for rejected updates

#### Scenario: Keep unaffected field attributes stable
- **WHEN** batch update payload omits a field attribute
- **THEN** system keeps existing attribute value unchanged
- **THEN** response reflects final persisted state without implicit resets

