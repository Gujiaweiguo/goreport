// 图表类型常量定义
export const CHART_TYPES = {
  // 基础图表
  bar: { key: 'bar', name: '柱状图', category: '基础图表', description: '用于横向对比不同维度的数据表现。\n适合展示销量、访问量和排行。', scene: '分类对比、指标排行、月度趋势对照', icon: 'Histogram' },
  line: { key: 'line', name: '折线图', category: '基础图表', description: '强调连续时间维度上的变化趋势。\n可用于观察峰值和转折点。', scene: '趋势分析、实时监控、同比环比', icon: 'TrendCharts' },
  pie: { key: 'pie', name: '饼图', category: '基础图表', description: '突出整体与部分之间的占比关系。\n数据类别不宜过多以保证可读性。', scene: '占比结构、渠道份额、成本构成', icon: 'PieChart' },
  scatter: { key: 'scatter', name: '散点图', category: '基础图表', description: '展示两个变量间的相关性与分布。\n便于发现异常点与聚类。', scene: '相关性分析、离群点检测、样本分布', icon: 'Grid' },
  
  // 饼图扩展
  'pie-extended': { key: 'pie-extended', name: '环形图', category: '饼图', description: '在饼图基础上保留中心留白区域。\n可放置关键 KPI 指标密度。', scene: '核心指标卡、占比+总量联合展示', icon: 'PieChart' },
  'donut': { key: 'donut', name: '南丁玫瑰图', category: '饼图', description: '使用半径或面积表达维度差异。\n视觉冲击更强，适合展示层级差。', scene: '品牌偏好、地域贡献、层级比较', icon: 'DataAnalysis' },
  
  // 仪表盘/数据看板
  gauge: { key: 'gauge', name: '仪表盘', category: '仪表盘', description: '适合单值指标的状态展示。\n可配阈值区间体现风险等级。', scene: '完成率、健康度、风险预警', icon: 'DataBoard' },
  
  // 高级图表
  'graph': { key: 'graph', name: '关系图', category: '关系图', description: '以节点和连线刻画实体关系网络。\n适合展示连接强度和传播链路。', scene: '组织关系、社交网络、调用链路', icon: 'Share' },
  'sankey': { key: 'sankey', name: '桑基图', category: '关系图', description: '展示流量在多个阶段间的流向。\n节点宽度体现体量变化。', scene: '转化漏斗、能耗流向、资金去向', icon: 'Connection' },
  'tree': { key: 'tree', name: '树图', category: '关系图', description: '用于层级结构和父子关系表达。\n支持自上而下的结构浏览。', scene: '组织架构、目录结构、权限树', icon: 'DataLine' },
  'relation': { key: 'relation', name: '关系图', category: '关系图', description: '刻画实体间的关联关系和流向。', scene: '关联分析、依赖追踪、影响分析', icon: 'Share' },
  
  // 地理图
  geo: { key: 'geo', name: '地理图', category: '地理图', items: [{ key: 'map', label: '地图', description: '支持区域数据对比和下钻。' }] },
  
  // 数据分析
  'heatmap': { key: 'heatmap', name: '热力图', category: '高级图表', description: '通过颜色深浅展示数据密度和分布。\n适用于大规模数据集的异常检测。', scene: '异常检测、数据密度、用户行为分析', icon: 'DataAnalysis' },
  'treemap': { key: 'treemap', name: '树图', category: '高级图表', description: '层级结构展示数据分布和占比。', scene: '目录分析、销售漏斗、用户路径', icon: 'DataLine' },
  'funnel': { key: 'funnel', name: '漏斗图', category: '高级图表', description: '展示数据在各个阶段的转化率。\n适用于转化漏斗、用户行为分析。', scene: '转化漏斗、流失分析、用户路径', icon: 'TrendCharts' },
  'wordcloud': { key: 'wordcloud', name: '词云', category: '高级图表', description: '展示文本词频统计和重要性。\n适用于关键词分析和文本挖掘。', scene: '关键词分析、舆情监控、文本挖掘', icon: 'DataAnalysis' },
  'parallel': { key: 'parallel', name: '平行坐标图', category: '高级图表', description: '在多维空间中展示数据相关性。\n适用于多维数据对比和雷达图。', scene: '多维度对比、绩效评估、能力矩阵', icon: 'DataAnalysis' },
  'radar': { key: 'radar', name: '雷达图', category: '高级图表', description: '以多维数据展示对比分析。\n适用于能力评估、绩效分析、多指标对比。', scene: '能力评估、绩效分析、多指标对比', icon: 'TrendCharts' },
  'boxplot': { key: 'boxplot', name: '箱线图', category: '高级图表', description: '展示数据分布、异常值和离群点。\n适用于统计分析、质量控制、异常检测。', scene: '统计分析、质量控制、异常检测', icon: 'DataAnalysis' },
  
  // 组合图表
  'combo': { key: 'combo', name: '组合图', category: '组合图表', description: '组合多种图表类型在同一画布。\n适用于综合数据分析和复杂场景展示。', scene: '综合看板、复杂报表、多维度分析', icon: 'TrendCharts' },
  'bar-line': { key: 'bar-line', name: '柱线图', category: '组合图表', description: '在柱状图基础上叠加折线展示趋势。\n适用于业务趋势对比。', scene: '趋势分析、多周期对比', icon: 'TrendCharts' },
  'line-area': { key: 'line-area', name: '面积图', category: '组合图表', description: '折线与区域填充结合。\n适用于趋势分析和数据对比。', scene: '趋势分析、多周期对比', icon: 'TrendCharts' },
  'multi-y': { key: 'multi-y', name: '多轴图', category: '组合图表', description: '多个 Y 轴共享 X 轴，展示多维度数据。', scene: '多维度对比、交叉分析', icon: 'TrendCharts' },
  'scatter-line': { key: 'scatter-line', name: '散点折线', category: '组合图表', description: '散点图与回归线结合。\n适用于趋势分析和相关性分析。', scene: '趋势分析、相关性分析', icon: 'TrendCharts' },
  'candlestick': { key: 'candlestick', name: 'K线图', category: '组合图表', description: '展示金融数据的开高低收和波动。\n适用于金融分析、趋势追踪。', scene: '金融分析、趋势追踪', icon: 'TrendCharts' },
  'effectscatter': { key: 'effectscatter', name: '效果散点', category: '组合图表', description: '带连线效果的散点图。\n适用于展示数据分布和聚类。', scene: '数据分布、聚类分析', icon: 'DataAnalysis' },
  
  // 3D 图表
  'map3d': { key: 'map3d', name: '3D地图', category: '地理图', items: [{ key: 'map', label: '地图' }], description: '在三维空间展示数据分布和地理位置。', scene: '地理分析、区域对比、位置追踪', icon: 'DataLine' },
  'mapbox3d': { key: 'mapbox3d', name: '3D 立方体图', category: '地理图', items: [{ key: 'mapbox', label: '3D 立方' }], description: '3D 立方图展示数据分布和比较。', scene: '地理分析、区域对比', icon: 'DataLine' }
} as const

export const SCENES = {
  '分类对比': '分类对比、指标排行、月度趋势对照',
  '趋势分析': '趋势分析、实时监控、同比环比',
  '占比结构': '占比结构、渠道份额、成本构成',
  '品牌偏好': '品牌偏好、地域贡献、层级比较',
  '组织关系': '组织架构、目录结构、权限树',
  '关联分析': '关联分析、依赖追踪、影响分析',
  '转化漏斗': '转化漏斗、流失分析、用户路径',
  '异常检测': '异常检测、数据密度、用户行为分析',
  '核心指标': '核心指标卡、占比+总量联合展示',
  '多维度对比': '多维度对比、交叉分析、绩效评估、能力矩阵',
  '综合看板': '综合看板、复杂报表、多维度分析'
} as const
