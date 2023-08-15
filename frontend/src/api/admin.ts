import { client } from '.'
import type {
  AdminConfigRequestDTO,
  AdminConfigResponseDTO,
  AdminGetAllUsersResponseDTO,
  AdminUserCreateRequestDTO,
  AdminUserCreateResponseDTO
} from './types'

export function adminGetAllUsers(): Promise<AdminGetAllUsersResponseDTO> {
  return client.get('/api/v1/admin/user/all').then((res) => res.data)
}

export function adminCreateUser(
  newUserData: AdminUserCreateRequestDTO
): Promise<AdminUserCreateResponseDTO> {
  return client.post('/api/v1/admin/user/create', newUserData).then((res) => res.data)
}

export function adminGetConfig(): Promise<AdminConfigResponseDTO> {
  return client.get('/api/v1/admin/config').then((res) => res.data)
}

export function adminSetConfig(newConfig: AdminConfigRequestDTO): Promise<AdminConfigRequestDTO> {
  return client.put('/api/v1/admin/config', newConfig).then((res) => res.data)
}