import axios from 'axios'

import type { ApiErrorResponse } from '@/types'

interface ProblemDetail {
  title?: string
  status?: number
  detail?: string
}

export function getApiErrorMessage(
  error: unknown,
  fallback = 'An unexpected error occurred',
): string {
  if (axios.isAxiosError(error)) {
    const data = error.response?.data as (ApiErrorResponse & ProblemDetail) | undefined

    if (data?.detail) return data.detail

    if (data?.error?.message) return data.error.message

    if (error.message) return error.message
  }

  if (error instanceof Error) return error.message

  return fallback
}
