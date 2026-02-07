#!/bin/bash

# Lighthouse 性能测试脚本
# 用于测试 goReport 前端性能

set -e

echo "================================"
echo "goReport 性能测试脚本"
echo "================================"
echo ""

# 配置
FRONTEND_URL="http://localhost"
BACKEND_URL="http://localhost:8085"
OUTPUT_DIR="./performance-reports"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

# 创建输出目录
mkdir -p "$OUTPUT_DIR/$TIMESTAMP"

echo "测试目标:"
echo "  前端: $FRONTEND_URL"
echo "  后端: $BACKEND_URL"
echo "  报告目录: $OUTPUT_DIR/$TIMESTAMP"
echo ""

# 检查服务是否可用
echo "检查服务可用性..."
if ! curl -s "$BACKEND_URL/health" > /dev/null 2>&1; then
    echo "⚠️  后端服务未启动或无法访问"
    echo "   请确保服务已启动: make build-prod"
    exit 1
fi

if ! curl -s "$FRONTEND_URL" > /dev/null 2>&1; then
    echo "⚠️  前端服务未启动或无法访问"
    echo "   请确保服务已启动: make build-prod"
    exit 1
fi

echo "✅ 服务检查通过"
echo ""

# 检查 Lighthouse 是否安装
if ! command -v lighthouse &> /dev/null; then
    echo "Lighthouse 未安装，正在安装..."
    npm install -g lighthouse
fi

echo "开始性能测试..."
echo ""

# 测试首页
echo "1. 测试首页..."
lighthouse "$FRONTEND_URL" \
    --output=html,json \
    --output-path="$OUTPUT_DIR/$TIMESTAMP/home" \
    --chrome-flags="--headless --no-sandbox" \
    --only-categories=performance,accessibility,best-practices,seo \
    --preset=desktop \
    --quiet

echo "   ✅ 首页测试完成"

# 测试报表管理页（需要登录，可能需要调整）
echo "2. 测试报表管理页..."
lighthouse "$FRONTEND_URL/datasource" \
    --output=html,json \
    --output-path="$OUTPUT_DIR/$TIMESTAMP/datasource" \
    --chrome-flags="--headless --no-sandbox" \
    --only-categories=performance,accessibility,best-practices \
    --preset=desktop \
    --quiet || echo "   ⚠️  数据源页面测试失败（可能需要登录）"

echo "   ✅ 数据源页面测试完成"

# 生成汇总报告
echo ""
echo "生成汇总报告..."

cat > "$OUTPUT_DIR/$TIMESTAMP/summary.md" << EOF
# goReport 性能测试报告

**测试时间**: $(date '+%Y-%m-%d %H:%M:%S')  
**测试版本**: v1.0.0  
**测试环境**: UAT

## 测试页面

| 页面 | 报告 |
|------|------|
| 首页 | [home.report.html](./home.report.html) |
| 数据源管理 | [datasource.report.html](./datasource.report.html) |

## 性能指标目标

| 指标 | 目标 | 说明 |
|------|------|------|
| FCP (First Contentful Paint) | < 1.8s | 首屏渲染时间 |
| LCP (Largest Contentful Paint) | < 2.5s | 最大内容绘制时间 |
| TTI (Time to Interactive) | < 3.8s | 可交互时间 |
| TBT (Total Blocking Time) | < 200ms | 总阻塞时间 |
| CLS (Cumulative Layout Shift) | < 0.1 | 累积布局偏移 |
| Speed Index | < 3.4s | 速度指数 |

## 查看详细报告

1. 打开 \`$OUTPUT_DIR/$TIMESTAMP/home.report.html\` 查看首页详细报告
2. 查看 JSON 数据: \`$OUTPUT_DIR/$TIMESTAMP/home.report.json\`

## 优化建议

根据 Lighthouse 报告，常见的优化方向：

1. **图片优化**
   - 使用 WebP 格式
   - 压缩图片大小
   - 使用懒加载

2. **代码优化**
   - 移除未使用的 JavaScript
   - 代码分割和懒加载
   - 压缩和混淆代码

3. **缓存优化**
   - 启用长期缓存（Cache-Control）
   - 使用 Service Worker

4. **字体优化**
   - 使用 font-display: swap
   - 预加载关键字体

EOF

echo ""
echo "================================"
echo "性能测试完成!"
echo "================================"
echo ""
echo "报告位置: $OUTPUT_DIR/$TIMESTAMP/"
echo ""
echo "查看报告:"
echo "  首页: $OUTPUT_DIR/$TIMESTAMP/home.report.html"
echo "  汇总: $OUTPUT_DIR/$TIMESTAMP/summary.md"
echo ""

# 输出评分摘要
echo "评分摘要:"
if [ -f "$OUTPUT_DIR/$TIMESTAMP/home.report.json" ]; then
    echo "  首页性能评分:"
    cat "$OUTPUT_DIR/$TIMESTAMP/home.report.json" | grep -o '"performance":[0-9]*' | head -1
fi

echo ""
echo "提示: 使用浏览器打开 HTML 报告查看详细分析"
