import { client } from '.'
import type { AccountChangePasswordRequestDTO } from './types'

export function accountChangePassword(body: AccountChangePasswordRequestDTO): Promise<string> {
  return client.put('/api/v1/account/change-password', body).then(res => res.data)
}
