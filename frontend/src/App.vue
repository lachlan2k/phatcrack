<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useToast } from 'vue-toastification'
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { RouterView, useRouter } from 'vue-router'
import { useAuthStore } from './stores/auth'

const AUTH_REFRESH_RATE = 0.5 * 60 * 1000 // Every 5 minutes
const authRefreshTimeout = ref(0)

const authStore = useAuthStore()

const { isLoggedIn, hasCompletedAuth, loggedInUser } = storeToRefs(authStore)

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
</script>

<template>
  <RouterView />
</template>
