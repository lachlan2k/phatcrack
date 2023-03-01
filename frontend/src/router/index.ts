import { createRouter, createWebHistory } from 'vue-router'
import DefaultLayout from '@/layouts/default.vue'

function withDefaultLayout(component: { (): Promise<any>; (): Promise<any> }) {
  return {
    component: DefaultLayout,
    children: [{ path: '', component }]
  }
}

function route(path, name, component) {
  return {
    path, name, ...withDefaultLayout(component)
  }
}

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/pages/Login.vue')
    },
    {
      path: '/',
      name: 'login-home',
      component: () => import('@/pages/Login.vue')
    },
    route('/agents', 'Agents', () => import('@/pages/Agents.vue')),
    route('/attacks', 'Attacks', () => import('@/pages/Attacks.vue')),
    route('/dashboard', 'Dashboard', () => import('@/pages/Dashboard.vue')),
    route('/potfile', 'Potfile', () => import('@/pages/Potfile.vue')),
    route('/projects', 'Projects', () => import('@/pages/projects/index.vue')),
    route('/wizard', 'Wizard', () => import('@/pages/Wizard.vue'))
  ]
})

export default router
