# Design: Custom goReport Frontend

## Context

goReport 官方前端是闭源的，包含以下核心功能：
1. 报表设计器 - 拖拽式报表设计，类似 Excel 操作
2. 报表渲染器 - 报表预览、分页、打印
3. 大屏设计器 - 可视化大屏设计，支持多种组件
4. 图表编辑器 - 28 种图表类型的配置和预览
5. 数据源管理 - 可视化数据源配置和测试
6. 导出配置 - Excel/PDF/Word/Image 导出参数设置

**Constraints**:
- 必须与 Go 后端 API 完全兼容
- 必须提供与官方版本相似的用户体验
- 开发时间控制在 3-6 个月内
- 必须支持主流浏览器（Chrome, Firefox, Edge, Safari）

**Stakeholders**:
- 产品团队 - 定义功能和用户体验
- UI/UX 设计师 - 设计界面和交互
- 前端开发团队 - 实现前端功能
- 后端开发团队 - 提供 API 支持
- 测试团队 - 验证功能和质量

## Goals / Non-Goals

**Goals**:
- 提供完整的报表设计、渲染、大屏设计功能
- 确保与 Go 后端 API 100% 兼容
- 提供良好的用户体验和性能
- 实现代码的可维护性和可扩展性
- 支持多租户和权限管理

**Non-Goals**:
- 完全复制 goReport 官方 UI（避免专利冲突）
- 实现 goReport 的所有高级功能（如 AI 辅助）
- 支持所有 goReport 的高级特性（如复杂的联动规则）
- 实现移动端适配（后续迭代）

## Decisions

### Decision 1: 前端技术栈

**Choice**: Vue 3 + TypeScript + Vite + Element Plus

**Why**:
- Vue 3 是现代前端框架，性能好，生态系统成熟
- TypeScript 提供类型安全，减少运行时错误
- Vite 提供快速的开发体验和构建速度
- Element Plus 是成熟的 UI 组件库，提供丰富的组件

**Alternatives Considered**:
- React + TypeScript - 生态更丰富，但学习曲线更陡
- Angular - 太重，不适合快速开发
- 纯 JavaScript - 缺少类型检查，不利于大型项目

**Configuration**:
```json
{
  "name": "jimureport-frontend",
  "version": "1.0.0",
  "dependencies": {
    "vue": "^3.4.0",
    "typescript": "^5.3.0",
    "vite": "^5.0.0",
    "element-plus": "^2.5.0",
    "echarts": "^5.4.0",
    "axios": "^1.6.0",
    "pinia": "^2.1.0",
    "vue-router": "^4.2.0",
    "monaco-editor": "^0.45.0"
  }
}
```

---

### Decision 2: 报表设计器实现方式

**Choice**: 基于虚拟滚动和 Canvas 的拖拽设计器

**Why**:
- Canvas 可以处理大量单元格，性能更好
- 虚拟滚动只渲染可见区域，减少 DOM 节点
- 拖拽和选择更容易实现

**Alternatives Considered**:
- 纯 DOM 表格 - 简单但性能差，单元格多时会卡顿
- WebAssembly - 性能最好但开发复杂度高

**Implementation Details**:
```typescript
// 设计器核心组件
<template>
  <canvas
    ref="canvasRef"
    @mousedown="handleMouseDown"
    @mousemove="handleMouseMove"
    @mouseup="handleMouseUp"
  />
  <!-- 属性面板 -->
  <PropertyPanel
    :selected-cell="selectedCell"
    @change="handlePropertyChange"
  />
</template>
```

**Features**:
- 拖拽单元格
- 合并/拆分单元格
- 调整行高/列宽
- 数据绑定配置
- 表达式编辑
- 样式设置（字体、颜色、边框等）

---

### Decision 3: 图表库选择

**Choice**: ECharts 5.x

**Why**:
- 官方 goReport 使用 ECharts，兼容性好
- 功能强大，支持 28 种图表类型
- 文档完善，社区活跃
- 性能优秀

**Alternatives Considered**:
- Chart.js - 简单但功能有限
- Highcharts - 商业授权成本高
- Plotly - 功能强但体积大

**Supported Chart Types**:
- 柱形图 (Bar)
- 折线图 (Line)
- 饼图 (Pie)
- 散点图 (Scatter)
- 漏斗图 (Funnel)
- 雷达图 (Radar)
- 仪表盘 (Gauge)
- 地图 (Map)
- 等等...

---

### Decision 4: 状态管理

**Choice**: Pinia

**Why**:
- Vue 3 官方推荐
- TypeScript 支持好
- API 简洁，易于使用
- 支持开发者工具调试

**Store Structure**:
```typescript
// stores/report.ts
export const useReportStore = defineStore('report', () => {
  const reports = ref<Report[]>([])
  const currentReport = ref<Report | null>(null)
  const isLoading = ref(false)

  const fetchReports = async () => {
    isLoading.value = true
    try {
      const response = await api.getReports()
      reports.value = response.data
    } finally {
      isLoading.value = false
    }
  }

  return {
    reports,
    currentReport,
    isLoading,
    fetchReports
  }
})
```

