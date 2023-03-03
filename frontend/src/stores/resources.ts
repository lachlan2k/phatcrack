import { loadHashTypes } from '@/api/resources'
import type { HashType } from '@/api/types'
import { defineStore } from 'pinia'

export type ResourceStore = {
  hashTypes: HashType[]
}

export const useResourcesStore = defineStore({
  id: 'resources-store',

  state: () =>
    ({
      hashTypes: []
    } as ResourceStore),

  actions: {
    async loadHashTypes() {
      if (this.hashTypes.length === 0) {
        const hashTypes = await loadHashTypes()
        this.hashTypes = Object.values(hashTypes.hashtypes).sort(
          (a: HashType, b: HashType) => a.id - b.id
        )
      }
    },

    loadAll() {
      this.loadHashTypes()
    }
  }
})