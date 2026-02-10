# 浏览器兼容性测试指南

本文档提供 goReport 浏览器兼容性测试的详细指南。

## 测试目标

确保 goReport 在主流浏览器和设备上正常运行。

## 测试环境

### 现代浏览器测试

| 浏览器 | 版本要求 | 测试重点 |
|--------|----------|---------|
| Chrome | 90+ | 全部功能 |
| Firefox | 88+ | 全部功能 |
| Safari | 14+ | 全部功能 |
| Edge | 90+ | 全部功能 |
| Opera | 76+ | 基础功能 |

### 移动端浏览器测试

| 平台 | 浏览器 | 版本要求 | 测试重点 |
|------|--------|----------|---------|
| iOS | Safari | 14+ | 触摸交互、响应式 |
| iOS | Chrome | 90+ | 触摸交互、响应式 |
| Android | Chrome | 90+ | 触摸交互、响应式 |
| Android | Firefox | 88+ | 触摸交互、响应式 |

## 功能测试清单

### 核心功能

#### 1. 用户认证

- [ ] 登录页面正常显示
- [ ] 用户名密码输入框可用
- [ ] 登录按钮可点击
- [ ] 登录成功后正确跳转
- [ ] 登录失败显示错误提示
- [ ] Token 存储正确
- [ ] 刷新页面保持登录状态

#### 2. 报表管理

- [ ] 报表列表正常加载
- [ ] 分页功能正常
- [ ] 搜索/筛选功能正常
- [ ] 新建报表按钮可用
- [ ] 编辑报表按钮可用
- [ ] 删除报表功能正常
- [ ] 报表预览正常显示

#### 3. 大屏设计

- [ ] 大屏列表正常显示
- [ ] 组件库正常加载
- [ ] 拖拽功能流畅
- [ ] 属性面板正常显示
- [ ] 图层管理功能正常
- [ ] 预览功能正常
- [ ] 全屏功能正常
- [ ] 保存功能正常

#### 4. 图表编辑

- [ ] 图表类型选择正常
- [ ] 图表预览正常显示
- [ ] 数据配置面板可用
- [ ] 属性配置面板可用
- [ ] 实时更新正常
- [ ] 图表交互正常（提示、缩放）

#### 5. 数据源管理

- [ ] 数据源列表正常显示
- [ ] 添加数据源表单正常
- [ ] 测试连接功能正常
- [ ] 元数据查询正常
- [ ] 编辑数据源功能正常
- [ ] 删除数据源功能正常

### 高级功能

#### Canvas 渲染

- [ ] 报表画布正常渲染
- [ ] 单元格选择正常
- [ ] 拖拽调整大小正常
- [ ] 合并单元格功能正常
- [ ] 样式设置正常
- [ ] 缩放功能正常

#### ECharts 图表

- [ ] 柱状图正常显示
- [ ] 折线图正常显示
- [ ] 饼图正常显示
- [ ] 地图正常显示
- [ ] 图表交互正常（hover、click）
- [ ] 图例切换正常
- [ ] 数据缩放正常

#### 文件导出

- [ ] Excel 导出功能正常
- [ ] PDF 导出功能正常
- [ ] 导出文件格式正确
- [ ] 导出文件命名正确
- [ ] 大数据量导出不超时

## CSS 兼容性测试

### 核心样式

- [ ] Flexbox 布局正常
- [ ] Grid 布局正常
- [ ] CSS 变量正常工作
- [ ] 动画效果流畅
- [ ] 过渡效果正常
- [ ] 阴影效果正常

### 响应式设计

- [ ] 桌面端（1920x1080）正常
- [ ] 笔记本（1366x768）正常
- [ ] 平板（768x1024）正常
- [ ] 手机（375x667）正常
- [ ] 横屏/竖屏切换正常

## JavaScript 兼容性测试

### ES6+ 特性

- [ ] 箭头函数正常
- [ ] async/await 正常
- [ ] Promise 正常工作
- [ ] 解构赋值正常
- [ ] 模板字符串正常
- [ ] 类语法正常

### Web APIs

- [ ] Fetch API 正常
- [ ] LocalStorage 正常
- [ ] SessionStorage 正常
- [ ] Clipboard API 正常
- [ ] WebSocket（如使用）正常
- [ ] Web Workers（如使用）正常

## 性能测试

### 加载性能

| 指标 | 目标 | 浏览器 |
|-------|------|--------|
| 首屏渲染 (FCP) | < 2s | ⬜ Chrome ⬜ Firefox ⬜ Safari ⬜ Edge |
| 最大内容绘制 (LCP) | < 3s | ⬜ Chrome ⬜ Firefox ⬜ Safari ⬜ Edge |
| 首次输入延迟 (FID) | < 100ms | ⬜ Chrome ⬜ Firefox ⬜ Safari ⬜ Edge |
| 累计布局偏移 (CLS) | < 0.1 | ⬜ Chrome ⬜ Firefox ⬜ Safari ⬜ Edge |

### 运行时性能

- [ ] 页面滚动流畅（60fps）
- [ ] 动画效果流畅
- [ ] 大数据量加载不卡顿
- [ ] 复杂计算不阻塞 UI
- [ ] 内存使用合理

## 已知兼容性问题

### CSS 兼容性

| 问题 | 影响范围 | 解决方案 |
|------|---------|---------|
| Grid 部分支持 | IE 11- | 使用 polyfill |
| CSS 变量 | IE 11- | 使用预处理器 |
| 滚动捕捉 | Firefox | 使用 event delegation |

### JavaScript 兼容性

| 问题 | 影响范围 | 解决方案 |
|------|---------|---------|
| 箭头函数 | IE 11- | 使用 babel 转译 |
| async/await | IE 11- | 使用 polyfill |
| optional chaining | 旧浏览器 | 不使用或转译 |

