import { getHashlist } from '@/api/project'
import type { HashlistDTO, HashlistHashDTO } from '@/api/types'

export enum ExportFormat {
  JSON,
  CSV,
  ColonSeparated
}

const hexToString = (hex: string) => {
  let str = ''
  for (let i = 0; i < hex.length; i += 2) {
    const hexValue = hex.substr(i, 2)
    const decimalValue = parseInt(hexValue, 16)
    str += String.fromCharCode(decimalValue)
  }
  return str
}

function getExportBlob(hashes: HashlistHashDTO[], format: ExportFormat, hasUsernames: boolean) {
  switch (format) {
    case ExportFormat.JSON:
      return new Blob([JSON.stringify(hashes)], { type: 'application/json' })
    case ExportFormat.CSV: {
      // Poor man's CSV, vulnerable to csv injection etc, but bleh
      // TODO: something proper csv generation
      const stringified = 
        hasUsernames ?
          [['username', 'hash', 'plaintext']]
            .concat(hashes.map((x) => [x.username, x.input_hash, hexToString(x.plaintext_hex)]))
            .map((x) => x.join(','))
            .join(',') :
          [['hash', 'plaintext']]
            .concat(hashes.map((x) => [x.input_hash, hexToString(x.plaintext_hex)]))
            .map((x) => x.join(','))
            .join(',')



      return new Blob([stringified], { type: 'text/csv' })
    }

    case ExportFormat.ColonSeparated: {
      const textBlob = hashes.map((x) =>
        hasUsernames ?  `${x.username}:${x.input_hash}:${hexToString(x.plaintext_hex)}` :  `${x.input_hash}:${hexToString(x.plaintext_hex)}`
      ).join('\n')
      return new Blob([textBlob], { type: 'text/plain' })
    }
  }
}

function getExportFilename(hashlist: HashlistDTO, format: ExportFormat) {
  switch (format) {
    case ExportFormat.JSON:
      return `${hashlist.name}.json`
    case ExportFormat.CSV:
      return `${hashlist.name}.csv`
    case ExportFormat.ColonSeparated:
      return `${hashlist.name}.txt`
  }
}

export async function exportResults(hashlistId: string, format: ExportFormat, crackedOnly: boolean) {
  const hashlist = await getHashlist(hashlistId)
  const filtered = crackedOnly ? hashlist.hashes.filter((x) => x.is_cracked) : hashlist.hashes

  const blob = getExportBlob(filtered, format, hashlist.has_usernames)
  const url = URL.createObjectURL(blob)

  const a = document.createElement('a')
  a.href = url
  a.download = getExportFilename(hashlist, format)
  a.click()
}
