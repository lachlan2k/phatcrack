import { loadHashTypes } from '@/api/hashcat'
import type { HashType } from '@/api/types'
import { defineStore } from 'pinia'

export type ResourceStore = {
  hashTypes: HashType[]
  loading: boolean
}

export const useResourcesStore = defineStore({
  id: 'resources-store',

  state: () =>
    ({
      hashTypes: [],
      loading: false
    }) as ResourceStore,

  actions: {
    async loadHashTypes() {
      if (this.hashTypes.length === 0 && !this.loading) {
        this.loading = true
        try {
          const hashTypes = await loadHashTypes()
          this.hashTypes = Object.values(hashTypes.hashtypes).sort((a: HashType, b: HashType) => a.id - b.id)
        } finally {
          this.loading = false
        }
      }
    },

    loadAll() {
      this.loadHashTypes()
    }
  },

  getters: {
    isHashTypesLoaded: state => state.hashTypes.length > 0,
    getHashTypeName: state => (hashId: number) => state.hashTypes.find(x => x.id == hashId)?.name ?? ''
  }
})
