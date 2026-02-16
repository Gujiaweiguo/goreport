import { computed, reactive, ref, watch, type Ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules, UploadUserFile } from 'element-plus'
import { datasourceApi, type DataSource } from '@/api/datasource'
import {
  formatBatchFieldErrors,
  getApiErrorMessage,
  datasetApi,
  type BatchUpdateFieldRequest,
  type CreateFieldRequest,
  type CreateDatasetRequest,
  type DatasetField,
  type DatasetSchema,
  type UpdateDatasetRequest
} from '@/api/dataset'
import {
  getBatchPrecheckMessage,
  resolveBatchResult,
  resolveRefreshLoadPlan,
  resolveSaveNavigationPlan,
  resolveTabSwitchLoadPlan
} from '../datasetEditWorkflow'

// Types
export type WorkflowStatus = 'idle' | 'loading' | 'success' | 'error'
export type SourceType = 'custom_sql' | 'table'
export type JoinType = 'inner' | 'left' | 'right' | 'full'

export interface WorkflowState {
  status: WorkflowStatus
  operation: 'init' | 'save' | 'save_and_return' | 'preview' | 'refresh' | 'tab_switch' | 'batch_update' | 'create_grouping_field' | ''
  message: string
}

export interface BatchConfig {
  type?: 'dimension' | 'measure'
  sortOrder?: 'asc' | 'desc' | 'none'
  displayNamePrefix: string
}

export interface GroupingFieldForm {
  name: string
  displayName: string
  dataType: 'string' | 'number' | 'date' | 'boolean'
  groupingRule: string
  groupingEnabled: boolean
}

export interface SelectedSource {
  key: string
  label: string
  type: SourceType
  table?: string
}

export interface SidebarTableItem {
  key: string
  label: string
  type: SourceType
  table?: string
}

export interface JoinConfigState {
  leftSourceKey: string
  rightSourceKey: string
  relationType: JoinType
  leftField: string
  rightField: string
}

export type CreateGroupingFieldRequest = CreateFieldRequest & {
  isGroupingField: boolean
  groupingRule: string
  groupingEnabled: boolean
}

// Constants
export const CUSTOM_SQL_KEY = 'custom_sql'
export const CUSTOM_SQL_LABEL = '自定义SQL'
export const MAX_SELECTED_SOURCES = 2

export const joinKeywordMap: Record<JoinType, string> = {
  inner: 'INNER JOIN',
  left: 'LEFT JOIN',
  right: 'RIGHT JOIN',
  full: 'FULL JOIN'
}

