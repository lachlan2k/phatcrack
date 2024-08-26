import { client } from '.'
import type { UsersGetAllResponseDTO } from './types'

export enum UserRole {
  Standard = 'standard',
  Admin = 'admin',
  ServiceAccount = 'service_account',
  MFAEnrolled = 'mfa_enrolled',
  MFAExempt = 'mfa_exempt',
  RequiresPasswordChange = 'requires_password_change'
}

export const userRoles = Object.values(UserRole)
export const userSignupRoles = [UserRole.Standard, UserRole.Admin]
export const userAssignableRoles = [...userSignupRoles, UserRole.RequiresPasswordChange, UserRole.MFAExempt]

export function getAllUsers(): Promise<UsersGetAllResponseDTO> {
  return client.get('/api/v1/user/all').then((res) => res.data)
}
