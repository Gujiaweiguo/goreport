<template>
  <div class="no-permission">
    <div class="permission-content">
      <el-icon :size="64" class="permission-icon">
        <Lock />
      </el-icon>
      <h3 class="permission-title">无访问权限</h3>
      <p class="permission-message">
        {{ message || '您没有访问该页面的权限，请联系管理员获取相应权限。' }}
      </p>
      <div class="permission-actions">
        <el-button type="primary" @click="handleGoHome">
          <el-icon><HomeFilled /></el-icon>
          返回首页
        </el-button>
        <el-button v-if="canLogin" @click="handleLogin">
          <el-icon><User /></el-icon>
          重新登录
        </el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { Lock, HomeFilled, User } from '@element-plus/icons-vue'

interface Props {
  message?: string
  canLogin?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  canLogin: true
})

const router = useRouter()

function handleGoHome() {
  router.push('/dashboard/designer')
}

function handleLogin() {
  router.push('/login')
}
</script>

<style scoped>
.no-permission {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.permission-content {
  max-width: 480px;
  background: #ffffff;
  border-radius: 16px;
  padding: 48px 40px;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
  text-align: center;
}

.permission-icon {
  color: #8b5cf6;
  margin-bottom: 24px;
}

.permission-title {
  margin: 0 0 16px;
  font-size: 24px;
  font-weight: 600;
  color: #374151;
}

.permission-message {
  margin: 0 0 32px;
  font-size: 14px;
  color: #64748b;
  line-height: 1.6;
}

.permission-actions {
  display: flex;
  gap: 12px;
  justify-content: center;
}
</style>
