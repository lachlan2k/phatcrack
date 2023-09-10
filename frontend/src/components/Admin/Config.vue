<script setup lang="ts">
import { AxiosError } from 'axios'
import { reactive, watch } from 'vue'
import { useToast } from 'vue-toastification'
import { useApi } from '@/composables/useApi'
import { adminGetConfig, adminSetConfig } from '@/api/admin'
import { useConfigStore } from '@/stores/config'

const configStore = useConfigStore()
const { data: settingsData, silentlyRefresh: reloadSettings, isLoading } = useApi(adminGetConfig)

const settings = reactive({
  is_mfa_required: false,
  is_maintenance_mode: false,
  auto_sync_listfiles: false,
  require_password_change_on_first_login: false,
  split_jobs_per_agent: 1,
  maximum_uploaded_file_size: 1,
  maximum_uploaded_file_line_scan_size: 1
})

watch(settingsData, (newSettings) => {
  if (newSettings == null) {
    return
  }

  settings.is_mfa_required = newSettings.is_mfa_required
  settings.is_maintenance_mode = newSettings.is_maintenance_mode
  settings.auto_sync_listfiles = newSettings.auto_sync_listfiles
  settings.require_password_change_on_first_login = newSettings.require_password_change_on_first_login
  settings.split_jobs_per_agent = newSettings.split_jobs_per_agent
  settings.maximum_uploaded_file_size = newSettings.maximum_uploaded_file_size
  settings.maximum_uploaded_file_line_scan_size = newSettings.maximum_uploaded_file_line_scan_size
})

const toast = useToast()

async function onSave() {
  try {
    const {
      is_mfa_required,
      is_maintenance_mode,
      require_password_change_on_first_login,
      auto_sync_listfiles,
      split_jobs_per_agent,
      maximum_uploaded_file_size,
      maximum_uploaded_file_line_scan_size
    } = settings

    await adminSetConfig({
      is_mfa_required,
      is_maintenance_mode,
      require_password_change_on_first_login,
      auto_sync_listfiles,
      split_jobs_per_agent,
      maximum_uploaded_file_size,
      maximum_uploaded_file_line_scan_size
    })
    configStore.load()
    toast.success('Settings saved')
  } catch (e: any) {
    let errorString = 'Unknown Error'
    if (e instanceof AxiosError) {
      errorString = e.response?.data?.message
    } else if (e instanceof Error) {
      errorString = e.message
    }

    toast.error('Failed to save settings: ' + errorString)
  } finally {
    reloadSettings()
  }
}
</script>

<template>
  <div v-if="isLoading" class="flex h-56 h-full w-56 w-full justify-center">
    <span class="loading loading-spinner loading-lg"></span>
  </div>
  <div v-else>
    <div class="form-control">
      <label class="label font-bold">
        <span class="label-text pr-3">Enable maintenance mode?</span>
        <input type="checkbox" v-model="settings.is_maintenance_mode" class="toggle" />
      </label>
    </div>
    <div class="form-control">
      <label class="label font-bold">
        <span class="label-text pr-3">Require MFA?</span>
        <input type="checkbox" v-model="settings.is_mfa_required" class="toggle" />
      </label>
    </div>
    <div>
      <label class="label font-bold">
        <span class="label-text pr-3">Require password change on first login?</span>
        <input type="checkbox" v-model="settings.require_password_change_on_first_login" class="toggle" />
      </label>
    </div>
    <div class="form-control">
      <label class="label font-bold">
        <span class="label-text pr-3">Automatically sync files to agents?</span>
        <input type="checkbox" v-model="settings.auto_sync_listfiles" class="toggle" />
      </label>
    </div>
    <div class="form-control">
      <label class="label font-bold">
        <span class="label-text pr-3">How many jobs per-agent to split (recommended: 1)?</span>
        <input type="number" v-model.number="settings.split_jobs_per_agent" class="input input-bordered input-sm w-16" />
      </label>
    </div>
    <div class="form-control">
      <label class="label font-bold">
        <span class="label-text pr-3">Max upload size (bytes)</span>
        <input type="number" v-model.number="settings.maximum_uploaded_file_size" class="input input-bordered input-sm w-40" />
      </label>
    </div>
    <div class="form-control">
      <label class="label font-bold">
        <span class="label-text pr-3">Max size uploaded file to scan (bytes)</span>
        <input type="number" v-model.number="settings.maximum_uploaded_file_line_scan_size" class="input input-bordered input-sm w-40" />
      </label>
    </div>
  </div>
  <div class="card-actions justify-end">
    <button class="btn btn-primary" @click="() => onSave()">Save</button>
  </div>
</template>
