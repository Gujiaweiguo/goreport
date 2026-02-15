# 测试覆盖率改进 - 详细测试计划

**创建时间**: 2026-02-15
**基于文档**: coverage-requirements.md

---

## 1. Repository 模块测试计划 (P0)

### 1.1 DatasetFieldRepository

| 测试用例 | 描述 | 优先级 |
|----------|------|--------|
| TestDatasetFieldRepository_Create_Success | 成功创建字段 | 高 |
| TestDatasetFieldRepository_Create_Duplicate | 重复创建处理 | 中 |
| TestDatasetFieldRepository_GetByID_Found | 成功获取字段 | 高 |
| TestDatasetFieldRepository_GetByID_NotFound | 字段不存在 | 高 |
| TestDatasetFieldRepository_List_Success | 成功列表查询 | 高 |
| TestDatasetFieldRepository_ListByType_Dimension | 按维度类型查询 | 中 |
| TestDatasetFieldRepository_ListByType_Measure | 按度量类型查询 | 中 |
| TestDatasetFieldRepository_Update_Success | 成功更新字段 | 高 |
| TestDatasetFieldRepository_Delete_Success | 成功删除字段 | 高 |
| TestDatasetFieldRepository_DeleteComputedFields_Success | 成功删除计算字段 | 中 |

### 1.2 DatasetRepository

| 测试用例 | 描述 | 优先级 |
|----------|------|--------|
| TestDatasetRepository_Create_Success | 成功创建数据集 | 高 |
| TestDatasetRepository_Create_WithFields | 带字段创建 | 中 |
| TestDatasetRepository_GetByID_Found | 成功获取数据集 | 高 |
| TestDatasetRepository_GetByID_NotFound | 数据集不存在 | 高 |
| TestDatasetRepository_GetByIDWithFields_Success | 获取数据集及字段 | 高 |
| TestDatasetRepository_List_Success | 成功列表查询 | 高 |
| TestDatasetRepository_List_Pagination | 分页查询 | 中 |
| TestDatasetRepository_Update_Success | 成功更新数据集 | 高 |
| TestDatasetRepository_Delete_Success | 成功删除数据集 | 高 |
| TestDatasetRepository_SoftDelete_Success | 成功软删除 | 中 |

### 1.3 DatasourceRepository

| 测试用例 | 描述 | 优先级 |
|----------|------|--------|
| TestDatasourceRepository_Create_Success | 成功创建数据源 | 高 |
| TestDatasourceRepository_GetByID_Found | 成功获取数据源 | 高 |
| TestDatasourceRepository_GetByID_NotFound | 数据源不存在 | 高 |
| TestDatasourceRepository_List_Success | 成功列表查询 | 高 |
| TestDatasourceRepository_List_ByTenant | 按租户查询 | 高 |
| TestDatasourceRepository_Update_Success | 成功更新数据源 | 高 |
| TestDatasourceRepository_Delete_Success | 成功删除数据源 | 高 |
| TestDatasourceRepository_Search_ByKeyword | 关键字搜索 | 中 |
| TestDatasourceRepository_Copy_Success | 成功复制数据源 | 中 |
| TestDatasourceRepository_Move_Success | 成功移动数据源 | 低 |
| TestDatasourceRepository_Rename_Success | 成功重命名数据源 | 低 |

### 1.4 UserRepository

| 测试用例 | 描述 | 优先级 |
|----------|------|--------|
| TestUserRepository_GetByID_Found | 成功获取用户 | 高 |
| TestUserRepository_GetByID_NotFound | 用户不存在 | 高 |
| TestUserRepository_GetByUsername_Found | 按用户名获取 | 高 |
| TestUserRepository_GetByUsername_NotFound | 用户名不存在 | 高 |
| TestUserRepository_Create_Success | 成功创建用户 | 高 |
| TestUserRepository_Update_Success | 成功更新用户 | 中 |
| TestUserRepository_Delete_Success | 成功删除用户 | 中 |

---

## 2. Database 模块测试计划 (P0)

| 测试用例 | 描述 | 优先级 |
|----------|------|--------|
| TestInitWithConfig_Success | 成功初始化 | 高 |
| TestInitWithConfig_InvalidDSN | 无效 DSN | 高 |
| TestInitWithConfig_ConnectionError | 连接错误 | 高 |
| TestInitWithConfig_MigrationError | 迁移错误 | 中 |
| TestEnsureSchemaCompatibility_Success | 成功兼容检查 | 高 |
| TestEnsureSchemaCompatibility_ColumnNotExists | 列不存在时添加 | 高 |
| TestEnsureSchemaCompatibility_BatchUpdate | 批量更新 | 中 |

---

## 3. Render 模块测试计划 (P1)

