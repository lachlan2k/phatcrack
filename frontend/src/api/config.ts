import { client } from '.'
import type { ConfigDTO } from './types'

export function getCurrentConfig(): Promise<ConfigDTO> {
  return client.get('/api/v1/config/current').then((res) => res.data)
}
