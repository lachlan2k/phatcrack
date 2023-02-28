import { defineStore } from 'pinia'
import { login as apiLogin, refreshAuth as apiRefreshAuth } from '@/api/auth'
import type { APILoggedInUserDetailsT } from '@/api/auth'
import type { AxiosError } from 'axios'


export type AuthState = {
    loggedInUser: APILoggedInUserDetailsT | null
    loginError: string | null
}

export const useAuthStore = defineStore({
    id: 'auth-store',
    
    state: () => ({
        loggedInUser: null,
        loginError: null
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
            }
        },

        logout() {

        },

        async refreshAuth() {
            try {
                const details = await apiRefreshAuth()
                console.log('refreshed', details.data)
                this.loggedInUser = details!.data?.user ?? null
                this.loginError = null
            } catch (err: any) {
                this.loggedInUser = null
                this.loginError = err.response.data.message
            }
        }
    },

    getters: {
        isLoggedIn: state => state.loggedInUser != null,
        username: state => state.loggedInUser?.username,
        error: state => state.loginError
    }
})