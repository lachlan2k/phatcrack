<script setup lang="ts">
import Modal from '@/components/Modal.vue'
import IconButton from '@/components/IconButton.vue'
import { formatDeviceName } from '@/util/formatDeviceName'
import { adminCreateAgent, adminDeleteAgent } from '@/api/admin'
import { useApi } from '@/composables/useApi'
import { useToast } from 'vue-toastification'
import { ref, computed } from 'vue'
import { getAllAgents } from '@/api/agent'
import { useToastError } from '@/composables/useToastError'
import ConfirmModal from '../ConfirmModal.vue'

const AgentStatusHealthy = 'AgentStatusHealthy'
const AgentStatusUnhealthyButConnected = 'AgentStatusUnhealthyButConnected'
const AgentStatusUnhealthyAndDisconnected = 'AgentStatusUnhealthyAndDisconnected'

const isAgentCreateOpen = ref(false)
const { data: agents, fetchData: fetchAgents, isLoading } = useApi(getAllAgents)

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
    <button class="btn btn-primary btn-sm" @click="() => (isAgentCreateOpen = true)">Create Agent</button>
  </div>

  <div v-if="isLoading" class="flex h-56 h-full w-56 w-full justify-center self-center">
    <span class="loading loading-spinner loading-lg"></span>
  </div>
  <table v-else class="table w-full">
    <thead>
      <tr>
        <th>Name</th>
        <th>Version</th>
        <th>Devices</th>
        <th>Status</th>
        <th>Actions</th>
      </tr>
    </thead>
    <tbody>
      <tr class="hover" v-for="agent in agents?.agents" :key="agent.id">
        <td>
          <strong>{{ agent.name }}</strong>
        </td>
        <td class="font-mono">{{ agent.agent_info.version }}</td>
        <td>
          <span v-for="device in agent.agent_devices" :key="device.device_id + device.device_name">
            <font-awesome-icon icon="fa-solid fa-memory" v-if="device.device_type == 'GPU'" />
            <font-awesome-icon icon="fa-solid fa-microchip" v-else />
            {{ formatDeviceName(device.device_name) }} ({{ device.temp }} °c)
            <br />
          </span>
        </td>

        <td class="text-center">
          <div class="badge badge-warning badge-sm m-auto block" title="Marked for maintenance" v-if="agent.is_maintenance_mode"></div>
          <div
            class="badge badge-accent badge-sm m-auto block"
            v-else-if="agent.agent_info.status == AgentStatusHealthy"
            title="Healthy"
          ></div>
          <div
            class="badge badge-warning badge-sm m-auto block"
            title="Unhealthy"
            v-else-if="
              agent.agent_info.status == AgentStatusUnhealthyAndDisconnected || agent.agent_info.status == AgentStatusUnhealthyButConnected
            "
          ></div>
          <div class="badge badge-ghost badge-sm m-auto block" title="Dead" v-else></div>
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
