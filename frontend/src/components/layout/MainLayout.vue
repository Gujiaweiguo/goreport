<template>
  <div class="layout">
    <NoPermission v-if="showNoPermission" />
    <aside v-else class="layout-sidebar">
      <div class="sidebar-header">
        <h2 class="app-title">JimuReport</h2>
      </div>
      <el-menu
        :default-active="activeMenu"
        class="sidebar-menu"
        @select="handleMenuSelect"
      >
        <el-menu-item index="/datasource">
          <el-icon><Setting /></el-icon>
          <template #title>数据源管理</template>
        </el-menu-item>
        <el-menu-item index="/report/designer">
          <el-icon><Edit /></el-icon>
          <template #title>报表设计器</template>
        </el-menu-item>
        <el-menu-item index="/report/preview">
          <el-icon><View /></el-icon>
          <template #title>报表预览</template>
        </el-menu-item>
        <el-menu-item index="/dashboard/designer">
          <el-icon><Monitor /></el-icon>
          <template #title>大屏设计器</template>
        </el-menu-item>
        <el-menu-item index="/chart/editor">
          <el-icon><PieChart /></el-icon>
          <template #title>图表编辑器</template>
        </el-menu-item>
      </el-menu>
    </aside>
    <div v-if="!showNoPermission" class="layout-main">
      <header class="layout-header">
        <div class="header-left"></div>
        <div class="header-right">
          <el-button text @click="handleLogout">退出登录</el-button>
        </div>
      </header>
      <main class="layout-content">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Setting, Edit, View, Monitor, PieChart } from '@element-plus/icons-vue'
import NoPermission from '@/components/common/NoPermission.vue'

const route = useRoute()
const router = useRouter()
const showNoPermission = ref(false)
const lastError = ref<any>(null)

const activeMenu = computed(() => route.path)

function handleMenuSelect(index: string) {
  showNoPermission.value = false
  if (index !== route.path) {
    router.push(index)
  }
}

function handleLogout() {
  localStorage.removeItem('token')
  router.push('/login')
}

watch(() => route.meta, (newMeta, oldMeta) => {
  if (newMeta?.error) {
    showNoPermission.value = true
    lastError.value = newMeta.error
  } else {
    showNoPermission.value = false
    lastError.value = null
  }
}, { immediate: true })
</script>

<style scoped>
.layout {
  display: flex;
  height: 100vh;
  overflow: hidden;
}

.layout-sidebar {
  width: 240px;
  background-color: #304156;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
}

.sidebar-header {
  padding: 20px;
  text-align: center;
  border-bottom: 1px solid #434a50;
}

.app-title {
  color: #fff;
  font-size: 20px;
  margin: 0;
}

.sidebar-menu {
  border: none;
  background-color: transparent;
  flex: 1;
}

.sidebar-menu .el-menu-item {
  color: #bfcbd9;
}

.sidebar-menu .el-menu-item:hover {
  background-color: #263445;
}

.sidebar-menu .el-menu-item.is-active {
  background-color: #409eff;
  color: #fff;
}

.layout-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.layout-header {
  height: 60px;
  border-bottom: 1px solid #dcdfe6;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  background-color: #fff;
}

.header-left {
  flex: 1;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.layout-content {
  flex: 1;
  overflow: auto;
  background-color: #f5f7fa;
}
</style>
