import { createRouter, createWebHistory } from 'vue-router'
import MainLayout from '@/components/layout/MainLayout.vue'
import Login from '@/views/Login.vue'
import DatasourceManage from '@/views/DatasourceManage.vue'
import ReportDesigner from '@/views/ReportDesigner.vue'
import ReportPreview from '@/views/ReportPreview.vue'
import DashboardDesigner from '@/views/DashboardDesigner.vue'
import ChartEditor from '@/views/ChartEditor.vue'

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
        path: 'report/designer',
        name: 'ReportDesigner',
        component: ReportDesigner,
        meta: { requiresAuth: true, title: '报表设计器', order: 4 }
      },
      {
        path: 'report/preview',
        name: 'ReportPreview',
        component: ReportPreview,
        meta: { requiresAuth: true, title: '报表预览', order: 5 }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')

  if (to.path === '/login' && token) {
    next('/dashboard/designer')
  } else if (to.meta.requiresAuth && !token) {
    next('/login')
  } else {
    next()
  }
})

router.onError((error) => {
  console.error('路由错误:', error)
  if (error.name === 'ChunkLoadError') {
    alert('资源加载失败，请刷新页面重试')
  }
})

export default router
