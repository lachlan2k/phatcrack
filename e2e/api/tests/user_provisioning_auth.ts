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

describe('User Provisioning & Credential Authentication', () => {
  describe('Admin Creating Users', () => {
    beforeAllSetupClientWithLogin(credsMap.admin)

    it('should be logged in as admin', async () => {
      const res = await api.refreshAuth()
      expect(res.user.username).toBe(credsMap.admin.username)
    })

    it('allows password change to be disabled', async () => {
      const config = await api.adminGetConfig()

      await api.adminSetConfig({
        auth: {
          general: {
            ...config.auth.general,
            require_password_change_on_first_login: false,
            enabled_methods: config.auth.general?.enabled_methods ?? [],
            is_mfa_required: config.auth.general?.is_mfa_required ?? false
          }
        }
      })

      const updatedConfig = await api.adminGetConfig()
      expect(updatedConfig.auth.general?.require_password_change_on_first_login).toBe(false)
    })

    it('allows alice user to be created with password (no change required)', async () => {
      const u = credsMap.alice

      const res = await api.adminCreateUser({
        gen_password: false,
        lock_password: false,
        password: u.password,
        username: u.username,
        roles: [api.UserRole.Standard]
      })

      expect(res.username).toBe(u.username)
      expect(res.roles).not.toContain(api.UserRole.RequiresPasswordChange)
    })

    it('allows password change to be enabled', async () => {
      const config = await api.adminGetConfig()

      await api.adminSetConfig({
        auth: {
          general: {
            ...config.auth.general,
            require_password_change_on_first_login: true,
            enabled_methods: config.auth.general?.enabled_methods ?? [],
            is_mfa_required: config.auth.general?.is_mfa_required ?? false
          }
        }
      })

      const updatedConfig = await api.adminGetConfig()
      expect(updatedConfig.auth.general?.require_password_change_on_first_login).toBe(true)
    })

    it('allows bob user to be created with password (change REQUIRED)', async () => {
      const u = credsMap.bob

      const res = await api.adminCreateUser({
        gen_password: false,
        lock_password: false,
        password: u.password,
        username: u.username,
        roles: [api.UserRole.Standard]
      })

      expect(res.username).toBe(u.username)
      expect(res.roles).toContain(api.UserRole.RequiresPasswordChange)
    })

    it('doesnt allow bob to be created a second time', async () => {
      const u = credsMap.bob

      await expect(
        api.adminCreateUser({
          gen_password: false,
          lock_password: false,
          password: u.password,
          username: u.username,
          roles: [api.UserRole.Standard]
        })
      ).rejects.toThrow(re40x)
    })

    it('doesnt allow a permutation of bobs username to be created', async () => {
      const u = credsMap.bob

      const startingUserCount = (await api.getAllUsers()).users.length
      expect(startingUserCount).toBeGreaterThan(0)

      const usernamesToTest = ['bOb', 'BOB', ' bob', 'bob ', ' BoB ']

      for (const username of usernamesToTest) {
        await expect(
          api.adminCreateUser({
            gen_password: false,
            lock_password: false,
            password: u.password,
            username: username,
            roles: [api.UserRole.Standard]
          })
        ).rejects.toThrow(re40x)
      }

      const finalUserCount = (await api.getAllUsers()).users.length
      expect(startingUserCount).toBe(finalUserCount)
    })
  })

  describe('User enumeration', () => {
    it('returns the same message whether username or password is invalid', async () => {
      const resInvalidUser = await api
        .loginWithCredentials('invaliduser', 'asdf')
        .catch(enrichAxiosError)
        .catch(x => x)
      const resInvalidPass = await api
        .loginWithCredentials(credsMap.admin.username, 'asdf')
        .catch(enrichAxiosError)
        .catch(x => x)

      const a = resInvalidUser.toString()
      const b = resInvalidPass.toString()

      expect(a).toMatch(/invalid credentials/i)
      expect(a).toBe(b)
    })
  })

  describe('Logging in as alice', () => {
    beforeAllSetupClientWithCookieJar()
    const u = credsMap.alice

    it('does not allow alice to log in with blank password', async () => {
      await expect(api.loginWithCredentials(u.username, '')).rejects.toThrow(re40x)
      await expect(api.refreshAuth()).rejects.toThrow(re401) // not logged in
    })

    it('does not allow alice to log in with incorrect password', async () => {
      await expect(api.loginWithCredentials(u.username, 'asdfasdfasdf').catch(enrichAxiosError)).rejects.toThrow(/invalid credentials/i)
      await expect(api.refreshAuth()).rejects.toThrow(re401) // not logged in
    })

    it('allows alice to log in without changing password', async () => {
      const res = await api.loginWithCredentials(u.username, u.password)
      expect(res.requires_password_change).toBe(false)
    })

    it('allows alice to hit endpoints after logging in', async () => {
      const res = await api.getAllProjects()
      expect(res).not.toBeNull()
      expect(res.projects).not.toBeNull()
      expect(Array.isArray(res.projects)).toBe(true)
    })
  })

  describe('Logging in as bob', () => {
    const u = credsMap.bob
    beforeAllSetupClientWithCookieJar()

    it('does not allow bob to log in with blank password', async () => {
      await expect(api.loginWithCredentials(u.username, '')).rejects.toThrow(re40x)
      await expect(api.refreshAuth()).rejects.toThrow(re401) // not logged in
    })

    it('does not allow bob to log in with incorrect password', async () => {
      await expect(api.loginWithCredentials(u.username, 'asdfasdfasdf').catch(enrichAxiosError)).rejects.toThrow(/invalid credentials/i)
      await expect(api.refreshAuth()).rejects.toThrow(re401) // not logged in
    })

    it('allows bob to log in, but requires password change', async () => {
      const res = await api.loginWithCredentials(u.username, u.password)
      expect(res.requires_password_change).toBe(true)
    })

    it('does not allow bob to hit endpoints after logging in without changing password', async () => {
      await expect(api.getAllProjects()).rejects.toThrow(re401)
    })

    it('does not allow bob to set the same password', async () => {
      await expect(
        api.changeTemporaryPassword({
          new_password: u.password,
          old_password: u.password
        })
      ).rejects.toThrow(re400)
    })

    it('does not allow bob to set a weak password', async () => {
      await expect(
        api.changeTemporaryPassword({
          new_password: 'a',
          old_password: u.password
        })
      ).rejects.toThrow(re400)
    })

    it('does not allow bob to set a password if old password is incorrect', async () => {
      await expect(
        api.changeTemporaryPassword({
          new_password: 'correct horse battery staple',
          old_password: 'incorrectpasswordhere'
        })
      ).rejects.toThrow(re400)
    })

    it('allows bob to set a new password', async () => {
      const res = await api.changeTemporaryPassword({
        new_password: u.password + '!',
        old_password: u.password
      })
      expect(res).toMatch(reOk)
    })

    it('allows bob to hit endpoints after changing password', async () => {
      const res = await api.getAllProjects()
      expect(res).not.toBeNull()
      expect(res.projects).not.toBeNull()
      expect(Array.isArray(res.projects)).toBe(true)
    })

    it('allows bob to change password', async () => {
      const res = await api.accountChangePassword({
        current_password: u.password + '!',
        new_password: u.password
      })
      expect(res).toMatch(reOk)
    })
  })

  describe('Logging in fresh as bob', () => {
    beforeAllSetupClientWithCookieJar()
    const u = credsMap.bob

    it('allows bob to log in and doesnt require password change', async () => {
      const res = await api.loginWithCredentials(u.username, u.password)
      expect(res.requires_password_change).toBe(false)
    })

    it('allows bob to hit endpoints after logging in', async () => {
      const res = await api.getAllProjects()
      expect(res).not.toBeNull()
      expect(res.projects).not.toBeNull()
      expect(Array.isArray(res.projects)).toBe(true)
    })
  })

  describe('Logout', () => {
    const u = credsMap.alice
    beforeAllSetupClientWithLogin(u)

    it('allows refreshing auth before logout', async () => {
      const res = await api.refreshAuth()
      expect(res.user.username).toBe(u.username)
    })

    it('allows logging out', async () => {
      const res = await api.logout()
      expect(res).toBe('Goodbye')
    })

    it('doesnt allow auth to be refreshed after logout', async () => {
      await expect(api.refreshAuth()).rejects.toThrow(re40x)
    })

    it('doesnt allow fetching endpoints after logging out', async () => {
      await expect(api.getAllProjects()).rejects.toThrow(re401)
    })
  })
})
