## 1. Database Schema Setup

- [x] 1.1 Create `datasets` table with columns: id, name, type (sql/api/file), tenant_id, datasource_id, config (JSON), created_at, updated_at, deleted_at
- [x] 1.2 Create `dataset_fields` table with columns: id, dataset_id, name, display_name, type (dimension/measure), data_type, is_computed, expression, is_sortable, is_groupable, default_sort_order, sort_index, created_at, updated_at
- [x] 1.3 Create `dataset_sources` table with columns: id, dataset_id, source_type, source_config (JSON), created_at, updated_at
- [x] 1.4 Add indexes on datasets.tenant_id, dataset_fields.dataset_id, dataset_sources.dataset_id
- [x] 1.5 Add foreign key constraints and cascade delete rules

## 2. Backend Model and Repository Layer

- [x] 2.1 Create Dataset model struct with GORM tags in `backend/internal/models/dataset.go`
- [x] 2.2 Create DatasetField model struct with GORM tags in `backend/internal/models/dataset_field.go`
- [x] 2.3 Create DatasetSource model struct with GORM tags in `backend/internal/models/dataset_source.go`
- [x] 2.4 Create DatasetRepository interface in `backend/internal/repository/dataset_repository.go`
- [x] 2.5 Implement CRUD methods in DatasetRepository: Create, GetByID, List, Update, Delete, SoftDelete
- [x] 2.6 Implement DatasetFieldRepository with CRUD methods
- [x] 2.7 Implement DatasetSourceRepository with CRUD methods
- [x] 2.8 Add tenant filtering in all repository query methods

## 3. Backend Service Layer - Dataset Management

- [x] 3.1 Create DatasetService struct in `backend/internal/service/dataset_service.go`
- [x] 3.2 Implement CreateDataset method: validate datasource belongs to tenant, execute SQL/API call to extract fields
- [x] 3.3 Implement GetDataset method with tenant validation
- [x] 3.4 Implement ListDatasets method with tenant filtering and pagination
- [x] 3.5 Implement UpdateDataset method: re-extract fields if SQL or source changed
- [x] 3.6 Implement DeleteDataset method: check for references to reports, dashboards, charts
- [x] 3.7 Implement PreviewDataset method: execute query with LIMIT 100
- [x] 3.8 Implement GetDatasetSchema method: return dimensions, measures, computed fields

## 4. Backend Service Layer - Field Management

- [x] 4.1 Implement CreateComputedField method: validate expression syntax, validate field references
- [x] 4.2 Implement UpdateFieldConfiguration method: update displayName, type, dataType, sort, grouping
- [x] 4.3 Implement DeleteComputedField method: validate is computed field before deletion
- [x] 4.4 Implement ListDimensions method: return dimension fields only
- [x] 4.5 Implement ListMeasures method: return measure fields only
- [x] 4.6 Implement GetDatasetFields method: return all fields with metadata

## 5. Backend Service Layer - Computed Fields Engine

- [x] 5.1 Create SQLExpressionBuilder in `backend/internal/service/expression_builder.go`
- [x] 5.2 Implement field reference substitution: [field_name] → actual column name
- [x] 5.3 Implement function translation: convert generic function names to database-specific SQL
- [x] 5.4 Implement expression validation: check syntax, field references, function usage
- [x] 5.5 Create APIExpressionBuilder for non-SQL datasets
- [x] 5.6 Implement in-memory expression evaluation using JavaScript for API datasets
- [x] 5.7 Implement computed field caching: cache compiled SQL expressions and JS evaluators
- [x] 5.8 Implement nested computed field resolution: calculate dependencies first

## 6. Backend Service Layer - Query Execution

- [x] 6.1 Create QueryExecutor in `backend/internal/service/query_executor.go`
- [x] 6.2 Implement SQL dataset query: construct SELECT clause with fields and computed fields
- [x] 6.3 Implement API dataset query: fetch data from API endpoint and apply transformations
- [x] 6.4 Implement file import dataset query: read file and parse data
- [x] 6.5 Implement filter application: translate filter operators to SQL WHERE clauses or JS filter
- [x] 6.6 Implement sorting: apply sortBy and sortOrder to query
- [x] 6.7 Implement pagination: apply LIMIT and OFFSET for SQL, array slice for API
- [x] 6.8 Implement field selection: SELECT only specified columns for performance
- [x] 6.9 Implement grouping and aggregation: apply GROUP BY and aggregation functions
- [x] 6.10 Add tenant isolation to queries: inject tenant filter conditions

## 7. Backend HTTP Handlers

