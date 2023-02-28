import { client } from '.'

export type APILoggedInUserDetailsT = {
    username: string
    id: string
    role: string
}

export type APILoginResponse = {
    user: APILoggedInUserDetailsT | null
}

export async function login(username: string, password: string) : Promise<APILoginResponse | null> {
    const res = await client.post('/api/v1/auth/login', {
        username, password
    })
    
    return res.data
}

export function refreshAuth(): Promise<APILoginResponse | null> {
    // return client.put('/api/v1/auth/refresh')
    return client.get('/api/v1/auth/whoami')
}

export function logout() : Promise<null> {
    return client.post('/api/v1/auth/logout')
}