---

### Decision 5: API 通信

**Choice**: Axios + 拦截器

**Why**:
- 成熟稳定，社区支持好
- 支持请求/响应拦截器
- 支持请求取消
- TypeScript 类型定义完整

**Implementation**:
```typescript
// api/client.ts
const client = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器 - 添加 Token
client.interceptors.request.use(
  (config) => {
    const token = getToken()
    if (token) {
      config.headers['X-Access-Token'] = token
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器 - 处理错误
client.interceptors.response.use(
  (response) => {
    return response.data
  },
  (error) => {
    if (error.response?.status === 401) {
      // 跳转登录页
      router.push('/login')
    }
    return Promise.reject(error)
  }
)
```

---

### Decision 6: 组件库

**Choice**: Element Plus

**Why**:
- Vue 3 官方推荐
- 组件丰富，满足需求
- 设计美观，主题定制方便
- 文档完善

**Custom Components**:
- 报表设计器组件
- 报表渲染器组件
- 大屏设计器组件
- 图表配置组件
- 数据源配置组件

---

## Architecture

### 项目结构

```
frontend/
├── public/
│   ├── favicon.ico
│   └── index.html
├── src/
│   ├── api/
│   │   ├── client.ts           # Axios 客户端
│   │   ├── auth.ts             # 认证 API
│   │   ├── report.ts           # 报表 API
│   │   ├── dashboard.ts        # 大屏 API
│   │   ├── export.ts           # 导出 API
│   │   └── datasource.ts      # 数据源 API
│   ├── assets/
│   │   ├── styles/
│   │   │   ├── variables.scss   # 样式变量
│   │   │   └── mixins.scss    # 样式混入
│   │   └── images/
│   ├── components/
│   │   ├── report/
│   │   │   ├── Designer.vue     # 报表设计器
│   │   │   ├── Renderer.vue     # 报表渲染器
│   │   │   ├── CellEditor.vue   # 单元格编辑器
│   │   │   └── PropertyPanel.vue # 属性面板
│   │   ├── dashboard/
│   │   │   ├── Designer.vue     # 大屏设计器
│   │   │   ├── ComponentList.vue # 组件列表
│   │   │   └── LayerPanel.vue  # 图层面板
│   │   ├── chart/
│   │   │   ├── ChartEditor.vue  # 图表编辑器
│   │   │   └── ChartPreview.vue # 图表预览
│   │   └── common/
│   │       ├── Layout.vue       # 布局组件
│   │       └── Header.vue      # 头部组件
│   ├── stores/
│   │   ├── auth.ts            # 认证状态
│   │   ├── report.ts          # 报表状态
│   │   ├── dashboard.ts       # 大屏状态
│   │   └── ui.ts             # UI 状态
│   ├── router/
│   │   └── index.ts          # 路由配置
│   ├── views/
│   │   ├── Login.vue          # 登录页
│   │   ├── ReportList.vue     # 报表列表
│   │   ├── ReportDesigner.vue # 报表设计
│   │   ├── DashboardList.vue  # 大屏列表
│   │   └── DashboardDesigner.vue # 大屏设计
│   ├── types/
│   │   ├── report.ts         # 报表类型
│   │   ├── dashboard.ts      # 大屏类型
│   │   └── api.ts           # API 类型
│   ├── utils/
│   │   ├── auth.ts          # 认证工具
│   │   ├── storage.ts       # 存储工具
│   │   └── request.ts      # 请求工具
│   ├── App.vue
│   └── main.ts
├── package.json
├── tsconfig.json
├── vite.config.ts
└── README.md
```

### 核心模块

#### 1. 报表设计器模块

**功能**:
- 虚拟滚动的 Canvas 设计器
- 单元格选择和编辑
- 拖拽调整行高/列宽
- 合并/拆分单元格
- 数据绑定配置
- 表达式编辑（集成 Monaco Editor）
- 样式设置

**API 调用**:
- 创建/保存报表: `POST /api/v1/jmreport/create`
- 更新报表: `POST /api/v1/jmreport/update`
- 预览报表: `POST /api/v1/jmreport/preview`
- 加载报表: `GET /api/v1/jmreport/get?id={id}`

#### 2. 报表渲染器模块

**功能**:
- 报表预览
- 分页显示
- 打印预览
- 参数查询
- 导出触发

**API 调用**:
- 渲染报表: `POST /api/v1/jmreport/render`
- 导出报表: `POST /api/v1/jmreport/export`

#### 3. 大屏设计器模块

**功能**:
- 组件拖拽
- 组件属性配置
- 图层管理
- 数据绑定
- 预览

**API 调用**:
- 创建/更新页面: `POST /api/v1/drag/page/create|update`
- 加载页面: `GET /api/v1/drag/page/get?id={id}`

#### 4. 图表编辑器模块

**功能**:
- 28 种图表类型选择
- 数据源配置
- 图表属性配置
- 实时预览

**API 调用**:
- 图表数据: `POST /api/v1/jmreport/render` (用于预览)

