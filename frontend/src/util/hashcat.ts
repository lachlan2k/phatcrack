import type { AttackSettingsT } from '@/components/Wizard/AttackSettings.vue'

import type { HashcatParams } from '@/api/types'

export enum AttackMode {
  Dictionary = 0,
  Combinator = 1,
  Mask = 3,
  HybridDM = 6,
  HybridMD = 7,
  Template = 999
}

export const attackModes = [
  { name: 'Wordlist', value: AttackMode.Dictionary },
  { name: 'Combinator', value: AttackMode.Combinator },
  { name: 'Brute-force/Mask', value: AttackMode.Mask },
  {
    name: 'Wordlist + Mask',
    value: AttackMode.HybridDM
  },
  {
    name: 'Mask + Wordlist',
    value: AttackMode.HybridMD
  },
  {
    name: 'From Template',
    value: AttackMode.Template
  }
]

export function getAttackModeName(id: number): string {
  return attackModes.find(x => x.value == id)?.name ?? ''
}

export function modeHasMask(value: number): boolean {
  return value == 1 || value == 3 || value == 6 || value == 7
}

export interface MaskInfo {
  mask: string
  charset: string
  description: string
}

export const maskCharsets: MaskInfo[] = [
  { mask: '?l', charset: 'abcdefghijklmnopqrstuvwxyz', description: 'Lowercase letters' },
  { mask: '?u', charset: 'ABCDEFGHIJKLMNOPQRSTUVWXYZ', description: 'Uppercase letters' },
  { mask: '?d', charset: '0123456789', description: 'Digits 0-9' },
  { mask: '?h', charset: '0123456789abcdef', description: 'Lowercase hex (0-f)' },
  { mask: '?H', charset: '0123456789ABCDEF', description: 'Uppercase hex (0-F)' },
  { mask: '?s', charset: '«space»!"#$%&\'()*+,-./:;<=>?@[]^_`{|}~', description: 'Symbols' },
  { mask: '?a', charset: '?l?u?d?s', description: 'Lowercase, uppercase, digits and symbols' },
  { mask: '?b', charset: '0x00 - 0xFF', description: 'All possible bytes' }
]

const hashrateUnits = ['H/s', 'KH/s', 'MH/s', 'GH/s', 'TH/s']

export function hashrateStr(hashrate: number): string {
  let n = 0,
    x = hashrate

  while (x > 1000 && n < hashrateUnits.length - 1) {
    n++
    x /= 1000
  }

  return `${x.toFixed(1)} ${hashrateUnits[n]}`
}

export function makeHashcatParams(hashType: number, attackSettings: AttackSettingsT): HashcatParams {
  // --loopback is only valid in wodlist attacks where there are valid rules to be looped back through
  const enable_loopback =
    attackSettings.enableLoopback && attackSettings.attackMode == AttackMode.Dictionary && attackSettings.selectedRulefiles.length > 0

  const baseParams: HashcatParams = {
    attack_mode: attackSettings.attackMode,
    hash_type: hashType,

    mask: '',
    mask_increment: false,
    mask_increment_min: 0,
    mask_increment_max: 0,
    mask_custom_charsets: [],
    mask_sharded_charset: '',

    wordlist_filenames: [],
    rules_filenames: [],

    optimized_kernels: attackSettings.optimizedKernels,
    slow_candidates: attackSettings.slowCandidates,
    enable_loopback,

    additional_args: [],
    skip: 0,
    limit: 0
  }

  switch (attackSettings.attackMode) {
    case AttackMode.Dictionary:
      return {
        ...baseParams,
        wordlist_filenames: attackSettings.selectedWordlists.slice(0, 1),
        rules_filenames: attackSettings.selectedRulefiles
      }

    case AttackMode.Combinator:
      return {
        ...baseParams,
        wordlist_filenames: attackSettings.selectedWordlists
      }

    case AttackMode.Mask:
      return {
        ...baseParams,
        mask: attackSettings.mask,
        mask_increment: attackSettings.maskIncrement
      }

    case AttackMode.HybridDM:
    case AttackMode.HybridMD:
      return {
        ...baseParams,
        mask: attackSettings.mask,
        mask_increment: attackSettings.maskIncrement,
        wordlist_filenames: attackSettings.selectedWordlists.slice(0, 1)
      }

    default:
      return baseParams
  }
}
