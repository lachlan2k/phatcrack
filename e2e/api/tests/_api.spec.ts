import * as api from '../../../frontend/src/api'
import axios, { AxiosError } from 'axios'
import tough from 'tough-cookie'
import { HttpsCookieAgent } from 'http-cookie-agent/http'

process.env.NODE_TLS_REJECT_UNAUTHORIZED = '0'

const reOk = /ok/i
const re401 = /code 401/
const re400 = /code 400/
const re40x = /code 40/

export const credsMap = {
  default: {
    username: 'admin',
    password: 'changeme'
  },

  admin: {
    username: 'admin',
    password: 'd#kdn19SeFiS0@3k'
  },

  alice: {
    username: 'alice',
    password: '84a30c9b39da5951'
  },

  bob: {
    username: 'bobby',
    password: '63e698d946dfbf01'
  },

  carol: {
    username: 'carol',
    password: '9b4b5282dc0cc8fb'
  }
}

const initialAdminPassword = 'changeme'

function setupClientWithCookieJar(_jar: tough.CookieJar | null = null) {
  const jar = _jar ?? new tough.CookieJar()

  api.setClient(
    axios.create({
      httpsAgent: new HttpsCookieAgent({
        cookies: {
          jar
        },
        rejectUnauthorized: false
      }),

      baseURL: 'https://localhost:8443'
    })
  )

  return jar
}

function beforeAllSetupClientWithCookieJar() {
  const jar = new tough.CookieJar()
  beforeAll(() => {
    setupClientWithCookieJar(jar)
  })
  return jar
}

function beforeAllSetupClientWithLogin({ username, password }: { username: string; password: string }) {
  const jar = new tough.CookieJar()

  beforeAll(async () => {
    setupClientWithCookieJar(jar)
    await api.loginWithCredentials(username, password)
  })

  return jar
}

const enrichAxiosError = (e: AxiosError) => {
  const data = e?.response?.data as any

  const newMessage = [e?.response?.status.toString(), e?.response?.statusText.toString(), data?.message?.toString()].join(': ')

  return Promise.reject(new Error(newMessage))
}

beforeAll(async () => {
  setupClientWithCookieJar()

  const wipeReq = await api.client.post('/api/v1/e2e/wipe', null, {
    headers: {
      'X-E2E-Key': 'authKeyForE2E'
    }
  })
  expect(wipeReq.status).toBe(200)
})

describe('Unauthenticated Tests', () => {
  describe('Allowed endpoints', () => {
    it('ping should return pong', async () => {
      const res = await api.ping()
      expect(res).toBe('pong')
    })

    it('should return config', async () => {
      const config = await api.getCurrentConfig()
      expect(config).not.toBeNull()
    })
  })

  describe('Disallowed Endpoints', () => {
    const t = (name: string, run: any) => ({ name, run })

    const tests = [
      t('accountChangePassword', () => api.accountChangePassword({} as any)),
      t('adminGetAllUsers', () => api.adminGetAllUsers()),
      t('adminCreateUser', () => api.adminCreateUser({} as any)),
      t('adminUpdateUser', () => api.adminUpdateUser('z', {} as any)),
      t('adminUpdateUserPassword', () => api.adminUpdateUserPassword('z', {} as any)),
      t('adminCreateServiceAccount', () => api.adminCreateServiceAccount({} as any)),
      t('adminDeleteUser', () => api.adminDeleteUser({} as any)),
      t('adminDeleteAgent', () => api.adminDeleteAgent({} as any)),
      t('adminCreateAgent', () => api.adminCreateAgent({} as any)),
      t('adminAgentSetMaintenance', () => api.adminAgentSetMaintenance('z', {} as any)),
      t('adminGetConfig', () => api.adminGetConfig()),
      t('adminSetConfig', () => api.adminSetConfig({} as any)),
      t('adminGetVersion', () => api.adminGetVersion()),
      t('getAllAgents', () => api.getAllAgents()),
      t('getAllAttackTemplates', () => api.getAllAttackTemplates()),
      t('deleteAttackTemplate', () => api.deleteAttackTemplate('z')),
      t('createAttackTemplate', () => api.createAttackTemplate({} as any)),
      t('createAttackTemplateSet', () => api.createAttackTemplateSet({} as any)),
      t('updateAttackTemplate', () => api.updateAttackTemplate('', {} as any)),
      t('startMFAEnrollment', () => api.startMFAEnrollment()),
      t('startMFAChallenge', () => api.startMFAChallenge()),
      t('changeTemporaryPassword', () => api.changeTemporaryPassword({} as any)),
      t('refreshAuth', () => api.refreshAuth()),
      t('loadHashTypes', () => api.loadHashTypes()),
      t('detectHashType', () => api.detectHashType({} as any)),
      t('getAllListfiles', () => api.getAllListfiles()),
      t('getListfilesForProject', () => api.getListfilesForProject('z')),
      t('deleteListfile', () => api.deleteListfile({} as any)),
      t('searchPotfile', () => api.searchPotfile({} as any)),
      t('createProject', () => api.createProject('', '')),
      t('deleteProject', () => api.deleteProject('z')),
      t('getAllProjects', () => api.getAllProjects()),
      t('getProject', () => api.getProject('z')),
      t('getProjectShares', () => api.getProjectShares('z')),
      t('addProjectShare', () => api.addProjectShare('z', {} as any)),
      t('deleteProjectShare', () => api.deleteProjectShare('z', {} as any)),
      t('createHashlist', () => api.createHashlist({} as any)),
      t('appendToHashlist', () => api.appendToHashlist('z', {} as any)),
      t('createAttack', () => api.createAttack({} as any)),
      t('deleteAttack', () => api.deleteAttack('z')),
      t('stopAttack', () => api.stopAttack('z')),
      t('startAttack', () => api.startAttack('z')),
      t('getHashlistsForProject', () => api.getHashlistsForProject('z')),
      t('getHashlist', () => api.getHashlist('z')),
      t('deleteHashlist', () => api.deleteHashlist('z')),
      t('getAttacksForHashlist', () => api.getAttacksForHashlist('z')),
      t('getAttacksWithJobsForHashlist', () => api.getAttacksWithJobsForHashlist('z', {} as any)),
      t('getAttacksInitialising', () => api.getAttacksInitialising()),
      t('getRunningJobs', () => api.getRunningJobs()),
      t('getJobCountPerUser', () => api.getJobCountPerUser()),
      t('getAllUsers', () => api.getAllUsers())
    ]

    for (const test of tests) {
      it('should return 401 for ' + test.name, async () => {
        await expect(test.run()).rejects.toThrow(re401)
      })
    }
  })
})

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

