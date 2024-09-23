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
import { dummyApiRequests } from './dummyRequests'

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
    for (const test of dummyApiRequests) {
      it('should return 401 for ' + test.name, async () => {
        await expect(test.run()).rejects.toThrow(re401)
      })
    }
  })
})
