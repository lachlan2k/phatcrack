import { defineStore } from 'pinia'

import { deleteAttackTemplate, getAllAttackTemplates } from '@/api/attackTemplate'
import type { AttackTemplateDTO } from '@/api/types'

export type AttackTemplatesStore = {
  templates: AttackTemplateDTO[]
  loading: boolean
}

export const useAttackTemplatesStore = defineStore({
  id: 'attack-templates-store',

  state: () =>
    ({
      templates: [],
      loading: false
    }) as AttackTemplatesStore,

  actions: {
    async load(forceRefresh: boolean = false) {
      if (this.loading) {
        return
      }

      if (this.templates.length > 0 && !forceRefresh) {
        return
      }

      try {
        this.loading = true
        const { attack_templates } = await getAllAttackTemplates()
        this.templates = attack_templates
      } finally {
        this.loading = false
      }
    },

    async delete(templateId: string) {
      const res = await deleteAttackTemplate(templateId)
      this.load()
      return res
    }
  },

  getters: {
    byId: state => (templateId: string) => state.templates.find(x => x.id === templateId)
  }
})
