import { describe, it, expect, beforeEach, vi } from 'vitest'
import { datasetApi } from './dataset'
import type { Dataset, DatasetField } from './dataset'

describe('datasetApi', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('list', () => {
    it('should have list method', () => {
      expect(datasetApi.list).toBeDefined()
      expect(typeof datasetApi.list).toBe('function')
    })
  })

  describe('create', () => {
    it('should have create method', () => {
      expect(datasetApi.create).toBeDefined()
      expect(typeof datasetApi.create).toBe('function')
    })
  })

  describe('getSchema', () => {
    it('should have getSchema method', () => {
      expect(datasetApi.getSchema).toBeDefined()
      expect(typeof datasetApi.getSchema).toBe('function')
    })
  })

  describe('update', () => {
    it('should have update method', () => {
      expect(datasetApi.update).toBeDefined()
      expect(typeof datasetApi.update).toBe('function')
    })
  })

  describe('delete', () => {
    it('should have delete method', () => {
      expect(datasetApi.delete).toBeDefined()
      expect(typeof datasetApi.delete).toBe('function')
    })
  })
})
