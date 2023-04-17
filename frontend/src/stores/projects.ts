import { getAllProjects } from '@/api/project'
import type { ProjectDTO } from '@/api/types'
import { defineStore } from 'pinia'

export type ProjectsState = {
  projects: ProjectDTO[]
}

export const useProjectsStore = defineStore({
  id: 'projects-store',

  state: () =>
    ({
      projects: []
    } as ProjectsState),

  actions: {
    async load() {
      this.projects = (await getAllProjects()).projects
    }
  }
})
