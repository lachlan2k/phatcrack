import { client } from '.'
import type { DetectHashTypeRequestDTO, DetectHashTypeResponseDTO, HashTypesDTO } from './types'

export function loadHashTypes(): Promise<HashTypesDTO> {
  return client.get('/api/v1/hashcat/hashtypes').then((res) => res.data)
}

const detectMemoMap = new Map<string, DetectHashTypeResponseDTO>()

export async function detectHashType(exampleHash: string): Promise<DetectHashTypeResponseDTO> {
  if (detectMemoMap.has(exampleHash)) {
    return detectMemoMap.get(exampleHash) as DetectHashTypeResponseDTO
  }

  const results = await client
    .post('/api/v1/hashcat/detect-hashtype', {
      test_hash: exampleHash
    } as DetectHashTypeRequestDTO)
    .then((res) => res.data as DetectHashTypeResponseDTO)

  detectMemoMap.set(exampleHash, results)

  return results
}
