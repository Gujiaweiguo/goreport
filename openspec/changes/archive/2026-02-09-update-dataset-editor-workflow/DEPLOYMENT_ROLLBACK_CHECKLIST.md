# 数据集编辑器工作流升级 - 部署与回滚清单

## 部署前准备

### 1. 功能验证
- [ ] 后端 API 测试通过
  - [ ] `POST /api/v1/datasets/:id/fields` - 批量字段更新
  - [ ] `POST /api/v1/datasets/:id/fields` - 创建分组字段
  - [ ] `GET /api/v1/datasets/:id/schema` - 返回分组元数据
  - [ ] `POST /api/v1/datasets/:id/data` - 支持 GROUP BY 和聚合
- [ ] 前端构建成功
  - [ ] `npm run build` 无错误
  - [ ] TypeScript 类型检查通过
- [ ] 前端功能测试
  - [ ] 保存按钮：成功后停留在编辑页
  - [ ] 保存并返回按钮：成功后导航到列表页
  - [ ] 数据预览/批量管理标签切换正常
  - [ ] 批量字段更新功能正常
  - [ ] 分组字段创建表单正常

### 2. 环境检查
- [ ] 后端依赖：Go 1.22+, Gin, GORM 已安装
- [ ] 前端依赖：Vue 3, TypeScript, Element Plus 已安装
- [ ] 数据库迁移已应用（如需要）
- [ ] Redis 缓存配置正确（如使用缓存）

### 3. 兼容性检查
- [ ] 现有数据集数据兼容新字段结构
- [ ] 现有 API 调用不会因新参数而失败
- [ ] 向后兼容性：旧客户端仍可使用数据集编辑

## 部署步骤

### 1. 后端部署
```bash
# 1. 备份当前版本
git tag backup-before-dataset-workflow-update

# 2. 构建新版本
cd backend
go build -o bin/server cmd/server/main.go

# 3. 部署到测试环境
# （根据部署流程进行）

# 4. 验证后端健康检查
curl http://localhost:8085/health
```

### 2. 前端部署
```bash
# 1. 备份当前版本
cd frontend
git tag backup-before-dataset-workflow-update

# 2. 构建生产版本
npm run build

# 3. 部署构建产物
# （将 dist/ 目录部署到静态服务器）

# 4. 验证前端可访问
curl http://localhost:3000/
```

### 3. 灰度发布（可选但推荐）
- [ ] 部署到 10% 用户群
- [ ] 监控错误日志和用户反馈 24 小时
- [ ] 如无问题，逐步扩大到 50% -> 100%

## 部署后验证

### 1. 功能验证
- [ ] 数据集编辑页正常加载
- [ ] 保存功能正常工作
- [ ] 保存并返回功能正常工作
- [ ] 批量管理标签页可用
- [ ] 分组字段创建功能可用
- [ ] 数据预览功能正常

### 2. 性能验证
- [ ] 页面加载时间 < 3 秒
- [ ] 批量更新响应时间 < 2 秒
- [ ] 缓存命中率保持良好水平

### 3. 错误监控
- [ ] 监控后端日志：无异常错误
- [ ] 监控前端控制台：无 JavaScript 错误
- [ ] 检查用户反馈渠道：无重大投诉

## 回滚计划

### 触发条件
以下情况之一发生时，应考虑回滚：
1. 发现关键功能故障（如保存失败）
2. 性能严重下降（> 50%）
3. 数据损坏或丢失
4. 用户反馈负面强烈（> 10% 投诉率）
5. 安全漏洞暴露

### 回滚步骤

#### 选项 1：前端回滚（通过特性开关）
```typescript
// 在 frontend/src/config.ts 中禁用新工作流
export const enableNewDatasetWorkflow = false
```
- [ ] 前端：修改配置开关，关闭新编辑页
- [ ] 重新部署前端
- [ ] 验证旧编辑页恢复可用

#### 选项 2：后端回滚（保留旧参数兼容）
```go
// 在 backend/internal/dataset/handler.go 中禁用新 API 端点
func (h *Handler) RegisterRoutes(r *gin.Engine) {
    // 注释掉新端点
    // r.POST("/datasets/:id/fields/batch", h.BatchUpdateFields)
    // r.POST("/datasets/:id/fields/grouping", h.CreateGroupingField)
}
```
- [ ] 后端：禁用新批量/分组 API 端点
- [ ] 重新部署后端
- [ ] 验证旧 API 行为恢复

#### 选项 3：完整回滚（Git 回滚）
```bash
# 1. 回滚到备份标签
git checkout backup-before-dataset-workflow-update

# 2. 重新构建
cd backend && go build -o bin/server cmd/server/main.go
cd frontend && npm run build

# 3. 重新部署
# （按照部署步骤进行）
```

### 回滚验证
- [ ] 旧版本功能正常恢复
- [ ] 无数据丢失或损坏
- [ ] 系统性能恢复正常
- [ ] 通知用户关于回滚

## 回滚后跟进

### 1. 问题分析
- [ ] 记录回滚原因和时间
- [ ] 分析错误日志和监控数据
- [ ] 识别根本原因

### 2. 修复与重新部署
- [ ] 在开发环境修复问题
- [ ] 进行充分的回归测试
- [ ] 准备新的部署计划
- [ ] 考虑更严格的灰度发布策略

### 3. 沟通
- [ ] 向用户解释回滚原因
- [ ] 提供修复时间表
- [ ] 更新文档和发布说明

## 监控指标

### 部署后 72 小时内重点监控

**功能指标**
- 数据集保存成功率 > 95%
- 批量更新成功率 > 90%
- 分组字段创建成功率 > 95%

**性能指标**
- API 响应时间 P95 < 1 秒
- 页面加载时间 P95 < 2 秒
- 错误率 < 0.5%

**用户体验指标**
- 用户满意度评分 > 4/5
- 支持工单增长率 < 20%

### 告警阈值
- 功能成功率 < 90% → 立即告警
- API 响应时间 P95 > 3 秒 → 立即告警
- 错误率 > 1% → 立即告警

## 联系信息

**部署联系人**
- 前端负责人：[待填写]
- 后端负责人：[待填写]
- DevOps 负责人：[待填写]

**应急联系人**
- 产品经理：[待填写]
- 技术负责人：[待填写]

## 附录

### 相关文档
- OpenSpec 变更提案：`/openspec/changes/update-dataset-editor-workflow/proposal.md`
- 设计文档：`/openspec/changes/update-dataset-editor-workflow/design.md`
- 任务清单：`/openspec/changes/update-dataset-editor-workflow/tasks.md`

### 已知问题
- 无

### 版本信息
- 前端版本：[待填写]
- 后端版本：[待填写]
- 部署日期：[待填写]
- 部署人员：[待填写]
