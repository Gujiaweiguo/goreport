## ADDED Requirements
### Requirement: Datasource Creation Wizard Workflow
The frontend datasource creation flow SHALL provide a two-step wizard that separates source-type selection from connection parameter configuration.

#### Scenario: Enter create wizard from datasource page
- **WHEN** user starts "创建数据源"
- **THEN** system opens step 1 for datasource type selection
- **THEN** wizard shows categories OLTP, OLAP, 数据湖, API数据, 文件

#### Scenario: Continue from type selection to configuration
- **WHEN** user selects a supported datasource template and clicks "下一步"
- **THEN** system opens step 2 with connection form fields for that datasource type
- **THEN** selected datasource type is reflected in form type value

#### Scenario: Unsupported template in wizard
- **WHEN** user selects a template marked as upcoming
- **THEN** system shows that template as "即将支持"
- **THEN** wizard prevents advancing to configuration for unsupported template

### Requirement: Datasource Create/Edit Connection Test Behavior
The system SHALL support deterministic connection testing in create and edit dialogs with credential-source aware behavior.

#### Scenario: Test connection in create dialog
- **WHEN** user clicks "测试连接" during create flow
- **THEN** system tests connectivity using current form input values
- **THEN** system returns success or actionable error message in dialog context

#### Scenario: Test connection in edit dialog without password change
- **WHEN** user edits datasource and does not provide a new password
- **THEN** system tests connectivity using saved datasource credentials
- **THEN** system does not require user to re-enter persisted password

#### Scenario: Test connection in edit dialog with password change
- **WHEN** user edits datasource and provides a new password
- **THEN** system tests connectivity using the new password from form input
- **THEN** result reflects the new credentials validity

### Requirement: Datasource Password Preservation on Edit
The system SHALL never expose persisted datasource password in plaintext to frontend and SHALL preserve stored password when edit request omits password.

#### Scenario: Open edit dialog for existing datasource
- **WHEN** user opens datasource edit dialog
- **THEN** password field is masked/placeholder state only
- **THEN** plaintext persisted password is not returned to client

#### Scenario: Save edit without changing password
- **WHEN** user submits datasource edit without entering new password
- **THEN** update request omits password field
- **THEN** backend keeps existing stored password unchanged

#### Scenario: Save edit with new password
- **WHEN** user submits datasource edit with a new password value
- **THEN** backend updates stored password to the new value
- **THEN** subsequent connection tests use updated password

### Requirement: Datasource Update Binding Robustness
Datasource update request binding SHALL accept handler-populated identity fields without requiring them from request body.

#### Scenario: Update request without id and tenantId in JSON body
- **WHEN** user submits update payload that does not include `id` and `tenantId`
- **THEN** handler fills `id` from route path and `tenantId` from auth context
- **THEN** request proceeds to service layer without generic "invalid request" binding failure

#### Scenario: Tenant context missing during update
- **WHEN** authenticated tenant context is absent
- **THEN** update request is rejected with tenant-specific authorization error

### Requirement: Datasource Create/Edit API Contract Consistency
The datasource UI workflow SHALL rely on stable API contracts for create test, edit test, and update operations.

#### Scenario: Create-flow test connection contract
- **WHEN** create dialog invokes connection test
- **THEN** system sends `POST /api/v1/datasources/test` with datasource form fields (`name`, `type`, `host`, `port`, `database`, `username`, optional `password`, optional `advanced`)
- **THEN** backend returns success flag and diagnostic message

#### Scenario: Edit-flow test connection without password change contract
- **WHEN** edit dialog invokes connection test and user does not provide a new password
- **THEN** system sends `POST /api/v1/datasources/:id/test`
- **THEN** backend validates using persisted datasource credentials

#### Scenario: Update contract with optional password
- **WHEN** edit dialog submits `PUT /api/v1/datasources/:id`
- **THEN** request body may omit `password` and still succeeds
- **THEN** backend updates non-password fields and keeps existing password when omitted

### Requirement: Datasource Wizard Data Model
The frontend datasource wizard SHALL maintain minimal state required for deterministic create flow and edit behavior.

#### Scenario: Wizard state model
- **WHEN** user enters create datasource flow
- **THEN** frontend keeps explicit step state (`step 1` select source, `step 2` configure)
- **THEN** frontend tracks selected category and selected template

#### Scenario: Template-to-backend type mapping
- **WHEN** user selects a supported template
- **THEN** frontend maps template to backend datasource type value supported by API contract
- **THEN** frontend applies template default connection values when defined (for example default port)

#### Scenario: Password state model in edit flow
- **WHEN** user opens edit dialog
- **THEN** frontend displays masked password placeholder state instead of plaintext password
- **THEN** frontend distinguishes unchanged-password state from user-entered new password state

### Requirement: Datasource Creation MVP Acceptance
Datasource management SHALL provide a minimum shippable create/edit workflow aligned with category-first source selection.

#### Scenario: Category coverage in create step
- **WHEN** user opens create step 1
- **THEN** UI presents all five category entries: OLTP, OLAP, 数据湖, API数据, 文件
- **THEN** selecting each category updates the template list content

#### Scenario: Supported template can be completed end-to-end
- **WHEN** user chooses a supported template and provides valid configuration
- **THEN** user can proceed to step 2, test connection, and create datasource successfully

#### Scenario: Upcoming template is visible but gated
- **WHEN** user chooses a template marked upcoming
- **THEN** UI clearly indicates "即将支持"
- **THEN** flow prevents create submission for that template
