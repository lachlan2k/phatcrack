import type { PublicConfigDTO } from './types'

import { client } from '.'

export const AuthMethodCredentials = 'method_credentials'
export const AuthMethodOIDC = 'method_oidc'

export function getCurrentConfig(): Promise<PublicConfigDTO> {
  return client.get('/api/v1/config/public').then(res => res.data)
}
