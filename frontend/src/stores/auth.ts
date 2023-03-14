import { defineStore } from 'pinia'
import { login as apiLogin, refreshAuth as apiRefreshAuth, logout as apiLogout } from '@/api/auth'
import type { AuthCurrentUserDTO } from '@/api/types'

export type AuthState = {
  loggedInUser: AuthCurrentUserDTO | null
  loginError: string | null
  hasTriedAuth: boolean
}

export const useAuthStore = defineStore({
  id: 'auth-store',

  state: () =>
    ({
      loggedInUser: null,
      loginError: null,
      hasTriedAuth: false // When the app first loads, we don't want to assume a session timeout, so we want to check auth at least once
    } as AuthState),

  actions: {
    async login(username: string, password: string) {
      try {
        const details = await apiLogin(username, password)
        this.loggedInUser = details?.user ?? null
        this.loginError = null
      } catch (err: any) {
        this.loggedInUser = null
        this.loginError = err.response.data.message
      } finally {
        this.hasTriedAuth = true
      }
    },

    async logout() {
      try {
        apiLogout()
      } finally {
        this.loggedInUser = null
      }
    },

    async refreshAuth() {
      try {
        const details = await apiRefreshAuth()
        this.loggedInUser = details.user
        this.loginError = null
      } catch (err: any) {
        // We were logged in before, and now we're not
        if (this.loggedInUser != null) {
          this.loginError = 'Session timeout'
        } else {
          this.loginError = err.response.data.message
        }

        this.loggedInUser = null
      } finally {
        this.hasTriedAuth = true
      }
    }
  },

  getters: {
    isLoggedIn: (state) => state.loggedInUser != null,
    username: (state) => state.loggedInUser?.username,
    error: (state) => state.loginError,
    isAdmin: (state) => state.loggedInUser?.role === 'admin'
  }
})
