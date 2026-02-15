# frontend/src/ AGENTS.md

**Generated:** 2026-02-14T11:22:56Z

## OVERVIEW

Vue 3 应用，支持报表设计器 (Canvas)、表达式编辑器 (Monaco)、图表渲染 (ECharts)。

## STRUCTURE

```
src/
├── api/            # HTTP 客户端 (每领域一个文件)
├── components/     # 组件按领域划分
│   ├── chart/      # ECharts 包装
│   ├── common/     # 通用组件
│   ├── dataset/    # 数据集相关
│   ├── report/     # 报表相关
│   └── dashboard/  # 仪表盘相关
├── canvas/         # Canvas 报表设计器
├── router/         # Vue Router 配置
├── stores/         # Pinia 状态管理
├── types/          # TypeScript 类型
├── views/          # 页面组件
└── tests/          # 测试配置和 setup
```

## WHERE TO LOOK

| Task | Location |
|------|----------|
| 添加新 API | `api/[domain].ts` |
| 页面组件 | `views/[Page].vue` |
| Canvas 设计器 | `canvas/` |
| 表达式编辑器 | Monaco 集成在相关组件中 |
| 测试 mock | `tests/setup.ts` |
| 路由配置 | `router/index.ts` |

## CONVENTIONS

### API 模式
每个领域一个文件:
```typescript
// api/dataset.ts
export const datasetApi = {
  list: () => client.get('/datasets'),
  create: (data) => client.post('/datasets', data),
}
```

### 组件命名
- 组件: `PascalCase.vue`
- 测试: `PascalCase.test.ts` (同目录)

## UNIQUE STYLES

### Canvas 2D Mock (测试)
`tests/setup.ts` 提供完整的 Canvas mock:
- 40+ 方法 (fillRect, measureText, createLinearGradient, 等)
- 支持 ReportDesigner 组件测试

### v-loading 指令 Mock
```typescript
const vLoading = {
  mounted: () => {},
  updated: () => {},
  unmounted: () => {},
}
// 在 mount 配置中: directives: { loading: vLoading }
```

### 浏览器 API Mock
`tests/setup.ts` 统一 mock:
- `window.location` (防止 jsdom 导航错误)
- `ResizeObserver`, `IntersectionObserver`
- `matchMedia`

### ECharts 集成
- `components/chart/EChartsComponent.vue` - 基础包装
- `components/chart/EChartsRenderer.vue` - 渲染器

## ANTI-PATTERNS

| 禁止 | 原因 |
|------|------|
| 直接修改 props | Vue 单向数据流 |
| 未 stub Element Plus 组件 | 测试会因 DOM 依赖失败 |
| 忘记 mock canvas | ReportDesigner 测试崩溃 |
