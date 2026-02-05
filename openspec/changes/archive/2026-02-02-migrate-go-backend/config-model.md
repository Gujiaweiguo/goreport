# 配置模型设计

## 1. 概述

本文档定义了 JimuReport Go 后端的配置模型，确保与现有 Java Spring Boot 配置的兼容性，并支持环境变量覆盖。

## 2. 配置文件结构

### 2.1 主配置文件 (config/config.yaml)

```yaml
# 服务配置
server:
  port: 8085              # 服务端口（保持与 Java 一致）
  mode: release          # 运行模式: debug, release
  read_timeout: 60s       # 读取超时
  write_timeout: 60s      # 写入超时

# 数据库配置
database:
  driver: mysql            # 数据库驱动: mysql, postgresql, sqlite
  host: ${DB_HOST:127.0.0.1}      # 数据库主机
  port: ${DB_PORT:3306}             # 数据库端口
  name: ${DB_NAME:jimureport}         # 数据库名
  username: ${DB_USERNAME:root}       # 用户名
  password: ${DB_PASSWORD:root}       # 密码
  charset: utf8mb4                    # 字符集
  parse_time: true                    # 解析时间
  loc: Local                           # 时区
  max_idle_conns: 10                 # 最大空闲连接数
  max_open_conns: 100                # 最大打开连接数
  conn_max_lifetime: 1h               # 连接最大生命周期

# JWT 配置
jwt:
  secret: ${JWT_SECRET:jimureport-secret-key-change-in-production}  # JWT 密钥
  expiration: ${JWT_EXPIRATION:2592000}  # Token 过期时间（秒），默认 30 天
  issuer: jimureport-go               # 签发者
  audience: jimureport-api             # 接收方
  algorithm: HS256                     # 签名算法

# Sa-Token 兼容配置
satoken:
  token_name: X-Access-Token            # Token 名称（与 Java 配置一致）
  timeout: 2592000                    # Token 有效期（秒），30 天
  active_timeout: -1                   # Token 最低活跃频率（秒），-1 表示不限制
  is_concurrent: true                  # 是否允许同一账号多地同时登录
  is_share: false                      # 是否共用 token
  token_style: uuid                    # Token 风格
  is_log: false                       # 是否输出操作日志
  is_print: false                     # 是否打印 banner

# 日志配置
logging:
  level: ${LOG_LEVEL:info}           # 日志级别: debug, info, warn, error, fatal
  format: json                        # 日志格式: json, text
  output: stdout                       # 输出目标: stdout, stderr, file
  file_path: /opt/jimureport/logs/app.log  # 日志文件路径（当 output=file 时）
  max_size: 100                       # 日志文件最大大小（MB）
  max_backups: 3                       # 保留的日志文件数量
  max_age: 28                          # 日志文件保留天数
  compress: true                      # 是否压缩旧日志

# CORS 配置
cors:
  enabled: true                       # 是否启用 CORS
  allowed_origins:
    - "*"                             # 允许的源
  allowed_methods:
    - GET
    - POST
    - PUT
    - DELETE
    - OPTIONS
  allowed_headers:
    - Content-Type
    - Authorization
    - X-Access-Token
    - X-Tenant-Key
    - X-Tenant-Id
  exposed_headers:
    - Content-Length
  allow_credentials: true             # 是否允许凭据
  max_age: 86400                       # 预检请求缓存时间（秒）

# 文件上传配置
upload:
  max_file_size: ${MAX_FILE_SIZE:10MB}  # 最大文件大小
  max_request_size: ${MAX_REQUEST_SIZE:10MB}  # 最大请求大小
  allowed_types:                       # 允许的文件类型
    - image/jpeg
    - image/png
    - image/gif
    - application/pdf
    - application/msword
    - application/vnd.openxmlformats-officedocument.wordprocessingml.document
    - application/vnd.ms-excel
    - application/vnd.openxmlformats-officedocument.spreadsheetml.sheet

# 报表配置
jmreport:
  # 签名密钥
  signature_secret: dd05f1c54d63749eda95f9fa6d49v442a

  # 防火墙配置
  firewall:
    data_source_safe: false            # 数据源安全
    sql_parse_safe: false               # SQL 解析安全
    low_code_mode: dev                 # 低代码模式: dev, prod
    sql_injection_level: basic          # SQL 注入检查级别: strict, basic, none

  # 展示配置
  col: 100                             # 展示列数
  row: 200                             # 展示行数

  # API 基础路径
  api_base_path: http://192.168.1.11:8085  # 自定义 API 接口前缀

  # 数据量最大限制
  max_data_rows: 100000               # 无分页模式和打印全部的最大记录数

  # 页面大小
  page_size:
    - 10
    - 20
    - 30
    - 40

  # 多租户模式
  saas_mode:  # created: 按照创建人隔离, tenant: 按照租户隔离, 空: 不隔离

  # 高德地图配置
  gao_de_api:
    api_key: ${GAODE_API_KEY:}       # 高德地图 API Key
    secret_key: ${GAODE_SECRET_KEY:}   # 高德地图秘钥

  # 邮件发送配置
  mail:
    enabled: false                      # 是否开启
    host: ${MAIL_HOST:}                # SMTP 服务器
    port: ${MAIL_PORT:25}              # SMTP 端口
    sender: ${MAIL_SENDER:}            # 发件人
    username: ${MAIL_USERNAME:}          # 用户名
    password: ${MAIL_PASSWORD:}          # 密码
    ssl: false                          # 是否使用 SSL
    from: ${MAIL_FROM:}                # 发件人地址

  # 自动化导出配置
  automate:
    export:
      enable_auto_export: true          # 是否开启自动导出
      expired: 30                      # 文件过期时间（天）
      jimu_view_path:                  # 积木报表 view 页面地址
      download_path: /opt/download      # 下载的报表存放目录

  # AI 配置
  ai:
    jeecg_host: ${JEECG_HOST:http://localhost:8080/jeecgboot/}

# 文件存储配置
storage:
  type: ${STORAGE_TYPE:local}          # 存储类型: local, minio, alioss
  
  # 本地存储
  local:
    upload_path: /opt/upload             # 文件上传路径

  # MinIO 对象存储
  minio:
    endpoint: ${MINIO_ENDPOINT:}       # MinIO 服务地址
    access_key: ${MINIO_ACCESS_KEY:} # Access Key
    secret_key: ${MINIO_SECRET_KEY:} # Secret Key
    bucket_name: ${MINIO_BUCKET:jimureport}  # 存储桶名称
    use_ssl: false                      # 是否使用 SSL

  # 阿里云 OSS
  alioss:
    endpoint: ${OSS_ENDPOINT:}           # OSS 端点
    access_key: ${OSS_ACCESS_KEY:}     # Access Key ID
    secret_key: ${OSS_SECRET_KEY:}     # Access Key Secret
    bucket_name: ${OSS_BUCKET:jimureport}  # 存储桶名称

# 租户配置
tenant:
  enabled: ${TENANT_ENABLED:false}      # 是否启用多租户
  default_id: 1                       # 默认租户 ID
  header_keys:
    tenant_key: X-Tenant-Key            # 租户 key header
    tenant_id: X-Tenant-Id             # 租户 ID header
  query_param: tenant_id               # 租户 query 参数名

# Redis 配置（可选）
redis:
  enabled: ${REDIS_ENABLED:false}       # 是否启用 Redis
  host: ${REDIS_HOST:127.0.0.1}      # Redis 主机
  port: ${REDIS_PORT:6379}             # Redis 端口
  password: ${REDIS_PASSWORD:}          # Redis 密码
  db: 1                               # 数据库索引
  pool_size: 10                         # 连接池大小
  min_idle_conns: 5                     # 最小空闲连接数

# 性能配置
performance:
  # 连接池
  db_pool:
    enabled: true
    max_idle: 10
    max_open: 100
    conn_max_lifetime: 3600  # 秒

  # 缓存
  cache:
    enabled: ${CACHE_ENABLED:false}     # 是否启用缓存
    type: ${CACHE_TYPE:redis}          # 缓存类型: redis, memory
    default_expiration: 3600           # 默认过期时间（秒）

  # 并发
  concurrent:
    max_workers: 100                   # 最大 worker 数量
    queue_size: 1000                  # 任务队列大小
```

