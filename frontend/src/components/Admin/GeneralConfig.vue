<script setup lang="ts">
import { reactive, watch } from 'vue'
import { useToast } from 'vue-toastification'
import { useApi } from '@/composables/useApi'
import { adminGetConfig, adminSetConfig } from '@/api/admin'
import { useConfigStore } from '@/stores/config'
import { useToastError } from '@/composables/useToastError'

const configStore = useConfigStore()
const { data: settingsData, silentlyRefresh: reloadSettings, isLoading } = useApi(adminGetConfig)

const generalSettings = reactive({
  is_maintenance_mode: false,
  maximum_uploaded_file_size: 1,
  maximum_uploaded_file_line_scan_size: 1
})

watch(settingsData, (newSettings) => {
  const general = newSettings?.general
  if (general == null) {
    return
  }

  generalSettings.is_maintenance_mode = general.is_maintenance_mode
  generalSettings.maximum_uploaded_file_size = general.maximum_uploaded_file_size
  generalSettings.maximum_uploaded_file_line_scan_size = general.maximum_uploaded_file_line_scan_size
})

const toast = useToast()
const { catcher } = useToastError()

async function onSave() {
  try {
    await adminSetConfig({
      general: {
        is_maintenance_mode: generalSettings.is_maintenance_mode,
        maximum_uploaded_file_size: generalSettings.maximum_uploaded_file_size,
        maximum_uploaded_file_line_scan_size: generalSettings.maximum_uploaded_file_line_scan_size
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
        <td>Enable maintenance mode</td>
        <td><input type="checkbox" class="toggle" v-model="generalSettings.is_maintenance_mode" /></td>
      </tr>

      <tr>
        <td>Maximum upload size (bytes)</td>
        <td>
          <input type="number" v-model.number="generalSettings.maximum_uploaded_file_size" class="input input-bordered input-sm w-40" />
        </td>
      </tr>

      <tr>
        <td>Maximum size of file to scan for line count (bytes)</td>
        <td>
          <input
            type="number"
            v-model.number="generalSettings.maximum_uploaded_file_line_scan_size"
            class="input input-bordered input-sm w-40"
          />
        </td>
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