- [x] 7.1 Create DatasetHandler in `backend/internal/httpserver/handlers/dataset_handler.go`
- [x] 7.2 Implement POST /api/v1/datasets (CreateDataset)
- [x] 7.3 Implement GET /api/v1/datasets (ListDatasets) with pagination query params
- [x] 7.4 Implement GET /api/v1/datasets/:id (GetDataset)
- [x] 7.5 Implement PUT /api/v1/datasets/:id (UpdateDataset)
- [x] 7.6 Implement DELETE /api/v1/datasets/:id (DeleteDataset)
- [x] 7.7 Implement GET /api/v1/datasets/:id/preview (PreviewDataset)
- [x] 7.8 Implement POST /api/v1/datasets/:id/data (QueryDataset) with filters, sorting, pagination
- [x] 7.9 Implement GET /api/v1/datasets/:id/dimensions (GetDimensions)
- [x] 7.10 Implement GET /api/v1/datasets/:id/measures (GetMeasures)
- [x] 7.11 Implement GET /api/v1/datasets/:id/schema (GetDatasetSchema)
- [x] 7.12 Implement POST /api/v1/datasets/:id/fields (CreateComputedField)
- [x] 7.13 Implement PUT /api/v1/datasets/:id/fields/:fieldId (UpdateField)
- [x] 7.14 Implement DELETE /api/v1/datasets/:id/fields/:fieldId (DeleteField)
- [x] 7.15 Implement POST /api/v1/datasets/validate (ValidateDatasetConfig)
- [x] 7.16 Add request DTO structs with validation tags
- [x] 7.17 Add response DTO structs with JSON serialization

## 8. Backend Routing and Middleware

- [x] 8.1 Add dataset routes to router in `backend/internal/httpserver/router.go`
- [x] 8.2 Apply JWT authentication middleware to all dataset routes
- [x] 8.3 Apply tenant extraction middleware to populate tenant context
- [x] 8.4 Add rate limiting for query endpoints
- [x] 8.5 Add request logging for audit trail

## 9. Frontend - Dataset Management Pages

- [x] 9.1 Create DatasetList page component in `frontend/src/views/dataset/DatasetList.vue`
- [x] 9.2 Implement dataset table with columns: name, type, datasource, field count, created_at, actions
- [x] 9.3 Implement create dataset button and dialog
- [x] 9.4 Implement pagination controls
- [x] 9.5 Implement search and filter functionality
- [x] 9.6 Create DatasetEdit page component in `frontend/src/views/dataset/DatasetEdit.vue`
- [x] 9.7 Implement dataset type selector (SQL, API, File Import)
- [x] 9.8 Implement SQL editor with syntax highlighting
- [x] 9.9 Implement API configuration form (URL, method, headers, body)
- [x] 9.10 Implement file upload component (Excel, CSV)
- [x] 9.11 Implement dataset preview table
- [x] 9.12 Implement save and cancel buttons with validation

## 10. Frontend - Field Configuration UI

- [x] 10.1 Create FieldList component in `frontend/src/components/dataset/FieldList.vue`
- [x] 10.2 Implement dimensions and measures tab view
- [x] 10.3 Implement field table with columns: name, display name, type, data type, is computed, actions
- [x] 10.4 Implement edit field dialog
- [x] 10.5 Implement field type toggle (dimension/measure)
- [x] 10.6 Implement data type selector (string, number, date, boolean)
- [x] 10.7 Implement sort order configuration (asc, desc, none)
- [x] 10.8 Implement grouping toggle
- [x] 10.9 Implement display name input

## 11. Frontend - Computed Field Editor

- [x] 11.1 Create ComputedFieldEditor component in `frontend/src/components/dataset/ComputedFieldEditor.vue`
- [x] 11.2 Implement expression editor with Monaco or custom text editor
- [x] 11.3 Implement field reference panel on left side
- [x] 11.4 Implement function library panel on right side
- [x] 11.5 Implement function signature display
- [x] 11.6 Implement syntax highlighting for fields, functions, operators
- [x] 11.7 Implement real-time expression validation with error display
- [x] 11.8 Implement auto-complete for field names
- [x] 11.9 Implement auto-complete for function names
- [x] 11.10 Implement field type selector (dimension/measure)
- [x] 11.11 Implement data type selector with auto-detect option
- [x] 11.12 Implement preview expression result on sample data

## 15. Frontend - API Client Layer

- [x] 15.1 Create dataset API client in `frontend/src/api/dataset.ts`
- [x] 15.2 Implement listDatasets function
- [x] 15.3 Implement getDataset function
- [x] 15.4 Implement createDataset function
- [x] 15.5 Implement updateDataset function
- [x] 15.6 Implement deleteDataset function
- [x] 15.7 Implement previewDataset function
- [x] 15.8 Implement queryDataset function with filters, sorting, pagination
- [x] 15.9 Implement getDimensions function
- [x] 15.10 Implement getMeasures function
- [x] 15.11 Implement getSchema function
- [x] 15.12 Implement createComputedField function
- [x] 15.13 Implement updateField function
- [x] 15.14 Implement deleteField function
- [x] 15.15 Implement validateDatasetConfig function

