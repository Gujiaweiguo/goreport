# Change: Build Custom goReport Frontend

## Why

goReport 官方前端是闭源的，无法通过开源方式获取。为了让 Go 后端成为完整的解决方案，需要开发自定义前端来替代官方前端，提供报表设计器、渲染器和大屏设计器等核心功能。

## What Changes

- [ ] 开发自定义的报表设计器（拖拽设计、属性编辑、数据绑定）
- [ ] 开发自定义的报表渲染器（预览、分页、打印）
- [ ] 开发自定义的大屏设计器（组件拖拽、图层管理）
- [ ] 实现 28 种图表类型的编辑器
- [ ] 实现导出配置界面
- [ ] 实现数据源管理界面

**Impact on Specs**:
- 新增 capability: `report-designer-ui`
- 新增 capability: `report-renderer-ui`
- 新增 capability: `bi-dashboard-ui`
- 新增 capability: `chart-editor-ui`
- 修改 capability: `embedding-integration` (添加前端集成文档)

## Impact

**Affected Specs**:
- `specs/report-designer/spec.md` - 新增 UI 需求
- `specs/report-rendering/spec.md` - 新增 UI 需求
- `specs/bi-dashboard/spec.md` - 新增 UI 需求
- `specs/embedding-integration/spec.md` - 新增前端集成指南

**Affected Code**:
- `jimureport-go/internal/httpserver/routes.go` - 添加前端静态资源路由
- `jimureport-go/static/` - 新增前端资源目录
- `frontend/` - 新增前端项目根目录

**Estimated Effort**:
- 前端开发: 400-600 人日
- UI/UX 设计: 50-100 人日
- 测试: 100-150 人日
- **总计**: 550-850 人日（约 2.5-4 人年）

**Risks**:
- UI 专利风险：可能与 goReport 有专利冲突
- 功能兼容性：难以保证与官方版本完全兼容
- 维护成本高：需要持续维护和更新
- 开发周期长：可能影响项目交付时间

**Opportunities**:
- 完全自主可控
- 可以根据需求定制
- 不依赖第三方资源
- 可以优化用户体验

## Dependencies

- 需要 goReport Go 后端 API 完整实现
- 需要 UI/UX 设计师参与
- 需要前端开发团队
- 需要测试团队

**Blocking Issues**:
- 目前导出文件生成功能未实现（需要先完成）
- 图表渲染需要依赖 ECharts 或类似库
