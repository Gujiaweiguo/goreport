## ADDED Requirements
### Requirement: Report Designer Access
The backend SHALL provide the report designer entry route `/jmreport/list` for users with access.

#### Scenario: Access report designer
- **WHEN** an authenticated user visits `/jmreport/list`
- **THEN** the report designer UI is available

### Requirement: Report Template CRUD
The backend SHALL provide APIs to create, update, list, and delete report templates used by the designer.

#### Scenario: Update report template
- **WHEN** a user submits an update to an existing report template
- **THEN** the template is persisted and subsequent views reflect the change
