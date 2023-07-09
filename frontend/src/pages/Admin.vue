<script setup lang="ts">
import { AxiosError } from 'axios'
import Modal from '@/components/Modal.vue'
import IconButton from '@/components/IconButton.vue'
import { adminCreateUser, adminGetAllUsers } from '@/api/admin'
import { useApi } from '@/composables/useApi'
import { useToast } from 'vue-toastification'
import { ref, computed } from 'vue'

const isUserCreateOpen = ref(false)
const { data: allUsers, fetchData: fetchUsers } = useApi(adminGetAllUsers)

const toast = useToast()

const possibleRoles = ['admin', 'standard']

const newUserUsername = ref('')
const newUserPassword = ref('')
const newUserRole = ref('standard')

const newUserValidationError = computed(() => {
  if (newUserUsername.value.length < 3) {
    return 'Username too short'
  }

  if (newUserPassword.value.length < 16) {
    return 'Password too short'
  }

  return null
})

async function onCreateUser() {
  try {
    const res = await adminCreateUser({
      username: newUserUsername.value,
      password: newUserPassword.value,
      roles: [newUserRole.value]
    })

    toast.success('Created new user: ' + res.username)
  } catch (e: any) {
    let errorString = 'Unknown Error'
    if (e instanceof AxiosError) {
      errorString = e.response?.data?.message
    } else if (e instanceof Error) {
      errorString = e.message
    }

    toast.error('Failed to create new user: ' + errorString)
  } finally {
    fetchUsers()
  }
}
// const newUserRoles = ref([])
</script>

<template>
  <main class="w-full p-4">
    <div class="prose">
      <h1>Admin</h1>
    </div>

    <div class="mt-6 flex flex-wrap gap-6">
      <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <div class="flex flex-row justify-between">
            <Modal v-model:isOpen="isUserCreateOpen">
              <h3 class="text-lg font-bold">Create a new user</h3>

              <div class="form-control">
                <label class="label font-bold">
                  <span class="label-text">Username</span>
                </label>
                <input
                  v-model="newUserUsername"
                  type="text"
                  placeholder="j.smith"
                  class="input-bordered input w-full max-w-xs"
                />
              </div>

              <div class="form-control">
                <label class="label font-bold">
                  <span class="label-text">Password</span>
                </label>
                <input
                  v-model="newUserPassword"
                  type="password"
                  placeholder="hunter2"
                  class="input-bordered input w-full max-w-xs"
                />
              </div>

              <div class="form-control">
                <label class="label font-bold">
                  <span class="label-text">Role</span>
                </label>
                <select class="select-bordered select" v-model="newUserRole">
                  <option v-for="role in possibleRoles" :value="role" :key="role">
                    {{ role }}
                  </option>
                </select>
              </div>

              <div class="form-control mt-3">
                <span class="tooltip" :data-tip="newUserValidationError">
                  <button
                    @click="onCreateUser"
                    :disabled="newUserValidationError != null"
                    class="btn-primary btn w-full"
                  >
                    Create
                  </button>
                </span>
              </div>
            </Modal>
            <h2 class="card-title">Users</h2>
            <button class="btn-primary btn-sm btn" @click="() => (isUserCreateOpen = true)">
              Create User
            </button>
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
              <tr class="hover" v-for="user in allUsers?.users" :key="user.id">
                <td>
                  <strong>{{ user.username }}</strong>
                </td>
                <td>
                  {{ user.roles.join(', ') }}
                </td>
                <td class="text-center">
                  <IconButton icon="fa-solid fa-trash" color="error" tooltip="Delete" />
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </main>
</template>
