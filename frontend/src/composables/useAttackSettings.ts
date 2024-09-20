import { reactive } from 'vue'

import type { AttackSettingsT } from '@/components/Wizard/AttackSettings.vue'

import type { HashcatParams } from '@/api/types'

import { AttackMode, makeHashcatParams } from '@/util/hashcat'

const startingAttackSettings: AttackSettingsT = {
  attackMode: AttackMode.Dictionary,

  selectedWordlists: [],
  selectedRulefiles: [],

  mask: '',
  maskIncrement: true,

  combinatorLeft: [],
  combinatorRight: [],

  optimizedKernels: false,
  slowCandidates: false,
  enableLoopback: true,
  isDistributed: true
}

export const useAttackSettings = () => {
  const attackSettings = reactive<AttackSettingsT>(startingAttackSettings)

  const resetAttackSettings = () => {
    Object.assign(attackSettings, { ...startingAttackSettings })
  }

  const asHashcatParams = (hashType: number) => {
    const xx = makeHashcatParams(hashType, { ...attackSettings })
    console.log(xx)
    return xx
  }

  const loadFromHashcatParams = (params: HashcatParams) => {
    resetAttackSettings()

    attackSettings.attackMode = params.attack_mode

    attackSettings.optimizedKernels = params.optimized_kernels
    attackSettings.slowCandidates = params.slow_candidates
    attackSettings.enableLoopback = params.enable_loopback

    // comb?
    switch (attackSettings.attackMode) {
      case AttackMode.Dictionary: {
        attackSettings.selectedWordlists = params.wordlist_filenames.slice(0, 1)
        attackSettings.selectedRulefiles = params.rules_filenames

        break
      }

      case AttackMode.Combinator: {
        attackSettings.selectedWordlists = params.wordlist_filenames
        if (attackSettings.selectedWordlists.length >= 2) {
          attackSettings.combinatorLeft = [attackSettings.selectedWordlists[0]]
          attackSettings.combinatorRight = [attackSettings.selectedWordlists[1]]
        }

        break
      }

      case AttackMode.Mask: {
        attackSettings.mask = params.mask
        attackSettings.maskIncrement = params.mask_increment
        break
      }

      case AttackMode.HybridDM:
      case AttackMode.HybridMD: {
        attackSettings.mask = params.mask
        attackSettings.maskIncrement = params.mask_increment
        attackSettings.selectedWordlists = params.wordlist_filenames
        break
      }
    }
  }

  return { attackSettings, resetAttackSettings, loadFromHashcatParams, asHashcatParams }
}
