# JimuReport 贡献指南

感谢您对 JimuReport 项目的关注！我们欢迎任何形式的贡献。

## 目录

- [如何贡献](#如何贡献)
- [开发流程](#开发流程)
- [代码规范](#代码规范)
- [提交指南](#提交指南)
- [Pull Request 流程](#pull-request-流程)
- [问题报告](#问题报告)

## 如何贡献

### 贡献方式

您可以通过以下方式贡献：

1. **代码贡献**：提交新功能、Bug 修复
2. **文档改进**：完善文档、翻译
3. **问题报告**：报告 Bug、提出建议
4. **测试**：参与测试、提供反馈
5. **推广**：分享项目、帮助其他用户

### 贡献前

在贡献代码前，请先：

1. 阅读 [开发指南](DEVELOPMENT_GUIDE.md)
2. 搜索 [现有 Issues](../../issues)
3. 创建 Issue 讨论您的想法
4. 等待维护者确认

## 开发流程

### 1. Fork 项目

```bash
# Fork 项目到您的 GitHub 账户
# 克隆您的 Fork
git clone https://github.com/your-username/jimureport-go.git
cd jimureport-go
```

### 2. 创建分支

```bash
# 添加上游仓库
git remote add upstream https://github.com/jeecg/jimureport-go.git

# 同步最新代码
git fetch upstream
git checkout main
git merge upstream/main

# 创建功能分支
git checkout -b feature/your-feature-name
```

### 3. 开发和测试

```bash
# 安装依赖
make install

# 运行开发环境
make dev

# 运行测试
make test
```

### 4. 提交代码

遵循 [提交指南](#提交指南)：

```bash
git add .
git commit -m "feat: add dashboard component library"
```

### 5. 推送和创建 PR

```bash
# 推送到您的 Fork
git push origin feature/your-feature-name

# 在 GitHub 上创建 Pull Request
```

## 代码规范

### Go 代码规范

#### 命名约定

- 包名：小写单词，不使用下划线：`package dashboard`
- 文件名：小写下划线：`repository.go`
- 函数/变量：驼峰命名：`getUserByID`
- 常量：大写下划线：`MAX_CONNECTIONS`
- 接口：驼峰：`UserService`

#### 注释规范

```go
// Package level comment
package dashboard

// UserService 用户服务接口
type UserService interface {
    // GetUserByID 根据 ID 获取用户
    // 返回用户信息或错误
    GetUserByID(ctx context.Context, id string) (*User, error)
}

func (s *service) GetUserByID(ctx context.Context, id string) (*User, error) {
    // Complex logic
}
```

#### 错误处理

```go
// 使用 fmt.Errorf 包装错误
return fmt.Errorf("failed to connect: %w", err)

// 创建自定义错误
var ErrNotFound = errors.New("not found")
```

### Vue/TypeScript 代码规范

#### 组件命名

```vue
<template>
  <MyComponent />
</template>

<script setup lang="ts">
import MyComponent from './MyComponent.vue'
</script>
```

#### 类型定义

```typescript
// 使用 interface 定义数据结构
interface User {
  id: string
  name: string
}

// 使用 type 定义联合类型
type Status = 'pending' | 'active' | 'disabled'
```

#### 样式规范

```vue
<style scoped lang="scss">
// 使用 scoped 避免污染全局样式
.component {
  padding: 16px;

  &__header {
    font-size: 18px;
  }
}
</style>
```

### 数据库规范

```sql
-- 表名：小写下划线
CREATE TABLE reports (
  -- 字段名：小写下划线
  id VARCHAR(36) PRIMARY KEY,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 索引名：idx_开头
INDEX idx_tenant_id (tenant_id);
```

## 提交指南

### 提交信息格式

使用 [Conventional Commits](https://www.conventionalcommits.org/) 格式：

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Type 类型

- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更新
- `style`: 代码格式（不影响功能）
- `refactor`: 重构
- `perf`: 性能优化
- `test`: 测试相关
- `chore`: 构建工具、依赖更新

### 示例

```
feat(dashboard): add drag and drop support

Implement drag and drop functionality for dashboard
components to improve user experience.

- Add drag handlers to components
- Implement drop zones on canvas
- Update component positions

Closes #123
```

## Pull Request 流程

### PR 描述模板

创建 PR 时使用以下模板：

```markdown
## 变更类型
- [ ] Bug 修复
- [x] 新功能
- [ ] 破坏性变更
- [ ] 文档更新

## 变更描述
简要描述此 PR 的目的和内容

## 相关 Issues
Closes #123

## 测试
描述如何测试此变更：

1. 步骤一
2. 步骤二
3. 预期结果

## 截图（如适用）
添加截图说明变更

## 检查清单
- [ ] 代码通过所有测试
- [ ] 遵循代码规范
- [ ] 添加必要的文档
- [ ] 更新 CHANGELOG
- [ ] 没有 `console.log`
- [ ] 没有合并冲突
```

### PR 审查流程

1. 提交 PR 后等待维护者审查
2. 根据反馈修改代码
3. 所有检查通过后合并
4. PR 合并后删除分支

### Code Review 要点

- 代码逻辑是否正确
- 是否有性能问题
- 是否有安全隐患
- 测试覆盖是否足够
- 文档是否完整

## 问题报告

### Bug 报告

使用 Bug 模板：

```markdown
## Bug 描述
简要描述 Bug

## 重现步骤
1. 前往 '...'
2. 点击 '....'
3. 滚动到 '....'
4. 看到错误

## 预期行为
描述预期行为

## 实际行为
描述实际发生的行为

## 截图
如果适用，添加截图

## 环境
- OS: [e.g. Ubuntu 20.04]
- Browser: [e.g. Chrome 120]
- Version: [e.g. v1.0.0]
```

### 功能请求

使用功能请求模板：

```markdown
## 功能描述
简要描述您希望添加的功能

## 问题背景
为什么需要这个功能

## 建议的解决方案
描述您希望的实现方式

## 替代方案
描述您考虑过的其他方案

## 额外上下文
添加任何其他上下文或截图
```

## 开发者认证

### 贡献者许可协议

提交代码即表示您同意：

1. 您的代码采用 LGPL-3.0 许可证
2. 您有权提交此代码
3. 您的代码不会侵犯第三方知识产权

### CLA 签署

对于重大贡献，可能需要签署 CLA（Contributor License Agreement）。

## 发布流程

### 版本号规则

使用 [语义化版本](https://semver.org/)：

- 主版本号：不兼容的 API 变更
- 次版本号：向后兼容的功能新增
- 修订号：向后兼容的问题修复

示例：`1.2.3`

### 发布步骤

1. 更新版本号
2. 更新 CHANGELOG
3. 创建 Git 标签
4. 发布到 GitHub
5. 构建并推送 Docker 镜像

## 获取帮助

如有问题：

- 查看 [文档](../../docs)
- 提交 [Issue](../../issues)
- 加入 [Discord](https://discord.gg/jimureport)

## 致谢

感谢所有贡献者！您可以在 [CONTRIBUTORS.md](CONTRIBUTORS.md) 查看完整列表。

特别感谢：

- 代码贡献者
- 文档贡献者
- 测试贡献者
- Bug 报告者

## 行为准则

### 尊重

- 尊重不同的观点和经验
- 接受建设性批评
- 专注于对社区最有利的事情

### 包容

- 欢迎不同背景的人
- 使用包容性语言
- 避免排他性术语

### 协作

- 合作解决问题
- 乐于帮助他人
- 倾听并考虑不同意见

### 专业

- 尊重其他人的时间
- 清楚说明意图
- 适当的反馈

---

再次感谢您的贡献！
