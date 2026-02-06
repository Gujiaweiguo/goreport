<template>
  <div class="error-state">
    <div class="error-icon" :class="`error-${type}`">
      <el-icon :size="48">
        <component :is="iconComponent" />
      </el-icon>
    </div>
    <h3 class="error-title">{{ title }}</h3>
    <p v-if="message" class="error-message">{{ message }}</p>
    <div v-if="details" class="error-details">
      <pre>{{ details }}</pre>
    </div>
    <div class="error-actions">
      <el-button v-if="canRetry" type="primary" @click="handleRetry">
        <el-icon><Refresh /></el-icon>
        重试
      </el-button>
      <el-button v-if="canGoBack" @click="handleGoBack">
        <el-icon><Back /></el-icon>
        返回
      </el-button>
      <el-button v-if="canGoHome" @click="handleGoHome">
        <el-icon><HomeFilled /></el-icon>
        返回首页
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { Warning, CircleCloseFilled, InfoFilled, Refresh, Back, HomeFilled } from '@element-plus/icons-vue'

interface Props {
  type?: 'network' | 'api' | 'config' | 'permission'
  title?: string
  message?: string
  details?: string
  canRetry?: boolean
  canGoBack?: boolean
  canGoHome?: boolean
  onRetry?: () => void
}

const props = withDefaults(defineProps<Props>(), {
  type: 'api',
  title: '操作失败',
  canRetry: true,
  canGoBack: true,
  canGoHome: false
})

const router = useRouter()

const iconComponent = computed(() => {
  switch (props.type) {
    case 'network':
      return CircleCloseFilled
    case 'api':
      return Warning
    case 'config':
      return InfoFilled
    case 'permission':
      return InfoFilled
    default:
      return CircleCloseFilled
  }
})

function handleRetry() {
  if (props.onRetry) {
    props.onRetry()
  }
}

function handleGoBack() {
  router.back()
}

function handleGoHome() {
  router.push('/dashboard/designer')
}
</script>

<style scoped>
.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 40px;
  min-height: 400px;
  text-align: center;
}

.error-icon {
  margin-bottom: 24px;
}

.error-network {
  color: #ef4444;
}

.error-api {
  color: #f59e0b;
}

.error-config {
  color: #3b82f6;
}

.error-permission {
  color: #8b5cf6;
}

.error-title {
  margin: 0 0 16px;
  font-size: 20px;
  font-weight: 600;
  color: #374151;
}

.error-message {
  margin: 0 0 24px;
  font-size: 14px;
  color: #64748b;
  max-width: 600px;
  line-height: 1.6;
}

.error-details {
  margin: 0 0 24px;
  padding: 16px;
  background: #f5f7fa;
  border-radius: 8px;
  text-align: left;
  max-width: 600px;
}

.error-details pre {
  margin: 0;
  font-size: 12px;
  color: #909399;
  white-space: pre-wrap;
  word-break: break-all;
}

.error-actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  justify-content: center;
}
</style>
