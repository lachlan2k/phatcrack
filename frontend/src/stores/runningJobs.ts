import { getRunningJobs } from '@/api/project'
import type { RunningJobForUserDTO } from '@/api/types'
import { defineStore } from 'pinia'

export type RunningJobsStore = {
  jobs: RunningJobForUserDTO[]
  loading: boolean
}

export const useRunningJobsStore = defineStore({
  id: 'running-jobs-store',

  state: () =>
    ({
      jobs: [],
      loading: false
    } as RunningJobsStore),

  actions: {
    async load() {
      if (this.loading) {
        return
      }

      try {
        this.loading = true
        const { jobs } = await getRunningJobs()

        this.jobs = jobs
      } finally {
        this.loading = false
      }
    }
  },

  getters: {
    forProject: (state) => (projectId: string) => state.jobs.filter((x) => x.project_id === projectId),
    forHashlist: (state) => (hashlistId: string) => state.jobs.filter((x) => x.hashlist_id === hashlistId)
  }
})
