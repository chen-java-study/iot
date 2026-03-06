import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'Query',
    component: () => import('../views/Query.vue')
  },
  {
    path: '/detail',
    name: 'Detail',
    component: () => import('../views/Detail.vue')
  },
  {
    path: '/recharge',
    name: 'Recharge',
    component: () => import('../views/Recharge.vue')
  }
]

const router = createRouter({
  // H5 部署在 /h5/ 子路径下，需要设置 base
  history: createWebHistory('/h5/'),
  routes
})

export default router
