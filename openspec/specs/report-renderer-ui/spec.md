# report-renderer-ui Specification

## Purpose
TBD - created by archiving change build-custom-frontend. Update Purpose after archive.
## Requirements
### Requirement: Report Preview

The system SHALL provide a report preview that displays the report with actual data populated from the data source.

#### Scenario: User previews report with sample data
- **WHEN** user clicks "Preview" button in designer
- **THEN** report is rendered with sample data from data source
- **AND** all cells display actual data values
- **AND** bound expressions are calculated and displayed
- **AND** report is displayed in read-only mode
- **AND** toolbar shows "Preview" mode

#### Scenario: User previews report with parameters
- **WHEN** user enters parameter values (e.g., date range, user ID) and clicks "Preview"
- **THEN** report is rendered with data filtered by parameters
- **AND** parameter values are displayed in header
- **AND** expressions are recalculated with parameter context

#### Scenario: User refreshes preview
- **WHEN** user clicks "Refresh" button in preview mode
- **THEN** report is re-rendered with fresh data from data source
- **AND** loading indicator is shown during refresh
- **AND** refresh completes within 2 seconds for small reports (<1000 rows)

#### Scenario: User exits preview
- **WHEN** user clicks "Back to Designer" button
- **THEN** preview mode is closed
- **AND** designer is reopened with original report design
- **AND** unsaved changes in preview are discarded

### Requirement: Report Pagination

The system SHALL support report pagination for reports with many rows.

#### Scenario: Report displays with pagination
- **WHEN** report has more than 50 rows
- **THEN** report is displayed in pages (default 50 rows per page)
- **AND** pagination controls appear at bottom of report
- **AND** controls show: "Page X of Y", "Previous", "Next"
- **AND** page 1 is displayed by default

#### Scenario: User navigates to next page
- **WHEN** user clicks "Next" button in pagination
- **THEN** next page of report is displayed
- **AND** "Next" button is disabled if on last page
- **AND** page indicator updates to show current page

#### Scenario: User navigates to previous page
- **WHEN** user clicks "Previous" button in pagination
- **THEN** previous page of report is displayed
- **AND** "Previous" button is disabled if on first page
- **AND** page indicator updates to show current page

#### Scenario: User jumps to specific page
- **WHEN** user enters a page number and presses Enter
- **THEN** specified page is displayed
- **AND** error message appears if page number is invalid
- **AND** page is clamped to valid range

#### Scenario: User changes page size
- **WHEN** user selects a different page size (25, 50, 100, 200)
- **THEN** page size is updated
- **AND** report is re-rendered with new page size
- **AND** pagination controls update (page count changes)
- **AND** user is returned to page 1

#### Scenario: Report displays without pagination (for small reports)
- **WHEN** report has 50 or fewer rows
- **THEN** pagination controls are hidden
- **AND** all rows are displayed on single page

### Requirement: Print Preview

The system SHALL provide a print preview that shows how the report will look when printed.

#### Scenario: User opens print preview
- **WHEN** user clicks "Print Preview" button
- **THEN** report is displayed in print preview mode
- **AND** page size matches configured paper size (A4 default)
- **AND** page orientation is applied (portrait or landscape)
- **AND** print preview toolbar appears (Print, Page Setup, Close)

#### Scenario: User changes paper size in print preview
- **WHEN** user selects a different paper size (A4, A3, Letter)
- **THEN** report is re-layout to fit new paper size
- **AND** page width and height adjust accordingly
- **AND** content is scaled or wrapped if needed

#### Scenario: User changes page orientation in print preview
- **WHEN** user switches between portrait and landscape
- **THEN** report is rotated to match orientation
- **AND** content is re-flowed to fit new orientation
- **AND** page count may change

#### Scenario: User sets margins in print preview
- **WHEN** user adjusts margin values (top, bottom, left, right)
- **THEN** printable area is reduced by margin size
- **AND** content is re-layout to fit printable area
- **AND** margins are shown as dotted lines on preview

#### Scenario: User prints from print preview
- **WHEN** user clicks "Print" button
- **THEN** browser print dialog is opened
- **AND** report is formatted for printing (hide UI elements)
- **AND** print settings (paper size, orientation, margins) are applied

### Requirement: Export Dialog

The system SHALL provide an export dialog for configuring export options before generating export files.

