import { client } from '.'
import type { PotfileSearchResponseDTO } from './types'

export function searchPotfile(hashes: string[]): Promise<PotfileSearchResponseDTO> {
  return client
    .post('/api/v1/potfile/search', {
      hashes
    })
    .then((res) => res.data)
}
