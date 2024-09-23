import * as api from '../../../frontend/src/api'
import { dummyApiRequests } from './dummyRequests'
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

process.env.NODE_TLS_REJECT_UNAUTHORIZED = '0'

beforeAll(async () => {
  setupClientWithCookieJar()

  const wipeReq = await api.client.post('/api/v1/e2e/wipe', null, {
    headers: {
      'X-E2E-Key': 'authKeyForE2E'
    }
  })
  expect(wipeReq.status).toBe(200)
})

import './unauth'
import './setup'
import './user_provisioning_auth'

import './admin_authz'

import './projects'
import './hashlists'
