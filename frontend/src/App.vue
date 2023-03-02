<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { RouterView, useRouter } from 'vue-router'
import { useAuthStore } from './stores/auth'

const AUTH_REFRESH_RATE = 0.5 * 60 * 1000 // Every 5 minutes
const authRefreshTimeout = ref(0)

const router = useRouter()

const authStore = useAuthStore()
const { isLoggedIn, hasTriedAuth } = storeToRefs(authStore)

onMounted(() => {
  authStore.refreshAuth()

  authRefreshTimeout.value = setInterval(() => {
    authStore.refreshAuth()
  }, AUTH_REFRESH_RATE)
})

onBeforeUnmount(() => {
  clearInterval(authRefreshTimeout.value)
})

watch(isLoggedIn, (newIsLoggedIn, prevIsLoggedIn) => {
  // We don't want to issue a redirect if we haven't even tried auth yet
  if (!hasTriedAuth.value) {
    console.log('remove me if you never see me')
    return
  }

  if (newIsLoggedIn && !prevIsLoggedIn && router.currentRoute.value.path === '/login') {
    console.log('Redirecting to dashboard')
    router.push('/dashboard')
  }

  if (!newIsLoggedIn) {
    console.log('Redirecting to login', router.currentRoute.value)
    router.push('/login')
  }
})
</script>

<template>
  <RouterView />
</template>