## 测试工具

### 自动化测试工具

1. **BrowserStack**
   - 网址：https://www.browserstack.com/
   - 支持浏览器：150+
   - 支持设备：2000+

2. **LambdaTest**
   - 网址：https://www.lambdatest.com/
   - 实时测试
   - 自动化测试

3. **Sauce Labs**
   - 网址：https://saucelabs.com/
   - 云端测试环境

4. **Playwright**
   - 本地自动化测试
   - 支持多浏览器

### 本地测试工具

1. **浏览器开发者工具**
   - Chrome DevTools
   - Firefox Developer Tools
   - Safari Web Inspector
   - Edge DevTools

2. **Lighthouse**
   ```bash
   npm install -g lighthouse
   lighthouse http://localhost:3000 --view
   ```

3. **WebPageTest**
   - 网址：https://www.webpagetest.org/
   - 多地测试
   - 性能报告

## 测试流程

### 1. 环境准备

```bash
# 启动开发环境
make dev

# 确认服务运行
curl http://localhost:3000
curl http://localhost:8085/health
```

### 2. 功能测试

1. 在每个浏览器中打开 `http://localhost:3000`
2. 登录系统
3. 按照功能测试清单逐项测试
4. 记录发现的问题

### 3. 性能测试

1. 使用 Lighthouse 运行性能测试
2. 记录各项指标
3. 对比不同浏览器性能
4. 识别性能瓶颈

### 4. 问题记录

使用模板记录问题：

```markdown
## 问题描述

**浏览器**：Chrome 120
**版本**：120.0.6099.109
**OS**：Windows 11

**重现步骤**
1. 步骤一
2. 步骤二

**预期行为**
描述预期结果

**实际行为**
描述实际结果

**截图**
添加截图

**控制台错误**
```
Console Errors:
Error 1: ...
Error 2: ...
```
```

## 问题优先级

| 优先级 | 定义 | 响应时间 |
|-------|------|---------|
| P0 | 阻塞性问题，无法使用 | 24h |
| P1 | 严重问题，主要功能受影响 | 72h |
| P2 | 中等问题，部分功能受限 | 1 周 |
| P3 | 轻微问题，美观性问题 | 2 周 |

## 测试报告模板

```markdown
# 浏览器兼容性测试报告

**测试日期**：YYYY-MM-DD
**测试人员**：姓名
**测试版本**：v1.0.0

## 测试环境

| 浏览器 | 版本 | OS | 结果 |
|--------|------|-----|------|
| Chrome | 120 | Windows 11 | ✅ 通过 |
| Firefox | 115 | Windows 11 | ⚠️ 部分 |
| Safari | 17 | macOS 14 | ❌ 失败 |
| Edge | 120 | Windows 11 | ✅ 通过 |

## 功能测试结果

### 核心功能
- 用户认证：✅ Chrome ✅ Firefox ⚠️ Safari ❌ Edge
- 报表管理：✅ Chrome ✅ Firefox ✅ Safari ✅ Edge
- 大屏设计：✅ Chrome ✅ Firefox ⚠️ Safari ✅ Edge
- 图表编辑：✅ Chrome ✅ Firefox ✅ Safari ✅ Edge
- 数据源管理：✅ Chrome ✅ Firefox ✅ Safari ✅ Edge

### 高级功能
- Canvas 渲染：✅ Chrome ✅ Firefox ❌ Safari ✅ Edge
- ECharts 图表：✅ Chrome ✅ Firefox ✅ Safari ✅ Edge
- 文件导出：✅ Chrome ✅ Firefox ✅ Safari ✅ Edge

## 性能测试结果

| 浏览器 | FCP | LCP | FID | CLS |
|--------|-----|-----|-----|-----|
| Chrome | 1.2s | 2.1s | 80ms | 0.05 |
| Firefox | 1.5s | 2.3s | 95ms | 0.08 |
| Safari | 1.8s | 2.5s | 110ms | 0.1 |

## 发现的问题

### P0 问题
1. [问题描述](链接)
   - 浏览器：Safari
   - 影响范围：大屏设计
   - 预计修复：2024-02-10

### P1 问题
2. [问题描述](链接)
   - 浏览器：Firefox
   - 影响范围：Canvas 渲染
   - 预计修复：2024-02-15

## 建议

1. 优化 Canvas 渲染性能
2. 添加 Safari polyfill
3. 改进 Firefox 兼容性
4. 优化移动端性能
```

## 常见问题处理

### Safari 特有问题

#### 1. IndexedDB 限制
- **问题**：Safari IndexedDB 配额较小
- **解决**：使用 IndexedDB 封装，自动降级到 LocalStorage

#### 2. Flexbox 旧版本
- **问题**：Safari 14 以下 flexbox 部分问题
- **解决**：添加 vendor prefix

### Firefox 特有问题

#### 1. WebGL 性能
- **问题**：Firefox WebGL 性能较差
- **解决**：检测性能，自动降级到 Canvas 2D

#### 2. 触摸事件
- **问题**：Firefox 触摸事件支持不完整
- **解决**：使用 Pointer Events API

### Chrome 特有问题

#### 1. 内存泄漏
- **问题**：Chrome 长时间使用内存增长
- **解决**：定期清理事件监听器，使用 WeakMap

## 获取帮助

如遇到兼容性问题：

- 查看 [浏览器兼容性文档](https://developer.mozilla.org/zh-CN/docs/Web/Compatibility)
- 提交 Issue：GitHub Issues (标签：browser-compatibility)
- 联系支持：`weiguogu@163.com`

---

**注意**：本指南应与每次发布前执行，确保跨浏览器兼容性。
