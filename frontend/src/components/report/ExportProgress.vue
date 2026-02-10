<template>
  <el-dialog
    v-model="visible"
    title="导出进度"
    width="600px"
    :close-on-click-modal="false"
    :show-close="false"
  >
    <div v-if="job.status === 'pending'" class="export-pending">
      <el-progress :percentage="0" :indeterminate="true" />
      <div class="export-status">等待中...</div>
    </div>

    <div v-else-if="job.status === 'processing'" class="export-processing">
      <el-progress :percentage="job.progress" />
      <div class="export-status">
        正在导出... {{ job.progress }}%
      </div>
      <div class="export-time">
        预计剩余时间: {{ estimatedTime }}
      </div>
    </div>

    <div v-else-if="job.status === 'completed'" class="export-completed">
      <el-progress :percentage="100" status="success" />
      <div class="export-status">
        导出完成！
      </div>
      <el-button type="primary" @click="downloadFile">
        <el-icon class="el-icon--left"><Download /></el-icon>
        下载文件
      </el-button>
      <div class="export-file-info">
        文件名: {{ job.fileName }}
        大小: {{ formatFileSize(job.fileSize) }}
      </div>
    </div>

    <div v-else-if="job.status === 'cancelled'" class="export-cancelled">
      <el-progress :percentage="100" status="warning" />
      <div class="export-status">
        已取消导出
      </div>
    </div>

    <div v-else-if="job.status === 'failed'" class="export-failed">
      <el-progress :percentage="100" status="exception" />
      <div class="export-status">
        导出失败
      </div>
      <div class="export-error">
        <div class="error-message">{{ job.error }}</div>
        <el-button @click="copyError" type="text" size="small">
          <el-icon><CopyDocument /></el-icon>
          复制错误信息
        </el-button>
      </div>
    </div>

    <template #footer>
      <el-button 
        v-if="job.status === 'processing'" 
        @click="handleCancel"
        :disabled="cancelling"
      >
        <el-icon class="el-icon--left"><Close /></el-icon>
        取消导出
      </el-button>
      <el-button @click="handleClose">
        {{ job.status === 'completed' ? '关闭' : '最小化' }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
  import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
  import { ElMessage } from 'element-plus'
  import { Download, Close, CopyDocument } from '@element-plus/icons-vue'

  interface ExportJob {
    id: string
    status: 'pending' | 'processing' | 'completed' | 'cancelled' | 'failed'
    progress: number
    fileName: string
    fileSize: number
    error: string
  }

  const props = defineProps<{
    visible: boolean
    jobId: string
  }>()

  const emit = defineEmits<{
    'update:visible': [value: boolean]
  }>()

  const job = ref<ExportJob>({
    id: '',
    status: 'pending',
    progress: 0,
    fileName: '',
    fileSize: 0,
    error: ''
  })

  const cancelling = ref(false)
  const pollTimer = ref<number | null>(null)

  const statusText: Record<string, string> = {
    pending: '等待中',
    processing: '处理中',
    completed: '已完成',
    cancelled: '已取消',
    failed: '失败'
  }

  const statusColor: Record<string, string> = {
    pending: '#409eff',
    processing: '#67c23a',
    completed: '#67c23a',
    cancelled: '#e6a23c',
    failed: '#f56c6c'
  }

  const estimatedTime = computed(() => {
    if (job.value.status === 'processing') {
      const remaining = 100 - job.value.progress
      if (remaining > 0) {
        const timePerPercent = 30 / 100
        const seconds = Math.ceil(remaining * timePerPercent)
        if (seconds < 60) {
          return `${seconds}秒`
        } else {
          const minutes = Math.ceil(seconds / 60)
          return `${minutes}分钟`
        }
      }
    }
    return '计算中...'
  })

  async function fetchJobStatus() {
    try {
      const response = await fetch(`/api/v1/jmreport/export/status/${props.jobId}`)
      const data = await response.json()
      if (data.success) {
        job.value = data.result
      }
    } catch (error: any) {
      console.error('Failed to fetch job status:', error)
    }
  }

  async function downloadFile() {
    try {
      const response = await fetch(`/api/v1/jmreport/export/download/${props.jobId}`)
      if (response.ok) {
        const blob = await response.blob()
        const url = window.URL.createObjectURL(blob)
        const link = document.createElement('a')
        link.href = url
        link.download = job.value.fileName
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
        window.URL.revokeObjectURL(url)
        ElMessage.success('下载成功')
      }
    } catch (error: any) {
      ElMessage.error(error.message || '下载失败')
    }
  }

  async function handleCancel() {
    cancelling.value = true
    try {
      const response = await fetch(`/api/v1/jmreport/export/cancel/${props.jobId}`, {
        method: 'POST'
      })
      const data = await response.json()
      if (data.success) {
        job.value = {
          ...job.value,
          status: 'cancelled'
        }
        ElMessage.success('已取消导出')
      }
    } catch (error: any) {
      ElMessage.error(error.message || '取消失败')
    } finally {
      cancelling.value = false
    }
  }

  function copyError() {
    if (job.value.error) {
      navigator.clipboard.writeText(job.value.error)
      ElMessage.success('已复制到剪贴板')
    }
  }

  function formatFileSize(bytes: number): string {
    if (bytes === 0) return '0 B'
    const k = 1024
    const m = 1024 * 1024
    const g = 1024 * 1024 * 1024
    if (bytes < k) return `${(bytes / k).toFixed(2)} KB`
    if (bytes < m) return `${(bytes / m).toFixed(2)} MB`
    return `${(bytes / g).toFixed(2)} GB`
  }

  function startPolling() {
    fetchJobStatus()
    pollTimer.value = window.setInterval(() => {
      fetchJobStatus()
      if (job.value.status === 'completed' || job.value.status === 'failed' || job.value.status === 'cancelled') {
        stopPolling()
      }
    }, 2000)
  }

  function stopPolling() {
    if (pollTimer.value !== null) {
      clearInterval(pollTimer.value)
      pollTimer.value = null
    }
  }

  function handleClose() {
    stopPolling()
    emit('update:visible', false)
  }

  onMounted(() => {
    job.value.id = props.jobId
    if (props.visible) {
      startPolling()
    }
  })

  onUnmounted(() => {
    stopPolling()
  })

  watch(() => props.visible, (newVal) => {
    if (newVal) {
      startPolling()
    } else {
      stopPolling()
    }
  })
</script>

<style scoped>
.export-pending,
.export-processing,
.export-completed,
.export-cancelled,
.export-failed {
  padding: 20px 0;
}

.export-status {
  text-align: center;
  font-size: 16px;
  font-weight: 500;
  margin: 16px 0;
  color: #333;
}

.export-time {
  text-align: center;
  font-size: 14px;
  color: #666;
  margin-bottom: 20px;
}

.export-file-info {
  text-align: center;
  font-size: 14px;
  color: #666;
  margin: 16px 0;
  padding: 16px;
  background: #f5f5f5;
  border-radius: 8px;
}

.export-error {
  text-align: center;
  margin-top: 20px;
}

.error-message {
  color: #f56c6c;
  background: #fef2f2f;
  padding: 12px;
  border-radius: 8px;
  margin-bottom: 12px;
  font-size: 14px;
}
</style>
