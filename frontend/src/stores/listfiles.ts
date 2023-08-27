import { getAllRulefiles, getAllWordlists } from '@/api/listfiles'
import type { ListfileDTO } from '@/api/types'
import { defineStore } from 'pinia'

export type ListfileStore = {
  wordlists: ListfileDTO[]
  rulefiles: ListfileDTO[]
  loading: boolean
}

export const useListfilesStore = defineStore({
  id: 'listfiles-store',

  state: () =>
    ({
      wordlists: [],
      rulefiles: [],
      loading: false
    } as ListfileStore),

  actions: {
    async load(forceRefetch = false) {
      if (this.loading) {
        return
      }

      if (forceRefetch || (this.wordlists.length === 0 && this.rulefiles.length == 0)) {
        this.loading = true
        try {
          const wordlistsReq = getAllWordlists()
          const rulefilesReq = getAllRulefiles()
          const [{ wordlists }, { rulefiles }] = await Promise.all([wordlistsReq, rulefilesReq])

          this.wordlists = wordlists
          this.rulefiles = rulefiles
        } finally {
          this.loading = false
        }
      }
    }
  },

  getters: {
    byId: (state) => (id: string) =>
      state.wordlists.find((x) => x.id == id) ?? state.rulefiles.find((x) => x.id == id)
  }
})
