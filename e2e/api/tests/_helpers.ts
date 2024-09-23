import axios, { AxiosError } from 'axios'
import * as tough from 'tough-cookie'
import { HttpsCookieAgent } from 'http-cookie-agent/http'
import * as api from '../../../frontend/src/api'

export const reOk = /ok/i
export const re401 = /code 401/
export const re400 = /code 400/
export const re40x = /code 40/

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

export const initialAdminPassword = 'changeme'

export function setupClientWithCookieJar(_jar: tough.CookieJar | null = null) {
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

export function beforeAllSetupClientWithCookieJar() {
  const jar = new tough.CookieJar()
  beforeAll(() => {
    setupClientWithCookieJar(jar)
  })
  return jar
}

export function beforeAllSetupClientWithLogin({ username, password }: { username: string; password: string }) {
  const jar = new tough.CookieJar()

  beforeAll(async () => {
    setupClientWithCookieJar(jar)
    await api.loginWithCredentials(username, password)
  })

  return jar
}

export const enrichAxiosError = (e: AxiosError) => {
  const data = e?.response?.data as any

  const newMessage = [e?.response?.status.toString(), e?.response?.statusText.toString(), data?.message?.toString()].join(': ')

  return Promise.reject(new Error(newMessage))
}
