export function pluralise(count: number, singular: string, plural: string = singular + 's', zero: string = plural): string {
  if (count === 0) {
    return 'No ' + zero
  }
  return count + ' ' + (count === 1 ? singular : plural)
}
