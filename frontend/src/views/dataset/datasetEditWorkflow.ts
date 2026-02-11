import type { BatchUpdateFieldsResponse, FieldError } from '@/api/dataset'

export interface SaveNavigationPlan {
  replaceRoute: boolean
  returnToList: boolean
}

export const resolveSaveNavigationPlan = (
  action: 'save' | 'save_and_return',
  isNewDataset: boolean,
  isSuccess: boolean
): SaveNavigationPlan => {
  if (!isSuccess) {
    return {
      replaceRoute: false,
      returnToList: false
    }
  }

  return {
    replaceRoute: action === 'save' && isNewDataset,
    returnToList: action === 'save_and_return'
  }
}

export interface TabSwitchLoadPlan {
  loadSchema: boolean
  loadPreview: boolean
}

export const resolveTabSwitchLoadPlan = (
  targetTab: string,
  datasetId: string,
  loadedSchemaDatasetId: string,
  loadedPreviewDatasetId: string
): TabSwitchLoadPlan => {
  if (!datasetId) {
    return {
      loadSchema: false,
      loadPreview: false
    }
  }

  return {
    loadSchema: targetTab === 'batch' && loadedSchemaDatasetId !== datasetId,
    loadPreview: targetTab === 'preview' && loadedPreviewDatasetId !== datasetId
  }
}

export interface RefreshLoadPlan {
  loadSchema: boolean
  loadPreview: boolean
}

export const resolveRefreshLoadPlan = (
  datasetId: string,
  activeTab: string,
  loadedPreviewDatasetId: string
): RefreshLoadPlan => {
  if (!datasetId) {
    return {
      loadSchema: false,
      loadPreview: false
    }
  }

  return {
    loadSchema: true,
    loadPreview: activeTab === 'preview' || loadedPreviewDatasetId === datasetId
  }
}

export const getBatchPrecheckMessage = (datasetId: string, selectedCount: number, hasPatch: boolean): string => {
  if (!datasetId) {
    return '请先保存数据集'
  }
  if (!selectedCount) {
    return '请先选择要批量更新的字段'
  }
  if (!hasPatch) {
    return '请至少配置一个批量更新项'
  }
  return ''
}

export interface BatchResultResolution {
  isPartialFailure: boolean
  failedFieldIds: string[]
}

export const resolveBatchResult = (result: BatchUpdateFieldsResponse): BatchResultResolution => {
  const errors: FieldError[] = Array.isArray(result?.errors) ? result.errors : []
  const failedFieldIds = errors
    .map((item) => item.fieldId)
    .filter((fieldId): fieldId is string => Boolean(fieldId))

  return {
    isPartialFailure: !result.success && failedFieldIds.length > 0,
    failedFieldIds
  }
}
