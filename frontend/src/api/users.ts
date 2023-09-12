import { client } from '.'
import type { UsersGetAllResponseDTO } from './types'

export function getAllUsers(): Promise<UsersGetAllResponseDTO> {
  return client.get('/api/v1/user/all').then((res) => res.data)
}
