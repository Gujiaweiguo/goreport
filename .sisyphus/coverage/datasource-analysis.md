# Datasource 模块覆盖率分析

**当前覆盖率**: 66.5%
**目标覆盖率**: 80%
**差距**: 13.5%
**优先级**: 🟡 P1

---

## 未覆盖函数列表

| 函数 | 文件 | 覆盖率 | 测试用例需求 |
|------|------|--------|-------------|
| GetTables | cached_metadata.go:40 | 0.0% | TestCachedMetadata_GetTables |
| GetFields | cached_metadata.go:73 | 0.0% | TestCachedMetadata_GetFields |
| BuildDSN | connection_builder.go:22 | 61.9% | 补充边界场景 |
| Connect | connection_builder.go:70 | 46.2% | TestConnect_* |
| TestConnection | connection_builder.go:138 | 25.0% | TestConnection_* |
| GetTables | metadata.go:17 | 0.0% | TestMetadata_GetTables |
| GetFields | metadata.go:33 | 0.0% | TestMetadata_GetFields |
| forwardConnections | ssh_tunnel.go:120 | 0.0% | SSH 集成测试 |
| copyData | ssh_tunnel.go:164 | 0.0% | SSH 集成测试 |
| Close | ssh_tunnel.go:193 | 0.0% | TestSSHTunnel_Close |
| LocalAddr | ssh_tunnel.go:215 | 0.0% | TestSSHTunnel_LocalAddr |

---

## 函数分析

### CachedMetadata (0.0%)
**需要测试用例**:
- TestCachedMetadata_GetTables_CacheHit
- TestCachedMetadata_GetTables_CacheMiss
- TestCachedMetadata_GetFields_CacheHit
- TestCachedMetadata_GetFields_CacheMiss
- TestCachedMetadata_CacheExpiration

### ConnectionBuilder
**需要补充测试用例**:
- TestBuildDSN_MySQL
- TestBuildDSN_PostgreSQL
- TestBuildDSN_WithSSL
- TestConnect_Timeout
- TestConnect_AuthFailure
- TestConnection_Success
- TestConnection_Failure

### Metadata (0.0%)
**需要测试用例**:
- TestMetadata_GetTables_Success
- TestMetadata_GetTables_Error
- TestMetadata_GetFields_Success
- TestMetadata_GetFields_Error

### SSHTunnel
**注意**: SSH 测试需要外部服务器，可使用 mock 或跳过

---

## 测试策略

### Mock 数据库连接
```go
type MockDBConnector struct {
    mock.Mock
}

func (m *MockDBConnector) Connect() (*sql.DB, error) {
    args := m.Called()
    return args.Get(0).(*sql.DB), args.Error(1)
}
```

### 预计测试用例数量
- CachedMetadata: 5-8 用例
- ConnectionBuilder 补充: 8-10 用例
- Metadata: 4-6 用例
- SSHTunnel (可选): 3-5 用例

**总计**: 20-30 个测试用例

---

## 实现优先级

1. **高优先级**:
   - CachedMetadata 缓存逻辑
   - ConnectionBuilder 核心连接

2. **中优先级**:
   - Metadata 元数据查询
   - 错误处理场景

3. **低优先级**:
   - SSHTunnel (需要外部依赖)
