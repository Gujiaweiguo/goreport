<template>
  <div class="datasource-manage">
    <h2>数据源管理</h2>
    <div class="page-actions">
      <el-button type="primary" @click="handleCreate">
        <el-icon><Plus /></el-icon>
        创建数据源
      </el-button>
      <el-button type="primary" plain @click="showCreateDialog">
        <span>批量导入</span>
      </el-button>
    </div>

    <el-table
      :data="datasources"
      @selection-change="handleSelectionChange"
      v-loading="loading"
      stripe
      border
      @row-click="handleRowClick"
      style="width: 100%"
    >
      <el-table-column prop="type" label="类型" width="80">
        <template #default="{ row }">
          <el-tag v-if="row.type === 'mysql'" type="success">MySQL</el-tag>
          <el-tag v-else-if="row.type === 'postgresql'" type="warning">PostgreSQL</el-tag>
          <el-tag v-else-if="row.type === 'sqlserver'" type="info">SQL Server</el-tag>
          <el-tag v-else-if="row.type === 'excel' || row.type === 'csv'" type="info">Excel/CSV</el-tag>
          <el-tag v-else type="warning">API 数据源</el-tag>
          <el-tag v-else>Unknown</el-tag>
        </template>
      </el-table-column>

      <el-table-column prop="name" label="名称" width="150">
        <template #default="{ row }">
          <span>{{ row.name || '-' }}</span>
        </template>
      </el-table-column>

      <el-table-column prop="host" label="主机" width="120">
        <template #default="{ row }">
          <span>{{ row.host || '-' }}</span>
        </template>
      </el-table-column>

      <el-table-column prop="port" label="端口" width="80">
        <template #default="{ row }">
          <span>{{ row.port || '-' }}</span>
        </template>
      </el-table-column>

      <el-table-column prop="database" label="数据库" width="100">
        <template #default="{ row }">
          <span>{{ row.database || '-' }}</span>
        </template>
      </el-table-column>

      <el-table-column prop="createdAt" label="创建时间" width="130">
        <template #default="{ row }">
          <span>{{ formatDate(row.createdAt) }}</span>
        </template>
      </el-table-column>

      <el-table-column prop="updatedAt" label="更新时间" width="130">
        <template #default="{ row }">
          <span>{{ formatDate(row.updatedAt) }}</span>
        </template>
      </el-table-column>

      <el-table-column prop="operations" label="操作" width="280">
        <template #default="{ row }">
          <el-button link type="primary" size="small" @click="handleEdit(row)">编辑</el-button>
          <el-button link type="success" size="small" @click="openTableDialog(row)">表/字段</el-button>
          <el-button link type="primary" size="small" @click="openTestConnection(row)">测试连接</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      background="#f5f7fa"
      layout="total, prev, pager, next"
      @current-change="currentPage"
      @size-change="pageSize"
    >
      <div class="page-controls">
        <span>共 {{ total }} 条数据源</span>
        <el-pagination
          layout="prev, pager, next"
          @current-page="currentPage"
          :page-size="pageSize"
        :total="total"
          :background="#f5f7fa"
        />
      </div>
    </el-pagination>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { datasourceApi, type DataSource } from '@/api/datasource'
import { ElMessage, ElMessageBox, ElNotification } from 'element-plus'

interface Datasource {
  id: string
  name: string
  type: string
  host: string
  port: number
  database: string
  username?: string
  password?: string
  tenantId: string
  createdAt?: string
  updatedAt?: string
}

interface OperationState {
  visible: boolean
  operation: string
  result?: any
  error?: string
}

