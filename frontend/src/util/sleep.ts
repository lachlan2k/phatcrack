export default function sleep(ms: number) {
  return new Promise((wake) => setTimeout(wake, ms))
}
