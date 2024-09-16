import type { AgentGetAllResponseDTO } from './types'
import { client } from '.'

export function getAllAgents(): Promise<AgentGetAllResponseDTO> {
  return client.get('/api/v1/agent/all').then(res => res.data)
}
