# UI 资源策略和静态路径

## 1. 概述

JimuReport 的前端 UI 主要由 `jimureport-spring-boot3-starter` 依赖提供，包括报表设计器、仪表盘设计器等。Go 后端需要正确处理静态资源的访问和服务。

## 2. 静态资源位置

### 2.1 本地静态资源

| 路径 | 说明 |
|------|------|
| `/static/login/login.html` | 登录页面 |
| `/static/login/cdn/*` | 登录页面依赖的 CDN 资源 |
| `/static/favicon.ico` | 网站图标 |

### 2.2 Starter 提供的静态资源

`jimureport-spring-boot3-starter` 会自动提供以下前端资源：

| 资源 | 路径 | 说明 |
|------|------|------|
| 报表工作台 | `/jmreport/list` | 报表设计器和列表页 |
| 仪表盘工作台 | `/drag/list` | 仪表盘/大屏设计器和列表页 |
| 报表预览 | `/jmreport/view/*` | 报表预览页面 |
| 仪表盘预览 | `/drag/page/view/*` | 仪表盘预览页面 |

## 3. Go 后端静态资源策略

### 3.1 策略选项

#### 选项 1: 直接嵌入静态资源 (推荐用于 MVP)

**优点**:
- 部署简单，单个二进制文件
- 无需外部依赖

**缺点**:
- 二进制文件较大
- 更新 UI 需要重新编译

**实现**:
```go
// 使用 embed.EmbedFS 将静态资源嵌入 Go 二进制
//go:embed all:static
var staticFiles embed.FS
```

#### 选项 2: 从外部目录服务 (推荐用于生产)

**优点**:
- UI 可以独立更新
- 二进制文件较小

**缺点**:
- 需要额外的部署步骤
- 需要管理静态资源目录

**实现**:
```go
http.Handle("/static/", http.FileServer(http.Dir("/path/to/static")))
```

#### 选项 3: 代理到现有 Java 服务 (过渡方案)

**优点**:
- 最小化初始工作
- 完全兼容现有 UI

**缺点**:
- 需要保持两个服务运行
- 性能开销

**实现**:
```go
// 使用 httputil.ReverseProxy 代理请求
proxy := httputil.NewSingleHostReverseProxy(&url.URL{
    Scheme: "http",
    Host:   "java-service:8085",
})
```

### 3.2 推荐方案

**阶段 1 (MVP)**: 使用选项 3 - 代理到 Java 服务
- 最小化初始开发工作量
- 确保 UI 完全兼容

**阶段 2 (过渡)**: 使用选项 2 - 从外部目录服务
- 将静态资源从 JAR 包中提取
- Go 服务直接提供静态资源

**阶段 3 (最终)**: 使用选项 1 - 直接嵌入静态资源
- 完全独立的 Go 服务
- 无外部依赖

## 4. 路由兼容性

### 4.1 需要保持的路由

以下路由必须保持不变，以确保 UI 兼容性：

| 路由 | 说明 | Go 实现 |
|------|------|----------|
| `/jmreport/list` | 报表工作台 | 代理或提供 HTML |
| `/jmreport/view/*` | 报表预览 | 代理或提供 HTML |
| `/drag/list` | 仪表盘工作台 | 代理或提供 HTML |
| `/drag/page/view/*` | 仪表盘预览 | 代理或提供 HTML |
| `/login/login.html` | 登录页面 | 提供 HTML |
| `/favicon.ico` | 网站图标 | 提供图标 |

### 4.2 API 路由

API 路由将由 Go 实现，不需要静态资源：

| 路由前缀 | 说明 |
|----------|------|
| `/jmreport/api/*` | 报表相关 API |
| `/drag/api/*` | 仪表盘相关 API |
| `/jmreport/export/*` | 报表导出 API |

## 5. 资源映射表

### 5.1 报表资源

| 原路径 | Go 实现 | 说明 |
|--------|----------|------|
| `/jmreport/list` | `/jmreport/list` | 代理到 Java 或提供新页面 |
| `/jmreport/view/{id}` | `/jmreport/view/{id}` | 代理到 Java 或提供新页面 |
| `/jmreport/export/*` | `/jmreport/export/*` | API 路由，由 Go 实现 |

