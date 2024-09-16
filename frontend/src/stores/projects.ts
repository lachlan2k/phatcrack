import { defineStore } from 'pinia'

import { getAllProjects } from '@/api/project'
import type { ProjectDTO } from '@/api/types'

export type ProjectsState = {
  projects: ProjectDTO[]
  isLoading: boolean
}

export const useProjectsStore = defineStore({
  id: 'projects-store',

  state: () =>
    ({
      projects: [],
      isLoading: false
    }) as ProjectsState,

  actions: {
    async load(forceRefetch: boolean = false) {
      if ((this.projects.length > 0 || this.isLoading) && !forceRefetch) {
        return
      }

      try {
        this.isLoading = true
        this.projects = (await getAllProjects()).projects
      } finally {
        this.isLoading = false
      }
    }
  },

  getters: {
    byId: state => (projId: string) => state.projects.find(x => x.id === projId) ?? null
  }
})
