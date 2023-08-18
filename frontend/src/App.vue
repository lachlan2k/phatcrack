<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useToast } from 'vue-toastification'
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { RouterView, useRouter } from 'vue-router'
import { useAuthStore } from './stores/auth'

const AUTH_REFRESH_RATE = 0.5 * 60 * 1000 // Every 5 minutes
const authRefreshTimeout = ref(0)

const router = useRouter()
const authStore = useAuthStore()

const { hasCompletedAuth, hasTriedAuth, loggedInUser, isLoggedIn } = storeToRefs(authStore)

onMounted(() => {
  authStore.refreshAuth()

  authRefreshTimeout.value = setInterval(() => {
    authStore.refreshAuth()
  }, AUTH_REFRESH_RATE)
})

onBeforeUnmount(() => {
  clearInterval(authRefreshTimeout.value)
})

const toast = useToast()

watch(hasCompletedAuth, (newHasCompletedAuth, prevHasCompletedAuth) => {
  if (newHasCompletedAuth && !prevHasCompletedAuth) {
    toast.success(`Welcome ${loggedInUser.value?.username}`)
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
  if (newHasTriedAuth && !oldHasTriedAuth) {
    if (!isLoggedIn.value) {
      router.push('/')
    }
  }
})
</script>

<template>
  <RouterView />
</template>
