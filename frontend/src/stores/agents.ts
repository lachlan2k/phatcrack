import { getAllAgents } from '@/api/agent'
import type { AgentDTO } from '@/api/types'
import { defineStore } from 'pinia'

export type AgentsStore = {
  agents: AgentDTO[]
  loading: boolean
}

export const useAgentsStore = defineStore({
  id: 'agents-store',

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