---

## Risks / Trade-offs

### Risk 1: UI 专利冲突

**Description**:
自定义前端可能与 goReport 官方 UI 存在专利冲突。

**Mitigation**:
- 不完全复制官方 UI 设计
- 采用不同的视觉风格
- 专注于功能实现而非 UI 相似
- 咨询法律顾问

### Risk 2: 功能兼容性问题

**Description**:
用户习惯使用官方 goReport，自定义前端可能不兼容。

**Mitigation**:
- 优先实现核心功能
- 提供迁移文档
- 收集用户反馈并快速迭代
- 保持 API 100% 兼容

### Risk 3: 开发周期长

**Description**:
完整的前端开发需要 3-6 个月，可能影响项目交付。

**Mitigation**:
- 采用分阶段交付策略
- 优先开发核心功能（报表设计器、渲染器）
- 后续迭代添加大屏设计器和高级功能

### Risk 4: 性能问题

**Description**:
报表设计器需要处理大量单元格，可能出现性能问题。

**Mitigation**:
- 使用 Canvas 而非 DOM
- 实现虚拟滚动
- 优化渲染逻辑
- 使用 Web Worker 处理复杂计算

---

## Migration Plan

### Phase 1: 基础架构（2-4 周）

**任务**:
- [ ] 搭建 Vue 3 + TypeScript + Vite 项目
- [ ] 配置 Element Plus 和 ECharts
- [ ] 实现路由和状态管理
- [ ] 实现认证模块（登录、登出、Token 管理）
- [ ] 实现基础布局（Header、Sidebar、Main）

**交付物**:
- 可运行的 Vue 项目
- 登录功能
- 基础布局

---

### Phase 2: 报表列表和数据管理（2-3 周）

**任务**:
- [ ] 实现报表列表页
- [ ] 实现报表 CRUD API 调用
- [ ] 实现数据源管理页
- [ ] 实现数据源配置和测试

**交付物**:
- 报表列表功能
- 数据源管理功能

---

### Phase 3: 报表设计器（8-12 周）

**任务**:
- [ ] 实现 Canvas 设计器基础框架
- [ ] 实现单元格选择和编辑
- [ ] 实现拖拽调整行高/列宽
- [ ] 实现合并/拆分单元格
- [ ] 实现数据绑定配置
- [ ] 实现表达式编辑器（Monaco Editor）
- [ ] 实现样式设置面板

**交付物**:
- 完整的报表设计器
- 支持基本的报表设计功能

---

### Phase 4: 报表渲染器（4-6 周）

**任务**:
- [ ] 实现报表预览功能
- [ ] 实现分页显示
- [ ] 实现打印预览
- [ ] 实现参数查询
- [ ] 集成导出功能

**交付物**:
- 完整的报表渲染器
- 支持预览、分页、打印、导出

---

### Phase 5: 大屏设计器（6-8 周）

**任务**:
- [ ] 实现大屏设计器基础框架
- [ ] 实现组件拖拽
- [ ] 实现组件属性配置
- [ ] 实现图层管理
- [ ] 实现数据绑定
- [ ] 实现预览

**交付物**:
- 完整的大屏设计器
- 支持基本的大屏设计功能

---

### Phase 6: 图表编辑器（4-6 周）

**任务**:
- [ ] 实现图表类型选择器
- [ ] 实现 ECharts 集成
- [ ] 实现数据源配置
- [ ] 实现图表属性配置
- [ ] 实现实时预览

**交付物**:
- 完整的图表编辑器
- 支持 28 种图表类型

---

### Phase 7: 测试和优化（4-6 周）

**任务**:
- [ ] 单元测试
- [ ] 集成测试
- [ ] 性能优化
- [ ] 兼容性测试（浏览器）
- [ ] 用户体验优化

**交付物**:
- 测试报告
- 性能优化报告
- 用户手册

---

### Phase 8: 部署和文档（2-3 周）

**任务**:
- [ ] 构建优化
- [ ] 部署配置
- [ ] 用户文档
- [ ] 开发者文档
- [ ] 迁移指南

**交付物**:
- 可部署的前端版本
- 完整的文档

---

## Open Questions

1. **UI 设计风格**: 是否需要完全模仿 goReport 官方 UI，还是采用自己的设计风格？

2. **功能优先级**: 是否所有功能都需要在第一阶段实现，还是可以分阶段交付？

3. **图表类型**: 是否需要支持全部 28 种图表类型，还是可以优先支持常用的 10-15 种？

4. **移动端适配**: 是否需要支持移动端访问？

5. **国际化**: 是否需要支持多语言（中英文）？

6. **主题定制**: 是否需要支持主题切换（亮色/暗色）？

7. **浏览器兼容性**: 需要支持哪些浏览器？最低版本要求？

8. **性能指标**: 报表设计器支持的最大单元格数量？预览加载时间要求？

9. **导出格式**: 导出功能需要支持哪些格式？优先级？

10. **权限系统**: 前端需要实现哪些权限控制？细粒度到字段级别吗？
