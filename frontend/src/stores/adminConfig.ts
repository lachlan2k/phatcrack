import { defineStore } from 'pinia'

import { adminGetConfig, adminSetConfig } from '@/api/admin'
import type { AdminConfigResponseDTO, AdminConfigRequestDTO } from '@/api/types'

export type ConfigStore = {
  config: AdminConfigResponseDTO | null
  loadingCount: number
}

export const useAdminConfigStore = defineStore('admin-config-store', {
  state: () => ({ config: null, loadingCount: 0 }) as ConfigStore,

  actions: {
    async load() {
      if (this.loadingCount > 0) {
        return
      }

      this.loadingCount++

      try {
        this.config = await adminGetConfig()
      } finally {
        this.loadingCount--
      }
    },

    async update(newConfig: AdminConfigRequestDTO) {
      this.loadingCount++

      try {
        this.config = await adminSetConfig(newConfig)
      } finally {
        this.loadingCount--
      }
    }
  },

  getters: {
    loading: state => state.loadingCount > 0
  }
})
