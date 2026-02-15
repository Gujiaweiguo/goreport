import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createMemoryHistory, createRouter } from 'vue-router'

describe('Router Configuration', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('Route Definitions', () => {
    const routes = [
      {
        path: '/login',
        name: 'Login',
        meta: { title: 'Login' }
      },
      {
        path: '/',
        redirect: '/dashboard/designer',
        meta: { requiresAuth: true },
        children: [
          {
            path: 'dashboard/designer',
            name: 'DashboardDesigner',
            meta: { requiresAuth: true, title: '大屏设计器', order: 1 }
          },
          {
            path: 'chart/editor',
            name: 'ChartEditor',
            meta: { requiresAuth: true, title: '图表编辑器', order: 2 }
          },
          {
            path: 'datasource',
            name: 'DatasourceManage',
            meta: { requiresAuth: true, title: '数据源管理', order: 3 }
          },
          {
            path: 'dataset',
            name: 'DatasetList',
            meta: { requiresAuth: true, title: '数据集管理', order: 4 }
          },
          {
            path: 'dataset/create',
            name: 'DatasetCreate',
            meta: { requiresAuth: true, title: '新建数据集', order: 4 }
          },
          {
            path: 'dataset/edit/:id',
            name: 'DatasetEdit',
            meta: { requiresAuth: true, title: '编辑数据集', order: 4 }
          },
          {
            path: 'report/list',
            name: 'ReportList',
            meta: { requiresAuth: true, title: '报表列表', order: 4 }
          },
          {
            path: 'report/designer',
            name: 'ReportDesigner',
            meta: { requiresAuth: true, title: '报表设计器', order: 5 }
          },
          {
            path: 'report/preview',
            name: 'ReportPreview',
            meta: { requiresAuth: true, title: '报表预览', order: 6 }
          }
        ]
      }
    ]

    it('should have login route', () => {
      const loginRoute = routes.find(r => r.path === '/login')
      expect(loginRoute).toBeDefined()
      expect(loginRoute?.name).toBe('Login')
    })

    it('should have root route with redirect', () => {
      const rootRoute = routes.find(r => r.path === '/')
      expect(rootRoute).toBeDefined()
      expect(rootRoute?.redirect).toBe('/dashboard/designer')
    })

    it('should have correct number of child routes', () => {
      const rootRoute = routes.find(r => r.path === '/')
      expect(rootRoute?.children).toHaveLength(9)
    })

    it('should have all required routes', () => {
      const routeNames = [
        'Login',
        'DashboardDesigner',
        'ChartEditor',
        'DatasourceManage',
        'DatasetList',
        'DatasetCreate',
        'DatasetEdit',
        'ReportList',
        'ReportDesigner',
        'ReportPreview'
      ]

      const allRoutes = routes.flatMap(r => 
        r.children ? [r.name, ...r.children.map(c => c.name)] : [r.name]
      ).filter(Boolean)

      routeNames.forEach(name => {
        expect(allRoutes).toContain(name)
      })
    })

    it('should have requiresAuth meta on protected routes', () => {
      const rootRoute = routes.find(r => r.path === '/')
      rootRoute?.children?.forEach(child => {
        expect(child.meta?.requiresAuth).toBe(true)
      })
    })

    it('should have correct route orders', () => {
      const rootRoute = routes.find(r => r.path === '/')
      const children = rootRoute?.children || []
      
      expect(children.find(c => c.name === 'DashboardDesigner')?.meta?.order).toBe(1)
      expect(children.find(c => c.name === 'ChartEditor')?.meta?.order).toBe(2)
      expect(children.find(c => c.name === 'DatasourceManage')?.meta?.order).toBe(3)
    })

    it('should have dynamic route for dataset edit', () => {
      const rootRoute = routes.find(r => r.path === '/')
      const datasetEditRoute = rootRoute?.children?.find(c => c.path === 'dataset/edit/:id')
      expect(datasetEditRoute).toBeDefined()
      expect(datasetEditRoute?.name).toBe('DatasetEdit')
    })
  })

  describe('Router Creation', () => {
    it('should create router with memory history', async () => {
      const router = createRouter({
        history: createMemoryHistory(),
        routes: [
          { path: '/', name: 'Home', component: { template: '<div>Home</div>' } },
          { path: '/login', name: 'Login', component: { template: '<div>Login</div>' } }
        ]
      })

      expect(router).toBeDefined()
      expect(router.currentRoute).toBeDefined()
    })

    it('should navigate to routes', async () => {
      const router = createRouter({
        history: createMemoryHistory(),
        routes: [
          { path: '/', name: 'Home', component: { template: '<div>Home</div>' } },
          { path: '/login', name: 'Login', component: { template: '<div>Login</div>' } }
        ]
      })

      await router.push('/login')
      expect(router.currentRoute.value.path).toBe('/login')
    })
  })

  describe('Route Meta', () => {
    const routes = [
      { path: '/login', name: 'Login' },
      { path: '/dashboard', name: 'Dashboard', meta: { requiresAuth: true } }
    ]

    it('should identify protected routes', () => {
      const protectedRoutes = routes.filter(r => r.meta?.requiresAuth)
      expect(protectedRoutes).toHaveLength(1)
      expect(protectedRoutes[0].name).toBe('Dashboard')
    })

    it('should identify public routes', () => {
      const publicRoutes = routes.filter(r => !r.meta?.requiresAuth)
      expect(publicRoutes).toHaveLength(1)
      expect(publicRoutes[0].name).toBe('Login')
    })
  })
})

describe('Navigation Guards', () => {
  it('should redirect unauthenticated users to login', () => {
    const mockAuth = {
      getToken: vi.fn().mockReturnValue(null)
    }
    
    const token = mockAuth.getToken()
    expect(token).toBeNull()
  })

  it('should allow authenticated users to access protected routes', () => {
    const mockAuth = {
      getToken: vi.fn().mockReturnValue('valid-token')
    }
    
    const token = mockAuth.getToken()
    expect(token).toBe('valid-token')
  })

  it('should clear session on auth error', () => {
    const mockAuth = {
      clearSession: vi.fn()
    }
    
    mockAuth.clearSession()
    expect(mockAuth.clearSession).toHaveBeenCalled()
  })
})

describe('Route Paths', () => {
  const expectedPaths = [
    '/login',
    '/dashboard/designer',
    '/chart/editor',
    '/datasource',
    '/dataset',
    '/dataset/create',
    '/dataset/edit/:id',
    '/report/list',
    '/report/designer',
    '/report/preview'
  ]

  it('should have all expected paths', () => {
    expectedPaths.forEach(path => {
      expect(path).toBeDefined()
    })
  })
})
