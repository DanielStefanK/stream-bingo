import { fetchJson, type PaginationResponse } from './api'

export const loadUsers = async (page: number, limit: number, term?: string) => {
  const resp = fetchJson<PaginationResponse<UserWithAllFields[]>>(
    `${import.meta.env.VITE_APP_SERVER_URL}/admin/user/list?page=${page}&limit=${limit}&term=${term}`
  )
  return resp
}

export const toggleUserActive = async (id: number, active: boolean) => {
  const resp = fetchJson<UserWithAllFields>(
    `${import.meta.env.VITE_APP_SERVER_URL}/admin/user/state/${id}`,
    {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ active }),
    }
  )
  return resp
}

export const deleteUser = async (id: number) => {
  const resp = fetchJson<UserWithAllFields>(
    `${import.meta.env.VITE_APP_SERVER_URL}/admin/user/delete/${id}`,
    {
      method: 'DELETE',
      headers: { 'Content-Type': 'application/json' },
    }
  )
  return resp
}

export interface UserWithAllFields {
  id: number
  name: string
  email: string
  avatarURL: string
  provider: string
  admin: boolean
  active: boolean
  createdAt: string
  updatedAt: string
  deleteAt?: string
  providerID: string
}
