<template>
  <div class="report-preview">
    <el-card class="preview-toolbar" shadow="never">
      <div class="toolbar-left">
        <span class="toolbar-title">报表预览</span>
        <el-tag size="small" effect="dark" type="info">ID: {{ reportId || '未指定' }}</el-tag>
      </div>
      <div class="toolbar-actions">
        <el-button @click="handleRefresh" :loading="loading">刷新</el-button>
        <el-button @click="handleExport">导出</el-button>
      </div>
    </el-card>

    <div class="preview-content" v-loading="loading">
      <div v-if="!renderedHTML" class="preview-empty">
        <el-empty description="暂无预览内容" />
      </div>
      <div v-else class="preview-html" v-html="renderedHTML"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { reportApi } from '@/api/report'

const route = useRoute()
const renderedHTML = ref('')
const loading = ref(false)

const reportId = computed(() => route.query.id as string | undefined)

async function handleRefresh() {
  if (!reportId.value) {
    ElMessage.warning('未指定报表 ID')
    return
  }

  loading.value = true
  try {
    const response = await reportApi.preview({ id: reportId.value })
    if (response.data.success) {
      renderedHTML.value = response.data.result?.html || ''
    } else {
      ElMessage.error(response.data.message || '预览加载失败')
    }
  } catch (error: any) {
    ElMessage.error('预览加载失败')
  } finally {
    loading.value = false
  }
}

function handleExport() {
  ElMessage.info('导出功能开发中')
}

onMounted(async () => {
  if (reportId.value) {
    await handleRefresh()
  }
})
</script>

<style scoped>
.report-preview {
  min-height: 100vh;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  background: radial-gradient(circle at top left, #eef2f6, #f8fafc 35%, #ffffff);
  font-family: "IBM Plex Sans", "Noto Sans", sans-serif;
}

.preview-toolbar {
  border-radius: 14px;
  border: 1px solid #e2e8f0;
  background: linear-gradient(145deg, #0f172a, #1f2937);
  color: #f8fafc;
}

.preview-toolbar :deep(.el-card__body) {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 16px;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.toolbar-title {
  font-size: 16px;
  font-weight: 600;
  letter-spacing: 0.6px;
}

.toolbar-actions {
  display: flex;
  gap: 10px;
}

.preview-content {
  flex: 1;
  border-radius: 14px;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  box-shadow: 0 12px 24px rgba(15, 23, 42, 0.08);
  padding: 16px;
  overflow: auto;
}

.preview-empty {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.preview-html {
  min-height: 400px;
}

.preview-html :deep(table) {
  border-collapse: collapse;
}

.preview-html :deep(td),
.preview-html :deep(th) {
  border: 1px solid #e2e8f0;
  padding: 6px 10px;
}
</style>