## 12. Frontend - Report Designer Integration

- [x] 12.1 Modify ReportDesigner property panel to add "Dataset" data source option
- [x] 12.2 Implement dataset selector dropdown
- [x] 12.3 Implement dimension and measure selector dropdowns
- [x] 12.4 Implement field binding display with dataset name and field name
- [x] 12.5 Implement aggregation function selector for measures (SUM, AVG, COUNT, MAX, MIN, none)
- [x] 12.6 Implement grouping toggle for dimensions
- [x] 12.7 Update report config serialization to include dataset references
- [x] 12.8 Implement migration wizard to convert existing SQL bindings to dataset bindings

## 13. Frontend - Dashboard Integration

- [x] 13.1 Modify DashboardDesigner widget configuration to add "Dataset" data source option
- [x] 13.2 Implement dataset selector dropdown in widget config
- [x] 13.3 Implement dimension and measure selectors for chart axes
- [x] 13.4 Implement multi-measure support for series charts
- [x] 13.5 Update dashboard config serialization to include dataset references
- [x] 13.6 Implement dashboard data query API call for dataset widgets

## 14. Frontend - Chart Editor Integration

- [x] 14.1 Modify ChartEditor data source configuration to add "Dataset" type
- [x] 14.2 Implement dataset selector dropdown
- [x] 14.3 Implement dimension mapping interface for X-axis, category, series
- [x] 14.4 Implement measure mapping interface for Y-axis, value, size, angle
- [x] 14.5 Implement computed field support in selectors
- [x] 14.6 Update chart preview to fetch data from dataset API
- [x] 14.7 Implement chart-dataset compatibility validation
- [x] 14.8 Handle dataset schema change notifications

## 15. Frontend - API Client Layer

- [x] 15.1 Create dataset API client in `frontend/src/api/dataset.ts`
- [x] 15.2 Implement listDatasets function
- [x] 15.3 Implement getDataset function
- [x] 15.4 Implement createDataset function
- [x] 15.5 Implement updateDataset function
- [x] 15.6 Implement deleteDataset function
- [x] 15.7 Implement previewDataset function
- [x] 15.8 Implement queryDataset function with filters, sorting, pagination
- [x] 15.9 Implement getDimensions function
- [x] 15.10 Implement getMeasures function
- [x] 15.11 Implement getSchema function
- [x] 15.12 Implement createComputedField function
- [x] 15.13 Implement updateField function
- [x] 15.14 Implement deleteField function
- [x] 15.15 Implement validateDatasetConfig function

## 16. Backend Testing

- [x] 16.1 Write unit tests for Dataset models
- [x] 16.2 Write unit tests for DatasetRepository methods
- [x] 16.3 Write unit tests for DatasetService CRUD methods
- [x] 16.4 Write unit tests for field management methods
- [x] 16.5 Write unit tests for computed field expression builder
- [x] 16.6 Write unit tests for query executor with filters, sorting, pagination
- [x] 16.7 Write integration tests for dataset CRUD API endpoints
- [x] 16.8 Write integration tests for dataset query API
- [x] 16.9 Write integration tests for field management API
- [x] 16.10 Write tests for tenant isolation enforcement
- [x] 16.11 Write tests for computed field expression validation
- [x] 16.12 Write performance tests for query execution

## 17. Frontend Testing

- [x] 17.1 Write unit tests for DatasetList component
- [x] 17.2 Write unit tests for DatasetEdit component
- [x] 17.3 Write unit tests for FieldList component
- [x] 17.4 Write unit tests for ComputedFieldEditor component
- [x] 17.5 Write unit tests for API client functions
- [x] 17.6 Write integration tests for dataset creation flow
- [x] 17.7 Write integration tests for computed field creation
- [x] 17.8 Write end-to-end tests for report designer with dataset
- [x] 17.9 Write end-to-end tests for dashboard with dataset
- [x] 17.10 Write end-to-end tests for chart editor with dataset

## 18. Documentation

- [x] 18.1 Write dataset management user guide
- [x] 18.2 Write computed field editor user guide
- [x] 18.3 Write API documentation for dataset endpoints
- [x] 18.4 Write developer guide for integrating dataset feature
- [x] 18.5 Write migration guide from direct SQL to dataset
- [x] 18.6 Update README with dataset feature description
- [x] 18.7 Update AGENTS.md with dataset conventions

## 19. Deployment and Monitoring

- [x] 19.1 Prepare database migration scripts
- [x] 19.2 Configure monitoring and logging for dataset queries
- [x] 19.3 Set up performance metrics collection
- [x] 19.4 Configure alerts for query failures and slow queries
- [x] 19.5 Prepare rollback procedure
- [x] 19.6 Create deployment checklist
