import { client } from '.'
import type { DetectHashTypeRequestDTO, DetectHashTypeResponseDTO, HashTypesDTO } from './types'

export function loadHashTypes(): Promise<HashTypesDTO> {
  return client.get('/api/v1/resources/hashtypes').then((res) => res.data)
}

export function detectHashType(exampleHash: string): Promise<DetectHashTypeResponseDTO> {
  return client
    .post('/api/v1/resources/detect_hashtype', {
      test_hash: exampleHash
    } as DetectHashTypeRequestDTO)
    .then((res) => res.data)
}
