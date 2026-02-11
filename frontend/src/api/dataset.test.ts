import { beforeEach, describe, expect, it, vi } from 'vitest'
import apiClient from './client'
import { datasetApi, formatBatchFieldErrors, getApiErrorMessage } from './dataset'

vi.mock('./client', () => {
  return {
    default: {
      get: vi.fn(),
      post: vi.fn(),
      put: vi.fn(),
      patch: vi.fn(),
      delete: vi.fn()
    }
  }
})

describe('datasetApi', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('calls batch update endpoint with fieldId in each item', () => {
    const payload = {
      fields: [
        {
          fieldId: 'field-1',
          type: 'dimension' as const,
          sortOrder: 'asc' as const
        }
      ]
    }

    datasetApi.batchUpdateFields('dataset-1', payload)

    expect(apiClient.patch).toHaveBeenCalledWith('/api/v1/datasets/dataset-1/fields', payload)
  })

  it('formats partial batch failures with field details', () => {
    const details = formatBatchFieldErrors([
      { fieldId: 'field-1', message: 'field does not belong to dataset' },
      { fieldId: 'field-2', message: 'invalid sort order' }
    ])

    expect(details).toBe('field-1: field does not belong to dataset；field-2: invalid sort order')
  })

  it('prefers backend business/guardrail message over fallback', () => {
    const message = getApiErrorMessage(
      {
        isAxiosError: true,
        response: {
          data: {
            message: 'query validation failed: disallowed SQL operation'
          }
        }
      },
      '批量更新失败'
    )

    expect(message).toBe('query validation failed: disallowed SQL operation')
  })

  it('returns transport fallback when request fails without response payload', () => {
    const message = getApiErrorMessage(
      {
        isAxiosError: true,
        request: {}
      },
      '批量更新失败'
    )

    expect(message).toBe('网络连接失败，请检查网络后重试')
  })
})
