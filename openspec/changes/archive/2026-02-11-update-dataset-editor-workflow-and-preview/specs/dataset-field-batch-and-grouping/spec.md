## ADDED Requirements

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
