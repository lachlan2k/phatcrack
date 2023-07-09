import { client } from '.'
import type { AgentGetAllResponseDTO } from './types'

export function getAllAgents(): Promise<AgentGetAllResponseDTO> {
  return client.get('/api/v1/agent/all').then((res) => res.data)
}
