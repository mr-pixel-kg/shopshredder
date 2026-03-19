import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { imagesApi } from '@/api'
import type { Image, CreateImageRequest } from '@/types'

export const useImagesStore = defineStore('images', () => {
  const images = ref<Image[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  const publicImages = computed(() => images.value.filter((i) => i.isPublic))

  async function fetchPublicImages() {
    loading.value = true
    error.value = null
    try {
      images.value = await imagesApi.listPublic()
    } catch (e: unknown) {
      error.value = e instanceof Error ? e.message : 'Fehler beim Laden'
    } finally {
      loading.value = false
    }
  }

  async function fetchAllImages() {
    loading.value = true
    error.value = null
    try {
      images.value = await imagesApi.listAll()
    } catch (e: unknown) {
      error.value = e instanceof Error ? e.message : 'Fehler beim Laden'
    } finally {
      loading.value = false
    }
  }

  async function createImage(req: CreateImageRequest): Promise<Image> {
    const image = await imagesApi.create(req)
    images.value.unshift(image)
    return image
  }

  async function deleteImage(id: string) {
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
  }
})
