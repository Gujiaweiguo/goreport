import { createRouter, createWebHistory } from 'vue-router'
import Home from '@/views/Home.vue'
import Login from '@/views/Login.vue'
import DatasourceManage from '@/views/DatasourceManage.vue'
import ReportDesigner from '@/views/ReportDesigner.vue'
import ReportPreview from '@/views/ReportPreview.vue'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: Login
  },
  {
    path: '/',
    name: 'Home',
    component: Home,
    meta: { requiresAuth: true }
  },
  {
    path: '/datasource',
    name: 'DatasourceManage',
    component: DatasourceManage,
    meta: { requiresAuth: true }
  },
  {
    path: '/report/designer',
    name: 'ReportDesigner',
    component: ReportDesigner,
    meta: { requiresAuth: true, title: '报表设计器' }
  },
  {
    path: '/report/preview',
    name: 'ReportPreview',
    component: ReportPreview,
    meta: { requiresAuth: true, title: '报表预览' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')

  if (to.meta.requiresAuth && !token) {
    next('/login')
  } else if (to.path === '/login' && token) {
    next('/')
  } else {
    next()
  }
})

export default router
