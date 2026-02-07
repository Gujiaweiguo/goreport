# goReport 路由清单

## 1. 发现的路由

### 1.1 认证相关路由 (本地实现)

| 路由 | 方法 | 描述 | 实现 |
|------|------|------|------|
| `/doLogin` | GET | 登录 | `LoginController.java` |
| `/` | GET | 首页，重定向到报表工作台 | `LoginController.java` |
| `/isLogin` | GET | 查询登录状态 | `LoginController.java` |
| `/logout` | GET | 退出登录 | `LoginController.java` |
| `/login/login.html` | GET | 登录页面 | 静态资源 |

### 1.2 goReport 报表路由 (由 starter 提供)

| 路由前缀 | 描述 | 功能 |
|----------|------|------|
| `/jmreport/list` | 报表工作台入口 | 主页面 |
| `/jmreport/*` | 报表相关路由 | 包括：<br>- 报表设计器<br>- 报表预览和渲染<br>- 报表导出/打印<br>- 报表模板 CRUD<br>- 数据源管理 |

### 1.3 JimuBI 大屏/仪表盘路由 (由 starter 提供)

| 路由前缀 | 描述 | 功能 |
|----------|------|------|
| `/drag/list` | 仪表盘/大屏工作台入口 | 主页面 |
| `/drag/*` | 仪表盘/大屏路由 | 包括：<br>- 大屏/仪表盘设计器<br>- 页面组件管理<br>- 数据集管理<br>- 数据源管理<br>- 分享功能 |

## 2. API 接口约定

### 2.1 认证接口 (JmReportTokenServiceI)

goReport 通过 `JmReportTokenServiceI` 接口进行认证集成：

| 方法 | 描述 |
|------|------|
| `getToken(HttpServletRequest request)` | 从请求中获取 token |
| `getUsername(String token)` | 根据 token 获取用户名 |
| `getRoles(String token)` | 获取用户角色 |
| `getPermissions(String token)` | 获取用户权限 |
| `verifyToken(String token)` | 验证 token 有效性 |
| `getTenantId()` | 获取租户 ID |

### 2.2 当前实现 (goReportTokenServiceImpl.java)

**Token 获取优先级**：
1. 从 Sa-Token 上下文获取 (`StpUtil.getTokenValue()`)
2. 从请求参数 `token` 获取

**Token 验证**：
- 如果 `security.enable=false`，则跳过验证
- 否则使用 `StpUtil.checkLogin()` 进行验证

**角色配置**：
- 默认返回：`admin`, `lowdeveloper`, `dbadeveloper`

**权限标识符**：
- `drag:datasource:testConnection` - 仪表盘数据库连接测试
- `onl:drag:clear:recovery` - 清空回收站
- `drag:analysis:sql` - SQL解析
- `drag:design:getTotalData` - 仪表盘对Online表单展示数据
- `drag:dataset:save` - 数据集保存
- `drag:dataset:delete` - 数据集删除
- `drag:datasource:saveOrUpate` - 数据源保存
- `drag:datasource:delete` - 数据源删除
- `onl:drag:page:delete` - 页面删除

**租户 ID 获取**：
1. Header: `X-Tenant-Key`
2. Header: `X-Tenant-Id`
3. Query Parameter: `tenant_id`

### 2.3 仪表盘外部服务 (IOnlDragExternalService)

| 方法 | 描述 |
|------|------|
| `getManyDictItems(List<String> codeList, List<JSONObject> tableDictList)` | 批量获取字典 |
| `getDictItems(String dictCode)` | 获取字典 |
| `addLog(DragLogDTO dragLogDTO)` | 添加日志 |

## 3. Token 传递方式

根据配置，支持以下 Token 传递方式：

1. **Sa-Token Cookie**: `X-Access-Token` (默认)
2. **URL 参数**: `?token=xxx`

## 4. 待进一步探索的路由

由于 jimureport-spring-boot3-starter 是外部依赖，以下路由需要通过运行时或文档进一步确认：

### 4.1 报表相关
- `/jmreport/view` - 报表预览
- `/jmreport/export` - 报表导出
- `/jmreport/save` - 报表保存
- `/jmreport/del` - 报表删除
- `/jmreport/dataSource` - 数据源管理
- `/jmreport/design` - 报表设计

### 4.2 仪表盘相关
- `/drag/page` - 页面管理
- `/drag/dataset` - 数据集管理
- `/drag/datasource` - 数据源管理
- `/drag/component` - 组件管理
- `/drag/preview` - 预览

## 5. 配置相关

### 5.1 Sa-Token 配置
```yaml
sa-token:
  token-name: X-Access-Token
  timeout: 2592000  # 30天
  is-concurrent: true
  is-share: false
  token-style: uuid
```

### 5.2 数据库配置
```yaml
spring:
  datasource:
    url: jdbc:mysql://host:port/db?characterEncoding=UTF-8&useUnicode=true&useSSL=false
    username: root
    password: root
```

### 5.3 文件上传配置
```yaml
spring:
  servlet:
    multipart:
      max-file-size: 10MB
      max-request-size: 10MB
```

## 6. 下一步行动

1. [ ] 启动 jimureport-example，访问 `/actuator/mappings` 查看完整的路由映射
2. [ ] 浏览 goReport 官方文档，获取完整的 API 文档
3. [ ] 通过抓包工具分析 UI 调用的实际 API 路由
