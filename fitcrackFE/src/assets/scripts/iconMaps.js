export function jobIcon (status) {
  const icons = {
    ready: 'mdi-play-circle-outline',
    running: 'mdi-dots-vertical-circle-outline',
    finished: 'mdi-checkbox-marked-circle-outline',
    exhausted: 'mdi-close-circle-outline',
    malformed: 'mdi-alert-circle-outline',
    finishing: 'mdi-progress-clock'
  }
  return icons[status] || 'mdi-help-circle-outline'
}

export function typeIcon (type) {
  const icons = {
    info: 'mdi-information-outline',
    success: 'mdi-checkbox-marked-circle-outline',
    error: 'mdi-close-circle-outline'
  }
  return icons[type] || 'mdi-help-circle-outline'
}

export function attackIcon (handler) {
  switch (handler) {
    case 'dictionary':
      return 'mdi-dictionary'
    case 'combinator':
      return 'mdi-vector-combine'
    case 'maskattack':
    case 'mask':
      return 'mdi-nuke'
    case 'hybridWordlistMask':
    case 'hybrid (Wordlist + Mask)':
      return 'mdi-vector-difference-ba'
    case 'hybridMaskWordlist':
    case 'hybrid (Mask + Wordlist)':
      return 'mdi-vector-difference-ab'
    case 'pcfgAttack':
    case 'pcfg':
      return 'mdi-ray-start-end'
    case 'princeAttack':
    case 'prince':
      return 'mdi-crown'
    default:
      return 'mdi-checkbox-blank-outline'
  }
}