import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { nextTick } from 'vue'
import DatasourceManage from './DatasourceManage.vue'

vi.mock('element-plus', () => ({
  ElMessage: {
    error: vi.fn(),
    success: vi.fn(),
    warning: vi.fn(),
    info: vi.fn()
  },
  ElMessageBox: {
    confirm: vi.fn()
  }
}))

vi.mock('@/api/datasource', () => ({
  datasourceApi: {
    list: vi.fn().mockResolvedValue({ data: { success: true, result: { datasources: [], total: 0 } } }),
    create: vi.fn().mockResolvedValue({ data: { success: true } }),
    update: vi.fn().mockResolvedValue({ data: { success: true } }),
    delete: vi.fn().mockResolvedValue({ data: { success: true } }),
    copy: vi.fn().mockResolvedValue({ data: { success: true } }),
    move: vi.fn().mockResolvedValue({ data: { success: true } }),
    rename: vi.fn().mockResolvedValue({ data: { success: true } }),
    test: vi.fn().mockResolvedValue({ data: { success: true } }),
    testById: vi.fn().mockResolvedValue({ data: { success: true } })
  }
}))

import { ElMessage, ElMessageBox } from 'element-plus'
import { datasourceApi } from '@/api/datasource'

const globalStubs = {
  'el-card': true,
  'el-button': true,
  'el-button-group': true,
  'el-input': true,
  'el-input-number': true,
  'el-select': true,
  'el-option': true,
  'el-table': true,
  'el-table-column': true,
  'el-tag': true,
  'el-pagination': true,
  'el-dialog': true,
  'el-form': true,
  'el-form-item': true,
  'el-checkbox': true,
  'el-radio-group': true,
  'el-radio': true,
  'el-divider': true,
  'el-result': true,
  'el-tooltip': true,
  'el-icon': true
}

