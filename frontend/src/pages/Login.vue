<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

const username = ref('')
const password = ref('')

function doLogin () {
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
          <button class="btn-primary btn" @click="doLogin">Login</button>
        </div>
      </div>
    </div>
  </main>
</template>
