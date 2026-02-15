# Repository 模块覆盖率分析

**当前覆盖率**: 4.1%
**目标覆盖率**: 80%
**差距**: 75.9%
**优先级**: 🔴 P0

---

## 未覆盖函数列表

### DatasetFieldRepository
| 函数 | 覆盖率 | 测试用例需求 |
|------|--------|-------------|
| Create | 0.0% | TestDatasetFieldRepository_Create |
| GetByID | 0.0% | TestDatasetFieldRepository_GetByID |
| List | 0.0% | TestDatasetFieldRepository_List |
| ListByType | 0.0% | TestDatasetFieldRepository_ListByType |
| Update | 0.0% | TestDatasetFieldRepository_Update |
| Delete | 0.0% | TestDatasetFieldRepository_Delete |
| DeleteComputedFields | 0.0% | TestDatasetFieldRepository_DeleteComputedFields |

### DatasetRepository
| 函数 | 覆盖率 | 测试用例需求 |
|------|--------|-------------|
| Create | 0.0% | TestDatasetRepository_Create |
| GetByID | 0.0% | TestDatasetRepository_GetByID |
| GetByIDWithFields | 0.0% | TestDatasetRepository_GetByIDWithFields |
| List | 0.0% | TestDatasetRepository_List |
| Update | 0.0% | TestDatasetRepository_Update |
| Delete | 0.0% | TestDatasetRepository_Delete |
| SoftDelete | 0.0% | TestDatasetRepository_SoftDelete |

### DatasetSourceRepository
| 函数 | 覆盖率 | 测试用例需求 |
|------|--------|-------------|
| Create | 0.0% | TestDatasetSourceRepository_Create |
| GetByID | 0.0% | TestDatasetSourceRepository_GetByID |
| List | 0.0% | TestDatasetSourceRepository_List |
| Update | 0.0% | TestDatasetSourceRepository_Update |
| Delete | 0.0% | TestDatasetSourceRepository_Delete |

### DatasourceRepository
| 函数 | 覆盖率 | 测试用例需求 |
|------|--------|-------------|
| Create | 0.0% | TestDatasourceRepository_Create |
| GetByID | 0.0% | TestDatasourceRepository_GetByID |
| List | 0.0% | TestDatasourceRepository_List |
| Update | 0.0% | TestDatasourceRepository_Update |
| Delete | 0.0% | TestDatasourceRepository_Delete |
| Search | 0.0% | TestDatasourceRepository_Search |
| Copy | 0.0% | TestDatasourceRepository_Copy |
| Move | 0.0% | TestDatasourceRepository_Move |
| Rename | 0.0% | TestDatasourceRepository_Rename |

### TenantRepository
| 函数 | 覆盖率 | 测试用例需求 |
|------|--------|-------------|
| GetByID | 0.0% | TestTenantRepository_GetByID |
| ListByUserID | 0.0% | TestTenantRepository_ListByUserID |

### UserRepository
| 函数 | 覆盖率 | 测试用例需求 |
|------|--------|-------------|
| GetByID | 部分覆盖 | TestUserRepository_GetByID |
| GetByUsername | 部分覆盖 | TestUserRepository_GetByUsername |
| Create | 0.0% | TestUserRepository_Create |
| Update | 0.0% | TestUserRepository_Update |
| Delete | 0.0% | TestUserRepository_Delete |

---

## 测试策略

### Mock 策略
由于 repository 层依赖 GORM DB，使用以下策略：
1. 使用 `gorm.io/driver/sqlite` 内存数据库进行测试
2. 或使用 `github.com/DATA-DOG/go-sqlmock` 进行 mock

### 预计测试用例数量
- DatasetFieldRepository: 7 函数 × 2-3 场景 = 15-20 用例
- DatasetRepository: 7 函数 × 2-3 场景 = 15-20 用例
- DatasetSourceRepository: 5 函数 × 2-3 场景 = 10-15 用例
- DatasourceRepository: 9 函数 × 2-3 场景 = 20-25 用例
- TenantRepository: 2 函数 × 2-3 场景 = 5-6 用例
- UserRepository: 5 函数 × 2-3 场景 = 10-15 用例

**总计**: 75-100 个测试用例

---

## 实现优先级

1. **高优先级** (核心 CRUD):
   - DatasourceRepository: Create, GetByID, List, Update, Delete
   - DatasetRepository: Create, GetByID, List, Update, Delete

2. **中优先级**:
   - DatasetFieldRepository: Create, List, Update, Delete
   - UserRepository: Create, GetByUsername

3. **低优先级**:
   - Search, Copy, Move, Rename 等辅助功能