describe('DatasourceManage.vue', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(datasourceApi.list).mockResolvedValue({ 
      data: { success: true, result: { datasources: [], total: 0 } } 
    } as any)
  })

  describe('Component Rendering', () => {
    it('renders correctly', async () => {
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.find('.datasource-manage').exists()).toBe(true)
    })
  })

  describe('Initial state', () => {
    it('currentPage starts at 1', async () => {
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.vm.currentPage).toBe(1)
    })

    it('pageSize starts at 10', async () => {
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.vm.pageSize).toBe(10)
    })

    it('total starts at 0', async () => {
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.vm.total).toBe(0)
    })

    it('dialogVisible starts false', async () => {
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.vm.dialogVisible).toBe(false)
    })

    it('isEdit starts false', async () => {
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.vm.isEdit).toBe(false)
    })

    it('loading state changes during data fetch', async () => {
      vi.mocked(datasourceApi.list).mockImplementation(() => new Promise(() => {}))
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.vm.loading).toBe(true)
    })
  })

  describe('needsDatabase computed', () => {
    it('returns true for mysql', async () => {
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.form.type = 'mysql'
      expect(wrapper.vm.needsDatabase).toBe(true)
    })

    it('returns true for postgresql', async () => {
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.form.type = 'postgresql'
      expect(wrapper.vm.needsDatabase).toBe(true)
    })

    it('returns true for sqlserver', async () => {
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.form.type = 'sqlserver'
      expect(wrapper.vm.needsDatabase).toBe(true)
    })

    it('returns true for mongodb', async () => {
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.form.type = 'mongodb'
      expect(wrapper.vm.needsDatabase).toBe(true)
    })

    it('returns false for excel', async () => {
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.form.type = 'excel'
      expect(wrapper.vm.needsDatabase).toBe(false)
    })

    it('returns false for csv', async () => {
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.form.type = 'csv'
      expect(wrapper.vm.needsDatabase).toBe(false)
    })

    it('returns false for api', async () => {
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.form.type = 'api'
      expect(wrapper.vm.needsDatabase).toBe(false)
    })
  })

  describe('needsAuth computed', () => {
    it('returns true for mysql', async () => {
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.form.type = 'mysql'
      expect(wrapper.vm.needsAuth).toBe(true)
    })

    it('returns true for postgresql', async () => {
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.form.type = 'postgresql'
      expect(wrapper.vm.needsAuth).toBe(true)
    })

    it('returns false for excel', async () => {
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.form.type = 'excel'
      expect(wrapper.vm.needsAuth).toBe(false)
    })
  })

  describe('supportsSSH computed', () => {
    it('returns true for mysql', async () => {
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.form.type = 'mysql'
      expect(wrapper.vm.supportsSSH).toBe(true)
    })

    it('returns true for postgresql', async () => {
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.form.type = 'postgresql'
      expect(wrapper.vm.supportsSSH).toBe(true)
    })

    it('returns true for mongodb', async () => {
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.form.type = 'mongodb'
      expect(wrapper.vm.supportsSSH).toBe(true)
    })

    it('returns false for sqlserver', async () => {
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.form.type = 'sqlserver'
      expect(wrapper.vm.supportsSSH).toBe(false)
    })

    it('returns false for excel', async () => {
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.form.type = 'excel'
      expect(wrapper.vm.supportsSSH).toBe(false)
    })
  })

  describe('formatDate function', () => {
    it('formats valid date string', async () => {
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      const result = wrapper.vm.formatDate('2024-01-15T10:30:00')
      expect(result).toContain('2024')
    })

    it('returns dash for empty string', async () => {
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.vm.formatDate('')).toBe('-')
    })

    it('returns dash for undefined', async () => {
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.vm.formatDate(undefined)).toBe('-')
    })

    it('returns dash for null', async () => {
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.vm.formatDate(null as any)).toBe('-')
    })
  })

  describe('handleSearch function', () => {
    it('resets currentPage to 1 and calls loadDatasources', async () => {
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      vi.clearAllMocks()
      
      wrapper.vm.currentPage = 5
      wrapper.vm.handleSearch()
      
      expect(wrapper.vm.currentPage).toBe(1)
      expect(datasourceApi.list).toHaveBeenCalled()
    })
  })

  describe('handlePageChange function', () => {
    it('updates currentPage and calls loadDatasources', async () => {
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      vi.clearAllMocks()
      
      wrapper.vm.handlePageChange(3)
      
      expect(wrapper.vm.currentPage).toBe(3)
      expect(datasourceApi.list).toHaveBeenCalled()
    })
  })

  describe('handleSizeChange function', () => {
    it('updates pageSize, resets currentPage, and calls loadDatasources', async () => {
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      vi.clearAllMocks()
      
      wrapper.vm.handleSizeChange(20)
      
      expect(wrapper.vm.pageSize).toBe(20)
      expect(wrapper.vm.currentPage).toBe(1)
      expect(datasourceApi.list).toHaveBeenCalled()
    })
  })

  describe('showCreateDialog function', () => {
    it('sets isEdit to false and dialogVisible to true', async () => {
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      
      wrapper.vm.showCreateDialog()
      
      expect(wrapper.vm.isEdit).toBe(false)
      expect(wrapper.vm.dialogTitle).toBe('创建数据源')
      expect(wrapper.vm.dialogVisible).toBe(true)
    })
  })

  describe('showEditDialog function', () => {
    it('populates form with row data and opens dialog', async () => {
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      
      const row = {
        id: 'ds-1',
        name: 'Test DS',
        type: 'mysql',
        host: 'localhost',
        port: 3306,
        database: 'testdb',
        username: 'root'
      } as any
      
      wrapper.vm.showEditDialog(row)
      
      expect(wrapper.vm.isEdit).toBe(true)
      expect(wrapper.vm.dialogTitle).toBe('编辑数据源')
      expect(wrapper.vm.form.name).toBe('Test DS')
      expect(wrapper.vm.form.type).toBe('mysql')
      expect(wrapper.vm.form.host).toBe('localhost')
      expect(wrapper.vm.form.port).toBe(3306)
      expect(wrapper.vm.currentEditId).toBe('ds-1')
      expect(wrapper.vm.dialogVisible).toBe(true)
    })
  })

  describe('handleCopy function', () => {
    it('calls datasourceApi.copy and shows success message', async () => {
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      
      await wrapper.vm.handleCopy({ id: 'ds-1', name: 'Test' } as any)
      
      expect(datasourceApi.copy).toHaveBeenCalledWith('ds-1')
      expect(ElMessage.success).toHaveBeenCalledWith('复制数据源成功')
    })

    it('shows error message on copy failure', async () => {
      vi.mocked(datasourceApi.copy).mockResolvedValueOnce({ 
        data: { success: false, message: 'Copy failed' } 
      } as any)
      
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      
      await wrapper.vm.handleCopy({ id: 'ds-1', name: 'Test' } as any)
      
      expect(ElMessage.error).toHaveBeenCalledWith('Copy failed')
    })

    it('shows error message on exception', async () => {
      vi.mocked(datasourceApi.copy).mockRejectedValueOnce(new Error('Network'))
      
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      
      await wrapper.vm.handleCopy({ id: 'ds-1', name: 'Test' } as any)
      
      expect(ElMessage.error).toHaveBeenCalledWith('复制数据源失败')
    })
  })

  describe('showMoveDialog function', () => {
    it('sets moveForm and opens dialog', async () => {
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      
      wrapper.vm.showMoveDialog({ id: 'ds-1', name: 'Test' } as any)
      
      expect(wrapper.vm.moveForm.id).toBe('ds-1')
      expect(wrapper.vm.moveForm.target).toBe('')
      expect(wrapper.vm.moveDialogVisible).toBe(true)
    })
  })

  describe('handleMoveSubmit function', () => {
    it('shows warning when target is empty', async () => {
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      
      wrapper.vm.moveForm.id = 'ds-1'
      wrapper.vm.moveForm.target = ''
      
      await wrapper.vm.handleMoveSubmit()
      
      expect(ElMessage.warning).toHaveBeenCalledWith('请输入目标位置')
      expect(datasourceApi.move).not.toHaveBeenCalled()
    })

    it('calls datasourceApi.move and shows success', async () => {
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      
      wrapper.vm.moveForm.id = 'ds-1'
      wrapper.vm.moveForm.target = 'folder-1'
      
      await wrapper.vm.handleMoveSubmit()
      
      expect(datasourceApi.move).toHaveBeenCalledWith('ds-1', 'folder-1')
      expect(ElMessage.success).toHaveBeenCalledWith('移动数据源成功')
      expect(wrapper.vm.moveDialogVisible).toBe(false)
    })

    it('shows error on move failure', async () => {
      vi.mocked(datasourceApi.move).mockResolvedValueOnce({ 
        data: { success: false, message: 'Move failed' } 
      } as any)
      
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      
      wrapper.vm.moveForm.id = 'ds-1'
      wrapper.vm.moveForm.target = 'folder-1'
      
      await wrapper.vm.handleMoveSubmit()
      
      expect(ElMessage.error).toHaveBeenCalledWith('Move failed')
    })
  })

  describe('showRenameDialog function', () => {
    it('sets renameForm and opens dialog', async () => {
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      
      wrapper.vm.showRenameDialog({ id: 'ds-1', name: 'Old Name' } as any)
      
      expect(wrapper.vm.renameForm.id).toBe('ds-1')
      expect(wrapper.vm.renameForm.name).toBe('Old Name')
      expect(wrapper.vm.renameDialogVisible).toBe(true)
    })
  })

  describe('handleRenameSubmit function', () => {
    it('shows warning when name is empty', async () => {
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      
      wrapper.vm.renameForm.id = 'ds-1'
      wrapper.vm.renameForm.name = '   '
      
      await wrapper.vm.handleRenameSubmit()
      
      expect(ElMessage.warning).toHaveBeenCalledWith('请输入新名称')
      expect(datasourceApi.rename).not.toHaveBeenCalled()
    })

    it('calls datasourceApi.rename and shows success', async () => {
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      
      wrapper.vm.renameForm.id = 'ds-1'
      wrapper.vm.renameForm.name = 'New Name'
      
      await wrapper.vm.handleRenameSubmit()
      
      expect(datasourceApi.rename).toHaveBeenCalledWith('ds-1', 'New Name')
      expect(ElMessage.success).toHaveBeenCalledWith('重命名成功')
      expect(wrapper.vm.renameDialogVisible).toBe(false)
    })

    it('shows error on rename failure', async () => {
      vi.mocked(datasourceApi.rename).mockResolvedValueOnce({ 
        data: { success: false, message: 'Rename failed' } 
      } as any)
      
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      
      wrapper.vm.renameForm.id = 'ds-1'
      wrapper.vm.renameForm.name = 'New Name'
      
      await wrapper.vm.handleRenameSubmit()
      
      expect(ElMessage.error).toHaveBeenCalledWith('Rename failed')
    })
  })

  describe('handleDelete function', () => {
    it('deletes datasource after confirmation', async () => {
      vi.mocked(ElMessageBox.confirm).mockResolvedValueOnce('confirm' as any)
      
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      
      await wrapper.vm.handleDelete({ id: 'ds-1', name: 'Test DS' } as any)
      
      expect(ElMessageBox.confirm).toHaveBeenCalled()
      expect(datasourceApi.delete).toHaveBeenCalledWith('ds-1')
      expect(ElMessage.success).toHaveBeenCalledWith('删除数据源成功')
    })

    it('does not delete when cancelled', async () => {
      vi.mocked(ElMessageBox.confirm).mockRejectedValueOnce('cancel')
      
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      
      await wrapper.vm.handleDelete({ id: 'ds-1', name: 'Test' } as any)
      
      expect(datasourceApi.delete).not.toHaveBeenCalled()
    })

    it('shows error on delete failure', async () => {
      vi.mocked(ElMessageBox.confirm).mockResolvedValueOnce('confirm' as any)
      vi.mocked(datasourceApi.delete).mockResolvedValueOnce({ 
        data: { success: false, message: 'Delete failed' } 
      } as any)
      
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      
      await wrapper.vm.handleDelete({ id: 'ds-1', name: 'Test' } as any)
      
      expect(ElMessage.error).toHaveBeenCalledWith('Delete failed')
    })
  })

  describe('openTestDialog function', () => {
    it('sets currentTestDatasourceId and opens dialog', async () => {
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      
      wrapper.vm.openTestDialog({ id: 'ds-1', name: 'Test' } as any)
      
      expect(wrapper.vm.currentTestDatasourceId).toBe('ds-1')
      expect(wrapper.vm.testDialogVisible).toBe(true)
      expect(wrapper.vm.testResult).toBeNull()
    })
  })

  describe('resetForm function', () => {
    it('resets form to default values', async () => {
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      
      wrapper.vm.form.name = 'Modified'
      wrapper.vm.form.type = 'postgresql'
      wrapper.vm.form.host = 'remote'
      wrapper.vm.form.port = 5432
      
      wrapper.vm.resetForm()
      
      expect(wrapper.vm.form.name).toBe('')
      expect(wrapper.vm.form.type).toBe('mysql')
      expect(wrapper.vm.form.host).toBe('localhost')
      expect(wrapper.vm.form.port).toBe(3306)
      expect(wrapper.vm.enableSSH).toBe(false)
    })
  })

  describe('handleDialogClose function', () => {
    it('closes dialog and resets form', async () => {
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      
      wrapper.vm.dialogVisible = true
      wrapper.vm.form.name = 'Test'
      
      wrapper.vm.handleDialogClose()
      
      expect(wrapper.vm.dialogVisible).toBe(false)
      expect(wrapper.vm.form.name).toBe('')
    })
  })

  describe('handleTypeChange function', () => {
    it('clears database when type does not need it', async () => {
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      
      wrapper.vm.form.type = 'mysql'
      wrapper.vm.form.database = 'testdb'
      wrapper.vm.form.type = 'excel'
      
      wrapper.vm.handleTypeChange()
      
      expect(wrapper.vm.form.database).toBe('')
    })

    it('clears auth fields when type does not need them', async () => {
      const wrapper = mount(DatasourceManage, { global: { stubs: globalStubs } })
      await nextTick()
      
      wrapper.vm.form.type = 'mysql'
      wrapper.vm.form.username = 'root'
      wrapper.vm.form.password = 'secret'
      wrapper.vm.form.type = 'excel'
      
      wrapper.vm.handleTypeChange()
      
      expect(wrapper.vm.form.username).toBe('')
      expect(wrapper.vm.form.password).toBe('')
    })
  })
})
