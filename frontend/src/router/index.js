import { createRouter, createWebHistory } from 'vue-router'
import store from '@/store'  
import Login from '@/views/Login.vue'
import { ElMessage, ElNotification } from 'element-plus'
import 'element-plus/es/components/notification/style/css'

const routes = [
  {
    path: '/',
    name: 'Login',
    component: Login
  },
  {
    path: '/user/create-report',
    name: 'CreateReport',
    component: () => import('../views/user/CreateReport.vue'),
    meta: { requiresAuth: true, role: 'user' }
  },
  {
    path: '/user/report-status',
    name: 'ReportStatus',
    component: () => import('../views/user/ReportStatus.vue'),
    meta: { requiresAuth: true, role: 'user' }
  },
  {
    path: '/user/manage-reports',
    name: 'ManageReports',
    component: () => import('../views/user/ManageReports.vue'),
    meta: { requiresAuth: true, role: 'user' }
  },
  {
    path: '/user/profile',
    name: 'Profile',
    component: () => import('../views/user/Profile.vue'),
    meta: { requiresAuth: true, role: 'user' }
  },
  {
    path: '/admin/approve-report',
    name: 'ApproveReport',
    component: () => import('../views/admin/ApproveReport.vue'),
    meta: { requiresAuth: true, role: 'admin' }
  },
  {
    path: '/admin/view-reports',
    name: 'ViewReports',
    component: () => import('../views/admin/ViewReports.vue'),
    meta: { requiresAuth: true, role: 'admin' }
  },
  {
    path: '/admin/manage-users',
    name: 'ManageUsers',
    component: () => import('../views/admin/ManageUsers.vue'),
    meta: { requiresAuth: true, role: 'admin' }
  },
  {
    path: '/admin/manage-templates',
    name: 'ManageTemplates',
    component: () => import('../views/admin/ManageTemplates.vue'),
    meta: { requiresAuth: true, role: 'admin' }
  },
  {
    path: '/admin/statistics',
    name: 'Statistics',
    component: () => import('../views/admin/Statistics.vue'),
    meta: { requiresAuth: true, role: 'admin' }
  },
  {
    path: '/admin/system',
    name: 'System',
    component: () => import('../views/admin/System.vue'),
    meta: { requiresAuth: true, role: 'admin' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const isLoggedIn = store.getters.isLoggedIn
  const userRole = store.getters.userRole
  const token = store.state.token
  
  console.log("🚦 路由守卫拦截:", {
    目标路由: to.path,
    需要认证: to.matched.some(record => record.meta.requiresAuth),
    需要角色: to.meta.role,
    当前登录状态: isLoggedIn,
    当前用户角色: userRole,
    当前token长度: token?.length,
    本地存储token: localStorage.getItem('token')
  })

  // 需要认证的页面
  if (to.matched.some(record => record.meta.requiresAuth)) {
    if (!isLoggedIn) {
      console.log("❌ 用户未登录，跳转到登录页")
      next('/')
    } else if (to.meta.role && to.meta.role !== userRole) {
      console.log(`❌ 权限不足，需要角色: ${to.meta.role}, 当前角色: ${userRole}`)
      ElMessage.error('权限不足')
      // 根据角色跳转到对应首页
      next(userRole === 'admin' ? '/admin/approve-report' : '/user/create-report')
    } else {
      console.log("✅ 权限验证通过，允许访问")
      next()
    }
  } else if (to.path === '/' && isLoggedIn) {
    // 已登录但访问登录页，跳转到首页
    console.log("📱 已登录，重定向到主页")
    next(userRole === 'admin' ? '/admin/approve-report' : '/user/create-report')
  } else {
    console.log("➡️  直接放行")
    next()
  }
})

export default router