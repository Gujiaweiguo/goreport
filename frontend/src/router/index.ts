import { createRouter, createWebHistory } from 'vue-router'
import MainLayout from '@/components/layout/MainLayout.vue'
import Login from '@/views/Login.vue'
import DatasourceManage from '@/views/DatasourceManage.vue'
import DatasetList from '@/views/dataset/DatasetList.vue'
import DatasetEdit from '@/views/dataset/DatasetEdit.vue'
import ReportDesigner from '@/views/ReportDesigner.vue'
import ReportList from '@/views/ReportList.vue'
import ReportPreview from '@/views/ReportPreview.vue'
import DashboardDesigner from '@/views/DashboardDesigner.vue'
import ChartEditor from '@/views/ChartEditor.vue'
import { auth } from '@/api/auth'
import { userApi } from '@/api/user'
import axios from 'axios'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: Login
  },
  {
    path: '/',
    component: MainLayout,
    meta: { requiresAuth: true },
    redirect: '/dashboard/designer',
    children: [
      {
        path: 'dashboard/designer',
        name: 'DashboardDesigner',
        component: DashboardDesigner,
        meta: { requiresAuth: true, title: '大屏设计器', order: 1 }
      },
      {
        path: 'chart/editor',
        name: 'ChartEditor',
        component: ChartEditor,
        meta: { requiresAuth: true, title: '图表编辑器', order: 2 }
      },
      {
        path: 'datasource',
        name: 'DatasourceManage',
        component: DatasourceManage,
        meta: { requiresAuth: true, title: '数据源管理', order: 3 }
      },
      {
        path: 'dataset',
        name: 'DatasetList',
        component: DatasetList,
        meta: { requiresAuth: true, title: '数据集管理', order: 4 }
      },
      {
        path: 'dataset/create',
        name: 'DatasetCreate',
        component: DatasetEdit,
        meta: { requiresAuth: true, title: '新建数据集', order: 4 }
      },
      {
        path: 'dataset/edit/:id',
        name: 'DatasetEdit',
        component: DatasetEdit,
        meta: { requiresAuth: true, title: '编辑数据集', order: 4 }
      },
      {
        path: 'report/list',
        name: 'ReportList',
        component: ReportList,
        meta: { requiresAuth: true, title: '报表列表', order: 4 }
      },
      {
        path: 'report/designer',
        name: 'ReportDesigner',
        component: ReportDesigner,
        meta: { requiresAuth: true, title: '报表设计器', order: 5 }
      },
      {
        path: 'report/preview',
        name: 'ReportPreview',
        component: ReportPreview,
        meta: { requiresAuth: true, title: '报表预览', order: 6 }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

let sessionValidationPromise: Promise<boolean> | null = null

async function validateSession(): Promise<boolean> {
  if (sessionValidationPromise) {
    return sessionValidationPromise
  }

  sessionValidationPromise = (async () => {
    const token = auth.getToken()
    if (!token) {
      return false
    }

    try {
      const response = await userApi.getMe()
      return !!response.data.success
    } catch (error) {
      if (axios.isAxiosError(error) && [401, 403].includes(error.response?.status ?? 0)) {
        auth.clearSession()
        return false
      }

      // Keep the local session for transient failures (network/CSP/proxy),
      // otherwise users can be stuck on login after successful auth.
      return true
    }
  })()

  try {
    return await sessionValidationPromise
  } finally {
    sessionValidationPromise = null
  }
}

router.beforeEach(async (to) => {
  const token = auth.getToken()

  if (to.path === '/login' && token) {
    const valid = await validateSession()
    return valid ? '/dashboard/designer' : true
  }

  if (to.meta.requiresAuth) {
    if (!token) {
      return '/login'
    }
    const valid = await validateSession()
    if (!valid) {
      return '/login'
    }
  }

  return true
})

router.onError((error) => {
  console.error('路由错误:', error)
  if (error.name === 'ChunkLoadError') {
    alert('资源加载失败，请刷新页面重试')
  }
})

export default router
