# Change: Add Comprehensive Test Coverage

## Why

当前系统已有 22 个能力规格定义的核心功能，但缺乏系统化的测试覆盖率要求和验证机制。为确保系统质量和稳定性，需要建立全面的测试覆盖体系，包括：

1. **测试覆盖率目标不明确** - 各模块测试覆盖率参差不齐，缺乏统一目标
2. **测试用例与规格未对齐** - 现有测试未系统性地覆盖所有 OpenSpec 定义的需求场景
3. **缺少 E2E 测试** - 端到端业务流程测试覆盖不足
4. **测试数据管理分散** - 缺乏统一的测试数据工厂和清理机制

## What Changes

### 新增测试能力规格
- 定义测试覆盖率要求（后端 ≥80%，前端 ≥70%）
- 定义测试类型层次（单元测试、集成测试、E2E 测试）
- 定义测试场景与 OpenSpec 需求的映射关系

### 后端测试增强
- **BREAKING** 添加测试覆盖率 CI 门禁（< 80% 时构建失败）
- 为所有 OpenSpec 定义的需求场景补充测试用例
- 完善 `testutil/` 测试工具包

### 前端测试增强
- 添加核心组件测试覆盖
- 添加 API 集成测试
- 添加关键用户流程 E2E 测试

### CI/CD 集成
- 测试覆盖率报告自动生成
- 覆盖率趋势追踪

## Impact

### Affected specs
- **新增**: `testing-coverage` - 测试覆盖率能力规格

### Affected code
- `backend/internal/*/` - 各领域测试用例补充
- `backend/internal/testutil/` - 测试工具增强
- `frontend/src/**/*.test.ts` - 前端测试用例补充
- `scripts/ci/` - CI 脚本更新
- `.github/workflows/` - GitHub Actions 配置

### Risks
- 测试覆盖率门禁可能导致初期构建失败（需要逐步提升覆盖率）
- E2E 测试增加 CI 运行时间

### Migration
- 分阶段实施，优先覆盖 P0 模块
- 提供测试覆盖率基线报告作为起点