#### Scenario: User opens export dialog
- **WHEN** user clicks "Export" button
- **THEN** export dialog is displayed
- **AND** user can select export format (Excel, PDF, Word, Image)
- **AND** user can configure export options for selected format

#### Scenario: User configures Excel export
- **WHEN** user selects "Excel" format
- **THEN** Excel-specific options are shown:
  - Include headers (yes/no)
  - Include styles (yes/no)
  - Sheet name
  - Encoding
- **AND** default values are pre-selected

#### Scenario: User configures PDF export
- **WHEN** user selects "PDF" format
- **THEN** PDF-specific options are shown:
  - Paper size (A4, A3, Letter)
  - Orientation (portrait, landscape)
  - Margins
  - Include page numbers (yes/no)
  - Quality (high, medium, low)
- **AND** default values are pre-selected

#### Scenario: User configures Word export
- **WHEN** user selects "Word" format
- **THEN** Word-specific options are shown:
  - Include headers (yes/no)
  - Include styles (yes/no)
  - Document template
- **AND** default values are pre-selected

#### Scenario: User configures Image export
- **WHEN** user selects "Image" format
- **THEN** Image-specific options are shown:
  - Image format (PNG, JPG, SVG)
  - Resolution (DPI)
  - Quality (1-100)
  - Background color (transparent, white)
- **AND** default values are pre-selected

#### Scenario: User starts export
- **WHEN** user configures options and clicks "Export"
- **THEN** export job is created via backend API
- **AND** progress indicator is displayed
- **AND** user can continue working while export runs in background
- **AND** user is notified when export completes

### Requirement: Export Progress

The system SHALL display export progress and allow users to download completed exports.

#### Scenario: Export in progress
- **WHEN** export job is running
- **THEN** progress indicator shows percentage complete
- **AND** progress indicator shows estimated time remaining
- **AND** user can cancel export by clicking "Cancel"

#### Scenario: Export completes successfully
- **WHEN** export job completes successfully
- **THEN** success message is displayed
- **AND** "Download" button appears
- **AND** export job ID is stored for future reference

#### Scenario: User downloads completed export
- **WHEN** user clicks "Download" button
- **THEN** browser downloads the export file
- **AND** file name includes report name and timestamp
- **AND** file format matches selected format

#### Scenario: Export fails with error
- **WHEN** export job fails
- **THEN** error message is displayed with error details
- **AND** error log is shown if available
- **AND** user can retry export
- **AND** error is logged for troubleshooting

### Requirement: Performance Requirements

The report renderer SHALL perform efficiently with large datasets.

#### Scenario: Render report with 10,000 rows
- **WHEN** user previews a report with 10,000 rows
- **THEN** report loads within 3 seconds
- **AND** pagination is applied (50 rows per page)
- **AND** user can navigate between pages smoothly

#### Scenario: Render report with complex expressions
- **WHEN** user previews a report with 500 cells containing expressions
- **THEN** all expressions are calculated within 2 seconds
- **AND** calculated values are displayed correctly
- **AND** no lag or freeze occurs

#### Scenario: Render report with charts
- **WHEN** user previews a report with 10 chart components
- **THEN** all charts render within 3 seconds
- **AND** charts are displayed with correct data
- **AND** charts are interactive (hover, tooltip)

#### Scenario: Switch between preview and designer
- **WHEN** user switches from preview to designer and back
- **THEN** switch completes within 500ms
- **AND** report state is preserved
- **AND** no data loss occurs

### Requirement: Responsive Design

The report renderer SHALL adapt to different screen sizes and devices.

#### Scenario: View report on large screen
- **WHEN** user views report on screen width > 1920px
- **THEN** report is displayed at full width
- **AND** content is optimally spaced
- **AND** no horizontal scrollbar appears

#### Scenario: View report on medium screen
- **WHEN** user views report on screen width 1024-1920px
- **THEN** report fits within screen width
- **AND** content is scaled if needed
- **AND** horizontal scrollbar appears only if content is too wide

#### Scenario: View report on small screen
- **WHEN** user views report on screen width < 1024px
- **THEN** report is optimized for small screen
- **AND** horizontal scrolling is enabled
- **AND** responsive layout is applied

#### Scenario: View report on mobile device
- **WHEN** user views report on mobile device (width < 768px)
- **THEN** report is displayed in mobile-friendly layout
- **AND** tables are scrollable horizontally
- **AND** pagination controls are optimized for touch
- **AND** export options are available in collapsible menu