describe('Projects', () => {
  let alicesProjectName = 'alices project'
  let alicesProjectId = ''

  describe('Alice creating projects', () => {
    const u = credsMap.alice
    beforeAllSetupClientWithLogin(u)

    it('allows alice to create a project', async () => {
      const res = await api.createProject(alicesProjectName, '')
      expect(res.name).toBe(alicesProjectName)
      alicesProjectId = res.id

      const aliceUserData = await api.refreshAuth()
      expect(res.owner_user_id).toBe(aliceUserData.user.id)
    })

    it('doesnt allow a project with a bad name to be created', async () => {
      const badNames = ['', 'a', 'a'.repeat(1000), '!@#$%^&*()-=_+"\'']

      for (const name of badNames) {
        await expect(api.createProject(name, '')).rejects.toThrow(re400)
      }
    })
  })

  describe('Access Control', () => {
    describe('Alice Perspective', () => {
      beforeAllSetupClientWithLogin(credsMap.alice)

      it('allows alice to see her project in list of all projects', async () => {
        const res = await api.getAllProjects()
        const ids = res.projects.map(x => x.id)

        expect(ids).toContain(alicesProjectId)
      })

      it('allows alice to get her project', async () => {
        const res = await api.getProject(alicesProjectId)
        expect(res.id).toBe(alicesProjectId)
        expect(res.name).toBe(alicesProjectName)
      })
    })

    describe('Bob Perspective', () => {
      beforeAllSetupClientWithLogin(credsMap.bob)

      it('doesnt include alices project in bobs list of all projects', async () => {
        const res = await api.getAllProjects()
        const ids = res.projects.map(x => x.id)

        expect(ids).not.toContain(alicesProjectId)
      })

      it('doesnt allow bob to see alices project', async () => {
        await expect(api.getProject(alicesProjectId)).rejects.toThrow(re40x)
      })

      it('doesnt allow bob to delete alices project', async () => {
        await expect(api.deleteProject(alicesProjectId)).rejects.toThrow(re40x)
      })
    })

    describe('Admin Perspective', () => {
      beforeAllSetupClientWithLogin(credsMap.admin)

      it("allows admin to see alice's project in list of all projects", async () => {
        const res = await api.getAllProjects()
        const ids = res.projects.map(x => x.id)

        expect(ids).toContain(alicesProjectId)
      })

      it("allows admin to get alice's project", async () => {
        const res = await api.getProject(alicesProjectId)
        expect(res.id).toBe(alicesProjectId)
        expect(res.name).toBe(alicesProjectName)
      })
    })
  })
})
