import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { sandboxesApi } from '@/api'
import type { Sandbox, CreateSandboxRequest } from '@/types'

export const useSandboxesStore = defineStore('sandboxes', () => {
  const sandboxes = ref<Sandbox[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  const activeSandboxes = computed(() =>
    sandboxes.value.filter((s) => s.status === 'running' || s.status === 'starting'),
  )

  const recentSandboxes = computed(() =>
    sandboxes.value.filter(
      (s) => s.status === 'stopped' || s.status === 'expired' || s.status === 'failed' || s.status === 'deleted',
    ),
  )

  async function fetchMySandboxes() {
    loading.value = true
    error.value = null
    try {
      sandboxes.value = await sandboxesApi.listMine()
    } catch (e: unknown) {
      error.value = e instanceof Error ? e.message : 'Fehler beim Laden'
    } finally {
      loading.value = false
    }
  }

  async function fetchAllSandboxes() {
    loading.value = true
    error.value = null
    try {
      sandboxes.value = await sandboxesApi.list()
    } catch (e: unknown) {
      error.value = e instanceof Error ? e.message : 'Fehler beim Laden'
    } finally {
      loading.value = false
    }
  }

  async function fetchGuestSandboxes() {
    loading.value = true
    error.value = null
    try {
      sandboxes.value = await sandboxesApi.listGuest()
    } catch (e: unknown) {
      error.value = e instanceof Error ? e.message : 'Fehler beim Laden'
    } finally {
      loading.value = false
    }
  }

  async function createSandbox(req: CreateSandboxRequest): Promise<Sandbox> {
    const sandbox = await sandboxesApi.create(req)
    sandboxes.value.unshift(sandbox)
    return sandbox
  }

  async function createPublicDemo(req: CreateSandboxRequest): Promise<Sandbox> {
    const sandbox = await sandboxesApi.createPublicDemo(req)
    sandboxes.value.unshift(sandbox)
    return sandbox
  }

  async function deleteSandbox(id: string) {
    await sandboxesApi.remove(id)
    sandboxes.value = sandboxes.value.filter((s) => s.id !== id)
  }

  return {
    sandboxes,
    loading,
    error,
    activeSandboxes,
    recentSandboxes,
    fetchMySandboxes,
    fetchAllSandboxes,
    fetchGuestSandboxes,
    createSandbox,
    createPublicDemo,
    deleteSandbox,
  }
})
