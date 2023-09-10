import { getCurrentConfig } from '@/api/config'
import type { ConfigDTO } from '@/api/types'
import { defineStore } from 'pinia'

export type ConfigStore = {
  config: ConfigDTO | null
}

export const useConfigStore = defineStore({
  id: 'config-store',

  state: () => ({ config: null } as ConfigStore),

  actions: {
    async load() {
      this.config = await getCurrentConfig()
    }
  }
})
