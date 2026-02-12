import { describe, expect, it } from 'vitest'
import {
  getBatchPrecheckMessage,
  resolveBatchResult,
  resolveRefreshLoadPlan,
  resolveSaveNavigationPlan,
  resolveTabSwitchLoadPlan
} from './datasetEditWorkflow'

describe('datasetEditWorkflow', () => {
  describe('save navigation', () => {
    it('keeps route context on save failure and only navigates on success', () => {
      expect(resolveSaveNavigationPlan('save', true, false)).toEqual({
        replaceRoute: false,
        returnToList: false
      })

      expect(resolveSaveNavigationPlan('save', true, true)).toEqual({
        replaceRoute: true,
        returnToList: false
      })

      expect(resolveSaveNavigationPlan('save_and_return', true, false)).toEqual({
        replaceRoute: false,
        returnToList: false
      })

      expect(resolveSaveNavigationPlan('save_and_return', false, true)).toEqual({
        replaceRoute: false,
        returnToList: true
      })
    })

    it('returns to list only when save_and_return succeeds', () => {
      expect(resolveSaveNavigationPlan('save', false, true).returnToList).toBe(false)
      expect(resolveSaveNavigationPlan('save_and_return', false, true).returnToList).toBe(true)
      expect(resolveSaveNavigationPlan('save_and_return', false, false).returnToList).toBe(false)
    })
  })

  describe('preview and batch tab transitions', () => {
    it('reuses loaded state on tab switch and only loads missing data', () => {
      expect(resolveTabSwitchLoadPlan('batch', 'dataset-1', 'dataset-1', 'dataset-1')).toEqual({
        loadSchema: false,
        loadPreview: false
      })

      expect(resolveTabSwitchLoadPlan('batch', 'dataset-1', '', 'dataset-1')).toEqual({
        loadSchema: true,
        loadPreview: false
      })

      expect(resolveTabSwitchLoadPlan('preview', 'dataset-1', 'dataset-1', '')).toEqual({
        loadSchema: false,
        loadPreview: true
      })
    })

    it('does not trigger loading when dataset context is missing', () => {
      expect(resolveTabSwitchLoadPlan('batch', '', 'dataset-1', 'dataset-1')).toEqual({
        loadSchema: false,
        loadPreview: false
      })
    })
  })

  describe('refresh behavior', () => {
    it('builds deterministic refresh loading plan for preview context', () => {
      expect(resolveRefreshLoadPlan('dataset-1', 'batch', 'dataset-1')).toEqual({
        loadSchema: true,
        loadPreview: true
      })

      expect(resolveRefreshLoadPlan('dataset-1', 'batch', '')).toEqual({
        loadSchema: true,
        loadPreview: false
      })

      expect(resolveRefreshLoadPlan('dataset-1', 'preview', '')).toEqual({
        loadSchema: true,
        loadPreview: true
      })
    })

    it('does not refresh when dataset context is missing', () => {
      expect(resolveRefreshLoadPlan('', 'batch', 'dataset-1')).toEqual({
        loadSchema: false,
        loadPreview: false
      })
    })
  })

  describe('batch workflow guards', () => {
    it('returns guard messages and extracts partial failure field ids', () => {
      expect(getBatchPrecheckMessage('', 1, true)).toBe('请先保存数据集')
      expect(getBatchPrecheckMessage('dataset-1', 0, true)).toBe('请先选择要批量更新的字段')
      expect(getBatchPrecheckMessage('dataset-1', 1, false)).toBe('请至少配置一个批量更新项')
      expect(getBatchPrecheckMessage('dataset-1', 1, true)).toBe('')

      const partial = resolveBatchResult({
        success: false,
        updatedFields: ['field-1'],
        errors: [
          { fieldId: 'field-2', message: 'invalid sort order' },
          { fieldId: 'field-3', message: 'field does not belong to dataset' }
        ]
      })

      expect(partial).toEqual({
        isPartialFailure: true,
        failedFieldIds: ['field-2', 'field-3']
      })
    })
  })
})
