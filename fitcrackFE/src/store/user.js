import router from '@/router'

import {axios} from '@bundled-es-modules/axios'
const api = window.serverAddress

const emptyUser = {
  id: null,
  username: null,
  mail: null,
  role: {
    'MANAGE_USERS': false,
    'ADD_NEW_JOB': false,
    'UPLOAD_DICTIONARIES': false,
    'VIEW_ALL_JOBS': false,
    'EDIT_ALL_JOBS': false,
    'OPERATE_ALL_JOBS': false,
    'ADD_USER_PERMISSIONS_TO_JOB': false
  }
}

export default {
  state: {
    loggedIn: false,
    userData: emptyUser
  },
  getters: {
    can: state => permission => state.userData.role[permission]
  },
  mutations: {
    _setUser (state, data) {
      if (!data) {
        state.loggedIn = false
        state.userData = emptyUser
      } else {
        state.loggedIn = true
        state.userData = data
      }
    }
  },
  actions: {
    async refreshUser ({ commit }) {
      const {user, loggedIn, token} = await axios.get(`${api}/user/isLoggedIn`).then(r => r.data)
      localStorage.setItem('jwt', token)
      commit('_setUser', loggedIn ? user : null)
      return user
    },
    async resume ({ state, dispatch }) {
      if (state.loggedIn) return true
      const user = await dispatch('refreshUser')
      return user.id != null
    },
    async signIn ({ commit }, credentials) {
      try {
        const {user, token} = await axios.post(`${api}/user/login`, credentials).then(r => r.data)
        localStorage.setItem('jwt', token)
        commit('_setUser', user)
        const redirect = sessionStorage.getItem('loginRedirect')
        sessionStorage.removeItem('loginRedirect')
        if (!redirect || redirect === '/login') {
          router.replace({ name: 'home' })
        } else {
          router.push(redirect)
        }
        return true
      } catch (e) {
        console.error(e)
        commit('_setUser', null)
        return false
      }
    },
    async signOut ({ commit }) {
      await axios.get(`${api}/user/logout`)
      commit('_setUser', null)
      router.push({ name: 'login' })
    },
    async hasPermission ({ state, dispatch }, name) {
      if (!state.loggedIn) await dispatch('resume')
      return state.userData.role[name]
    }
  }
}