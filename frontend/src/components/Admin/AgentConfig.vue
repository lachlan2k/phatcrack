<script setup lang="ts">
import { reactive, watch } from 'vue'
import { useToast } from 'vue-toastification'
import { useApi } from '@/composables/useApi'
import { adminGetConfig, adminSetConfig } from '@/api/admin'
import { useConfigStore } from '@/stores/config'
import { useToastError } from '@/composables/useToastError'

const configStore = useConfigStore()
const { data: settingsData, silentlyRefresh: reloadSettings, isLoading } = useApi(adminGetConfig)

const agentSettings = reactive({
  auto_sync_listfiles: false,
  split_jobs_per_agent: 0
})

watch(settingsData, (newSettings) => {
  const agent = newSettings?.agent
  if (agent == null) {
    return
  }

  agentSettings.auto_sync_listfiles = agent.auto_sync_listfiles
  agentSettings.split_jobs_per_agent = agent.split_jobs_per_agent
})

const toast = useToast()
const { catcher } = useToastError()

async function onSave() {
  try {
    await adminSetConfig({
      agent: {
        auto_sync_listfiles: agentSettings.auto_sync_listfiles,
        split_jobs_per_agent: agentSettings.split_jobs_per_agent
      }
    })
    configStore.load()
    toast.success('Settings saved')
  } catch (e: any) {
    catcher(e)
  } finally {
    reloadSettings()
  }
}
</script>

<template>
  <table class="compact-table table-first-col-bold table">
    <thead>
      <tr>
        <th>Setting</th>
        <th>Value</th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td>Automatically sync listfiles to agents</td>
        <td><input type="checkbox" class="toggle" v-model="agentSettings.auto_sync_listfiles" /></td>
      </tr>

      <tr>
        <td>Number of jobs per agent for each job (recommended: 1)</td>
        <td><input type="number" v-model.number="agentSettings.split_jobs_per_agent" class="input input-bordered input-sm w-40" /></td>
      </tr>

      <tr>
        <td></td>
        <td>
          <button class="btn btn-primary btn-sm" @click="() => onSave()">Save</button>
        </td>
      </tr>
    </tbody>
  </table>
</template>
