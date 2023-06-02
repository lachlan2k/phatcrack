export default function decodeHex(input: string): string {
  let out = ''

  for (let i = 0; i < input.length; i += 2) {
    out += String.fromCharCode(parseInt(input.substring(i, i + 2), 16))
  }

  return out
}
