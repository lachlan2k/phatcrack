import Vue from 'vue'
import Vuex from 'vuex'

import user from './user'
import jobForm from './job-form'
import binInterface from './bin-interface'
import transfer from './transfer'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    darkAppearance: false,
    project: window.projectName || 'Fitcrack'
  },
  modules: {
    user,
    jobForm,
    binInterface,
    transfer
  }
})

export function twoWayMap (props) {
  const map = {}
  props.forEach(p => {
    map[p] = `${p}Mut`
  })
  return map
}

// Creates a generic Vuex mutation for every prop in 'obj'

export function supermutator9000 (obj) {
  const mut = prop => (state, val) => { state[prop] = val }
  const muts = {}
  Object.keys(obj).forEach(key => {
    muts[`${key}Mut`] = mut(key)
  })
  return muts
}