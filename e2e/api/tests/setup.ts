import * as api from '../../../frontend/src/api'
import {
  beforeAllSetupClientWithCookieJar,
  beforeAllSetupClientWithLogin,
  enrichAxiosError,
  setupClientWithCookieJar,
  initialAdminPassword,
  credsMap,
  re400,
  re401,
  re40x,
  reOk
} from './_helpers'

describe('First Login Flow', () => {
  const cookieJar = beforeAllSetupClientWithCookieJar()

  it('should deny log in with incorrect credentials', async () => {
    await expect(api.loginWithCredentials('admin', 'wrongpassword')).rejects.toThrow(re401)
  })

  it('should allow log in with admin:chanageme', async () => {
    const res = await api.loginWithCredentials('admin', initialAdminPassword)
    await expect(res.user.username).toBe('admin')
  })

  it('should require password to be changed', async () => {
    const res = await api.loginWithCredentials('admin', initialAdminPassword)
    await expect(res.requires_password_change).toBe(true)
  })

  it('should deny requests before password is changed', async () => {
    await expect(api.getAllListfiles()).rejects.toThrow(re401)
  })

  it('should not allow password to be changed with standard endpoint', async () => {
    await expect(
      api.accountChangePassword({
        new_password: 'newpassword',
        current_password: initialAdminPassword
      })
    ).rejects.toThrow(re401)
  })

  it('should not allow password to be changed with wrong current password', async () => {
    await expect(
      api
        .changeTemporaryPassword({
          new_password: 'password1234123412341234',
          old_password: 'wrongpassword'
        })
        .catch(enrichAxiosError)
    ).rejects.toThrow(/old password was incorrect/i)
  })

  it('should not allow password to be changed to a weak password', async () => {
    await expect(
      api.changeTemporaryPassword({
        new_password: 'password1',
        old_password: initialAdminPassword
      })
    ).rejects.toThrow(re400)
  })

  it('should allow password to be changed to a strong password', async () => {
    const res = await api.changeTemporaryPassword({
      new_password: credsMap.admin.password,
      old_password: initialAdminPassword
    })
    await expect(res).toMatch(reOk)
  })

  it('should allow log in with new password', async () => {
    cookieJar!.removeAllCookiesSync()

    const res = await api.loginWithCredentials('admin', credsMap.admin.password)
    await expect(res.user.username).toBe('admin')
  })

  it('should not require password to be changed', async () => {
    const res = await api.loginWithCredentials('admin', credsMap.admin.password)
    await expect(res.requires_password_change).toBe(false)
  })

  it('should allow requests after password is changed', async () => {
    await expect(api.getAllListfiles()).resolves.not.toThrow()
  })
})
