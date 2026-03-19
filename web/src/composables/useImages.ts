import { onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useImagesStore } from '@/stores/images.store'
import { useAuthStore } from '@/stores/auth.store'

export function useImages(mode: 'public' | 'all' = 'public') {
  const store = useImagesStore()
  const authStore = useAuthStore()
  const { images, publicImages, loading, error } = storeToRefs(store)

  function fetch() {
    if (mode === 'all' && authStore.isAuthenticated) {
      return store.fetchAllImages()
    }
    return store.fetchPublicImages()
  }

  onMounted(() => {
    fetch()
  })

  return {
    images,
    publicImages,
    loading,
    error,
    refresh: fetch,
    createImage: store.createImage,
    deleteImage: store.deleteImage,
  }
}
