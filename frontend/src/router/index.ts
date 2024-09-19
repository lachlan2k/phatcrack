import { createRouter, createWebHistory } from 'vue-router'

import DefaultLayout from '@/layouts/default.vue'

function withDefaultLayout(component: () => any, name: string) {
  return {
    component: DefaultLayout,
    children: [{ path: '', name: `${name} Layout`, component }]
  }
}

function route(path: string, name: string, component: () => any) {
  return {
    path,
    name,
    ...withDefaultLayout(component, name)
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
    {
      path: '/oidc-callback',
      name: 'oidc-callback',
      component: () => import('@/pages/LoginOIDCCallback.vue')
    },
    route('/agents', 'Agents', () => import('@/pages/Agents.vue')),
    route('/dashboard', 'Dashboard', () => import('@/pages/projects/index.vue')),
    route('/hash-search', 'Hash Search', () => import('@/pages/HashSearch.vue')),

    route('/project/:id', 'Project', () => import('@/pages/projects/project.vue')),

    route('/hashlist/:id', 'Hashlist', () => import('@/pages/Hashlist.vue')),

    route('/listfiles', 'Listfiles', () => import('@/pages/Listfiles.vue')),

    route('/attack-templates', 'Attack Templates', () => import('@/pages/AttackTemplates.vue')),

    route('/wizard', 'Wizard', () => import('@/pages/Wizard.vue')),

    route('/account', 'Account', () => import('@/pages/Account.vue')),

    route('/utilisation', 'Utilisation', () => import('@/pages/Utilisation.vue')),

    route('/admin/general', 'General Settings', () => import('@/pages/admin/Configuration.vue')),
    route('/admin/users', 'User Management', () => import('@/pages/admin/Users.vue')),
    route('/admin/agents', 'Agent Management', () => import('@/pages/admin/Agents.vue'))
  ]
})

export default router
