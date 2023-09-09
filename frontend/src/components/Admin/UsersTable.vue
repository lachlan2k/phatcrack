<script setup lang="ts">
import Modal from '@/components/Modal.vue'
import ConfirmModal from '../ConfirmModal.vue'
import IconButton from '@/components/IconButton.vue'
import PaginationControls from '@/components/PaginationControls.vue'
import { adminCreateServiceAccount, adminCreateUser, adminDeleteUser, adminGetAllUsers } from '@/api/admin'
import { useApi } from '@/composables/useApi'
import { useToast } from 'vue-toastification'
import { ref, computed } from 'vue'
import { usePagination } from '@/composables/usePagination'
import { useToastError } from '@/composables/useToastError'
import { useAuthStore } from '@/stores/auth'
import { storeToRefs } from 'pinia'

const isUserCreateOpen = ref(false)
const isServiceAccountCreateOpen = ref(false)

const { data: allUsers, fetchData: fetchUsers } = useApi(adminGetAllUsers)

const usersToPaginate = computed(() => allUsers.value?.users ?? [])

const {
  next: nextPage,
  prev: prevPage,
  totalPages: totalPages,
  currentItems: paginatedUsers,
  activePage
} = usePagination(usersToPaginate, 10)

const possibleRoles = ['admin', 'standard']

const newUserUsername = ref('')
const newUserPassword = ref('')
const newUserRole = ref('standard')

const serviceAccountValidationError = computed(() => {
  if (newUserUsername.value.length < 3) {
    return 'Username too short'
  }
  return null
})

const newUserValidationError = computed(() => {
  if (newUserUsername.value.length < 3) {
    return 'Username too short'
  }

  if (newUserPassword.value.length < 16) {
    return 'Password too short'
  }

  return null
})

const toast = useToast()
const { catcher } = useToastError()

async function onCreateUser() {
  try {
    const res = await adminCreateUser({
      username: newUserUsername.value,
      password: newUserPassword.value,
      roles: [newUserRole.value]
    })

    toast.success('Created new user: ' + res.username)
  } catch (e: any) {
    catcher(e)
  } finally {
    fetchUsers()
  }
}

async function onCreateServiceAccount() {
  try {
    const res = await adminCreateServiceAccount({
      username: newUserUsername.value,
      roles: [newUserRole.value]
    })

    toast.info(`Created new service account ${res.username}.\n\n API Key (note this down, won't be displayed again):\n${res.api_key}`, {
      // force user to dismiss this
      timeout: false,
      closeOnClick: false,
      draggable: false
    })
  } catch (e: any) {
    catcher(e)
  } finally {
    fetchUsers()
  }
}

const authStore = useAuthStore()
const { loggedInUser } = storeToRefs(authStore)

async function onDeleteUser(id: string) {
  if (loggedInUser.value?.id === id) {
    toast.error("You can't delete your own user")
    return
  }

  try {
    await adminDeleteUser(id)
    toast.info('Deleted user')
  } catch (e: any) {
    catcher(e)
  } finally {
    fetchUsers()
  }
}
</script>

<template>
  <div class="flex flex-row justify-between">
    <Modal v-model:isOpen="isUserCreateOpen">
      <h3 class="text-lg font-bold">Create a new user</h3>

      <div class="form-control">
        <label class="label font-bold">
          <span class="label-text">Username</span>
        </label>
        <input v-model="newUserUsername" type="text" placeholder="j.smith" class="input input-bordered w-full max-w-xs" />
      </div>

      <div class="form-control">
        <label class="label font-bold">
          <span class="label-text">Password</span>
        </label>
        <input v-model="newUserPassword" type="password" placeholder="hunter2" class="input input-bordered w-full max-w-xs" />
      </div>

      <div class="form-control">
        <label class="label font-bold">
          <span class="label-text">Role</span>
        </label>
        <select class="select select-bordered" v-model="newUserRole">
          <option v-for="role in possibleRoles" :value="role" :key="role">
            {{ role }}
          </option>
        </select>
      </div>

      <div class="form-control mt-3">
        <span class="tooltip" :data-tip="newUserValidationError">
          <button @click="onCreateUser" :disabled="newUserValidationError != null" class="btn btn-primary w-full">Create</button>
        </span>
      </div>
    </Modal>
    <Modal v-model:isOpen="isServiceAccountCreateOpen">
      <h3 class="text-lg font-bold">Create a new service account</h3>

      <div class="form-control">
        <label class="label font-bold">
          <span class="label-text">Service Account Name</span>
        </label>
        <input v-model="newUserUsername" type="text" placeholder="mr.roboto" class="input input-bordered w-full max-w-xs" />
      </div>

      <div class="form-control">
        <label class="label font-bold">
          <span class="label-text">Role</span>
        </label>
        <select class="select select-bordered" v-model="newUserRole">
          <option v-for="role in possibleRoles" :value="role" :key="role">
            {{ role }}
          </option>
        </select>
      </div>

      <div class="form-control mt-3">
        <span class="tooltip" :data-tip="serviceAccountValidationError">
          <button @click="onCreateServiceAccount" :disabled="serviceAccountValidationError != null" class="btn btn-primary w-full">
            Create
          </button>
        </span>
      </div>
    </Modal>
    <h2 class="card-title">Users</h2>
    <div>
      <button class="btn btn-primary btn-sm mr-1" @click="() => (isServiceAccountCreateOpen = true)">Create Service Account</button>
      <button class="btn btn-primary btn-sm" @click="() => (isUserCreateOpen = true)">Create User</button>
    </div>
  </div>

  <table class="table w-full">
    <thead>
      <tr>
        <th>Username</th>
        <th>Role</th>
        <th>Actions</th>
      </tr>
    </thead>
    <tbody>
      <tr class="hover" v-for="user in paginatedUsers" :key="user.id">
        <td>
          <strong>{{ user.username }}</strong>
        </td>
        <td>
          {{ user.roles.join(', ') }}
        </td>
        <td class="text-center">
          <ConfirmModal @on-confirm="() => onDeleteUser(user.id)">
            <IconButton icon="fa-solid fa-trash" color="error" tooltip="Delete" />
          </ConfirmModal>
        </td>
      </tr>
    </tbody>
  </table>
  <div class="mt-2 w-full text-center">
    <PaginationControls @next="() => nextPage()" @prev="() => prevPage()" :current-page="activePage" :total-pages="totalPages" />
  </div>
</template>
