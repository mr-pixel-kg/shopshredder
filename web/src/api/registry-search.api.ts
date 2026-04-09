import { apiClient } from './client'

export interface ImageSuggestion {
  name: string
  description: string
  starCount: number
  isOfficial: boolean
}

export interface TagSuggestion {
  name: string
  fullSize?: number
  lastUpdated?: string
}

export const registrySearchApi = {
  async searchImages(query: string): Promise<ImageSuggestion[]> {
    const { data } = await apiClient.get<{ results: ImageSuggestion[] }>(
      '/api/registry/images/search',
      { params: { q: query } },
    )
    return data.results
  },

  async searchTags(image: string, query: string): Promise<TagSuggestion[]> {
    const { data } = await apiClient.get<{ results: TagSuggestion[] }>('/api/registry/tags', {
      params: { image, q: query },
    })
    return data.results
  },
}
