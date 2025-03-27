import type { User } from '../store/userStore'
import { fetchJson } from './api'

export const providers = [
  // { name: 'twitch', id: 'twitch' },
  // { name: 'google', id: 'google' },
  { name: 'github', id: 'github' },
]

export const getLoginUrl = (provider: string) => {
  return `${import.meta.env.VITE_APP_SERVER_URL}/auth/oauth/${provider}`
}

export const loginLocal = async (
  email: string,
  password: string
): Promise<LoginResponse> => {
  return fetchJson<LoginResponse>(
    `${import.meta.env.VITE_APP_SERVER_URL}/auth/login`,
    {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password }),
    }
  ).then((res) => {
    if (res.token) {
      localStorage.setItem('token', res.token)
    }
    return res
  })
}

export const loginOAuth = async (
  provider: string,
  code: string
): Promise<LoginResponse> => {
  return fetchJson<LoginResponse>(
    `${import.meta.env.VITE_APP_SERVER_URL}/auth/oauth/${provider}/callback?code=${code}`
  ).then((res) => {
    if (res.token) {
      localStorage.setItem('token', res.token)
    }
    return res
  })
}

export const register = async (
  name: string,
  email: string,
  password: string
): Promise<LoginResponse> => {
  return fetchJson<LoginResponse>(
    `${import.meta.env.VITE_APP_SERVER_URL}/auth/register`,
    {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name, email, password }),
    }
  )
}

export const fetchUser = async (): Promise<User> => {
  return fetchJson<User>(`${import.meta.env.VITE_APP_SERVER_URL}/auth/me`)
}

interface LoginResponse {
  token: string
  user: User
}
