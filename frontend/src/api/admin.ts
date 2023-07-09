import { client } from '.'
import type {
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
