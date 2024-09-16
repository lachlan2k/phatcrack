import type {
  AuthLoginResponseDTO,
  AuthWhoamiResponseDTO,
  AuthChangePasswordRequestDTO,
  AuthWebAuthnStartEnrollmentResponseDTO,
  AuthWebAuthnStartChallengeResponseDTO
} from './types'
import { client } from '.'

export function loginWithCredentials(username: string, password: string): Promise<AuthLoginResponseDTO> {
  return client
    .post('/api/v1/auth/login/credentials', {
      username,
      password
    })
    .then(res => res.data)
}

export function loginWithOIDCCallback(querystring: string): Promise<AuthLoginResponseDTO> {
  return client.post('/api/v1/auth/login/oidc/callback' + querystring).then(res => res.data)
}

function urlSafeB64Encode(value: ArrayBuffer) {
  return btoa(String.fromCharCode.apply(null, new Uint8Array(value) as unknown as number[]))
    .replace(/\+/g, '-')
    .replace(/\//g, '_')
    .replace(/=/g, '')
}

export function startMFAEnrollment(): Promise<AuthWebAuthnStartEnrollmentResponseDTO> {
  return client.post('/api/v1/auth/mfa/start-enrollment?method=MFATypeWebAuthn').then(res => res.data)
}

export function finishMFAEnrollment(cred: PublicKeyCredential): Promise<string> {
  const attestationResponse = cred.response as AuthenticatorAttestationResponse

  return client
    .post('/api/v1/auth/mfa/finish-enrollment?method=MFATypeWebAuthn', {
      id: cred.id,
      rawId: urlSafeB64Encode(cred.rawId),
      type: cred.type,
      response: {
        clientDataJSON: urlSafeB64Encode(cred.response.clientDataJSON),
        attestationObject: urlSafeB64Encode(attestationResponse.attestationObject)
      }
    })
    .then(res => res.data)
}

export function startMFAChallenge(): Promise<AuthWebAuthnStartChallengeResponseDTO> {
  return client.post('/api/v1/auth/mfa/start-challenge?method=MFATypeWebAuthn').then(res => res.data)
}

export function finishMFAChallenge(cred: PublicKeyCredential): Promise<string> {
  const assertion = cred.response as AuthenticatorAssertionResponse

  return client
    .post('/api/v1/auth/mfa/finish-challenge?method=MFATypeWebAuthn', {
      id: cred.id,
      rawId: urlSafeB64Encode(cred.rawId),
      type: cred.type,
      response: {
        authenticatorData: urlSafeB64Encode(assertion.authenticatorData),
        clientDataJSON: urlSafeB64Encode(assertion.clientDataJSON),
        signature: urlSafeB64Encode(assertion.signature),
        userHandle: urlSafeB64Encode(assertion.userHandle as ArrayBuffer)
      }
    })
    .then(res => res.data)
}

export function changeTemporaryPassword(body: AuthChangePasswordRequestDTO): Promise<string> {
  return client.post('/api/v1/auth/change-temporary-password', body).then(res => res.data)
}

export function refreshAuth(): Promise<AuthWhoamiResponseDTO> {
  return client.put('/api/v1/auth/refresh').then(res => res.data)
}

export function logout(): Promise<null> {
  return client.post('/api/v1/auth/logout').then(res => res.data)
}
