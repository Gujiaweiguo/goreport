import type { DatasetSchema } from '@/api/dataset'

type RowData = Record<string, any>

export type PreviewMode = 'schema_chart' | 'fallback_count' | 'fallback_index'

export interface PreviewChartModel {
  mode: PreviewMode
  categories: string[]
  values: number[]
  categoryLabel: string
  valueLabel: string
  fallbackMessage?: string
}

const toNumber = (value: unknown): number | null => {
  if (typeof value === 'number' && Number.isFinite(value)) {
    return value
  }

  if (typeof value === 'string') {
    const parsed = Number(value)
    if (Number.isFinite(parsed)) {
      return parsed
    }
  }

  return null
}

const getAllKeys = (rows: RowData[]): string[] => {
  const set = new Set<string>()
  rows.forEach((row) => {
    Object.keys(row || {}).forEach((key) => set.add(key))
  })
  return Array.from(set).sort((a, b) => a.localeCompare(b))
}

const firstSchemaField = (schema: DatasetSchema | null | undefined, section: Array<'dimensions' | 'measures' | 'computed'>, keys: Set<string>, predicate?: (fieldName: string) => boolean): string => {
  if (!schema) {
    return ''
  }

  for (const part of section) {
    const fields = schema[part] || []
    for (const field of fields) {
      const fieldName = field?.name || ''
      if (!fieldName || !keys.has(fieldName)) {
        continue
      }
      if (!predicate || predicate(fieldName)) {
        return fieldName
      }
    }
  }

  return ''
}

const hasNumericValue = (rows: RowData[], key: string): boolean => {
  return rows.some((row) => toNumber(row?.[key]) !== null)
}

const hasCategoricalValue = (rows: RowData[], key: string): boolean => {
  return rows.some((row) => {
    const value = row?.[key]
    return value !== null && value !== undefined && String(value).trim().length > 0
  })
}

export const buildPreviewChartModel = (
  schema: DatasetSchema | null | undefined,
  data: RowData[]
): PreviewChartModel => {
  const rows = Array.isArray(data) ? data : []
  const keys = getAllKeys(rows)

  if (!rows.length || !keys.length) {
    return {
      mode: 'fallback_count',
      categories: [],
      values: [],
      categoryLabel: '分类',
      valueLabel: '数量',
      fallbackMessage: '暂无可用预览数据'
    }
  }

  const keySet = new Set(keys)
  const numericKeys = keys.filter((key) => hasNumericValue(rows, key))
  const categoricalKeys = keys.filter((key) => hasCategoricalValue(rows, key))
  const nonNumericCategoricalKeys = categoricalKeys.filter((key) => !numericKeys.includes(key))

  const schemaCategoryKey = firstSchemaField(schema, ['dimensions', 'computed'], keySet, (fieldName) => !numericKeys.includes(fieldName))
  const schemaValueKey = firstSchemaField(schema, ['measures', 'computed'], keySet, (fieldName) => numericKeys.includes(fieldName))

  const categoryKey = schemaCategoryKey || nonNumericCategoricalKeys[0] || ''
  const valueKey = schemaValueKey || numericKeys[0] || ''

  if (categoryKey && valueKey) {
    return {
      mode: 'schema_chart',
      categories: rows.map((row) => String(row?.[categoryKey] ?? '')),
      values: rows.map((row) => toNumber(row?.[valueKey]) ?? 0),
      categoryLabel: categoryKey,
      valueLabel: valueKey
    }
  }

  if (categoryKey) {
    const countMap = new Map<string, number>()
    rows.forEach((row) => {
      const category = String(row?.[categoryKey] ?? '空值')
      countMap.set(category, (countMap.get(category) || 0) + 1)
    })

    const categories = Array.from(countMap.keys())
    return {
      mode: 'fallback_count',
      categories,
      values: categories.map((category) => countMap.get(category) || 0),
      categoryLabel: categoryKey,
      valueLabel: '记录数',
      fallbackMessage: '未找到可用指标字段，已切换为分类计数预览'
    }
  }

  if (valueKey) {
    return {
      mode: 'fallback_index',
      categories: rows.map((_, index) => `第${index + 1}行`),
      values: rows.map((row) => toNumber(row?.[valueKey]) ?? 0),
      categoryLabel: '行序号',
      valueLabel: valueKey,
      fallbackMessage: '未找到可用维度字段，已按行序号展示指标'
    }
  }

  return {
    mode: 'fallback_count',
    categories: [],
    values: [],
    categoryLabel: '分类',
    valueLabel: '数量',
    fallbackMessage: '当前数据结构暂不支持图表展示'
  }
}
