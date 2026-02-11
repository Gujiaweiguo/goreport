## ADDED Requirements

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
