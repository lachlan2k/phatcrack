<script setup lang="ts">
import Modal from '@/components/Modal.vue'
import IconButton from '@/components/IconButton.vue'
import { adminCreateAgent, adminDeleteAgent } from '@/api/admin'
import { useApi } from '@/composables/useApi'
import { useToast } from 'vue-toastification'
import { ref, computed } from 'vue'
import { getAllAgents } from '@/api/agent'
import { useToastError } from '@/composables/useToastError'
import ConfirmModal from '../ConfirmModal.vue'

const isAgentCreateOpen = ref(false)
const { data: agents, fetchData: fetchAgents } = useApi(getAllAgents)

const toast = useToast()
const { catcher } = useToastError()

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
      name: newAgentName.value
    })
  
    toast.info(`Created new agent ${res.name}.\n\nNew agent's auth key (note this down, won't be displayed again):\n${res.key}`, {
      // force user to dismiss this
      timeout: false,
      closeOnClick: false,
      draggable: false
    })
  } catch (e: any) {
    catcher(e)
  } finally {
    fetchAgents()
  }
}

async function onDeleteAgent(id: string) {
  try {
    await adminDeleteAgent(id)
    toast.info('Deleted agent')
  } catch (e: any) {
    catcher(e)
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
        <input v-model="newAgentName" type="text" placeholder="crack01" class="input input-bordered w-full max-w-xs" />
      </div>

      <div class="form-control mt-3">
        <span class="tooltip" :data-tip="newAgentValidationError">
          <button @click="onCreateAgent" :disabled="newAgentValidationError != null" class="btn btn-primary w-full">Create</button>
        </span>
      </div>
    </Modal>
    <h2 class="card-title">Agents</h2>
    <button class="btn btn-primary btn-sm ml-12" @click="() => (isAgentCreateOpen = true)">Create Agent</button>
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
          <ConfirmModal @on-confirm="() => onDeleteAgent(agent.id)">
            <IconButton icon="fa-solid fa-trash" color="error" tooltip="Delete" />
          </ConfirmModal>
        </td>
      </tr>
    </tbody>
  </table>
</template>
