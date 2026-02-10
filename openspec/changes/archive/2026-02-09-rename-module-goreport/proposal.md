# Change: 重构模块名从 jimureport-go 到 goreport

## Why
- 仓库从组织 github.com/jeecg 迁移到个人 github.com/gujiaweiguo
- 简化模块名，从 jimureport-go 改为 goreport
- 统一品牌标识

## What Changes
- 更新 go.mod 模块声明
- 更新所有内部 import 路径
- 更新文档中的模块引用
- 更新配置文件中的路径引用（如有）

## Impact
- 所有 import 路径需要更新
- 外部依赖需要更新引用
- 这是一个 breaking change

## Migration Guide
对于外部使用者：
```bash
# 旧路径
go get github.com/gujiaweiguo/goreport

# 新路径
go get github.com/gujiaweiguo/goreport
```

内部代码中所有 import 需要替换：
```go
// 旧
import "github.com/gujiaweiguo/goreport/internal/..."

// 新
import "github.com/gujiaweiguo/goreport/internal/..."
```
