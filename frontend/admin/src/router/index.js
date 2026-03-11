import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue')
  },
  {
    path: '/',
    component: () => import('../views/Layout.vue'),
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: () => import('../views/Dashboard.vue')
      },
      {
        path: 'cards',
        name: 'CardManage',
        component: () => import('../views/CardManage.vue')
      },
      {
        path: 'recharges',
        name: 'RechargeRecord',
        component: () => import('../views/RechargeRecord.vue')
      },
      {
        path: 'config',
        name: 'SystemConfig',
        component: () => import('../views/SystemConfig.vue')
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('admin_token')
  // 如果在登录页，清除可能存在的无效 token
  if (to.path === '/login') {
    if (token) {
      localStorage.removeItem('admin_token')
      localStorage.removeItem('admin_user')
    }
    next()
  } else if (!token) {
    next('/login')
  } else {
    next()
  }
})

export default router
