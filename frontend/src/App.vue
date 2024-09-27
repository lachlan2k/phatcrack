<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useToast } from 'vue-toastification'
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { RouterView, useRouter } from 'vue-router'
import { AxiosError } from 'axios'

import { checkCors } from './api'

import { useAuthStore } from '@/stores/auth'
import { useConfigStore } from '@/stores/config'
import { useProjectsStore } from '@/stores/projects'
import { useListfilesStore } from '@/stores/listfiles'
import { useAgentsStore } from '@/stores/agents'
import { useUsersStore } from '@/stores/users'
import { useActiveAttacksStore } from '@/stores/activeAttacks'
import { useAttackTemplatesStore } from '@/stores/attackTemplates'

const AUTH_REFRESH_RATE = 0.5 * 60 * 1000 // Every 5 minutes
const authRefreshInterval = ref(0)
const activeAttacksRefreshInterval = ref(0)

const router = useRouter()
const authStore = useAuthStore()
const configStore = useConfigStore()
const projectStore = useProjectsStore()
const listfileStore = useListfilesStore()
const agentsStore = useAgentsStore()
const usersStore = useUsersStore()
const activeAttacksStore = useActiveAttacksStore()
const attackTemplatesStore = useAttackTemplatesStore()

const { hasCompletedAuth, hasTriedAuth, loggedInUser, isLoggedIn } = storeToRefs(authStore)

const toast = useToast()

onMounted(async () => {
  await router.isReady()

  try {
    await checkCors()
  } catch (e) {
    // check if axios error
    if (e instanceof AxiosError) {
      if (e.response?.data?.message.toLowerCase() === 'origin not allowed') {
        toast.error(
          'The request origin is not allowed. Phatcrack will not work.\n\nPlease ensure your administrator sets BASE_URL correctly. Alternatively, consider setting INSECURE_ORIGIN=1 if this is a local development instance.',
          {
            // force user to dismiss this
            timeout: false,
            closeOnClick: false,
            draggable: false
          }
        )
      }
    }
  }

  router.beforeEach(() => {
    authStore.refreshAuth()
    configStore.load()
  })

  authStore.refreshAuth()
  configStore.load()

  authRefreshInterval.value = setInterval(() => {
    authStore.refreshAuth()
    configStore.load()
  }, AUTH_REFRESH_RATE)

  activeAttacksRefreshInterval.value = setInterval(() => {
    if (isLoggedIn.value) {
      activeAttacksStore.load()
    }
  }, 15 * 1000)
})

onBeforeUnmount(() => {
  clearInterval(authRefreshInterval.value)
  clearInterval(activeAttacksRefreshInterval.value)
})

watch(hasCompletedAuth, (newHasCompletedAuth, prevHasCompletedAuth) => {
  if (newHasCompletedAuth && !prevHasCompletedAuth) {
    toast.success(`Welcome ${loggedInUser.value?.username}`)
    listfileStore.load(true)
    projectStore.load(true)
    agentsStore.load(true)
    usersStore.load(true)
    attackTemplatesStore.load(true)
    activeAttacksStore.load()
  }
})

watch(isLoggedIn, (newIsLoggedIn, oldIsLoggedIn) => {
  // We've been logged out
  if (oldIsLoggedIn && !newIsLoggedIn) {
    router.push('/')
  }
})

// When we try to refresh auth for the first time (after page load)
// check to see if we're logged in our not
// if we're not logged in, go to login page
watch(hasTriedAuth, (newHasTriedAuth, oldHasTriedAuth) => {
  if (newHasTriedAuth && !oldHasTriedAuth && router.currentRoute.value.path != '/oidc-callback') {
    if (!isLoggedIn.value) {
      router.push('/')
    }
  }
})
</script>

<template>
  <RouterView />
</template>
