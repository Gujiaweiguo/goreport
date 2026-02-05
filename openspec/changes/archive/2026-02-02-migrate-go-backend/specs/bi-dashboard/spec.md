## ADDED Requirements
### Requirement: Dashboard Designer Access
The backend SHALL provide the dashboard designer entry route `/drag/list` for authorized users.

#### Scenario: Access dashboard designer
- **WHEN** an authenticated user visits `/drag/list`
- **THEN** the dashboard designer UI is available

### Requirement: Dashboard CRUD
The backend SHALL provide APIs to create, update, list, and delete dashboard definitions.

#### Scenario: Update dashboard
- **WHEN** a user updates an existing dashboard
- **THEN** the updated dashboard is persisted and available for preview
