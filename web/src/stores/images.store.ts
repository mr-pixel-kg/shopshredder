import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { imagesApi } from '@/api'
import type { Image, CreateImageRequest } from '@/types'

export const useImagesStore = defineStore('images', () => {
  const images = ref<Image[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)
  const sseConnections = new Map<string, EventSource>()

  const publicImages = computed(() => images.value.filter((i) => i.isPublic))

  function mergeImages(fetched: Image[]) {
    const fetchedMap = new Map(fetched.map((img) => [img.id, img]))
    for (let i = images.value.length - 1; i >= 0; i--) {
      const existing = images.value[i]
      const updated = fetchedMap.get(existing.id)
      if (updated) {
        Object.assign(existing, updated)
        fetchedMap.delete(existing.id)
      } else {
        images.value.splice(i, 1)
      }
    }

    for (const img of fetched) {
      if (fetchedMap.has(img.id)) {
        images.value.push(img)
      }
    }
  }

  function subscribePullProgress(imageId: string) {
    if (sseConnections.has(imageId)) return

    const baseURL = import.meta.env.WEB_API_URL || ''
    const es = new EventSource(`${baseURL}/api/images/${imageId}/progress`)
    sseConnections.set(imageId, es)

    es.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data) as { percent?: number; status?: string; error?: string }
        const image = images.value.find((i) => i.id === imageId)
        if (!image) {
          closeSseConnection(imageId)
          return
        }

        const pct = data.percent ?? 0
        if (pct > (image.pullProgress ?? 0)) {
          image.pullProgress = pct
        }

        if (data.status === 'complete' || data.status === 'ready') {
          image.status = 'ready'
          image.pullProgress = 100
          closeSseConnection(imageId)
        } else if (data.status === 'failed') {
          image.status = 'failed'
          image.errorMessage = data.error || 'Pull fehlgeschlagen'
          closeSseConnection(imageId)
        }
      } catch {
        // ignore parse errors
      }
    }

    es.onerror = () => {
      closeSseConnection(imageId)
    }
  }

  function closeSseConnection(imageId: string) {
    const es = sseConnections.get(imageId)
    if (es) {
      es.close()
      sseConnections.delete(imageId)
    }
  }

  function subscribeAllPulling() {
    for (const image of images.value) {
      if (image.status === 'pulling') {
        subscribePullProgress(image.id)
      }
    }
  }

  function unsubscribeAll() {
    for (const [id, es] of sseConnections) {
      es.close()
      sseConnections.delete(id)
    }
  }

  function $reset() {
    unsubscribeAll()
    images.value = []
    loading.value = true
    error.value = null
  }

  async function fetchPublicImages() {
    const isInitial = images.value.length === 0
    if (isInitial) loading.value = true
    error.value = null
    try {
      mergeImages(await imagesApi.listPublic())
      subscribeAllPulling()
    } catch (e: unknown) {
      error.value = e instanceof Error ? e.message : 'Fehler beim Laden'
    } finally {
      loading.value = false
    }
  }

  async function fetchAllImages() {
    const isInitial = images.value.length === 0
    if (isInitial) loading.value = true
    error.value = null
    try {
      mergeImages(await imagesApi.listAll())
      subscribeAllPulling()
    } catch (e: unknown) {
      error.value = e instanceof Error ? e.message : 'Fehler beim Laden'
    } finally {
      loading.value = false
    }
  }

  async function createImage(req: CreateImageRequest): Promise<Image> {
    const image = await imagesApi.create(req)
    images.value.unshift(image)
    if (image.status === 'pulling') {
      subscribePullProgress(image.id)
    }
    return image
  }

  async function deleteImage(id: string) {
    closeSseConnection(id)
    await imagesApi.remove(id)
    images.value = images.value.filter((i) => i.id !== id)
  }

  return {
    images,
    loading,
    error,
    publicImages,
    fetchPublicImages,
    fetchAllImages,
    createImage,
    deleteImage,
    unsubscribeAll,
    $reset,
  }
})
