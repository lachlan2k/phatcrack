// Hashcat often outputs device names like "cpu-Intel(R) Xeon(R) CPU E5-2630 v3 @ 2.40GHz", so let's tidy that up
export function formatDeviceName(deviceName: string): string {
  return deviceName
    .replace(/^cpu-/, '') // remove cpu- prefix
    .replace('Intel(R)', 'Intel')
    .replace('(R) CPU', ' CPU')
}
