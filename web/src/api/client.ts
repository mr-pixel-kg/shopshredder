import axios from 'axios'

import { getToken } from '@/utils/storage'

import type { AxiosInstance } from 'axios'

const CLIENT_ID_KEY = 'client_id'

function getOrCreateClientId(): string {
  let id = localStorage.getItem(CLIENT_ID_KEY)
  if (!id) {
    id = crypto.randomUUID()
    localStorage.setItem(CLIENT_ID_KEY, id)
  }
  return id
}

export function getClientId(): string {
  return getOrCreateClientId()
}

export const AUTH_LOGOUT_EVENT = 'auth:session-expired'

export const apiClient: AxiosInstance = axios.create({
  baseURL: import.meta.env.WEB_API_URL,
  timeout: 30_000,
  headers: { 'Content-Type': 'application/json' },
})

apiClient.interceptors.request.use((config) => {
  const token = getToken()
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }

  config.headers['X-Client-Id'] = getOrCreateClientId()

  return config
})

apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (axios.isAxiosError(error) && error.response?.status === 401) {
      const url = error.config?.url ?? ''
      if (!url.includes('/auth/login') && !url.includes('/auth/register')) {
        window.dispatchEvent(new CustomEvent(AUTH_LOGOUT_EVENT))
      }
    }
    return Promise.reject(error)
  },
)
