# dataset-computed-fields Specification

## Purpose
TBD - created by archiving change add-dataset-feature. Update Purpose after archive.
## Requirements
### Requirement: Computed Field Creation
The frontend and backend SHALL support creating computed fields with expressions that reference dataset fields, functions, and arithmetic operators.

#### Scenario: Create computed field with field references
- **WHEN** user creates a new field in dataset field editor
- **THEN** user can click existing dataset fields to reference them
- **THEN** referenced fields are added to field expression
- **THEN** expression editor shows field names wrapped in brackets [field_name]

#### Scenario: Create computed field with arithmetic operators
- **WHEN** user creates a new field in dataset field editor
- **THEN** user can add arithmetic operators (+, -, *, /) to expression
- **THEN** operators can be typed manually or added via operator buttons
- **THEN** expression validation checks for valid operator syntax

#### Scenario: Create computed field with database functions
- **WHEN** user creates a new field in dataset field editor
- **THEN** user can access function library compatible with dataset's database
- **THEN** functions include aggregation (SUM, AVG, COUNT, MAX, MIN)
- **THEN** functions include string (CONCAT, SUBSTRING, LENGTH)
- **THEN** functions include date (DATE_FORMAT, DATE_ADD, DATEDIFF)
- **THEN** functions include math (ROUND, CEIL, FLOOR, ABS)
- **THEN** function signature and parameters are displayed

#### Scenario: Create computed field with complex expression
- **WHEN** user creates a field with expression `ROUND([price] * [quantity], 2)`
- **THEN** expression references price and quantity fields
- **THEN** expression uses ROUND function with precision parameter
- **THEN** expression uses multiplication operator
- **THEN** computed field is calculated when dataset is queried

### Requirement: Computed Field Type Configuration
The frontend and backend SHALL support configuring computed field type (dimension/measure) and data type based on expression result.

#### Scenario: Set computed field as dimension
- **WHEN** user creates computed field
- **THEN** user can select field type as "dimension"
- **THEN** dimension type is used for categorical data
- **THEN** field appears in dimension list for visualizations

#### Scenario: Set computed field as measure
- **WHEN** user creates computed field
- **THEN** user can select field type as "measure"
- **THEN** measure type is used for numerical data
- **THEN** field appears in measure list for visualizations
- **THEN** field can be aggregated in visualizations

#### Scenario: Configure computed field data type
- **WHEN** user creates computed field
- **THEN** user can manually set data type (string, number, date)
- **THEN** data type is used for type validation
- **THEN** user can enable auto-detection based on expression result

#### Scenario: Auto-detect computed field data type
- **WHEN** user enables auto-detection for computed field
- **THEN** system executes expression on sample data
- **THEN** system infers data type from result (string, number, date)
- **THEN** system sets field data type automatically
- **THEN** user can override auto-detected type

### Requirement: Computed Field Expression Editor
The frontend SHALL provide a rich expression editor for creating and editing computed field expressions.

#### Scenario: Expression editor with field reference
- **WHEN** user opens computed field expression editor
- **THEN** editor shows field list on the left side
- **THEN** clicking a field adds it to expression at cursor position
- **THEN** fields are displayed with their data types

#### Scenario: Expression editor with function library
- **WHEN** user opens computed field expression editor
- **THEN** editor shows function library on the right side
- **THEN** functions are organized by category (aggregation, string, date, math)
- **THEN** clicking a function shows signature and adds it to expression

#### Scenario: Expression editor with syntax highlighting
- **WHEN** user types in expression editor
- **THEN** fields are highlighted in one color
- **THEN** functions are highlighted in another color
- **THEN** operators and literals are highlighted
- **THEN** invalid syntax is underlined in red

#### Scenario: Expression editor with validation
- **WHEN** user types in expression editor
- **THEN** real-time validation checks syntax
- **THEN** validation errors are displayed below editor
- **THEN** errors include expected and actual tokens
- **THEN** validation button can be clicked to check entire expression

#### Scenario: Expression editor with auto-complete
- **WHEN** user types field name prefix
- **THEN** auto-complete shows matching fields
- **THEN** selecting from auto-complete inserts field reference
- **WHEN** user types function name prefix
- **THEN** auto-complete shows matching functions
- **THEN** selecting from auto-complete inserts function call with parameters

### Requirement: Computed Field Expression Execution
The backend SHALL execute computed field expressions when dataset data is queried.

#### Scenario: Execute computed field in SQL query
- **WHEN** dataset contains computed field with expression
- **THEN** backend wraps expression in SQL SELECT clause
- **THEN** backend substitutes field references with actual column names
- **THEN** backend executes query with computed field expression
- **THEN** results include computed field values

#### Scenario: Execute computed field for API dataset
- **WHEN** dataset contains computed field with expression
- **THEN** backend fetches data from API
- **THEN** backend applies expression to each row in-memory
- **THEN** results include computed field values

#### Scenario: Execute computed field with nested references
- **WHEN** computed field expression references another computed field
- **THEN** backend resolves nested field references
- **THEN** backend calculates dependent fields first
- **THEN** results include all computed field values

#### Scenario: Handle computed field errors
- **WHEN** computed field expression evaluation fails for a row
- **THEN** system logs error with row identifier
- **THEN** system returns NULL for that field value
- **THEN** system continues processing other rows

### Requirement: Computed Field Caching
The backend SHALL cache computed field expressions to improve query performance.

#### Scenario: Cache computed field SQL
- **WHEN** computed field is created or updated
- **THEN** backend generates and caches SQL SELECT expression
- **THEN** cached expression is reused for subsequent queries
- **THEN** cache is invalidated when field is modified

#### Scenario: Cache computed field for API datasets
- **WHEN** computed field is created for API dataset
- **THEN** backend caches compiled expression evaluator
- **THEN** cached evaluator is reused for subsequent queries
- **THEN** cache is invalidated when field is modified

