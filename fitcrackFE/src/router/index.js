/*
   * Author : see AUTHORS
   * Licence: MIT, see LICENSE
*/

import Vue from 'vue'
import Router from 'vue-router'

import store from '@/store'

const home = () => import('@/components/dashboard/dashboard.vue')
const jobs = () => import('@/components/job/jobsView.vue')
const jobDetail = () => import('@/components/jobDetail/jobDetailView.vue')
const batchList = () => import('@/components/batch/batchList.vue')
const batchDetail = () => import('@/components/batch/batchDetail.vue')
const login = () => import('@/components/user/loginView.vue')
const mainPage = () => import('@/components/mainPage.vue')
const addJob = () => import('@/components/job/addJobView.vue')
const templates = () => import('@/components/jobTemplate/templateList.vue')
const hosts = () => import('@/components/host/hostsView.vue')
const hashes = () => import('@/components/hashcache/hashCacheView.vue')
const dictionaries = () => import('@/components/dictionary/dictionariesView.vue')
const dictionaryDetail = () => import('@/components/dictionary/dictionaryDetailView.vue')
const pcfg = () => import('@/components/pcfg/pcfgView.vue')
const pcfgBrowser = () => import('@/components/pcfg/pcfgBrowser.vue')
const manageUsers = () => import('@/components/user/manageUsersView.vue')
const myAccount = () => import('@/components/myAccount/myAccountView.vue')
const hostDetail = () => import('@/components/host/hostDetailView.vue')
const rules = () => import('@/components/rule/rulesView.vue')
const ruleDetail = () => import('@/components/rule/ruleDetailView.vue')
const masks = () => import('@/components/mask/masksView.vue')
const charsets = () => import('@/components/charset/charsetsView.vue')
const charsetDetail = () => import('@/components/charset/charsetDetailView.vue')
const markovChains = () => import('@/components/markovChains/markovView.vue')
const maskDetail = () => import('@/components/mask/maskDetailView.vue')
const testView = () => import('@/components/test/test.vue')
const PageNotFound = () => import('@/components/pageNotFound.vue')
const EncryptedFiles = () => import('@/components/encryptedFile/encryptedFilesView.vue')
const Server = () => import('@/components/server/serverMonitor.vue')
const Settings = () => import('@/components/settings/settingsView.vue')
const Transfer = () => import('@/components/settings/dataTransfer.vue')
const UnauthorizedError = () => import('@/components/errorPages/unauthorized.vue')

Vue.use(Router);

