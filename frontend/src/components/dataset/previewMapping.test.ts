import { describe, expect, it } from 'vitest'
import type { DatasetSchema } from '@/api/dataset'
import { buildPreviewChartModel } from './previewMapping'

const schema: DatasetSchema = {
  dimensions: [
    {
      id: 'field-region',
      datasetId: 'dataset-1',
      name: 'region_name',
      displayName: '区域',
      type: 'dimension',
      dataType: 'string',
      isComputed: false,
      isSortable: true,
      isGroupable: true,
      defaultSortOrder: 'none',
      sortIndex: 0,
      createdAt: '',
      updatedAt: ''
    }
  ],
  measures: [
    {
      id: 'field-amount',
      datasetId: 'dataset-1',
      name: 'order_amount',
      displayName: '金额',
      type: 'measure',
      dataType: 'number',
      isComputed: false,
      isSortable: true,
      isGroupable: false,
      defaultSortOrder: 'none',
      sortIndex: 1,
      createdAt: '',
      updatedAt: ''
    }
  ],
  computed: []
}

describe('previewMapping', () => {
  it('maps chart axis from schema field names instead of fixed keys', () => {
    const model = buildPreviewChartModel(schema, [
      { region_name: '华北', order_amount: 1200 },
      { region_name: '华南', order_amount: 900 }
    ])

    expect(model.mode).toBe('schema_chart')
    expect(model.categoryLabel).toBe('region_name')
    expect(model.valueLabel).toBe('order_amount')
    expect(model.categories).toEqual(['华北', '华南'])
    expect(model.values).toEqual([1200, 900])
  })

  it('falls back deterministically to category count when no numeric measure exists', () => {
    const model = buildPreviewChartModel(undefined, [
      { city: '上海', tag: 'A' },
      { city: '上海', tag: 'B' },
      { city: '北京', tag: 'C' }
    ])

    expect(model.mode).toBe('fallback_count')
    expect(model.categoryLabel).toBe('city')
    expect(model.valueLabel).toBe('记录数')
    expect(model.categories).toEqual(['上海', '北京'])
    expect(model.values).toEqual([2, 1])
    expect(model.fallbackMessage).toContain('未找到可用指标字段')
  })

  it('falls back deterministically to row index when no category exists', () => {
    const model = buildPreviewChartModel(undefined, [
      { score: 88 },
      { score: 92 }
    ])

    expect(model.mode).toBe('fallback_index')
    expect(model.categoryLabel).toBe('行序号')
    expect(model.valueLabel).toBe('score')
    expect(model.categories).toEqual(['第1行', '第2行'])
    expect(model.values).toEqual([88, 92])
    expect(model.fallbackMessage).toContain('未找到可用维度字段')
  })
})