export function useDatasetEdit(datasetId: string) {
  const router = useRouter()

  // Form refs
  const formRef = ref<FormInstance>()
  const groupingFieldFormRef = ref<FormInstance>()

  // Route state
  const currentDatasetId = ref<string>(datasetId)
  const activeTab = ref('preview')

  // Form data
  const formData = ref<CreateDatasetRequest>({
    name: '',
    type: 'sql',
    datasourceId: '',
    config: {}
  })

  const formRules = reactive<FormRules>({
    name: [{ required: true, message: '请输入数据集名称', trigger: 'blur' }],
    type: [{ required: true, message: '请选择数据集类型', trigger: 'change' }],
    datasourceId: [{ required: true, message: '请选择数据源', trigger: 'change' }]
  })

  // Workflow state
  const workflowState = ref<WorkflowState>({ status: 'idle', operation: '', message: '' })

  // Datasource state
  const datasources = ref<DataSource[]>([])
  const datasourcesLoading = ref(false)
  const tables = ref<string[]>([])
  const tablesLoading = ref(false)

  // SQL workspace state
  const activeSidebarTableKey = ref<string>(CUSTOM_SQL_KEY)
  const activeSourceKey = ref<string>('')
  const selectedSources = ref<SelectedSource[]>([])
  const leftJoinFieldOptions = ref<string[]>([])
  const rightJoinFieldOptions = ref<string[]>([])
  const joinFieldLoading = reactive({ left: false, right: false })
  const joinConfig = ref<JoinConfigState>({
    leftSourceKey: '',
    rightSourceKey: '',
    relationType: 'inner',
    leftField: '',
    rightField: ''
  })

  // File upload state
  const fileList = ref<UploadUserFile[]>([])

  // Preview & schema state
  const previewData = ref<Record<string, any>[]>([])
  const schema = ref<DatasetSchema>({ dimensions: [], measures: [], computed: [] })
  const loadedSchemaDatasetId = ref('')
  const loadedPreviewDatasetId = ref('')

  // Batch management state
  const selectedFields = ref<DatasetField[]>([])
  const batchConfig = ref<BatchConfig>({
    displayNamePrefix: ''
  })

  // Grouping field dialog state
  const groupingFieldDialogVisible = ref(false)
  const groupingFieldForm = ref<GroupingFieldForm>({
    name: '',
    displayName: '',
    dataType: 'string',
    groupingRule: '',
    groupingEnabled: false
  })
  const groupingFieldRules = reactive<FormRules>({
    name: [{ required: true, message: '请输入字段名', trigger: 'blur' }],
    groupingRule: [{ required: true, message: '请输入分组规则', trigger: 'blur' }],
    dataType: [{ required: true, message: '请选择数据类型', trigger: 'change' }]
  })

  // Computed properties
  const allFields = computed(() => [...schema.value.dimensions, ...schema.value.measures, ...schema.value.computed])
  const hasCustomSqlSource = computed(() => selectedSources.value.some((source) => source.type === 'custom_sql'))
  const showSqlEditor = computed(() => formData.value.type === 'sql' && hasCustomSqlSource.value)
  const sidebarTableItems = computed<SidebarTableItem[]>(() => {
    const baseItems: SidebarTableItem[] = [{ key: CUSTOM_SQL_KEY, label: CUSTOM_SQL_LABEL, type: 'custom_sql' }]
    if (!formData.value.datasourceId) {
      return baseItems
    }
    const tableItems = tables.value.map((table) => ({
      key: `table:${table}`,
      label: table,
      type: 'table' as const,
      table
    }))
    return [...baseItems, ...tableItems]
  })

  // Error handling
  const getErrorMessage = (error: unknown, fallback: string): string => {
    return getApiErrorMessage(error, fallback)
  }

  // Workflow state runner
  const runWithState = async <T>(
    operation: WorkflowState['operation'],
    runner: () => Promise<T>,
    options?: { successMessage?: string; errorMessage?: string; notifySuccess?: boolean }
  ): Promise<T> => {
    workflowState.value = { status: 'loading', operation, message: '' }
    try {
      const result = await runner()
      workflowState.value = { status: 'success', operation, message: options?.successMessage || '' }
      if (options?.successMessage && options.notifySuccess !== false) {
        ElMessage.success(options.successMessage)
      }
      return result
    } catch (error) {
      const message = getErrorMessage(error, options?.errorMessage || '操作失败')
      workflowState.value = { status: 'error', operation, message }
      ElMessage.error(message)
      throw error
    }
  }

  // Data loading functions
  const loadDatasources = async () => {
    datasourcesLoading.value = true
    try {
      const response = await datasourceApi.list()
      if (response.data.success) {
        datasources.value = response.data.result?.datasources || []
        return
      }
      throw new Error(response.data.message || '加载数据源失败')
    } catch (error) {
      ElMessage.error(getErrorMessage(error, '加载数据源失败'))
    } finally {
      datasourcesLoading.value = false
    }
  }

  const loadTables = async (datasourceId: string) => {
    if (!datasourceId) {
      tables.value = []
      return
    }
    try {
      tablesLoading.value = true
      const response = await datasourceApi.getTables(datasourceId)
      if (response.data.success) {
        tables.value = response.data.result || []
        return
      }
      tables.value = []
    } catch {
      ElMessage.error('加载表格列表失败')
      tables.value = []
    } finally {
      tablesLoading.value = false
    }
  }

  const loadDataset = async (id: string) => {
    const response = await datasetApi.get(id)
    if (!response.data.success) {
      throw new Error(response.data.message || '加载数据集失败')
    }
    const dataset = response.data.result
    formData.value = {
      name: dataset.name,
      type: dataset.type,
      datasourceId: dataset.datasourceId,
      config: typeof dataset.config === 'string' ? JSON.parse(dataset.config) : dataset.config || {}
    }
  }

  const loadSchema = async (id: string) => {
    const response = await datasetApi.getSchema(id)
    if (!response.data.success) {
      throw new Error(response.data.message || '加载字段结构失败')
    }
    schema.value = response.data.result || { dimensions: [], measures: [], computed: [] }
    loadedSchemaDatasetId.value = id
  }

  const fetchPreview = async (id: string) => {
    const response = await datasetApi.preview(id)
    if (!response.data.success) {
      throw new Error(response.data.message || '预览失败')
    }
    previewData.value = response.data.result || []
    loadedPreviewDatasetId.value = id
  }

  // SQL workspace functions
  const ensureSqlConfig = () => {
    if (!formData.value.config || typeof formData.value.config !== 'object') {
      formData.value.config = {}
    }
    if (!Array.isArray(formData.value.config.sources)) {
      formData.value.config.sources = []
    }
    if (!formData.value.config.join || typeof formData.value.config.join !== 'object') {
      formData.value.config.join = {
        leftSourceKey: '',
        rightSourceKey: '',
        relationType: 'inner',
        leftField: '',
        rightField: ''
      }
    }
    if (typeof formData.value.config.query !== 'string') {
      formData.value.config.query = ''
    }
  }

  const syncSourcesToConfig = () => {
    ensureSqlConfig()
    formData.value.config.sources = selectedSources.value.map((source, index) => ({
      key: source.key,
      label: source.label,
      type: source.type,
      table: source.table,
      sortIndex: index
    }))
  }

  const syncJoinToConfig = () => {
    ensureSqlConfig()
    formData.value.config.join = { ...joinConfig.value }
  }

  const resetJoinFieldOptions = () => {
    leftJoinFieldOptions.value = []
    rightJoinFieldOptions.value = []
    joinFieldLoading.left = false
    joinFieldLoading.right = false
  }

  const resetSqlWorkspace = () => {
    selectedSources.value = []
    activeSourceKey.value = ''
    joinConfig.value = {
      leftSourceKey: '',
      rightSourceKey: '',
      relationType: 'inner',
      leftField: '',
      rightField: ''
    }
    resetJoinFieldOptions()
    if (formData.value.type === 'sql') {
      ensureSqlConfig()
      formData.value.config.sources = []
      formData.value.config.join = { ...joinConfig.value }
      formData.value.config.query = ''
    }
  }

  const normalizeFieldOptions = (fields: any[]): string[] => {
    return fields
      .map((field) => {
        if (typeof field === 'string') return field
        if (field && typeof field.name === 'string') return field.name
        if (field && typeof field.columnName === 'string') return field.columnName
        return ''
      })
      .filter((field): field is string => !!field)
  }

  const loadJoinFieldsBySource = async (sourceKey: string, side: 'left' | 'right') => {
    const source = selectedSources.value.find((item) => item.key === sourceKey)
    if (!source || source.type !== 'table' || !source.table || !formData.value.datasourceId) {
      if (side === 'left') {
        leftJoinFieldOptions.value = []
      } else {
        rightJoinFieldOptions.value = []
      }
      return
    }

    joinFieldLoading[side] = true
    try {
      const response = await datasourceApi.getFields(formData.value.datasourceId, source.table)
      const options = response.data.success ? normalizeFieldOptions(response.data.result || []) : []
      if (side === 'left') {
        leftJoinFieldOptions.value = options
      } else {
        rightJoinFieldOptions.value = options
      }
    } catch {
      if (side === 'left') {
        leftJoinFieldOptions.value = []
      } else {
        rightJoinFieldOptions.value = []
      }
      ElMessage.error('加载关联字段失败')
    } finally {
      joinFieldLoading[side] = false
    }
  }

  const syncAutoQueryFromSources = () => {
    if (formData.value.type !== 'sql') return
    ensureSqlConfig()
    if (hasCustomSqlSource.value) return

    const tableSources = selectedSources.value.filter((source) => source.type === 'table' && source.table)
    if (!tableSources.length) {
      formData.value.config.query = ''
      return
    }

    if (tableSources.length === 1) {
      formData.value.config.query = `SELECT * FROM ${tableSources[0].table}`
      return
    }

    const leftSource = selectedSources.value.find((source) => source.key === joinConfig.value.leftSourceKey)
    const rightSource = selectedSources.value.find((source) => source.key === joinConfig.value.rightSourceKey)
    const fallbackLeft = tableSources[0]
    const fallbackRight = tableSources[1]
    const actualLeft = leftSource?.type === 'table' ? leftSource : fallbackLeft
    const actualRight = rightSource?.type === 'table' ? rightSource : fallbackRight

    if (!actualLeft?.table || !actualRight?.table) {
      formData.value.config.query = `SELECT * FROM ${fallbackLeft.table}`
      return
    }

    if (!joinConfig.value.leftField || !joinConfig.value.rightField) {
      formData.value.config.query = `SELECT * FROM ${actualLeft.table}`
      return
    }

    formData.value.config.query = `SELECT * FROM ${actualLeft.table} t1 ${joinKeywordMap[joinConfig.value.relationType]} ${actualRight.table} t2 ON t1.${joinConfig.value.leftField} = t2.${joinConfig.value.rightField}`
  }

  const applyDefaultJoinSources = () => {
    if (selectedSources.value.length < 2) {
      joinConfig.value.leftSourceKey = ''
      joinConfig.value.rightSourceKey = ''
      joinConfig.value.leftField = ''
      joinConfig.value.rightField = ''
      return
    }

    if (!selectedSources.value.some((source) => source.key === joinConfig.value.leftSourceKey)) {
      joinConfig.value.leftSourceKey = selectedSources.value[0].key
      joinConfig.value.leftField = ''
    }

    if (!selectedSources.value.some((source) => source.key === joinConfig.value.rightSourceKey)) {
      const fallbackRight = selectedSources.value[1] || selectedSources.value[0]
      joinConfig.value.rightSourceKey = fallbackRight.key
      joinConfig.value.rightField = ''
    }
  }

  const parseSelectedTable = (query: unknown): string => {
    if (typeof query !== 'string') return ''
    const match = query.match(/^\s*SELECT\s+\*\s+FROM\s+([`"\[]?[\w.]+[`"\]]?)\s*;?\s*$/i)
    return match?.[1]?.replace(/[`"\[\]]/g, '') || ''
  }

  const restoreSqlWorkspaceFromConfig = async () => {
    if (formData.value.type !== 'sql') return

    ensureSqlConfig()
    selectedSources.value = []

    if (Array.isArray(formData.value.config.sources) && formData.value.config.sources.length) {
      selectedSources.value = formData.value.config.sources
        .slice(0, MAX_SELECTED_SOURCES)
        .map((source: any) => ({
          key: source.key || (source.type === 'custom_sql' ? CUSTOM_SQL_KEY : `table:${source.table}`),
          label: source.label || (source.type === 'custom_sql' ? CUSTOM_SQL_LABEL : source.table || ''),
          type: source.type === 'custom_sql' ? 'custom_sql' : 'table',
          table: source.table
        }))
        .filter((source: SelectedSource) => !!source.label)
    }

    if (!selectedSources.value.length) {
      const parsedTable = parseSelectedTable(formData.value.config.query)
      if (parsedTable) {
        selectedSources.value = [{ key: `table:${parsedTable}`, label: parsedTable, type: 'table', table: parsedTable }]
      } else if (formData.value.config.query) {
        selectedSources.value = [{ key: CUSTOM_SQL_KEY, label: CUSTOM_SQL_LABEL, type: 'custom_sql' }]
      }
    }

    const configJoin = formData.value.config.join || {}
    joinConfig.value = {
      leftSourceKey: typeof configJoin.leftSourceKey === 'string' ? configJoin.leftSourceKey : '',
      rightSourceKey: typeof configJoin.rightSourceKey === 'string' ? configJoin.rightSourceKey : '',
      relationType: ['inner', 'left', 'right', 'full'].includes(configJoin.relationType) ? configJoin.relationType : 'inner',
      leftField: typeof configJoin.leftField === 'string' ? configJoin.leftField : '',
      rightField: typeof configJoin.rightField === 'string' ? configJoin.rightField : ''
    }

    activeSourceKey.value = selectedSources.value[0]?.key || ''
    applyDefaultJoinSources()
    syncSourcesToConfig()
    syncJoinToConfig()
    syncAutoQueryFromSources()

    await Promise.all([
      loadJoinFieldsBySource(joinConfig.value.leftSourceKey, 'left'),
      loadJoinFieldsBySource(joinConfig.value.rightSourceKey, 'right')
    ])
  }

  const handleSidebarItemDblClick = async (item: SidebarTableItem) => {
    if (formData.value.type !== 'sql') return

    ensureSqlConfig()

    const existingIndex = selectedSources.value.findIndex((source) => source.key === item.key)
    if (existingIndex >= 0) {
      const existing = selectedSources.value[existingIndex]
      selectedSources.value.splice(existingIndex, 1)
      selectedSources.value.push(existing)
    } else {
      if (selectedSources.value.length >= MAX_SELECTED_SOURCES) {
        selectedSources.value.shift()
        ElMessage.info('最多配置两个来源，已替换最早添加的来源')
      }
      selectedSources.value.push({
        key: item.key,
        label: item.label,
        type: item.type,
        table: item.table
      })
    }

    activeSourceKey.value = item.key

    if (item.type === 'custom_sql' && !formData.value.config.query) {
      formData.value.config.query = ''
    }

    applyDefaultJoinSources()
    syncSourcesToConfig()
    syncJoinToConfig()
    syncAutoQueryFromSources()
    await Promise.all([
      loadJoinFieldsBySource(joinConfig.value.leftSourceKey, 'left'),
      loadJoinFieldsBySource(joinConfig.value.rightSourceKey, 'right')
    ])
  }

  // CRUD operations
  const ensureDatasetPersisted = async () => {
    if (currentDatasetId.value) {
      return currentDatasetId.value
    }

    const response = await datasetApi.create({
      name: formData.value.name,
      type: formData.value.type,
      datasourceId: formData.value.datasourceId,
      config: formData.value.config
    })

    if (!response.data.success || !response.data.result?.id) {
      throw new Error(response.data.message || '创建数据集失败')
    }

    const newId = response.data.result.id
    currentDatasetId.value = newId
    await router.replace(`/dataset/edit/${newId}`)
    return newId
  }

  // Event handlers
  const handleTypeChange = () => {
    formData.value.config = {}
    fileList.value = []
    previewData.value = []
    loadedPreviewDatasetId.value = ''
    activeSidebarTableKey.value = CUSTOM_SQL_KEY
    resetSqlWorkspace()

    if (formData.value.type === 'sql' && formData.value.datasourceId) {
      ensureSqlConfig()
      loadTables(formData.value.datasourceId)
      return
    }

    tables.value = []
  }

  const handleDatasourceChange = async (datasourceId: string) => {
    if (formData.value.type === 'sql') {
      resetSqlWorkspace()
    }

    if (formData.value.type === 'sql' && datasourceId) {
      ensureSqlConfig()
      await loadTables(datasourceId)
      return
    }

    tables.value = []
  }

  const handleFileChange = (file: UploadUserFile) => {
    fileList.value = [file]
    formData.value.config.file = (file as UploadUserFile & { raw?: File }).raw
  }

  const previewDataset = async () => {
    if (!formRef.value) return

    await formRef.value.validate()
    try {
      await runWithState(
        'preview',
        async () => {
          const id = await ensureDatasetPersisted()
          await fetchPreview(id)
        },
        { successMessage: '预览成功', errorMessage: '预览失败' }
      )
    } catch {
      // Preview failed, keep existing data
    }
  }

  const handleSaveAction = async (action: 'save' | 'save_and_return') => {
    if (!formRef.value) return

    await formRef.value.validate()

    try {
      await runWithState(
        action,
        async () => {
          const isNewDataset = !currentDatasetId.value
          let persistedId = currentDatasetId.value

          if (!isNewDataset) {
            const updatePayload: UpdateDatasetRequest = {
              name: formData.value.name,
              config: formData.value.config,
              action
            }
            const updateResponse = await datasetApi.update(currentDatasetId.value, updatePayload)
            if (!updateResponse.data.success) {
              throw new Error(updateResponse.data.message || '保存失败')
            }
          } else {
            const createResponse = await datasetApi.create({
              name: formData.value.name,
              type: formData.value.type,
              datasourceId: formData.value.datasourceId,
              config: formData.value.config
            })
            if (!createResponse.data.success || !createResponse.data.result?.id) {
              throw new Error(createResponse.data.message || '保存失败')
            }
            persistedId = createResponse.data.result.id
            currentDatasetId.value = persistedId
          }

          await loadSchema(persistedId)

          const navigationPlan = resolveSaveNavigationPlan(action, isNewDataset, true)
          if (navigationPlan.replaceRoute) {
            await router.replace(`/dataset/edit/${persistedId}`)
          }
          if (navigationPlan.returnToList) {
            goBack()
          }
        },
        {
          successMessage: action === 'save' ? '保存成功' : '保存成功，已返回列表',
          errorMessage: '保存失败'
        }
      )
    } catch {
      // Save failed, keep form context
    }
  }

  const handleRefreshData = async () => {
    try {
      await runWithState(
        'refresh',
        async () => {
          if (!currentDatasetId.value) {
            throw new Error('请先保存数据集后再刷新数据')
          }

          const refreshPlan = resolveRefreshLoadPlan(
            currentDatasetId.value,
            activeTab.value,
            loadedPreviewDatasetId.value
          )

          const tasks: Promise<void>[] = []
          if (refreshPlan.loadSchema) {
            tasks.push(loadSchema(currentDatasetId.value))
          }
          if (refreshPlan.loadPreview) {
            tasks.push(fetchPreview(currentDatasetId.value))
          }
          await Promise.all(tasks)
        },
        { successMessage: '数据刷新成功', errorMessage: '数据刷新失败' }
      )
    } catch {
      // Refresh failed
    }
  }

  const handleTabChange = async (name: string | number) => {
    const targetTab = String(name)
    if (targetTab === activeTab.value) return

    const tabPlan = resolveTabSwitchLoadPlan(
      targetTab,
      currentDatasetId.value,
      loadedSchemaDatasetId.value,
      loadedPreviewDatasetId.value
    )

    activeTab.value = targetTab

    if (!tabPlan.loadSchema && !tabPlan.loadPreview) return

    try {
      await runWithState(
        'tab_switch',
        async () => {
          const tasks: Promise<void>[] = []
          if (tabPlan.loadSchema && currentDatasetId.value) {
            tasks.push(loadSchema(currentDatasetId.value))
          }
          if (tabPlan.loadPreview && currentDatasetId.value) {
            tasks.push(fetchPreview(currentDatasetId.value))
          }
          await Promise.all(tasks)
        },
        { errorMessage: '切换标签失败', notifySuccess: false }
      )
    } catch {
      // Tab switch failed
    }
  }

  const handleSelectionChange = (fields: DatasetField[]) => {
    selectedFields.value = fields
  }

  const submitBatchUpdate = async () => {
    const hasPatch = !!(batchConfig.value.type || batchConfig.value.sortOrder || batchConfig.value.displayNamePrefix.trim())
    const precheckMessage = getBatchPrecheckMessage(currentDatasetId.value, selectedFields.value.length, hasPatch)
    if (precheckMessage) {
      ElMessage.warning(precheckMessage)
      return
    }

    try {
      await runWithState(
        'batch_update',
        async () => {
          const fields: BatchUpdateFieldRequest[] = selectedFields.value.map((field) => {
            const payload: BatchUpdateFieldRequest = { fieldId: field.id }
            if (batchConfig.value.type) {
              payload.type = batchConfig.value.type
            }
            if (batchConfig.value.sortOrder) {
              payload.sortOrder = batchConfig.value.sortOrder
            }
            if (batchConfig.value.displayNamePrefix.trim()) {
              payload.displayName = `${batchConfig.value.displayNamePrefix.trim()}${field.displayName || field.name}`
            }
            return payload
          })

          const response = await datasetApi.batchUpdateFields(currentDatasetId.value, { fields })
          if (!response.data.success) {
            throw new Error(response.data.message || '批量更新失败')
          }

          const result = response.data.result
          const resolution = resolveBatchResult(result)

          await loadSchema(currentDatasetId.value)

          if (resolution.isPartialFailure) {
            const failedSet = new Set(resolution.failedFieldIds)
            selectedFields.value = allFields.value.filter((field) => failedSet.has(field.id))
            const details = formatBatchFieldErrors(result.errors || [])
            throw new Error(details ? `以下字段更新失败：${details}` : response.data.message || '批量更新失败')
          }

          selectedFields.value = []
        },
        { successMessage: '批量更新成功', errorMessage: '批量更新失败' }
      )
    } catch {
      // Batch update failed
    }
  }

  // Grouping field handlers
  const resetGroupingFieldForm = () => {
    groupingFieldForm.value = {
      name: '',
      displayName: '',
      dataType: 'string',
      groupingRule: '',
      groupingEnabled: false
    }
    groupingFieldFormRef.value?.clearValidate()
  }

  const openGroupingFieldDialog = () => {
    groupingFieldDialogVisible.value = true
  }

  const submitGroupingField = async () => {
    if (!currentDatasetId.value) {
      ElMessage.warning('请先保存数据集')
      return
    }

    if (!groupingFieldFormRef.value) return

    await groupingFieldFormRef.value.validate()

    await runWithState(
      'create_grouping_field',
      async () => {
        const payload: CreateGroupingFieldRequest = {
          name: groupingFieldForm.value.name.trim(),
          displayName: groupingFieldForm.value.displayName.trim() || undefined,
          type: 'dimension',
          dataType: groupingFieldForm.value.dataType,
          isGroupingField: true,
          groupingRule: groupingFieldForm.value.groupingRule.trim(),
          groupingEnabled: groupingFieldForm.value.groupingEnabled
        }

        const response = await datasetApi.createField(currentDatasetId.value, payload)
        if (!response.data.success) {
          throw new Error(response.data.message || '创建分组字段失败')
        }

        await loadSchema(currentDatasetId.value)
        groupingFieldDialogVisible.value = false
        resetGroupingFieldForm()
      },
      { successMessage: '分组字段创建成功', errorMessage: '创建分组字段失败' }
    )
  }

  const closePreview = () => {
    previewData.value = []
    loadedPreviewDatasetId.value = ''
  }

  const goBack = () => {
    router.push('/dataset')
  }

  // Initialize
  const initialize = async () => {
    await loadDatasources()

    if (!currentDatasetId.value) return

    try {
      await runWithState(
        'init',
        async () => {
          await Promise.all([loadDataset(currentDatasetId.value), loadSchema(currentDatasetId.value)])

          if (formData.value.type === 'sql' && formData.value.datasourceId) {
            await loadTables(formData.value.datasourceId)
            await restoreSqlWorkspaceFromConfig()
          }
        },
        { notifySuccess: false, errorMessage: '初始化页面失败' }
      )
    } catch {
      // Init failed
    }
  }

  // Watchers
  watch(
    () => [joinConfig.value.leftSourceKey, joinConfig.value.rightSourceKey],
    async () => {
      if (formData.value.type !== 'sql') return
      await Promise.all([
        loadJoinFieldsBySource(joinConfig.value.leftSourceKey, 'left'),
        loadJoinFieldsBySource(joinConfig.value.rightSourceKey, 'right')
      ])
    }
  )

  watch(
    () => joinConfig.value,
    () => {
      if (formData.value.type !== 'sql') return
      syncJoinToConfig()
      syncAutoQueryFromSources()
    },
    { deep: true }
  )

  watch(
    selectedSources,
    () => {
      if (formData.value.type !== 'sql') return
      syncSourcesToConfig()
      applyDefaultJoinSources()
      syncJoinToConfig()
      syncAutoQueryFromSources()
    },
    { deep: true }
  )

  return {
    // Refs
    formRef,
    groupingFieldFormRef,
    currentDatasetId,
    activeTab,
    formData,
    formRules,
    workflowState,
    datasources,
    datasourcesLoading,
    tables,
    tablesLoading,
    activeSidebarTableKey,
    activeSourceKey,
    selectedSources,
    leftJoinFieldOptions,
    rightJoinFieldOptions,
    joinFieldLoading,
    joinConfig,
    fileList,
    previewData,
    schema,
    loadedSchemaDatasetId,
    loadedPreviewDatasetId,
    selectedFields,
    groupingFieldDialogVisible,
    groupingFieldForm,
    groupingFieldRules,
    batchConfig,

    // Computed
    allFields,
    hasCustomSqlSource,
    showSqlEditor,
    sidebarTableItems,

    // Methods
    handleSidebarItemDblClick,
    handleTypeChange,
    handleDatasourceChange,
    handleFileChange,
    previewDataset,
    handleSaveAction,
    handleRefreshData,
    handleTabChange,
    handleSelectionChange,
    submitBatchUpdate,
    resetGroupingFieldForm,
    openGroupingFieldDialog,
    submitGroupingField,
    closePreview,
    goBack,
    initialize
  }
}
