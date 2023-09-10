import { client } from '.'
import type {
  AdminAgentCreateRequestDTO,
  AdminAgentCreateResponseDTO,
  AdminAgentSetMaintanceRequestDTO,
  AdminConfigRequestDTO,
  AdminConfigResponseDTO,
  AdminGetAllUsersResponseDTO,
  AdminServiceAccountCreateRequestDTO,
  AdminServiceAccountCreateResponseDTO,
  AdminUserCreateRequestDTO,
  AdminUserCreateResponseDTO
} from './types'

export function adminGetAllUsers(): Promise<AdminGetAllUsersResponseDTO> {
  return client.get('/api/v1/admin/user/all').then((res) => res.data)
}

export function adminCreateUser(newUserData: AdminUserCreateRequestDTO): Promise<AdminUserCreateResponseDTO> {
  return client.post('/api/v1/admin/user/create', newUserData).then((res) => res.data)
}

export function adminCreateServiceAccount(
  newAccountData: AdminServiceAccountCreateRequestDTO
): Promise<AdminServiceAccountCreateResponseDTO> {
  return client.post('/api/v1/admin/user/create-service-account', newAccountData).then((res) => res.data)
}

export function adminDeleteUser(id: string): Promise<string> {
  return client.delete('/api/v1/admin/user/' + id).then((res) => res.data)
}

export function adminDeleteAgent(id: string): Promise<string> {
  return client.delete('/api/v1/admin/agent/' + id).then((res) => res.data)
}

export function adminCreateAgent(newAgentData: AdminAgentCreateRequestDTO): Promise<AdminAgentCreateResponseDTO> {
  return client.post('/api/v1/admin/agent/create', newAgentData).then((res) => res.data)
}

export function adminAgentSetMaintenance(id: string, body: AdminAgentSetMaintanceRequestDTO): Promise<string> {
  return client.put(`/api/v1/admin/agent/${id}/set-maintenance-mode`, body).then((res) => res.data)
}

export function adminGetConfig(): Promise<AdminConfigResponseDTO> {
  return client.get('/api/v1/admin/config').then((res) => res.data)
}

export function adminSetConfig(newConfig: AdminConfigRequestDTO): Promise<AdminConfigRequestDTO> {
  return client.put('/api/v1/admin/config', newConfig).then((res) => res.data)
}
