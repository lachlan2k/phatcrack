import { defineStore } from 'pinia'

import { AuthMethodCredentials, AuthMethodOIDC, getCurrentConfig } from '@/api/config'
import type { PublicConfigDTO } from '@/api/types'

export type ConfigStore = {
  config: PublicConfigDTO | null
  loading: boolean
}

export const useConfigStore = defineStore('config-store', {
  state: () => ({ config: null, loading: false }) as ConfigStore,

  actions: {
    async load() {
      if (this.loading) {
        return
      }

      this.loading = true

      try {
        this.config = await getCurrentConfig()
      } finally {
        this.loading = false
      }
    }
  },

  getters: {
    isCredentialAuthEnabled: state => state.config?.auth.enabled_methods?.includes(AuthMethodCredentials) ?? false,
    isOIDCAuthEnabled: state => state.config?.auth.enabled_methods?.includes(AuthMethodOIDC) ?? false
  }
})
