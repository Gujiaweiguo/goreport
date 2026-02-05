<template>
  <div class="datasource-manage">
    <el-card class="page-header">
      <el-button type="primary" @click="showCreateDialog">创建数据源</el-button>
    </el-card>

    <el-card class="table-container">
      <el-table :data="datasources" style="width: 100%" v-loading="loading">
        <el-table-column prop="name" label="名称" width="180" />
        <el-table-column prop="type" label="类型" width="100" />
        <el-table-column label="连接地址" width="250">
          <template #default="{ row }">
            {{ row.host }}:{{ row.port }}
          </template>
        </el-table-column>
        <el-table-column prop="database" label="数据库" width="150" />
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column label="操作" width="280">
          <template #default="{ row }">
            <el-button size="small" @click="handleTest(row)">测试连接</el-button>
            <el-button size="small" @click="showEditDialog(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
      @close="handleDialogClose"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="formRules"
        label-width="100px"
      >
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入数据源名称" />
        </el-form-item>

        <el-form-item label="类型" prop="type">
          <el-select v-model="form.type" placeholder="请选择数据源类型">
            <el-option label="MySQL" value="mysql" />
            <el-option label="PostgreSQL" value="postgresql" />
          </el-select>
        </el-form-item>

        <el-form-item label="主机" prop="host">
          <el-input v-model="form.host" placeholder="请输入主机地址" />
        </el-form-item>

        <el-form-item label="端口" prop="port">
          <el-input-number v-model="form.port" :min="1" :max="65535" />
        </el-form-item>

        <el-form-item label="数据库" prop="database">
          <el-input v-model="form.database" placeholder="请输入数据库名称" />
        </el-form-item>

        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" placeholder="请输入用户名" />
        </el-form-item>

        <el-form-item label="密码" prop="password">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="请输入密码"
            show-password
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="handleDialogClose">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">
          {{ isEdit ? '更新' : '创建' }}
        </el-button>
      </template>
    </el-dialog>

    <el-dialog
      v-model="testDialogVisible"
      title="测试连接"
      width="400px"
    >
      <div v-loading="testing">
        <p v-if="!testResult">点击下方按钮测试连接...</p>
        <div v-else>
          <el-result
            :icon="testResult.success ? 'success' : 'error'"
            :title="testResult.success ? '连接成功' : '连接失败'"
            :sub-title="testResult.message"
          />
        </div>
      </div>
      <template #footer>
        <el-button @click="testDialogVisible = false">关闭</el-button>
        <el-button type="primary" :loading="testing" @click="runTest">
          测试
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { datasourceApi, type DataSource, type CreateDataSourceRequest } from '@/api/datasource'

const dialogVisible = ref(false)
const testDialogVisible = ref(false)
const dialogTitle = ref('创建数据源')
const loading = ref(false)
const submitting = ref(false)
const testing = ref(false)
const isEdit = ref(false)
const currentEditId = ref('')
const formRef = ref<FormInstance>()

const datasources = ref<DataSource[]>([])
const testResult = ref<{ success: boolean; message: string } | null>(null)

interface DataSourceForm {
  name: string
  type: string
  host: string
  port: number
  database: string
  username: string
  password: string
}

const form = reactive<DataSourceForm>({
  name: '',
  type: 'mysql',
  host: 'localhost',
  port: 3306,
  database: '',
  username: '',
  password: ''
})

const formRules = reactive<FormRules<DataSourceForm>>({
  name: [{ required: true, message: '请输入数据源名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择数据源类型', trigger: 'change' }],
  host: [{ required: true, message: '请输入主机地址', trigger: 'blur' }],
  port: [{ required: true, message: '请输入端口号', trigger: 'blur' }],
  database: [{ required: true, message: '请输入数据库名称', trigger: 'blur' }],
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
})

async function loadDatasources() {
  loading.value = true
  try {
    const response = await datasourceApi.list()
    if (response.data.success) {
      datasources.value = response.data.result || []
    } else {
      ElMessage.error(response.data.message || '加载数据源失败')
    }
  } catch (error: any) {
    ElMessage.error('加载数据源失败')
  } finally {
    loading.value = false
  }
}

function showCreateDialog() {
  isEdit.value = false
  dialogTitle.value = '创建数据源'
  resetForm()
  dialogVisible.value = true
}

function showEditDialog(row: DataSource) {
  isEdit.value = true
  dialogTitle.value = '编辑数据源'
  form.name = row.name
  form.type = row.type
  form.host = row.host
  form.port = row.port
  form.database = row.database
  form.username = row.username
  form.password = ''
  currentEditId.value = row.id
  dialogVisible.value = true
}

function resetForm() {
  form.name = ''
  form.type = 'mysql'
  form.host = 'localhost'
  form.port = 3306
  form.database = ''
  form.username = ''
  form.password = ''
}

function handleDialogClose() {
  dialogVisible.value = false
  resetForm()
  if (formRef.value) {
    formRef.value.clearValidate()
  }
}

async function handleSubmit() {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
  } catch {
    return
  }

  submitting.value = true

  try {
    if (isEdit.value) {
      const response = await datasourceApi.update(currentEditId.value, form)
      if (response.data.success) {
        ElMessage.success('更新数据源成功')
        await loadDatasources()
        handleDialogClose()
      } else {
        ElMessage.error(response.data.message || '更新数据源失败')
      }
    } else {
      const response = await datasourceApi.create(form)
      if (response.data.success) {
        ElMessage.success('创建数据源成功')
        await loadDatasources()
        handleDialogClose()
      } else {
        ElMessage.error(response.data.message || '创建数据源失败')
      }
    }
  } catch (error: any) {
    ElMessage.error('操作失败')
  } finally {
    submitting.value = false
  }
}

async function handleDelete(row: DataSource) {
  try {
    await ElMessageBox.confirm(`确认删除数据源 "${row.name}" 吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    const response = await datasourceApi.delete(row.id)
    if (response.data.success) {
      ElMessage.success('删除数据源成功')
      await loadDatasources()
    } else {
      ElMessage.error(response.data.message || '删除数据源失败')
    }
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error('删除数据源失败')
    }
  }
}

async function handleTest(row: DataSource) {
  testDialogVisible.value = true
  testResult.value = null

  const testData: CreateDataSourceRequest = {
    name: row.name,
    type: row.type,
    host: row.host,
    port: row.port,
    database: row.database,
    username: row.username,
    password: ''
  }

  testing.value = true
  try {
    const response = await datasourceApi.test(testData)
    testResult.value = {
      success: response.data.success,
      message: response.data.message || (response.data.success ? '连接成功' : '连接失败')
    }
  } catch (error: any) {
    testResult.value = {
      success: false,
      message: error.message || '连接测试失败'
    }
  } finally {
    testing.value = false
  }
}

function runTest() {
  const currentDataSource = datasources.value.find(ds => ds.id === currentEditId.value)
  if (currentDataSource) {
    handleTest(currentDataSource)
  }
}

onMounted(() => {
  loadDatasources()
})
</script>

<style scoped>
.datasource-manage {
  padding: 20px;
}

.page-header {
  margin-bottom: 20px;
}

.table-container {
  min-height: 400px;
}
</style>
