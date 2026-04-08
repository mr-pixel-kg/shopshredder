import { apiClient } from './client'

import type { AddWhitelistRequest, User } from '@/types'

export const whitelistApi = {
  async list(): Promise<User[]> {
    const { data } = await apiClient.get<User[]>('/api/whitelist')
    return data
  },

  async add(req: AddWhitelistRequest): Promise<User> {
    const { data } = await apiClient.post<User>('/api/whitelist', req)
    return data
  },

  async remove(id: string): Promise<void> {
    await apiClient.delete(`/api/whitelist/${id}`)
  },
}