| 测试用例 | 描述 | 优先级 |
|----------|------|--------|
| TestFetchCellValue_StaticValue | 静态值 | 高 |
| TestFetchCellValue_Expression | 表达式计算 | 高 |
| TestFetchCellValue_EmptyCell | 空单元格 | 中 |
| TestFetchCellValue_NilBinding | 空绑定 | 中 |
| TestFetchCellValueFromDB_Success | 成功获取 | 高 |
| TestFetchCellValueFromDB_QueryError | 查询错误 | 高 |
| TestFetchCellValueFromDB_EmptyResult | 空结果 | 中 |
| TestFetchCellValueFromDB_WithParams | 带参数 | 中 |
| TestEngine_Render_ComplexConfig | 复杂配置 | 中 |
| TestEngine_Render_LargeDataset | 大数据集 | 低 |

---

## 4. Datasource 模块测试计划 (P1)

| 测试用例 | 描述 | 优先级 |
|----------|------|--------|
| TestCachedMetadata_GetTables_CacheHit | 缓存命中 | 高 |
| TestCachedMetadata_GetTables_CacheMiss | 缓存未命中 | 高 |
| TestCachedMetadata_GetFields_CacheHit | 缓存命中 | 高 |
| TestCachedMetadata_GetFields_CacheMiss | 缓存未命中 | 高 |
| TestCachedMetadata_CacheExpiration | 缓存过期 | 中 |
| TestBuildDSN_MySQL | MySQL DSN | 高 |
| TestBuildDSN_PostgreSQL | PostgreSQL DSN | 中 |
| TestConnect_Success | 成功连接 | 高 |
| TestConnect_Timeout | 连接超时 | 中 |
| TestConnect_AuthFailure | 认证失败 | 高 |
| TestConnection_Success | 测试连接成功 | 高 |
| TestConnection_Failure | 测试连接失败 | 高 |
| TestMetadata_GetTables_Success | 获取表成功 | 高 |
| TestMetadata_GetFields_Success | 获取字段成功 | 高 |

---

## 5. Testutil 模块测试计划 (P1)

| 测试用例 | 描述 | 优先级 |
|----------|------|--------|
| TestFixturesAuth_Setup_Success | 设置成功 | 高 |
| TestFixturesAuth_Cleanup_Success | 清理成功 | 高 |
| TestFixturesDatasource_Setup_Success | 设置成功 | 高 |
| TestFixturesDatasource_Cleanup_Success | 清理成功 | 高 |
| TestTestHelper_SetupMySQLTestDB_Success | 设置测试 DB | 高 |
| TestTestHelper_EnsureTenants_Success | 确保租户存在 | 高 |
| TestTestHelper_CleanupTenantData_Success | 清理租户数据 | 高 |
| TestTestHelper_NewTenantTestContext_Success | 创建测试上下文 | 高 |

---

## 6. Dataset 模块测试计划 (P2)

| 测试用例 | 描述 | 优先级 |
|----------|------|--------|
| TestQueryExecutor_Query_Success | 成功查询 | 高 |
| TestQueryExecutor_Query_WithFilter | 带过滤查询 | 中 |
| TestQueryExecutor_Query_WithSort | 带排序查询 | 中 |
| TestQueryExecutor_QuerySQLDataset_Success | SQL 查询成功 | 高 |
| TestQueryExecutor_QuerySQLDataset_Error | SQL 查询错误 | 高 |
| TestHandler_Create_ValidationError | 创建验证错误 | 中 |
| TestHandler_Update_ValidationError | 更新验证错误 | 中 |

---

## 7. Handlers 模块测试计划 (P2)

| 测试用例 | 描述 | 优先级 |
|----------|------|--------|
| TestAuthHandler_Login_Success | 登录成功 | 高 |
| TestAuthHandler_Login_InvalidCredentials | 无效凭据 | 高 |
| TestAuthHandler_Login_MissingFields | 缺少字段 | 中 |
| TestAuthHandler_Logout_Success | 登出成功 | 高 |
| TestAuthHandler_Logout_NoToken | 无 Token | 中 |
| TestDatasourceHandler_TestConnection_Success | 测试连接成功 | 高 |
| TestDatasourceHandler_TestConnection_Failure | 测试连接失败 | 高 |
| TestHealthHandler_Check_Success | 健康检查成功 | 高 |
| TestHealthHandler_Check_DBError | 数据库错误 | 中 |

---

## 8. 执行时间表

| 阶段 | 模块 | 预计用例数 | 预计时间 |
|------|------|-----------|----------|
| Day 1 上午 | repository (核心 CRUD) | 30-40 | 4h |
| Day 1 下午 | repository (辅助功能) | 30-40 | 4h |
| Day 2 上午 | database | 15-20 | 2h |
| Day 2 下午 | render | 16-25 | 3h |
| Day 3 上午 | datasource | 20-30 | 4h |
| Day 3 下午 | testutil | 15-20 | 2h |
| Day 4 上午 | dataset, handlers | 20-30 | 3h |
| Day 4 下午 | 验证和修复 | - | 2h |
| Day 5 | 验证报告和归档 | - | 4h |

---

## 9. 验收标准

### 覆盖率目标
- repository: ≥80%
- database: ≥80%
- render: ≥80%
- datasource: ≥80%
- testutil: ≥80%
- dataset: ≥80%
- handlers: ≥80%

### 质量标准
- 所有测试通过
- 无不稳定测试
- 测试命名规范
- 测试代码覆盖率可追溯
