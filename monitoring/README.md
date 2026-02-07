# JimuReport 生产环境监控配置

本目录包含生产环境监控相关的配置文件。

## 监控栈

- **Prometheus**: 指标收集和存储
- **Grafana**: 可视化和告警
- **Alertmanager**: 告警管理
- **Node Exporter**: 系统指标收集
- **cAdvisor**: 容器指标收集

## 快速开始

```bash
# 启动监控栈
docker compose -f docker-compose.monitoring.yml up -d

# 访问 Grafana
open http://localhost:3000
# 默认账号: admin/admin
```

## 配置文件

### Prometheus

- `prometheus.yml`: Prometheus 主配置
- `alerts.yml`: 告警规则
- `targets/`: 目标配置

### Grafana

- `grafana.ini`: Grafana 配置
- `dashboards/`: 仪表板 JSON
- `datasources/`: 数据源配置

## 默认端口

| 服务 | 端口 | 说明 |
|------|------|------|
| Prometheus | 9090 | 指标查询和存储 |
| Grafana | 3000 | 可视化界面 |
| Alertmanager | 9093 | 告警管理 |
| Node Exporter | 9100 | 系统指标 |
| cAdvisor | 8080 | 容器指标 |

## 告警规则

### 系统告警

- CPU 使用率 > 80%
- 内存使用率 > 85%
- 磁盘使用率 > 90%
- 网络异常

### 应用告警

- 服务宕机
- 响应时间 > 2s
- 错误率 > 5%
- 数据库连接池耗尽

### 业务告警

- 用户登录失败率 > 10%
- 报表生成失败
- 数据导出超时

## 自定义仪表板

参考 `dashboards/jimureport-dashboard.json` 创建自定义仪表板。

## 维护

```bash
# 查看监控服务状态
docker compose -f docker-compose.monitoring.yml ps

# 查看日志
docker compose -f docker-compose.monitoring.yml logs -f prometheus

# 重启服务
docker compose -f docker-compose.monitoring.yml restart
```
