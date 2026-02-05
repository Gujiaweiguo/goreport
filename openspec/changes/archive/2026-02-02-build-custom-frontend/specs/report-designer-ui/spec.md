# Spec: Report Designer UI

## ADDED Requirements

### Requirement: Canvas-based Report Designer

The system SHALL provide a canvas-based report designer that allows users to design reports by dragging and manipulating cells in an Excel-like interface.

#### Scenario: User opens blank report designer
- **WHEN** user creates a new report
- **THEN** a blank canvas designer is displayed with default grid size (100 columns x 50 rows)
- **AND** standard toolbar is visible (font, color, border, alignment)
- **AND** property panel is visible on the right side

#### Scenario: User selects single cell
- **WHEN** user clicks on a cell
- **THEN** the cell is highlighted with a selection border
- **AND** property panel shows cell properties (value, font, color, borders)
- **AND** user can edit cell value in place or via property panel

#### Scenario: User selects multiple cells via drag
- **WHEN** user drags mouse across multiple cells
- **THEN** all cells in the selection area are highlighted
- **AND** property panel shows properties for the selection
- **AND** user can apply styles to all selected cells

#### Scenario: User adjusts column width
- **WHEN** user drags column header border
- **THEN** column width changes in real-time
- **AND** adjacent columns shift accordingly
- **AND** width snaps to grid (8px increments)

#### Scenario: User adjusts row height
- **WHEN** user drags row header border
- **THEN** row height changes in real-time
- **AND** adjacent rows shift accordingly
- **AND** height snaps to grid (4px increments)

#### Scenario: User merges cells
- **WHEN** user selects multiple cells and clicks "Merge" button
- **THEN** cells are merged into a single cell
- **AND** cell value from top-left cell is preserved
- **AND** merged cell spans the selected area
- **AND** merge operation is saved to report design

#### Scenario: User unmerges cell
- **WHEN** user selects a merged cell and clicks "Unmerge" button
- **THEN** cell is unmerged to original cell grid
- **AND** cell value is retained in top-left cell
- **AND** unmerge operation is saved to report design

#### Scenario: User applies cell style
- **WHEN** user selects cells and applies font style (bold, italic, size)
- **THEN** style is applied to all selected cells
- **AND** style changes are visible immediately on canvas
- **AND** style changes are saved to report design

#### Scenario: User applies cell border
- **WHEN** user selects cells and applies border (top, bottom, left, right, all)
- **THEN** border is applied to all selected cells
- **AND** border changes are visible immediately on canvas
- **AND** border changes are saved to report design

#### Scenario: User applies cell color
- **WHEN** user selects cells and applies background color or text color
- **THEN** color is applied to all selected cells
- **AND** color changes are visible immediately on canvas
- **AND** color changes are saved to report design

### Requirement: Cell Data Binding

The system SHALL allow users to bind cells to data sources, enabling dynamic data population in reports.

#### Scenario: User binds cell to database field
- **WHEN** user selects a cell and opens data binding dialog
- **THEN** user can select a data source
- **AND** user can select a table from the data source
- **AND** user can select a field from the table
- **AND** cell displays the field name as placeholder (e.g., "${field_name}")
- **AND** binding is saved to report design

#### Scenario: User binds cell with expression
- **WHEN** user enters an expression in data binding dialog (e.g., "${SUM(field)}")
- **THEN** expression is validated for syntax errors
- **AND** cell displays the expression as placeholder
- **AND** expression is saved to report design
- **AND** expression syntax is highlighted in Monaco Editor

#### Scenario: User removes data binding
- **WHEN** user opens data binding dialog for a bound cell and clicks "Clear"
- **THEN** data binding is removed
- **AND** cell reverts to showing static value
- **AND** unbinding is saved to report design

#### Scenario: User tests data binding
- **WHEN** user clicks "Test" button in data binding dialog
- **THEN** system fetches sample data from data source
- **AND** cell displays sample value on canvas
- **AND** binding is confirmed to work correctly

### Requirement: Expression Editor

The system SHALL provide an expression editor with syntax highlighting and auto-completion for building complex data expressions.

#### Scenario: User opens expression editor
- **WHEN** user clicks expression editor button in property panel
- **THEN** Monaco Editor is opened in a modal dialog
- **AND** current expression is loaded into editor
- **AND** expression syntax is highlighted

#### Scenario: User types expression with auto-completion
- **WHEN** user starts typing a function name (e.g., "SU")
- **THEN** auto-completion suggestions appear (e.g., "SUM", "SUBSTR", "UPPER")
- **AND** user can select from suggestions with Tab or Enter
- **AND** selected function is inserted with cursor inside parentheses

#### Scenario: User validates expression
- **WHEN** user clicks "Validate" button in expression editor
- **THEN** expression syntax is checked
- **AND** validation result is displayed (valid or error message)
- **AND** error location is highlighted if there's an error

