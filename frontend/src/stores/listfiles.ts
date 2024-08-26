import { getAllListfiles, LISTFILE_TYPE_RULEFILE, LISTFILE_TYPE_WORDLIST } from '@/api/listfiles'
import type { ListfileDTO } from '@/api/types'
import { defineStore } from 'pinia'

export type ListfileStore = {
  listfiles: ListfileDTO[]
  loading: boolean
}

export const useListfilesStore = defineStore({
  id: 'listfiles-store',

  state: () =>
    ({
      listfiles: [],
      loading: false
    } as ListfileStore),

  actions: {
    async load(forceRefetch = false) {
      if (this.loading) {
        return
      }

      if (forceRefetch || (this.listfiles.length === 0)) {
        this.loading = true
        try {
          const { listfiles } = await getAllListfiles()
          this.listfiles = listfiles
        } finally {
          this.loading = false
        }
      }
    }
  },

  getters: {
    byId: (state) => (id: string) => state.listfiles.find((x) => x.id == id),
    wordlists: (state) => state.listfiles.filter(x => x.file_type === LISTFILE_TYPE_WORDLIST),
    rulefiles: (state) => state.listfiles.filter(x => x.file_type === LISTFILE_TYPE_RULEFILE)
  }
})
