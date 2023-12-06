<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { storeToRefs } from 'pinia'
import { finishMFAChallenge, startMFAEnrollment, startMFAChallenge, finishMFAEnrollment, changeTemporaryPassword } from '@/api/auth'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import { useToast } from 'vue-toastification'
import { useToastError } from '@/composables/useToastError'

const toast = useToast()
const authStore = useAuthStore()
const router = useRouter()

const { hasCompletedAuth, isAwaitingMFA, requiresPasswordChange, requiresMFAEnrollment, loginError, isLoginLoading, loggedInUser } =
  storeToRefs(authStore)

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

if (hasCompletedAuth.value) {
  router.push('/dashboard')
}

watch(hasCompletedAuth, (hasCompletedAuth) => {
  if (hasCompletedAuth) {
    router.push('/dashboard')
  }
})

const username = ref('')
const password = ref('')
const newPassword = ref('')

async function doLogin(event: Event) {
  if (event) {
    event.preventDefault()
  }

  authStore.login(username.value, password.value)
}

const { catcher } = useToastError()
const isPasswordChangeLoading = ref(false)

async function doPasswordChange(event: Event) {
  if (event) {
    event.preventDefault()
  }

  try {
    isPasswordChangeLoading.value = true
    const res = await changeTemporaryPassword({
      old_password: password.value,
      new_password: newPassword.value
    })
    if (res === 'Ok') {
      toast.success('Password changed successfully')
    } else {
      toast.warning('Unexpected API response: ' + res)
    }

    authStore.refreshAuth()
  } catch (e: any) {
    catcher(e, 'Failed to change temporary password. ')
  } finally {
    isPasswordChangeLoading.value = false
  }
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
            <input type="text" placeholder="john.doe" class="input input-bordered" v-model="username" />
          </div>
          <div class="form-control">
            <label class="label">
              <span class="label-text">Password</span>
            </label>
            <input type="password" placeholder="hunter2" class="input input-bordered" v-model="password" />
          </div>
          <div v-if="loginError != null" class="mt-4 text-center text-red-500">
            <p>{{ loginError }}</p>
          </div>
          <div class="form-control mt-6">
            <button type="submit" class="btn btn-primary" :disabled="isLoginLoading">
              <span class="loading loading-spinner loading-md" v-if="isLoginLoading"></span>
              Login
            </button>
          </div>
        </form>

        <div v-if="activeScreen == ActiveScreens.MFAVerification" class="text-center">
          <p>We need to verify your identity</p>
          <div class="cursor-pointer" @click="verifyKey">
            <font-awesome-icon icon="fa-solid fa-key" class="my-8" style="font-size: 5rem" />
          </div>
          <div>
            <button class="btn btn-ghost" @click="verifyKey">Verify</button>
          </div>
        </div>

        <div v-if="activeScreen == ActiveScreens.MFAEnrollment" class="text-center">
          <p>You are required to enroll a security key</p>
          <div class="cursor-pointer" @click="enrollKey">
            <font-awesome-icon icon="fa-solid fa-key" class="my-8" style="font-size: 5rem" />
          </div>
          <div>
            <button class="btn btn-primary" @click="enrollKey">Enroll Key</button>
          </div>
        </div>

        <div v-if="activeScreen == ActiveScreens.PasswordChange">
          <p class="text-center">You are required to change your password</p>
          <form @submit="doPasswordChange">
            <div class="form-control">
              <label class="label">
                <span class="label-text">Old Password</span>
              </label>
              <input type="password" placeholder="hunter2" class="input input-bordered" v-model="password" />
            </div>
            <div class="form-control">
              <label class="label">
                <span class="label-text">New Password</span>
              </label>
              <input type="password" placeholder="hunter2" class="input input-bordered" v-model="newPassword" />
            </div>
            <div v-if="loginError != null" class="mt-4 text-center text-red-500">
              <p>{{ loginError }}</p>
            </div>
            <div class="form-control mt-6">
              <button type="submit" class="btn btn-primary" :disabled="isPasswordChangeLoading">
                <span class="loading loading-spinner loading-md" v-if="isLoginLoading"></span>
                Change Password
              </button>
            </div>
          </form>
        </div>

        <div v-if="activeScreen == ActiveScreens.Done" class="text-center">
          <p>Welcome</p>
          <font-awesome-icon icon="fa-solid fa-check" class="my-8 text-success" style="font-size: 5rem" />
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
