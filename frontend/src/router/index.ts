import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/login', name: 'login', component: () => import('../pages/LoginPage.vue'), meta: { public: true } },
    {
      path: '/',
      component: () => import('../layouts/MainLayout.vue'),
      children: [
        { path: '', redirect: '/dashboard' },
        { path: 'dashboard', name: 'dashboard', component: () => import('../pages/DashboardPage.vue'), meta: { title: '仓库总览' } },
        { path: 'inbound', name: 'inbound', component: () => import('../pages/InboundPage.vue'), meta: { title: '入库管理' } },
        { path: 'outbound', name: 'outbound', component: () => import('../pages/OutboundPage.vue'), meta: { title: '出库管理' } },
        { path: 'bin-locations', name: 'bin-locations', component: () => import('../pages/BinLocationsPage.vue'), meta: { title: '库位管理' } },
        { path: 'products', name: 'products', component: () => import('../pages/ProductsPage.vue'), meta: { title: '商品管理' } },
        { path: 'owners', name: 'owners', component: () => import('../pages/OwnersPage.vue'), meta: { title: '货主管理' } },
        { path: 'owners/:id', name: 'owner-detail', component: () => import('../pages/OwnerDetailPage.vue'), meta: { title: '货主详情' } },
        { path: 'audit', name: 'audit', component: () => import('../pages/AuditPage.vue'), meta: { title: '操作日志' } },
      ],
    },
    { path: '/:pathMatch(.*)*', name: 'not-found', component: () => import('../pages/NotFoundPage.vue'), meta: { public: true } },
  ],
})

router.beforeEach((to) => {
  const auth = useAuthStore()
  if (!to.meta.public && !auth.isLoggedIn) {
    return { path: '/login', query: { redirect: to.fullPath } }
  }
  if (to.path === '/login' && auth.isLoggedIn) {
    return { path: '/dashboard' }
  }
  return true
})

export default router
