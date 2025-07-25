import { defineStore } from 'pinia'

import { getAllAgents } from '@/api/agent'
import type { AgentDTO } from '@/api/types'

export type AgentsStore = {
  agents: AgentDTO[]
  loading: boolean
}

export const useAgentsStore = defineStore('agents-store', {
  state: () =>
    ({
      agents: [],
      loading: false
    }) as AgentsStore,

  actions: {
    async load(forceRefetch = false) {
      if (this.loading) {
        return
      }

      if (forceRefetch || this.agents.length === 0) {
        this.loading = true
        try {
          this.agents = (await getAllAgents()).agents
        } finally {
          this.loading = false
        }
      }
    }
  },

  getters: {
    byId: state => (id: string) => state.agents.find(x => x.id == id)
  }
})
