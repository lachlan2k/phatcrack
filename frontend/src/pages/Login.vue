<script setup lang="ts">
import { ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'

const authStore = useAuthStore()
const router = useRouter()

const { isLoggedIn } = storeToRefs(authStore)

watch(isLoggedIn, newIsLoggedIn => {
  if (newIsLoggedIn) {
    router.push('/dashboard')
  }
})

const username = ref('')
const password = ref('')

function doLogin (event) {
  if (event) {
    event.preventDefault()
  }
  authStore.login(username.value, password.value)
}
</script>

<template>
  <main class="z-10 flex min-h-screen items-center justify-center self-center bg-neutral">
    <div class="card w-96 bg-base-100 shadow-xl">
      <div class="card-body">
        <div class="card-title justify-center">
          <h2>Login to Phatcrack</h2>
        </div>

        <form @submit="doLogin">
          <div class="form-control">
            <label class="label">
              <span class="label-text">Username</span>
            </label>
            <input type="text" placeholder="john.doe" class="input-bordered input" v-model="username" />
          </div>
          <div class="form-control">
            <label class="label">
              <span class="label-text">Password</span>
            </label>
            <input type="password" placeholder="hunter2" class="input-bordered input" v-model="password" />
          </div>
          <div class="form-control mt-6">
            <button type="submit" class="btn-primary btn" @click="doLogin">Login</button>
          </div>
        </form>
      </div>
    </div>
  </main>
</template>
