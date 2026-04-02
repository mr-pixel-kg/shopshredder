import { computed, onMounted, ref } from 'vue'

import { usersApi, whitelistApi } from '@/api'

import type { CreateUserRequest, UpdateUserRequest, User } from '@/types'

export function useUsers() {
  const users = ref<User[]>([])
  const loading = ref(false)
  const initialized = ref(false)
  const error = ref<string | null>(null)
  const busyIds = ref(new Set<string>())

  const activeUsers = computed<User[]>(() => users.value.filter((user) => !user.isPending))

  const invitedUsers = computed<User[]>(() => users.value.filter((user) => user.isPending))

  async function fetch() {
    if (!initialized.value) loading.value = true
    error.value = null
    try {
      users.value = await usersApi.list()
      initialized.value = true
    } catch (e: unknown) {
      error.value = e instanceof Error ? e.message : 'Fehler beim Laden'
    } finally {
      loading.value = false
    }
  }

  async function createUser(req: CreateUserRequest): Promise<User> {
    const user = await usersApi.create(req)
    users.value.unshift(user)
    return user
  }

  async function inviteUser(req: CreateUserRequest): Promise<User> {
    const user = await whitelistApi.add(req)
    users.value.unshift(user)
    return user
  }

  async function updateUser(id: string, req: UpdateUserRequest): Promise<User> {
    const user = await usersApi.update(id, req)
    users.value = users.value.map((existing) => (existing.id === id ? user : existing))
    return user
  }

  async function deleteUser(id: string): Promise<void> {
    await usersApi.remove(id)
    users.value = users.value.filter((u) => u.id !== id)
  }

  async function deleteInvite(id: string): Promise<void> {
    await whitelistApi.remove(id)
    users.value = users.value.filter((u) => u.id !== id)
  }

  onMounted(() => {
    void fetch()
  })

  return {
    users,
    activeUsers,
    invitedUsers,
    loading,
    error,
    busyIds,
    createUser,
    inviteUser,
    updateUser,
    deleteUser,
    deleteInvite,
    refresh: fetch,
  }
}
