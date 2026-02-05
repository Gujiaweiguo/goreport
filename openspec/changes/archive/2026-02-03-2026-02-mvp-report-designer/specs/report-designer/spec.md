## ADDED Requirements

### Requirement: Report Data Binding

The frontend SHALL allow users to bind datasource, table, and field to report cells, and the backend SHALL render report with actual data from the database.

#### Scenario: User creates data-bound report

- **WHEN** user selects a cell in report designer
- **THEN** user can choose a datasource from available datasources
- **THEN** user can select a database table from the chosen datasource
- **THEN** user can select a field from the chosen table
- **WHEN** user saves the report with data bindings
- **THEN** report configuration includes datasource ID, table name, and field name
- **WHEN** user views the report preview
- **THEN** backend queries the datasource for data
- **THEN** backend renders the report with actual data from database
- **THEN** frontend displays the report with data-filled cells

#### Scenario: Report preview with parameters

- **WHEN** user opens report preview
- **THEN** system loads the report configuration
- **THEN** backend extracts all data bindings from the report configuration
- **WHEN** data bindings exist
- **THEN** backend queries the database for each unique datasource and table combination
- **THEN** backend passes the data to the rendering engine
- **THEN** backend generates HTML with data
- **THEN** frontend displays the rendered report

#### Scenario: Report with no data binding

- **WHEN** user creates a report without data binding
- **THEN** cells display static text values
- **WHEN** user saves report
- **THEN** static values are persisted
- **WHEN** user previews report
- **THEN** report displays with static text values

---

### Requirement: Datasource Metadata API

The backend SHALL provide APIs to retrieve datasource metadata including tables and fields, enabling the report designer to present data binding options.

#### Scenario: Load available datasources

- **WHEN** user opens report designer
- **THEN** system loads available datasources from `/datasource/list`
- **THEN** user can see datasource names, types, and connection status

#### Scenario: Load datasource tables

- **WHEN** user selects a datasource in property panel
- **THEN** user can request tables from the datasource
- **THEN** backend queries the datasource's database for available tables
- **THEN** backend returns list of table names
- **THEN** user can select a table from the dropdown

#### Scenario: Load table fields

- **WHEN** user selects a table
- **THEN** user can request fields from the table
- **THEN** backend queries the database schema for the table
- **THEN** backend returns list of field names with their types
- **THEN** user can select a field from the dropdown

---

### Requirement: Cell Data Binding Configuration

The frontend property panel SHALL provide UI for configuring cell data binding, including datasource, table, field, and optional aggregate function.

#### Scenario: Configure text cell data binding

- **WHEN** user clicks a cell in report designer
- **THEN** property panel displays cell properties
- **THEN** property panel includes "Data Binding" section
- **THEN** user can select datasource from dropdown
- **THEN** user can select table from dropdown (disabled until datasource is selected)
- **THEN** user can select field from dropdown (disabled until table is selected)
- **THEN** user can optionally select aggregate function (sum, avg, count, max, min, none)
- **THEN** user can enter display label for the binding
- **THEN** configuration is saved when user saves cell or selects another cell

#### Scenario: Clear data binding

- **WHEN** user wants to remove data binding from a cell
- **THEN** user can click "Clear Data Binding" button in property panel
- **THEN** all data binding fields are reset to empty
- **THEN** cell type changes from `bound` back to `text`
- **THEN** cell displays static text value instead of data binding indicator

---

### Requirement: Report Preview with Data Rendering

The frontend SHALL provide a report preview page that calls the backend render API and displays the rendered HTML with actual data.

#### Scenario: Preview report with data

- **WHEN** user clicks "Preview" button in report designer
- **THEN** user is navigated to preview page
- **THEN** preview page loads report configuration from backend
- **THEN** preview page calls `/api/v1/jmreport/preview` API with report ID
- **THEN** backend renders the report with data from database
- **THEN** backend returns HTML with data-filled cells
- **THEN** frontend displays the rendered HTML in preview container
- **THEN** user can see actual data from database in the report

#### Scenario: Refresh preview data

- **WHEN** user clicks "Refresh" button in preview page
- **THEN** system re-queries the database for latest data
- **THEN** preview page re-renders the report
- **THEN** updated data is displayed to user

---

### Requirement: Report Persistence with Data Bindings

The frontend and backend SHALL support saving and loading report configurations that include cell data bindings.

#### Scenario: Save report with data bindings

- **WHEN** user clicks "Save" button in report designer
- **THEN** system collects all cell configurations including data bindings
- **THEN** system constructs report JSON configuration with cells containing datasource, tableName, and fieldName
- **THEN** system sends configuration to `/api/v1/jmreport/create` or `/api/v1/jmreport/update`
- **THEN** backend persists the configuration in `jimu_report` table's `json_str` column
- **THEN** report is saved with data binding information

#### Scenario: Load report with data bindings

- **WHEN** user opens an existing report from the report list
- **THEN** system loads the report configuration
- **THEN** system parses the configuration including cell data bindings
- **THEN** report designer populates the canvas with cells
- **THEN** cells with data bindings display the datasource, table, and field information
- **THEN** property panel shows current data binding configuration
- **THEN** user can edit the report or preview it
