import { defineStore } from 'pinia'
import { login as apiLogin, refreshAuth as apiRefreshAuth, logout as apiLogout } from '@/api/auth'
import type {
  AuthCurrentUserDTO,
  AuthLoginResponseDTO,
  AuthRefreshResponseDTO,
  AuthWhoamiResponseDTO
} from '@/api/types'

export type AuthState = {
  whoamiDetails: AuthLoginResponseDTO | AuthWhoamiResponseDTO | AuthRefreshResponseDTO | null

  isLoginLoading: boolean
  loginError: string | null
  hasTriedAuth: boolean
}

export const useAuthStore = defineStore({
  id: 'auth-store',

  state: () =>
    ({
      whoamiDetails: null,

      loginError: null,
      isLoginLoading: false,

      // When the app first loads, we don't want to assume a session timeout, so we want to check auth at least once
      hasTriedAuth: false
    } as AuthState),

  actions: {
    async login(username: string, password: string): Promise<boolean> {
      this.isLoginLoading = true
      try {
        const details = await apiLogin(username, password)
        this.whoamiDetails = details
        this.loginError = null
      } catch (err: any) {
        this.whoamiDetails = null
        this.loginError = err.response.data.message
      } finally {
        this.hasTriedAuth = true
        this.isLoginLoading = false
      }

      return this.loggedInUser != null
    },

    async logout() {
      try {
        apiLogout()
      } finally {
        this.whoamiDetails = null
      }
    },

    async refreshAuth() {
      try {
        const details = await apiRefreshAuth()
        this.whoamiDetails = details
        this.loginError = null
      } catch (err: any) {
        // We were logged in before, and now we're not
        if (this.loggedInUser != null) {
          this.loginError = 'Session timeout'
        } else {
          this.loginError = err.response.data.message
        }

        this.whoamiDetails = null
      } finally {
        this.hasTriedAuth = true
      }
    }
  },

  getters: {
    isLoggedIn: (state) => state.whoamiDetails?.user != null,
    loggedInUser: (state) => state.whoamiDetails?.user,

    hasCompletedAuth: (state) =>
      state.whoamiDetails?.user != null &&
      !state.whoamiDetails?.is_awaiting_mfa &&
      !state.whoamiDetails?.requires_password_change &&
      !state.whoamiDetails?.requires_mfa_enrollment,

    isAwaitingMFA: state => state?.whoamiDetails?.is_awaiting_mfa ?? false,
    requiresPasswordChange: state => state?.whoamiDetails?.requires_password_change ?? false,
    requiresMFAEnrollment: state => state?.whoamiDetails?.requires_mfa_enrollment,

    username: (state) => state.whoamiDetails?.user.username,

    error: (state) => state.loginError,
    
    isAdmin: (state) => state.whoamiDetails?.user.roles.includes('admin') ?? false,
    hasRole: (state) => (role: string) => state.whoamiDetails?.user.roles.includes(role) ?? false
  }
})
