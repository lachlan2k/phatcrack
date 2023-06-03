export const attackModes = [
  { name: 'Wordlist', value: 0 },
  { name: 'Combinator', value: 1 },
  { name: 'Brute-force/Mask', value: 3 },
  {
    name: 'Wordlist + Mask',
    value: 6
  },
  {
    name: 'Mask + Wordlist',
    value: 7
  }
]

export function getAttackModeName(id: number): string {
  return attackModes.find((x) => x.value == id)?.name ?? ''
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
