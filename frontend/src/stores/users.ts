import type { UserMinimalDTO } from '@/api/types'
import { getAllUsers } from '@/api/users'
import { defineStore } from 'pinia'

export type UsersState = {
  users: UserMinimalDTO[]
  isLoading: boolean
}

export const useUsersStore = defineStore({
  id: 'users-store',

  state: () =>
    ({
      users: [],
      isLoading: false
    }) as UsersState,

  actions: {
    async load(forceRefetch: boolean = false) {
      if ((this.users.length > 0 || this.isLoading) && !forceRefetch) {
        return
      }

      try {
        this.isLoading = true
        this.users = (await getAllUsers()).users
      } finally {
        this.isLoading = false
      }
    }
  },

  getters: {
    byId: (state) => (userId: string) => state.users.find((x) => x.id === userId) ?? null
  }
})