const formatDate = (date?: string) => {
  if (!date) return '-'
  const d = new Date(date)
  return d.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const dialogVisible = ref(false)
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(10)

const datasources = ref<DataSource[]>([])
const total = ref(0)
const searchKeyword = ref('')
const filter = reactive<{}>({ keyword: '', type: '', page: 1, pageSize: 10 })

const operationState = reactive<OperationState>({
  visible: false,
  operation: '',
  result: null,
  error: ''
})

// Datasource form
const createForm = reactive({
  name: '',
  type: 'mysql',
  host: '',
  port: 3306,
  database: '',
  username: '',
  password: ''
})

const editForm = reactive({
  id: '',
  name: '',
  type: '',
  host: '',
  port: 3306,
  database: '',
  username: '',
  password: ''
})

const deleteForm = reactive({
  id: '',
  name: ''
})

const rules = {
  name: [{ required: true, message: '请输入数据源名称', trigger: 'blur' }]
}

const loadDatasources = async () => {
  loading.value = true
  try {
    const { data, total } = await datasourceApi.list(currentPage.value, pageSize.value)
    datasources.value = data || []
    total.value = total || 0
  } catch (error) {
    ElMessage.error('加载数据源失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = async () => {
  currentPage.value = 1
  loadDatasources()
}

const handleCreate = async () => {
  try {
    if (!createForm.name || createForm.type === '') {
      ElMessage.warning('请先选择数据源类型')
      return
    }
    if (createForm.port !== undefined && (createForm.port < 1 || createForm.port > 65535)) {
      ElMessage.warning('端口号范围应在 1-65535 之间')
      return
    }
    if ((createForm.type === 'mysql' || createForm.type === 'postgresql' || createForm.type === 'sqlserver') && !createForm.database) {
      ElMessage.warning('请输入数据库连接名称')
      return
    }

    loading.value = true
    const response = await datasourceApi.create(createForm as any)
    loading.value = false

    if (response.data.success) {
      ElMessage.success('数据源创建成功')
      dialogVisible.value = false
      createForm.id = response.data.result.id
      createForm.name = ''
      createForm.type = ''
      createForm.host = ''
      createForm.port = ''
      createForm.database = ''
      await loadDatasources()
    } else {
      ElMessage.error(response.data.message || '创建失败')
      dialogVisible.value = false
    }
  } catch (error) {
    ElMessage.error('创建失败')
    loading.value = false
  }
}

const handleEdit = async (ds: Datasource) => {
  editForm.id = ds.id
  editForm.name = ds.name
  editForm.type = ds.type
  editForm.host = ds.host
  editForm.port = ds.port
  editForm.database = ds.database
  editForm.username = ds.username || ''
  editForm.password = ds.password || ''
  dialogVisible.value = false

  await loadDatasources()
}

const handleUpdate = async () => {
  if (!editForm.id) {
    ElMessage.warning('请先选择数据源')
    return
  }

  loading.value = true
  const response = await datasourceApi.update(editFormRef.value as any)
  loading.value = false

  if (response.data.success) {
    ElMessage.success('数据源更新成功')
    editFormRef.value = {}
    dialogVisible.value = false
    await loadDatasources()
  } else {
    ElMessage.error(response.data.message || '更新失败')
    dialogVisible.value = false
  }
}

const handleDelete = async (ds: Datasource) => {
  if (!ds.id) {
    ElMessage.warning('请先选择数据源')
    return
  }

  await ElMessageBox.confirm(`确认删除数据源"${ds.name}" 吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })

  try {
    const response = await datasourceApi.delete(ds.id)
    loading.value = true
    if (response.data.success) {
      ElMessage.success('数据源删除成功')
      await loadDatasources()
    } else {
      ElMessage.error(response.data.message || '删除失败')
    }
  } catch (error) {
    ElMessage.error('删除失败')
  }
    loading.value = false
  }
}

const openCreateDialog = () => {
  dialogVisible.value = true
  Object.assign(createForm, {
    type: 'mysql',
    port: 3306,
    database: '',
    username: '',
    password: ''
  })
}

const openEditDialog = (ds: Datasource) => {
  dialogVisible.value = true
  editForm.id = ds.id
  editForm.name = ds.name
  editForm.type = ds.type
  editForm.host = ds.host
  editForm.port = ds.port
  editForm.database = ds.database
  editForm.username = ds.username || ''
  editForm.password = ds.password || ''

  dialogVisible.value = true
}

const openTableDialog = (ds: Datasource) => {
  tableDatasourceId.value = ds.id
  tableDialogVisible.value = true
}

const closeTableDialog = () => {
  tableDialogVisible.value = false
}

const showCreateDialog = () => {
  dialogVisible.value = true
}

const closeDialog = () => {
  dialogVisible.value = false
}
</script>

<style scoped>
.datasource-manage {
  padding: 20px;
}

.page-actions {
  margin-bottom: 20px;
  display: flex;
  gap: 16px;
}

.table-container {
  min-height: 400px;
  background: #f5f7fa;
  border: 1px solid #e0e0e0e;
  border-radius: 8px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
}

.el-button {
  margin-right: 8px;
}

.page-title {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.el-dialog {
  width: 600px;
}

.el-table {
  font-size: 13px;
}

.el-form {
  max-width: 400px;
}

.el-pagination {
  margin-top: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-controls {
  display: flex;
  gap: 12px;
  align-items: center;
}

.el-empty-text {
  color: #909399;
}
</style>
