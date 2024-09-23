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

describe('Admin authz', () => {
  const u = credsMap.alice
  beforeAllSetupClientWithLogin(u)

  const adminEndpionts = dummyApiRequests.filter(x => x.name.toLowerCase().includes('admin'))
  for (const test of adminEndpionts) {
    it('should return access 401 for alice accessing ' + test.name, async () => {
      await expect(test.run()).rejects.toThrow(re401)
    })
  }
})
