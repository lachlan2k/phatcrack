import { client } from '.'
import type { AuthLoginResponseDTO, AuthWhoamiResponseDTO } from './types'

export function login(username: string, password: string): Promise<AuthLoginResponseDTO> {
  return client
    .post('/api/v1/auth/login', {
      username,
      password
    })
    .then((res) => res.data)
}

export function refreshAuth(): Promise<AuthWhoamiResponseDTO> {
  return client.put('/api/v1/auth/refresh').then((res) => res.data)
}

export function logout(): Promise<null> {
  return client.post('/api/v1/auth/logout').then((res) => res.data)
}
