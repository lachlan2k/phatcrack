<script setup lang="ts">
import { AxiosError } from 'axios'
import { reactive, watch } from 'vue'
import { useToast } from 'vue-toastification'
import { useApi } from '@/composables/useApi'
import { adminGetConfig, adminSetConfig } from '@/api/admin'

const { data: settingsData, fetchData: reloadSettings } = useApi(adminGetConfig)

const settings = reactive({
  is_mfa_required: false,
  require_password_change_on_first_login: false
})

watch(settingsData, (newSettings) => {
  if (newSettings == null) {
    return
  }

  settings.is_mfa_required = newSettings.is_mfa_required
  settings.require_password_change_on_first_login =
    newSettings.require_password_change_on_first_login
})

const toast = useToast()

async function onSave() {
  try {
    const { is_mfa_required, require_password_change_on_first_login } = settings
    await adminSetConfig({ is_mfa_required, require_password_change_on_first_login })
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
  <div>
    <div class="form-control">
      <label class="label font-bold">
        <span class="label-text pr-3">Require MFA?</span>
        <input type="checkbox" v-model="settings.is_mfa_required" class="toggle" />
      </label>
    </div>
    <div>
      <label class="label font-bold">
        <span class="label-text pr-3">Require password change on first login?</span>
        <input
          type="checkbox"
          v-model="settings.require_password_change_on_first_login"
          class="toggle"
        />
      </label>
    </div>
  </div>
  <div class="card-actions justify-end">
    <button class="btn-primary btn" @click="() => onSave()">Save</button>
  </div>
</template>
