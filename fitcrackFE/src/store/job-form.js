import Vue from 'vue'
import { supermutator9000 } from '@/store'

// Vuex module for the Add Job view

// Props that don't save to templates
const base = {
  selectedTemplate: 0, // empty
  step: 1, // Form stepper step
  validatedHashes: [],
  // info
  name: '',
  // input
  inputMethod: 'multipleHashes',
  hashList: '',
  hashType: null,
  ignoreHashes: false,
  // hosts
  hosts: [],
  // other
  startDate: null, // init!
  endDate: null // init!,
}

// Basic empty template
export const empty = {
  template: 'Empty',
  attackSettingsTab: undefined, // Selected attack mode
  // info
  comment: '',
  // settings
  leftDicts: [],
  rightDicts: [],
  rules: [],
  ruleLeft: '',
  ruleRight: '',
  masks: [''],
  hybridMask: '',
  charset: [],
  submode: 0,
  distributionMode: 0,
  markov: [], // init!
  markovThresh: NaN,
  pcfg: [],
  keyspaceLimit: 0,
  // prince
  checkDuplicates: true,
  casePermute: false,
  minPasswordLen: 1,
  maxPasswordLen: 8,
  minElemInChain: 1,
  maxElemInChain: 8,
  generateRandomRules: 0,
  // other
  optimized: true,
  startNow: true,
  endNever: true,
  timeForJob: undefined
}

export default {
  namespaced: true,
  state: {
    ...base,
    ...empty
  },
  getters: {
    attackSettings (state) {
      const current = attacks.find(a => a.handler === state.attackSettingsTab)
      if (!current) return false
      return {
        'attack_mode': current.id,
        'attack_name': current.serverName,
        'rules': state.rules.length > 0 ? state.rules[0] : null,
        'left_dictionaries': state.leftDicts,
        'right_dictionaries': state.rightDicts,
        'rule_left': state.ruleLeft,
        'rule_right': state.ruleRight,
        'masks': state.masks,
        'attack_submode': state.submode,
        'distribution_mode': parseInt(state.distributionMode),
        'markov_treshold': state.submode == 0 ? null : parseInt(state.markovThresh),
        'markov': state.submode == 0 ? null : state.markov[0],
        'charset': state.charset,
        'mask': state.hybridMask,
        'pcfg_grammar': state.pcfg[0],
        'keyspace_limit': state.keyspaceLimit,
        // prince
        'check_duplicates': state.checkDuplicates,
        'case_permute': state.casePermute,
        'min_password_len': parseInt(state.minPasswordLen),
        'max_password_len': parseInt(state.maxPasswordLen),
        'min_elem_in_chain': parseInt(state.minElemInChain),
        'max_elem_in_chain': parseInt(state.maxElemInChain),
        'generate_random_rules': parseInt(state.generateRandomRules),
         // other
        'optimized': state.optimized,
      }
    },
    jobSettings (state, { attackSettings }) {
      return {
        "name": state.name,
        "comment": state.comment,
        "hosts_ids": state.hosts.map(h => h.id),
        "seconds_per_job": parseInt(state.timeForJob),
        "time_start": (state.startNow ? '' : state.startDate),
        "time_end": (state.endNever ? '' : state.endDate),
        'attack_settings': attackSettings,
        "hash_settings": {
          "hash_type": state.hashType ? state.hashType.code : null,
          "hash_list": state.validatedHashes,
          "valid_only": !state.ignoreHashes
        }
      }
    },
    validAttackSpecificSettings (state) {
      switch (state.attackSettingsTab) {
        case 'dictionary':
          return state.leftDicts.length > 0
        case 'combinator':
          return state.leftDicts.length > 0 && state.rightDicts.length > 0
        case 'maskattack':
          if (state.submode > 0 && state.markov.length == 0) return false 
          return state.masks.filter(m => m !== '').length > 0
        case 'hybridWordlistMask':
          return state.leftDicts.length > 0 && state.hybridMask !== ''
        case 'hybridMaskWordlist':
          return state.rightDicts.length > 0 && state.hybridMask !== ''
        case 'pcfgAttack':
          return state.pcfg.length > 0
        case 'princeAttack':
          if (state.leftDicts.length == 0)
            return false
          if (state.minPasswordLen <= 0)
            return false;
          if (state.maxPasswordLen <= 0)
            return false;
          if (state.maxPasswordLen > 32)
            return false;
          if (state.maxElemInChain > 16)
            return false;
          if (state.minElemInChain <= 0)
            return false;
          if (state.maxElemInChain <= 0)
            return false;
          if (state.minPasswordLen > state.maxPasswordLen)
            return false;
          if (state.minElemInChain > state.maxElemInChain)
            return false;
          if (state.maxElemInChain > state.maxPasswordLen)
            return false;
          if (state.keyspaceLimit < 0)
            return false;
          if (state.generateRandomRules < 0)
            return false;
          if (state.rules.length > 0 && state.generateRandomRules > 0)
            return false;

          // All ok!
          return true
        default:
          return true
      }
    },
    valid (state, { attackSettings, validAttackSpecificSettings }) {
      if (
        !state.attackSettingsTab ||
        !attackSettings ||
        state.hashType == null ||
        state.validatedHashes.length == 0 ||
        state.timeForJob < 10 ||
        state.name === ''
      ) {
        return false
      } else {
        return validAttackSpecificSettings
      }
    },
    keyspaceKnown ({ attackSettingsTab }, { validAttackSpecificSettings }) {
      return attackSettingsTab && validAttackSpecificSettings
    },
    template (state) {
      const keys = Object.keys(empty)
      return Object.keys(state)
        .filter(key => keys.includes(key))
        .reduce((obj, key) => {
          if (JSON.stringify(state[key]) !== JSON.stringify(empty[key]) || ( key === 'masks' && state.masks[0] !== '' )) {
            obj[key] = state[key]
          }
          return obj
        }, {})
    },
    validTemplate (state, { template }) {
      return Object.keys(template).length > 0
    },
    makeTemplate (state, { template }) {
      return function (name) {
        const t = template
        t.template = name
        return t
      }
    }
  },
  mutations: {
    ...supermutator9000({
      ...base,
      ...empty
    }),
    update (state, changes) {
      Object.assign(state, changes)
    },
    applyTemplate (state, template = '{}') {
      Object.assign(state, {
        ...empty,
        ...template
      })
    },
    //
    addMask ({ masks }, val = '') {
      masks.push(val)
    },
    deleteMask ({ masks }, index) {
      if (masks.length <= 1) return
      masks.splice(index, 1)
    },
    updateMask ({ masks }, {index, val}) {
      Vue.set(masks, index, val)
    },
    mergeMasks (state, arr) {
      state.masks = [...state.masks, ...arr]
      if (state.masks[0] == '') state.masks.shift()
    }
  }
}

export const attacks = [
  {handler: 'dictionary', name: 'Dictionary', id: 0, serverName: 'dict'},
  {handler: 'combinator', name: 'Combination', id: 1, serverName: 'combinator'},
  {handler: 'maskattack', name: 'Brute-force', id: 3, serverName: 'mask'},
  {handler: 'hybridWordlistMask', name: 'Hybrid Wordlist + Mask', id: 6, serverName: 'Hybrid wordlist+mask'},
  {handler: 'hybridMaskWordlist', name: 'Hybrid Mask + Wordlist', id: 7, serverName: 'Hybrid mask+wordlist'},
  {handler: 'princeAttack', name: 'PRINCE', id: 8, serverName: 'prince'},
  {handler: 'pcfgAttack', name: 'PCFG', id: 9, serverName: 'pcfg'}
]