### 2.2 开发环境配置 (config/config.dev.yaml)

```yaml
server:
  port: 8085
  mode: debug

database:
  host: 127.0.0.1
  port: 3306
  name: jimureport
  username: root
  password: root

jwt:
  secret: dev-secret-key-do-not-use-in-production
  expiration: 86400  # 1 day in dev

logging:
  level: debug
  format: text
  output: stdout

jmreport:
  firewall:
    low_code_mode: dev
```

### 2.3 生产环境配置 (config/config.prod.yaml)

```yaml
server:
  port: 8085
  mode: release

database:
  host: ${DB_HOST}
  port: ${DB_PORT}
  name: ${DB_NAME}
  username: ${DB_USERNAME}
  password: ${DB_PASSWORD}

jwt:
  secret: ${JWT_SECRET}
  expiration: 2592000

logging:
  level: ${LOG_LEVEL:info}
  format: json
  output: stdout

cors:
  allowed_origins:
    - ${ALLOWED_ORIGINS:*}
```

## 3. 配置模型 (Go 代码)

```go
package config

import "time"

// Config 应用配置
type Config struct {
    Server    ServerConfig    `mapstructure:"server"`
    Database  DatabaseConfig  `mapstructure:"database"`
    JWT       JWTConfig       `mapstructure:"jwt"`
    SaToken   SaTokenConfig   `mapstructure:"satoken"`
    Logging   LoggingConfig   `mapstructure:"logging"`
    CORS       CORSConfig       `mapstructure:"cors"`
    Upload     UploadConfig     `mapstructure:"upload"`
    JimuReport JimuReportConfig `mapstructure:"jmreport"`
    Storage    StorageConfig    `mapstructure:"storage"`
    Tenant     TenantConfig     `mapstructure:"tenant"`
    Redis      RedisConfig      `mapstructure:"redis"`
    Performance PerformanceConfig `mapstructure:"performance"`
}

// ServerConfig 服务配置
type ServerConfig struct {
    Port          int           `mapstructure:"port"`
    Mode          string        `mapstructure:"mode"`  // debug, release
    ReadTimeout   time.Duration `mapstructure:"read_timeout"`
    WriteTimeout  time.Duration `mapstructure:"write_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
    Driver           string        `mapstructure:"driver"`
    Host             string        `mapstructure:"host"`
    Port             int           `mapstructure:"port"`
    Name             string        `mapstructure:"name"`
    Username         string        `mapstructure:"username"`
    Password         string        `mapstructure:"password"`
    Charset          string        `mapstructure:"charset"`
    ParseTime        bool          `mapstructure:"parse_time"`
    Loc              time.Location `mapstructure:"loc"`
    MaxIdleConns     int           `mapstructure:"max_idle_conns"`
    MaxOpenConns     int           `mapstructure:"max_open_conns"`
    ConnMaxLifetime  time.Duration `mapstructure:"conn_max_lifetime"`
}

