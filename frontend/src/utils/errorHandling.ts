export function getErrorMessage(error: unknown, fallback = '操作失败，请重试'): string {
  if (error instanceof Error) {
    return error.message
  }
  if (typeof error === 'string') {
    return error
  }
  if (error && typeof error === 'object') {
    const axiosError = error as { response?: { data?: { message?: string } }; message?: string }
    if (axiosError.response?.data?.message) {
      return axiosError.response.data.message
    }
    if (axiosError.message) {
      return axiosError.message
    }
  }
  return fallback
}

export function isAxiosError(error: unknown): error is { response: { status: number; data?: { message?: string } } } {
  return (
    error !== null &&
    typeof error === 'object' &&
    'response' in error &&
    error.response !== null &&
    typeof error.response === 'object'
  )
}

export function getErrorStatus(error: unknown): number | null {
  if (isAxiosError(error)) {
    return error.response.status
  }
  return null
}
