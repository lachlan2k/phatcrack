import { defineStore } from 'pinia'
import { loginWithCredentials as apiLogin, refreshAuth as apiRefreshAuth, logout as apiLogout, loginWithOIDCCallback } from '@/api/auth'
import type { AuthLoginResponseDTO, AuthRefreshResponseDTO, AuthWhoamiResponseDTO } from '@/api/types'

export type AuthState = {
  whoamiDetails: AuthLoginResponseDTO | AuthWhoamiResponseDTO | AuthRefreshResponseDTO | null

  isLoginLoading: boolean
  loginError: string | null
  hasTriedAuth: boolean
  hasLoggedOut: boolean
  isRefreshing: boolean
}

export const useAuthStore = defineStore({
  id: 'auth-store',

  state: () =>
    ({
      whoamiDetails: null,

      loginError: null,
      isLoginLoading: false,

      // When the app first loads, we don't want to assume a session timeout, so we want to check auth at least once
      hasTriedAuth: false,
      hasLoggedOut: false,
      isRefreshing: false
    }) as AuthState,

  actions: {
    async handleOIDCCallback(querystring: string) {
      this.isLoginLoading = true

      try {
        const details = await loginWithOIDCCallback(querystring)
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
        await apiLogout()
      } finally {
        this.hasLoggedOut = true
        this.loginError = ''
        this.whoamiDetails = null
      }
    },

    async refreshAuth() {
      if (this.isRefreshing) {
        return
      }
      this.isRefreshing = true

      try {
        const details = await apiRefreshAuth()
        this.whoamiDetails = details
        this.loginError = null
      } catch (err: any) {
        // We were logged in before, and now we're not
        if (this.loggedInUser != null) {
          if (this.hasLoggedOut) {
            // Did we click logout? if so, reset and don't show a session timeout
            this.hasLoggedOut = false
            this.loginError = ''
          } else {
            // Otherwise, probably a session timeout
            this.loginError = 'Session timeout'
          }
        } else if (err?.response?.data?.message == 'Login required') {
          // this "Login required" error is generic and pointless so we ignore it
        } else {
          this.loginError = err?.response?.data?.message || 'Unknown Error'
        }

        this.whoamiDetails = null
      } finally {
        this.hasTriedAuth = true
        this.isRefreshing = false
      }
    }
  },

  getters: {
    isLoggedIn: state => state.whoamiDetails?.user != null,
    loggedInUser: state => state.whoamiDetails?.user,

    hasCompletedAuth: state =>
      state.whoamiDetails?.user != null &&
      !state.whoamiDetails?.is_awaiting_mfa &&
      !state.whoamiDetails?.requires_password_change &&
      !state.whoamiDetails?.requires_mfa_enrollment,

    isAwaitingMFA: state => state?.whoamiDetails?.is_awaiting_mfa ?? false,
    requiresPasswordChange: state => state?.whoamiDetails?.requires_password_change ?? false,
    requiresMFAEnrollment: state => state?.whoamiDetails?.requires_mfa_enrollment,

    username: state => state.whoamiDetails?.user.username,

    error: state => state.loginError,

    isAdmin: state => state.whoamiDetails?.user.roles.includes('admin') ?? false,
    hasRole: state => (role: string) => state.whoamiDetails?.user.roles.includes(role) ?? false
  }
})
