# Database 模块覆盖率分析

**当前覆盖率**: 10.8%
**目标覆盖率**: 80%
**差距**: 69.2%
**优先级**: 🔴 P0

---

## 未覆盖函数列表

| 函数 | 文件 | 覆盖率 | 测试用例需求 |
|------|------|--------|-------------|
| InitWithConfig | database.go:20 | 13.6% | TestInitWithConfig_* |
| ensureDatasourceSchemaCompatibility | database.go:61 | 0.0% | TestEnsureDatasourceSchemaCompatibility |

---

## 函数分析

### InitWithConfig (13.6%)
**功能**: 使用配置初始化数据库连接

**未覆盖场景**:
- 配置验证失败
- 连接失败
- 迁移执行失败
- 重复调用处理

**需要测试用例**:
- TestInitWithConfig_Success
- TestInitWithConfig_InvalidConfig
- TestInitWithConfig_ConnectionError
- TestInitWithConfig_MigrationError

### ensureDatasourceSchemaCompatibility (0.0%)
**功能**: 确保数据源表结构兼容

**未覆盖场景**:
- 正常迁移流程
- 列不存在时的处理
- 批量更新

**需要测试用例**:
- TestEnsureDatasourceSchemaCompatibility_Success
- TestEnsureDatasourceSchemaCompatibility_ColumnExists
- TestEnsureDatasourceSchemaCompatibility_BatchUpdate

---

## 测试策略

### 使用 SQLite 内存数据库
```go
import (
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    require.NoError(t, err)
    return db
}
```

### 预计测试用例数量
- InitWithConfig: 5-8 用例
- ensureDatasourceSchemaCompatibility: 3-5 用例
- 其他辅助函数: 5-8 用例

**总计**: 15-20 个测试用例

---

## 实现优先级

1. **高优先级**:
   - InitWithConfig 完整覆盖
   - 连接池配置测试

2. **中优先级**:
   - Schema 兼容性检查
   - 迁移逻辑测试
