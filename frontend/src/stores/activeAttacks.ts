import { defineStore } from 'pinia'

import { getAttacksInitialising, getRunningJobs } from '@/api/project'
import type { AttackIDTreeDTO, RunningJobForUserDTO } from '@/api/types'

export type ActiveAttacksStore = {
  jobs: RunningJobForUserDTO[]
  initialisingAttacks: AttackIDTreeDTO[]
  loading: boolean
}

export const useActiveAttacksStore = defineStore({
  id: 'active-attacks-store',

  state: () =>
    ({
      jobs: [],
      initialisingAttacks: [],
      loading: false
    }) as ActiveAttacksStore,

  actions: {
    async load() {
      if (this.loading) {
        return
      }

      try {
        this.loading = true
        const { jobs } = await getRunningJobs()
        const { attacks } = await getAttacksInitialising()

        this.jobs = jobs
        this.initialisingAttacks = attacks
      } finally {
        this.loading = false
      }
    }
  },

  getters: {
    jobsForProject: state => (projectId: string) => state.jobs.filter(x => x.project_id === projectId),
    initialisingAttacksForHashlist: state => (hashlistId: string) => state.initialisingAttacks.filter(x => x.hashlist_id == hashlistId),
    initialisingAttacksForProject: state => (projectId: string) => state.initialisingAttacks.filter(x => x.project_id == projectId),
    jobsForHashlist: state => (hashlistId: string) => state.jobs.filter(x => x.hashlist_id === hashlistId)
  }
})