### 5.2 仪表盘资源

| 原路径 | Go 实现 | 说明 |
|--------|----------|------|
| `/drag/list` | `/drag/list` | 代理到 Java 或提供新页面 |
| `/drag/page/view/{id}` | `/drag/page/view/{id}` | 代理到 Java 或提供新页面 |
| `/drag/api/*` | `/drag/api/*` | API 路由，由 Go 实现 |

### 5.3 其他资源

| 原路径 | Go 实现 | 说明 |
|--------|----------|------|
| `/login/login.html` | `/login/login.html` | 由 Go 提供登录页面 |
| `/static/login/*` | `/static/login/*` | 登录页面资源 |
| `/favicon.ico` | `/favicon.ico` | 由 Go 提供 |

## 6. 文件上传配置

### 6.1 上传配置 (从 application-dev.yml)

```yaml
spring:
  servlet:
    multipart:
      max-file-size: 10MB
      max-request-size: 10MB
```

### 6.2 Go 实现

```go
const (
    maxUploadSize = 10 << 20 // 10 MB
)

func UploadLimitMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
        next.ServeHTTP(w, r)
    })
}
```

## 7. CORS 配置

### 7.1 当前 CORS 配置

从 `CustomCorsConfiguration.java`:
- 允许跨域请求
- 支持自定义 headers

### 7.2 Go CORS 中间件

```go
func CORSMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Access-Token")

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}
```

## 8. 资源压缩

### 8.1 Gzip 压缩

建议为静态资源启用 Gzip 压缩：

```go
import "github.com/NYTimes/gziphandler"

func main() {
    fs := http.FileServer(http.Dir("./static"))
    gzippedFs := gziphandler.GzipHandler(fs)

    http.Handle("/static/", gzippedFs)
}
```

### 8.2 缓存策略

为静态资源设置适当的缓存头：

```go
func StaticCacheMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 设置缓存时间
        w.Header().Set("Cache-Control", "public, max-age=3600")
        next.ServeHTTP(w, r)
    })
}
```

## 9. 部署结构

### 9.1 开发环境

```
jimureport-go/
├── main.go
├── static/              # 静态资源
│   ├── login/
│   ├── favicon.ico
│   └── ...
├── templates/           # 模板文件
│   └── ...
└── go.mod
```

### 9.2 生产环境

```
/opt/jimureport/
├── jimureport-go        # Go 二进制
├── static/              # 静态资源目录
└── config/             # 配置文件
    ├── config.yaml
    └── ...
```

## 10. 从 JAR 提取静态资源

如果选择从 Java JAR 包中提取静态资源：

```bash
# 解压 JAR 包
unzip jimureport-spring-boot3-starter-2.3.0.jar -d extracted

# 复制静态资源
cp -r extracted/BOOT-INF/classes/static/* /opt/jimureport/static/
```

## 11. 实现检查清单

### 11.1 基础功能

- [ ] 实现静态资源文件服务
- [ ] 实现路由兼容性（保持 `/jmreport/*` 和 `/drag/*`）
- [ ] 实现登录页面
- [ ] 实现文件上传限制
- [ ] 实现 CORS 支持

### 11.2 性能优化

- [ ] 启用 Gzip 压缩
- [ ] 配置静态资源缓存
- [ ] 配置 ETag
- [ ] 实现静态资源版本控制

### 11.3 安全性

- [ ] 验证文件类型上传
- [ ] 限制文件上传大小
- [ ] 实现 XSS 防护
- [ ] 配置 CSP (Content Security Policy)

## 12. 监控和日志

### 12.1 静态资源访问日志

```go
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
    })
}
```

### 12.2 错误处理

为静态资源访问实现友好的错误页面：

```go
func ErrorHandler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if r := recover(); r != nil {
                log.Printf("Panic: %v", r)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            }
        }()
        next.ServeHTTP(w, r)
    })
}
```