const appRoutes = [
  {
    path: '/',
    name: 'home',
    component: home,
    meta: {
      title: 'Dashboard',
      icon: 'mdi-view-dashboard'
    }
  },
  {
    path: '/test',
    name: 'test',
    component: testView
  },
  {
    path: '/batches',
    name: 'batches',
    component: batchList,
    meta: {
      title: 'Batches',
      icon: 'mdi-tray-full',
      navtab: 0
    }
  },
  {
    path: '/batches/:id',
    name: 'batch',
    component: batchDetail,
    meta: {
      title: 'Batch Detail'
    }
  },
  {
    path: '/jobs',
    name: 'jobs',
    component: jobs,
    meta: {
      title: 'Jobs',
      icon: 'mdi-briefcase',
      navtab: 0
    }
  },
  {
    path: '/bins/:id',
    name: 'bins',
    component: jobs,
    meta: {
      title: 'Bins',
      icon: 'mdi-folder',
      navtab: 0
    }
  },
  {
    path: '/bins/trash',
    name: 'trash',
    component: jobs,
    meta: {
      title: 'Trash',
      icon: 'mdi-delete',
      navtab: 0
    }
  },
  {
    path: '/jobs/add',
    name: 'addJob',
    component: addJob,
    meta: {
      guard: 'ADD_NEW_JOB',
      title: 'Create a job',
      icon: 'mdi-briefcase-plus',
      navtab: 0
    }
  },
  {
    path: '/jobs/:id',
    name: 'jobDetail',
    component: jobDetail,
    meta: {
      title: 'Job detail'
    }
  },
  {
    path: '/templates',
    name: 'templates',
    component: templates,
    meta: {
      title: 'Job templates',
      icon: 'mdi-file',
      navtab: 2
    }
  },
  {
    path: '/hosts',
    name: 'hosts',
    component: hosts,
    meta: {
      title: 'Hosts',
      icon: 'mdi-desktop-classic',
      navtab: 2
    }
  },
  {
    path: '/hosts/hidden',
    name: 'hiddenHosts',
    component: hosts
  },
  {
    path: '/hosts/:id',
    name: 'hostDetail',
    component: hostDetail,
    meta: {
      title: 'Host info'
    }
  },
  {
    path: '/hashes',
    name: 'hashes',
    component: hashes,
    meta: {
      title: 'Hashes',
      icon: 'mdi-book-lock',
      navtab: 1
    }
  },
  {
    path: '/dictionaries',
    name: 'dictionaries',
    component: dictionaries,
    meta: {
      title: 'Dictionaries',
      icon: 'mdi-dictionary',
      navtab: 1
    }
  },
  {
    path: '/dictionaries/:id',
    name: 'dictionaryDetail',
    component: dictionaryDetail
  },
  {
    path: '/pcfg',
    name: 'pcfg',
    component: pcfg,
    meta: {
      title: 'PCFG',
      icon: 'mdi-ray-start-end',
      navtab: 1
    }
  },
  {
    path: '/pcfg/:id',
    name: 'pcfgBrowser',
    component: pcfgBrowser,
    meta: {
      title: 'Inspect PCFG'
    }
  },
  {
    path: '/rules',
    name: 'rules',
    component: rules,
    meta: {
      title: 'Rules',
      icon: 'mdi-gavel',
      navtab: 1
    }
  },
  {
    path: '/rules/:id',
    name: 'ruleDetail',
    component: ruleDetail
  },
  {
    path: '/charsets',
    name: 'charsets',
    component: charsets,
    meta: {
      title: 'Charsets',
      icon: 'mdi-alphabetical',
      navtab: 1
    }
  },
  {
    path: '/charsets/:id',
    name: 'charsetDetail',
    component: charsetDetail
  },
  {
    path: '/masks',
    name: 'masks',
    component: masks,
    meta: {
      title: 'Masks',
      icon: 'mdi-guy-fawkes-mask',
      navtab: 1
    }
  },
  {
    path: '/masks/:id',
    name: 'maskDetail',
    component: maskDetail
  },
  {
    path: '/markovChains',
    name: 'markovChains',
    component: markovChains,
    meta: {
      title: 'Markov chains',
      icon: 'mdi-matrix',
      navtab: 1
    }
  },
  {
    path: '/files',
    name: 'files',
    component: EncryptedFiles,
    meta: {
      title: 'Encrypted Files',
      icon: 'mdi-file-lock',
      navtab: 1
    }
  },
  {
    path: '/user/manageUsers',
    name: 'manageUsers',
    component: manageUsers,
    meta: {
      guard: 'MANAGE_USERS',
      title: 'Manage users',
      icon: 'mdi-folder-account',
      navtab: 2
    }
  },
  {
    path: '/myAccount',
    name: 'myAccount',
    component: myAccount,
    meta: {
      title: 'My account',
      icon: 'mdi-account'
    }
  },
  {
    path: '/server',
    name: 'server',
    component: Server,
    meta: {
      title: 'Server monitor',
      icon: 'mdi-gauge',
      navtab: 2
    }
  },
  {
    path: '/settings',
    name: 'settings',
    component: Settings,
    meta: {
      title: 'Settings',
      icon: 'mdi-cogs',
      navtab: 2
    }
  },
  {
    path: '/transfer',
    name: 'transfer',
    component: Transfer,
    meta: {
      title: 'Data Transfer',
      icon: 'mdi-dolly',
      navtab: 2
    }
  },
  {
    path: '/error/unauthorized',
    name: 'unauthorized',
    component: UnauthorizedError,
    meta: {
      title: 'Access Denied',
      icon: 'mdi-alert'
    }
  },
  {
    path: "*",
    component: PageNotFound
  }
]

const rt = new Router({
  mode: 'history',
  routes: [
    {
      path: '/login',
      name: 'login',
      component: login
    },
    {
      path: '',
      component: mainPage,
      children: appRoutes
    }
  ]
})

rt.beforeResolve(async (to, from, next) => {
  const perm = to.meta.guard
  if (!perm) {
    next()
    return
  }
  //
  const hasAccess = await store.dispatch('hasPermission', perm)
  if (hasAccess) {
    next()
  } else {
    next({
      name: 'unauthorized',
      query: {
        accessing: to.path,
        missing: perm
      }
    })
  }
})

rt.afterEach((to, from) => {
  const project = rt.app.$store.state.project || 'Fitcrack'
  if (to.meta.title) {
    document.title = `${to.meta.title} – ${project}`
  } else {
    document.title = project
  }
})

export default rt

export function routeIcon(name) {
  const route = appRoutes.find(r => r.name === name)
  if (route.meta && route.meta.icon) {
    return route.meta.icon
  } else {
    return null
  }
}