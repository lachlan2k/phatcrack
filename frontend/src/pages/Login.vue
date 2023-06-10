<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { storeToRefs } from 'pinia'
import { finishMFAChallenge, startMFAEnrollment } from '@/api/auth'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import { finishMFAEnrollment } from '@/api/auth'
import { startMFAChallenge } from '@/api/auth'

const authStore = useAuthStore()
const router = useRouter()

const {
  hasCompletedAuth,
  isAwaitingMFA,
  requiresPasswordChange,
  requiresMFAEnrollment,
  loginError,
  isLoginLoading,
  loggedInUser
} = storeToRefs(authStore)

enum ActiveScreens {
  Credentials,
  PasswordChange,
  MFAEnrollment,
  MFAVerification,
  Done
}

const activeScreen = computed(() => {
  if (loggedInUser.value == null) {
    return ActiveScreens.Credentials
  }

  if (requiresPasswordChange.value) {
    return ActiveScreens.PasswordChange
  }

  if (requiresMFAEnrollment.value) {
    return ActiveScreens.MFAEnrollment
  }

  if (isAwaitingMFA.value) {
    return ActiveScreens.MFAVerification
  }

  return ActiveScreens.Done
})

watch(hasCompletedAuth, (hasCompletedAuth) => {
  if (hasCompletedAuth) {
    router.push('/dashboard')
  }
})

const username = ref('')
const password = ref('')

async function doLogin(event: Event) {
  if (event) {
    event.preventDefault()
  }

  authStore.login(username.value, password.value)
}

function urlSafeB64Decode(value: string) {
  return Uint8Array.from(atob(value.replace(/_/g, '/').replace(/-/g, '+')), (c) => c.charCodeAt(0))
}

async function enrollKey() {
  const response = await startMFAEnrollment()
  const challenge = {
    ...response,
    publicKey: {
      ...response.publicKey,
      challenge: urlSafeB64Decode(response.publicKey.challenge as unknown as string), // type codegen is wrong, its a base64 encoded string once marshalled, not a []byte
      user: {
        ...response.publicKey.user,
        id: urlSafeB64Decode(response.publicKey.user.id as string)
      },
      excludeCredentials: response.publicKey.excludeCredentials?.map((cred) => ({
        ...cred,
        id: urlSafeB64Decode(cred.id as unknown as string)
      })),
      attestation: 'none'
    } as PublicKeyCredentialCreationOptions
  }

  const newCred = await navigator.credentials.create({
    publicKey: challenge.publicKey
  })

  await finishMFAEnrollment(newCred as PublicKeyCredential)
  await authStore.refreshAuth()
}

async function verifyKey() {
  const response = await startMFAChallenge()
  const challenge = {
    ...response,
    publicKey: {
      ...response.publicKey,
      challenge: urlSafeB64Decode(response.publicKey.challenge as unknown as string), // type codegen is wrong, its a base64 encoded string once marshalled, not a []byte
      allowCredentials: response.publicKey.allowCredentials?.map((cred) => ({
        ...cred,
        id: urlSafeB64Decode(cred.id as unknown as string)
      }))
    } as PublicKeyCredentialRequestOptions
  }

  const assertion = await navigator.credentials.get({
    publicKey: challenge.publicKey
  })

  await finishMFAChallenge(assertion as PublicKeyCredential)
  await authStore.refreshAuth()
}

watch(activeScreen, (newActiveScreen) => {
  if (newActiveScreen == ActiveScreens.MFAVerification) {
    verifyKey()
  }
})

const cardTitle = computed(() => {
  switch (activeScreen.value) {
    case ActiveScreens.Credentials:
      return 'Login to Phatcrack'

    case ActiveScreens.PasswordChange:
      return 'Set a new password'

    case ActiveScreens.MFAEnrollment:
      return 'Plug in your security key'

    case ActiveScreens.MFAVerification:
      return 'Plug in your security key'

    case ActiveScreens.Done:
      return 'You have successfully logged in!'

    default:
      return ''
  }
})
</script>

<template>
  <main class="z-10 flex min-h-screen items-center justify-center self-center bg-neutral">
    <div class="card w-96 bg-base-100 shadow-xl">
      <div class="card-body">
        <div class="card-title justify-center">
          <h2>{{ cardTitle }}</h2>
        </div>

        <form @submit="doLogin" v-if="activeScreen == ActiveScreens.Credentials">
          <div class="form-control">
            <label class="label">
              <span class="label-text">Username</span>
            </label>
            <input
              type="text"
              placeholder="john.doe"
              class="input-bordered input"
              v-model="username"
            />
          </div>
          <div class="form-control">
            <label class="label">
              <span class="label-text">Password</span>
            </label>
            <input
              type="password"
              placeholder="hunter2"
              class="input-bordered input"
              v-model="password"
            />
          </div>
          <div v-if="loginError != null" class="mt-4 text-center text-red-500">
            <p>{{ loginError }}</p>
          </div>
          <div class="form-control mt-6">
            <button type="submit" class="btn-primary btn" :disabled="isLoginLoading">Login</button>
          </div>
        </form>

        <div v-if="activeScreen == ActiveScreens.MFAVerification" class="text-center">
          <p>We need to verify your identity</p>
          <div class="cursor-pointer" @click="verifyKey">
            <font-awesome-icon icon="fa-solid fa-key" class="my-8" style="font-size: 5rem" />
          </div>
          <div>
            <button class="btn-ghost btn" @click="verifyKey">Verify</button>
          </div>
        </div>

        <div v-if="activeScreen == ActiveScreens.MFAEnrollment">
          <p>You are required to enroll a security key</p>
          <div class="cursor-pointer" @click="enrollKey">
            <font-awesome-icon icon="fa-solid fa-key" class="my-8" style="font-size: 5rem" />
          </div>
          <div>
            <button class="btn-primary btn" @click="enrollKey">Enroll Key</button>
          </div>
        </div>
      </div>
    </div>
  </main>
</template>

<style scoped>
main {
  /* font-size: 1.25rem; */
}
</style>
