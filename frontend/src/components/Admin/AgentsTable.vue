<script setup lang="ts">
import Modal from '@/components/Modal.vue'
import IconButton from '@/components/IconButton.vue'
import { adminCreateAgent, adminGetAllUsers } from '@/api/admin'
import { useApi } from '@/composables/useApi'
import { useToast } from 'vue-toastification'
import { ref, computed } from 'vue'
import { AxiosError } from 'axios'
import { getAllAgents } from '@/api/agent'

const isAgentCreateOpen = ref(false)
const { data: agents, fetchData: fetchAgents } = useApi(getAllAgents)

const toast = useToast()

const newAgentName = ref('')
const newAgentValidationError = computed(() => {
  if (newAgentName.value.length < 4) {
    return 'Name too short'
  }

  if (newAgentName.value.length > 30) {
    return 'Name too long'
  }

  return null
})

async function onCreateAgent() {
  try {
    const res = await adminCreateAgent({
        name: newAgentName.value,
    })

    toast.success('Created new agent: ' + res.name)
    alert(`New agent's auth key: ${res.key} (this won't be displayed again)`)
  } catch (e: any) {
    let errorString = 'Unknown Error'
    if (e instanceof AxiosError) {
      errorString = e.response?.data?.message
    } else if (e instanceof Error) {
      errorString = e.message
    }

    toast.error('Failed to create new agent: ' + errorString)
  } finally {
    fetchAgents()
  }
}
</script>

<template>
  <div class="flex flex-row justify-between">
    <Modal v-model:isOpen="isAgentCreateOpen">
      <h3 class="text-lg font-bold">Create a new agent</h3>

      <div class="form-control">
        <label class="label font-bold">
          <span class="label-text">Agent Name</span>
        </label>
        <input
          v-model="newAgentName"
          type="text"
          placeholder="crack01"
          class="input-bordered input w-full max-w-xs"
        />
      </div>

      <div class="form-control mt-3">
        <span class="tooltip" :data-tip="newAgentValidationError">
          <button
            @click="onCreateAgent"
            :disabled="newAgentValidationError != null"
            class="btn-primary btn w-full"
          >
            Create
          </button>
        </span>
      </div>
    </Modal>
    <h2 class="card-title">Agents</h2>
    <button class="btn-primary btn-sm btn ml-12" @click="() => (isAgentCreateOpen = true)">
      Create Agent
    </button>
  </div>

  <table class="table w-full">
    <thead>
      <tr>
        <th>Name</th>
        <th>Actions</th>
      </tr>
    </thead>
    <tbody>
      <tr class="hover" v-for="agent in agents?.agents" :key="agent.id">
        <td>
          <strong>{{ agent.name }}</strong>
        </td>
        <td class="text-center">
          <IconButton icon="fa-solid fa-trash" color="error" tooltip="Delete" />
        </td>
      </tr>
    </tbody>
  </table>
</template>
