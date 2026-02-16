<template>
  <div class="report-list">
    <el-card class="page-header">
      <el-button type="primary" @click="handleCreate">新建报表</el-button>
    </el-card>

    <el-card class="table-container">
      <LoadingState v-if="loading" :loading="loading" text="加载报表列表..." />
      <el-empty v-else-if="reports.length === 0" description="暂无报表，点击上方按钮创建" />
      <el-table v-else :data="reports" style="width: 100%">
        <el-table-column prop="name" label="报表名称" width="200" />
        <el-table-column prop="code" label="报表编码" width="150" />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="row.type === 'dashboard' ? 'success' : 'primary'" size="small">
              {{ row.type === 'dashboard' ? '仪表盘' : '报表' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="浏览次数" width="100">
          <template #default="{ row }">
            {{ row.viewCount }}
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="更新时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.updatedAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" @click="handlePreview(row)">预览</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { reportApi, type Report } from '@/api/report'
import LoadingState from '@/components/common/LoadingState.vue'
import { getErrorMessage } from '@/utils/errorHandling'

const router = useRouter()
const loading = ref(false)
const reports = ref<Report[]>([])

async function loadReports() {
  loading.value = true
  try {
    const response = await reportApi.list()
    if (response.data.success) {
      reports.value = response.data.result || []
    } else {
      ElMessage.error(response.data.message || '加载报表列表失败')
    }
  } catch (error: unknown) {
    ElMessage.error(getErrorMessage(error, '加载报表列表失败'))
  } finally {
    loading.value = false
  }
}

function handleCreate() {
  router.push({ name: 'ReportDesigner' })
}

function handleEdit(row: Report) {
  router.push({
    name: 'ReportDesigner',
    query: { id: row.id }
  })
}

function handlePreview(row: Report) {
  router.push({
    name: 'ReportPreview',
    query: { id: row.id }
  })
}

async function handleDelete(row: Report) {
  try {
    await ElMessageBox.confirm(`确认删除报表 "${row.name}" 吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    const response = await reportApi.delete(row.id)
    if (response.data.success) {
      ElMessage.success('删除报表成功')
      await loadReports()
    } else {
      ElMessage.error(response.data.message || '删除报表失败')
    }
  } catch (error: unknown) {
    if (error !== 'cancel') {
      ElMessage.error('删除报表失败')
    }
  }
}

function formatDate(dateStr: string) {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

onMounted(() => {
  loadReports()
})
</script>

<style scoped>
.report-list {
  padding: 20px;
}

.page-header {
  margin-bottom: 20px;
}

.table-container {
  min-height: 400px;
}
</style>
