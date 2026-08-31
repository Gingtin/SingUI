import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import MainLayout from '@/layouts/MainLayout.vue'
import Login from '@/views/Login.vue'
import Dashboard from '@/views/Dashboard.vue'
import Inbounds from '@/views/Inbounds.vue'
import Outbounds from '@/views/Outbounds.vue'
import Routing from '@/views/Routing.vue'
import Subscriptions from '@/views/Subscriptions.vue'
import Logs from '@/views/Logs.vue'
import Settings from '@/views/Settings.vue'
import ApiDocs from '@/views/ApiDocs.vue'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { public: true },
  },
  {
    path: '/',
    component: MainLayout,
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: Dashboard,
      },
      {
        path: 'inbounds',
        name: 'Inbounds',
        component: Inbounds,
      },
      {
        path: 'outbounds',
        name: 'Outbounds',
        component: Outbounds,
      },
      {
        path: 'routing',
        name: 'Routing',
        component: Routing,
      },
      {
        path: 'subscriptions',
        name: 'Subscriptions',
        component: Subscriptions,
      },
      {
        path: 'logs',
        name: 'Logs',
        component: Logs,
      },
      {
        path: 'settings',
        name: 'Settings',
        component: Settings,
      },
      {
        path: 'docs',
        name: 'ApiDocs',
        component: ApiDocs,
      },
    ],
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/dashboard',
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem('token')
  if (!to.meta.public && !token) {
    next('/login')
  } else if (to.path === '/login' && token) {
    next('/dashboard')
  } else {
    next()
  }
})

export default router