// GetDSN 获取数据库连接字符串
func (c *DatabaseConfig) GetDSN() string {
    switch c.Driver {
    case "mysql":
        return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
            c.Username, c.Password, c.Host, c.Port, c.Name,
            c.Charset, c.ParseTime, c.Loc.String())
    case "postgresql":
        return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
            c.Host, c.Port, c.Username, c.Password, c.Name)
    case "sqlite":
        return c.Name
    default:
        return ""
    }
}

// JWTConfig JWT 配置
type JWTConfig struct {
    Secret     string `mapstructure:"secret"`
    Expiration int64  `mapstructure:"expiration"`  // seconds
    Issuer     string `mapstructure:"issuer"`
    Audience   string `mapstructure:"audience"`
    Algorithm  string `mapstructure:"algorithm"`
}

// SaTokenConfig Sa-Token 兼容配置
type SaTokenConfig struct {
    TokenName      string `mapstructure:"token_name"`
    Timeout        int64  `mapstructure:"timeout"`
    ActiveTimeout  int64  `mapstructure:"active_timeout"`
    IsConcurrent   bool   `mapstructure:"is_concurrent"`
    IsShare        bool   `mapstructure:"is_share"`
    TokenStyle     string `mapstructure:"token_style"`
    IsLog          bool   `mapstructure:"is_log"`
    IsPrint        bool   `mapstructure:"is_print"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
    Level      string `mapstructure:"level"`
    Format     string `mapstructure:"format"`
    Output     string `mapstructure:"output"`
    FilePath   string `mapstructure:"file_path"`
    MaxSize    int    `mapstructure:"max_size"`
    MaxBackups int    `mapstructure:"max_backups"`
    MaxAge     int    `mapstructure:"max_age"`
    Compress   bool   `mapstructure:"compress"`
}

// CORSConfig CORS 配置
type CORSConfig struct {
    Enabled         bool     `mapstructure:"enabled"`
    AllowedOrigins []string `mapstructure:"allowed_origins"`
    AllowedMethods []string `mapstructure:"allowed_methods"`
    AllowedHeaders []string `mapstructure:"allowed_headers"`
    ExposedHeaders []string `mapstructure:"exposed_headers"`
    AllowCredentials bool    `mapstructure:"allow_credentials"`
    MaxAge         int      `mapstructure:"max_age"`
}

// UploadConfig 上传配置
type UploadConfig struct {
    MaxFileSize    int64    `mapstructure:"max_file_size"`
    MaxRequestSize int64    `mapstructure:"max_request_size"`
    AllowedTypes   []string `mapstructure:"allowed_types"`
}

// JimuReportConfig 报表配置
type JimuReportConfig struct {
    SignatureSecret string              `mapstructure:"signature_secret"`
    Firewall       FirewallConfig       `mapstructure:"firewall"`
    Col            int                 `mapstructure:"col"`
    Row            int                 `mapstructure:"row"`
    APIBasePath    string              `mapstructure:"api_base_path"`
    MaxDataRows   int                 `mapstructure:"max_data_rows"`
    PageSize       []int               `mapstructure:"page_size"`
    SaaSMode       string              `mapstructure:"saas_mode"`
    GaoDeAPI       GaoDeAPIConfig      `mapstructure:"gao_de_api"`
    Mail            MailConfig          `mapstructure:"mail"`
    Automate        AutomateConfig      `mapstructure:"automate"`
    AI              AIConfig            `mapstructure:"ai"`
}

// FirewallConfig 防火墙配置
type FirewallConfig struct {
    DataSourceSafe      bool `mapstructure:"data_source_safe"`
    SQLParseSafe         bool `mapstructure:"sql_parse_safe"`
    LowCodeMode         string `mapstructure:"low_code_mode"`  // dev, prod
    SQLInjectionLevel  string `mapstructure:"sql_injection_level"`  // strict, basic, none
}

// GaoDeAPIConfig 高德地图配置
type GaoDeAPIConfig struct {
    APIKey     string `mapstructure:"api_key"`
    SecretKey string `mapstructure:"secret_key"`
}

// MailConfig 邮件配置
type MailConfig struct {
    Enabled  bool   `mapstructure:"enabled"`
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    Sender   string `mapstructure:"sender"`
    Username string `mapstructure:"username"`
    Password string `mapstructure:"password"`
    SSL      bool   `mapstructure:"ssl"`
    From     string `mapstructure:"from"`
}

// AutomateConfig 自动化导出配置
type AutomateConfig struct {
    Export AutomateExportConfig `mapstructure:"export"`
}

// AutomateExportConfig 自动化导出配置
type AutomateExportConfig struct {
    EnableAutoExport bool   `mapstructure:"enable_auto_export"`
    Expired          int    `mapstructure:"expired"`
    JimuViewPath     string `mapstructure:"jimu_view_path"`
    DownloadPath      string `mapstructure:"download_path"`
}

// AIConfig AI 配置
type AIConfig struct {
    JeecgHost string `mapstructure:"jeecg_host"`
}

// StorageConfig 存储配置
type StorageConfig struct {
    Type   string       `mapstructure:"type"`  // local, minio, alioss
    Local  LocalConfig  `mapstructure:"local"`
    MinIO  MinIOConfig  `mapstructure:"minio"`
    AliOSS AliOSSConfig `mapstructure:"alioss"`
}

// LocalConfig 本地存储配置
type LocalConfig struct {
    UploadPath string `mapstructure:"upload_path"`
}

// MinIOConfig MinIO 配置
type MinIOConfig struct {
    Endpoint   string `mapstructure:"endpoint"`
    AccessKey string `mapstructure:"access_key"`
    SecretKey string `mapstructure:"secret_key"`
    BucketName string `mapstructure:"bucket_name"`
    UseSSL     bool   `mapstructure:"use_ssl"`
}

// AliOSSConfig 阿里云 OSS 配置
type AliOSSConfig struct {
    Endpoint   string `mapstructure:"endpoint"`
    AccessKey string `mapstructure:"access_key"`
    SecretKey string `mapstructure:"secret_key"`
    BucketName string `mapstructure:"bucket_name"`
}

// TenantConfig 租户配置
type TenantConfig struct {
    Enabled     bool   `mapstructure:"enabled"`
    DefaultID   string `mapstructure:"default_id"`
    HeaderKeys  TenantHeaderKeys `mapstructure:"header_keys"`
    QueryParam  string `mapstructure:"query_param"`
}

// TenantHeaderKeys 租户 Header 配置
type TenantHeaderKeys struct {
    TenantKey string `mapstructure:"tenant_key"`
    TenantID string `mapstructure:"tenant_id"`
}

// RedisConfig Redis 配置
type RedisConfig struct {
    Enabled        bool   `mapstructure:"enabled"`
    Host           string `mapstructure:"host"`
    Port           int    `mapstructure:"port"`
    Password       string `mapstructure:"password"`
    DB             int    `mapstructure:"db"`
    PoolSize       int    `mapstructure:"pool_size"`
    MinIdleConns   int    `mapstructure:"min_idle_conns"`
}

// PerformanceConfig 性能配置
type PerformanceConfig struct {
    DBPool    PoolConfig    `mapstructure:"db_pool"`
    Cache     CacheConfig   `mapstructure:"cache"`
    Concurrent ConcurrentConfig `mapstructure:"concurrent"`
}

// PoolConfig 连接池配置
type PoolConfig struct {
    Enabled          bool          `mapstructure:"enabled"`
    MaxIdle          int           `mapstructure:"max_idle"`
    MaxOpen          int           `mapstructure:"max_open"`
    ConnMaxLifetime  time.Duration `mapstructure:"conn_max_lifetime"`
}

// CacheConfig 缓存配置
type CacheConfig struct {
    Enabled           bool   `mapstructure:"enabled"`
    Type              string `mapstructure:"type"`  // redis, memory
    DefaultExpiration int64  `mapstructure:"default_expiration"`  // seconds
}

// ConcurrentConfig 并发配置
type ConcurrentConfig struct {
    MaxWorkers int `mapstructure:"max_workers"`
    QueueSize  int `mapstructure:"queue_size"`
}
```

## 4. 配置加载

```go
package config

import (
    "github.com/spf13/viper"
    "strings"
)

// Load 加载配置
func Load(configPath string) (*Config, error) {
    v := viper.New()

    // 设置配置文件名
    v.SetConfigName("config")
    v.SetConfigType("yaml")

    // 添加配置文件路径
    v.AddConfigPath(configPath)
    v.AddConfigPath("./config")
    v.AddConfigPath("/etc/jimureport")
    v.AddConfigPath("$HOME/.jimureport")

    // 环境变量替换
    v.AutomaticEnv()
    v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

    // 读取配置
    if err := v.ReadInConfig(); err != nil {
        return nil, err
    }

    // 解析配置
    var cfg Config
    if err := v.Unmarshal(&cfg); err != nil {
        return nil, err
    }

    return &cfg, nil
}

// LoadWithProfile 加载指定 profile 的配置
func LoadWithProfile(configPath, profile string) (*Config, error) {
    v := viper.New()

    // 设置配置文件名
    if profile != "" {
        v.SetConfigName("config." + profile)
    } else {
        v.SetConfigName("config")
    }
    v.SetConfigType("yaml")

    // 添加配置文件路径
    v.AddConfigPath(configPath)
    v.AddConfigPath("./config")
    v.AddConfigPath("/etc/jimureport")

    // 环境变量替换
    v.AutomaticEnv()
    v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

    // 读取配置
    if err := v.ReadInConfig(); err != nil {
        return nil, err
    }

    // 解析配置
    var cfg Config
    if err := v.Unmarshal(&cfg); err != nil {
        return nil, err
    }

    return &cfg, nil
}
```

## 5. 配置验证

```go
package config

import "errors"

// Validate 验证配置
func (c *Config) Validate() error {
    // 验证必需的配置
    if c.JWT.Secret == "" || c.JWT.Secret == "dev-secret-key-do-not-use-in-production" && c.Server.Mode == "release" {
        return errors.New("JWT_SECRET must be set in production")
    }

    if c.Database.Password == "" && c.Server.Mode == "release" {
        return errors.New("database password must be set in production")
    }

    // 验证端口号
    if c.Server.Port <= 0 || c.Server.Port > 65535 {
        return errors.New("invalid server port")
    }

    // 验证数据库配置
    if c.Database.Host == "" {
        return errors.New("database host must be set")
    }

    // 验证存储配置
    if c.Storage.Type != "local" && c.Storage.Type != "minio" && c.Storage.Type != "alioss" {
        return errors.New("invalid storage type")
    }

    return nil
}
```

## 6. 环境变量映射表

| 环境变量 | 对应配置项 | 默认值 | 说明 |
|-----------|-----------|--------|------|
| SERVER_PORT | server.port | 8085 | 服务端口 |
| DB_HOST | database.host | 127.0.0.1 | 数据库主机 |
| DB_PORT | database.port | 3306 | 数据库端口 |
| DB_NAME | database.name | jimureport | 数据库名 |
| DB_USERNAME | database.username | root | 数据库用户名 |
| DB_PASSWORD | database.password | root | 数据库密码 |
| JWT_SECRET | jwt.secret | - | JWT 密钥（必填） |
| JWT_EXPIRATION | jwt.expiration | 2592000 | Token 过期时间 |
| LOG_LEVEL | logging.level | info | 日志级别 |
| MAX_FILE_SIZE | upload.max_file_size | 10MB | 最大文件大小 |
| MAX_REQUEST_SIZE | upload.max_request_size | 10MB | 最大请求大小 |
| STORAGE_TYPE | storage.type | local | 存储类型 |
| TENANT_ENABLED | tenant.enabled | false | 是否启用多租户 |
| REDIS_ENABLED | redis.enabled | false | 是否启用 Redis |
| CACHE_ENABLED | performance.cache.enabled | false | 是否启用缓存 |
| MINIO_ENDPOINT | storage.minio.endpoint | - | MinIO 端点 |
| MINIO_ACCESS_KEY | storage.minio.access_key | - | MinIO Access Key |
| MINIO_SECRET_KEY | storage.minio.secret_key | - | MinIO Secret Key |

## 7. 与 Java 配置的兼容性

### 7.1 端口保持一致

Go 后端默认端口保持为 `8085`，与 Java 版本一致。

### 7.2 Sa-Token 配置兼容

保持 Sa-Token 的关键配置项：
- `token_name`: X-Access-Token
- `timeout`: 30 天
- 其他行为通过 JWT 模拟

### 7.3 数据库连接字符串

Go 版本使用相同的连接参数：
- `characterEncoding=UTF-8`
- `useUnicode=true`
- `useSSL=false`
- `serverTimezone=Asia/Shanghai`

### 7.4 日志配置

保持相同的日志级别和配置方式。

## 8. 配置热重载（可选）

```go
package config

import (
    "github.com/fsnotify/fsnotify"
    "log"
)

// WatchConfig 监听配置文件变化
func WatchConfig(v *viper.Viper, onChange func(*Config)) {
    v.WatchConfig()
    v.OnConfigChange(func(e fsnotify.Event) {
        log.Printf("Config file changed: %s", e.String())
        
        var cfg Config
        if err := v.Unmarshal(&cfg); err != nil {
            log.Printf("Failed to unmarshal config: %v", err)
            return
        }
        
        onChange(&cfg)
    })
}
```

## 9. 配置最佳实践

1. **敏感信息**: 使用环境变量而非配置文件存储敏感信息
2. **默认值**: 为所有必需的配置项提供合理的默认值
3. **配置验证**: 启动时验证配置的有效性
4. **文档化**: 保持配置文件中的注释清晰和最新
5. **环境隔离**: 为不同环境（dev, staging, prod）使用不同的配置文件

## 10. 迁移检查清单

从 Java 配置迁移到 Go 配置时，检查以下项目：

- [ ] 服务器端口一致 (8085)
- [ ] 数据库连接参数一致
- [ ] JWT 密钥已设置
- [ ] 文件上传大小限制一致 (10MB)
- [ ] CORS 配置一致
- [ ] 日志级别配置一致
- [ ] 租户 header 配置一致
- [ ] Sa-Token token 名称一致
- [ ] 报表配置参数一致（signature_secret, page_size 等）
