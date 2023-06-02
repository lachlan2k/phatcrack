import { createRouter, createWebHistory } from 'vue-router'
import DefaultLayout from '@/layouts/default.vue'

function withDefaultLayout(component: () => any) {
  return {
    component: DefaultLayout,
    children: [{ path: '', component }]
  }
}

function route(path: string, name: string, component: () => any) {
  return {
    path,
    name,
    ...withDefaultLayout(component)
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
    route('/jobs', 'All Running Jobs', () => import('@/pages/Attacks.vue')),
    route('/dashboard', 'Dashboard', () => import('@/pages/Dashboard.vue')),
    route('/potfile', 'Potfile', () => import('@/pages/Potfile.vue')),

    route('/projects', 'All Projects', () => import('@/pages/projects/index.vue')),
    route('/project/:id', 'Project', () => import('@/pages/projects/project.vue')),

    route('/hashlist/:id', 'Hashlist', () => import('@/pages/Hashlist.vue')),

    route('/wizard', 'Wizard', () => import('@/pages/Wizard.vue'))
  ]
})

export default router
