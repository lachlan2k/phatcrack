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
    return attackModes.find(x => x.value == id)?.name ?? ''
}