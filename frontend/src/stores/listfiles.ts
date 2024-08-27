import { getAllListfiles, getListfilesForProject, LISTFILE_TYPE_RULEFILE, LISTFILE_TYPE_WORDLIST } from '@/api/listfiles'
import type { ListfileDTO } from '@/api/types'
import { defineStore } from 'pinia'

export type ListfileStore = {
  listfiles: ListfileDTO[]
  projectListfileMap: { [key: string]: ListfileDTO[] }
  projectLoadingMap: { [key: string]: boolean }
  loading: boolean
}

export const useListfilesStore = defineStore({
  id: 'listfiles-store',

  state: () =>
    ({
      listfiles: [],
      projectListfileMap: {},
      projectLoadingMap: {},
      loading: false
    } as ListfileStore),

  actions: {
    async load(forceRefetch = false) {
      if (this.loading) {
        return
      }

      if (forceRefetch || this.listfiles.length === 0) {
        this.loading = true
        try {
          const { listfiles } = await getAllListfiles()
          this.listfiles = listfiles
        } finally {
          this.loading = false
        }
      }
    },

    async loadForProject(projectId: string) {
      if (this.projectLoadingMap[projectId]) {
        return
      }

      this.projectLoadingMap[projectId] = true

      try {
        const { listfiles } = await getListfilesForProject(projectId)
        this.projectListfileMap[projectId] = listfiles
      } finally {
        this.projectLoadingMap[projectId] = false
      }
    }
  },

  getters: {
    byId: (state) => (id: string) => state.listfiles.find((x) => x.id == id),

    wordlists: (state) => state.listfiles.filter((x) => x.file_type === LISTFILE_TYPE_WORDLIST),
    rulefiles: (state) => state.listfiles.filter((x) => x.file_type === LISTFILE_TYPE_RULEFILE),

    groupedByType: (state) => {
      // map to { wordlists: [...], rulefiles: [...], etc... }
      const map = {} as { [key: string]: ListfileDTO[] }

      const sortedInsert = (arr: ListfileDTO[] | null, val: ListfileDTO): ListfileDTO[] => {
        if (arr == null || arr.length == 0) {
          return [val]
        }

        // Find the first one bigger
        const index = arr.findIndex((x) => x.lines > val.lines)
        if (index == -1) {
          return [...arr, val]
        }

        return [...arr.slice(0, index), val, ...arr.slice(index)]
      }

      const grouped = state.listfiles.reduce(
        (acc, obj) => ({
          ...acc,
          [obj.file_type]: sortedInsert(acc[obj.file_type], obj)
        }),
        map
      )

      return grouped
    },

    attachedToProject: (state) => (projectId: string) => state.projectListfileMap[projectId] ?? [],

    // Returns all listfiles a project can use, filtered by type. Places project-specific ones first
    availableToProjectByType: (state) => (projectId: string, type: string) =>
      (state.projectListfileMap[projectId] ?? [])
        .filter((x) => x.file_type === type)
        .concat(state.listfiles.filter((x) => x.file_type === type))
  }
})