#### Scenario: User saves expression
- **WHEN** user clicks "Save" button in expression editor
- **THEN** expression is saved to cell binding
- **AND** editor modal is closed
- **AND** cell displays expression as placeholder

### Requirement: Report Saving and Loading

The system SHALL allow users to save report designs and load them for editing.

#### Scenario: User saves new report
- **WHEN** user clicks "Save" button for a new report
- **THEN** user is prompted to enter report name
- **AND** report design is sent to backend via API
- **AND** success message is displayed
- **AND** report ID is assigned and stored

#### Scenario: User saves existing report
- **WHEN** user clicks "Save" button for an existing report
- **THEN** report design is sent to backend via API
- **AND** success message is displayed
- **AND** "Saved" indicator appears in status bar

#### Scenario: User loads report for editing
- **WHEN** user selects a report from report list and clicks "Edit"
- **THEN** report design is fetched from backend via API
- **AND** canvas designer displays the loaded report
- **AND** all cells, styles, and bindings are restored

#### Scenario: Auto-save triggers
- **WHEN** user makes changes and 30 seconds pass without save
- **THEN** auto-save is triggered
- **AND** report design is sent to backend via API
- **AND** "Auto-saved" indicator appears in status bar

#### Scenario: User opens unsaved report with changes
- **WHEN** user tries to open another report without saving current changes
- **THEN** confirmation dialog appears ("You have unsaved changes. Save?")
- **AND** user can choose to Save, Don't Save, or Cancel

### Requirement: Performance Requirements

The report designer SHALL perform smoothly with large reports.

#### Scenario: Designer loads large report (1000+ cells)
- **WHEN** user opens a report with 1000+ cells
- **THEN** canvas loads within 3 seconds
- **AND** all cells are rendered
- **AND** user can interact with cells immediately

#### Scenario: Designer handles scrolling with large report
- **WHEN** user scrolls through a report with 10000+ cells
- **THEN** scrolling is smooth (60 FPS)
- **AND** only visible cells are rendered (virtual scrolling)
- **AND** no lag or freeze occurs

#### Scenario: Designer handles cell selection with large report
- **WHEN** user selects 100+ cells in a large report
- **THEN** selection operation completes within 100ms
- **AND** all selected cells are highlighted
- **AND** property panel updates immediately

#### Scenario: Designer handles copy-paste with large report
- **WHEN** user copies and pastes 50+ cells in a large report
- **THEN** copy-paste operation completes within 500ms
- **AND** all cells are pasted correctly
- **AND** no memory leak occurs

### Requirement: Keyboard Shortcuts

The system SHALL provide keyboard shortcuts for common operations to improve productivity.

#### Scenario: User uses Ctrl+C to copy cells
- **WHEN** user selects cells and presses Ctrl+C
- **THEN** cells are copied to clipboard
- **AND** "Copied" tooltip appears

#### Scenario: User uses Ctrl+V to paste cells
- **WHEN** user has cells in clipboard and presses Ctrl+V
- **THEN** cells are pasted at current cursor position
- **AND** all cell properties (styles, bindings) are preserved

#### Scenario: User uses Ctrl+Z to undo
- **WHEN** user presses Ctrl+Z
- **THEN** last operation is undone
- **AND** canvas reverts to previous state
- **AND** undo can be repeated for multiple operations

#### Scenario: User uses Ctrl+Y to redo
- **WHEN** user presses Ctrl+Y
- **THEN** last undone operation is redone
- **AND** canvas reverts to redone state
- **AND** redo can be repeated for multiple undo operations

#### Scenario: User uses Delete to remove cell content
- **WHEN** user selects cells and presses Delete key
- **THEN** cell content is cleared
- **AND** cell styles are preserved
- **AND** cell bindings are preserved

#### Scenario: User uses Ctrl+S to save report
- **WHEN** user presses Ctrl+S
- **THEN** report is saved via API
- **AND** success message is displayed
- **AND** operation is faster than clicking "Save" button

### Requirement: Context Menu

The system SHALL provide a context menu for quick access to common cell operations.

#### Scenario: User right-clicks on cell
- **WHEN** user right-clicks on a cell
- **THEN** context menu appears at cursor position
- **AND** menu contains: Cut, Copy, Paste, Delete, Format Cells, Data Binding, Merge Cells, Unmerge Cells

#### Scenario: User right-clicks on merged cell
- **WHEN** user right-clicks on a merged cell
- **THEN** context menu contains "Unmerge Cells" instead of "Merge Cells"
- **AND** all other options are available

#### Scenario: User right-clicks on empty area
- **WHEN** user right-clicks on empty canvas area
- **THEN** context menu contains: Paste, Insert Row, Insert Column, Delete Row, Delete Column

#### Scenario: User selects context menu item
- **WHEN** user clicks an item in context menu
- **THEN** corresponding operation is performed
- **AND** context menu is closed
- **AND** changes are saved to report design
