import type { FetchOptions } from 'ofetch'

const handleAdminAuthError = async () => {
  const adminAuth = useAdminAuth()
  const route = useRoute()
  adminAuth.logout()
  if (route.path !== '/admin/login') {
    await navigateTo('/admin/login')
  }
}

export const useAdminApi = () => {
  const config = useRuntimeConfig()
  const token = useState<string>('admin:token', () => '')
  const base = config.public.apiBase || '/api'

  const request = async <T>(path: string, options: FetchOptions<'json'> = {}) => {
    const headers: Record<string, string> = {
      Accept: 'application/json',
      ...(options.headers as Record<string, string>),
    }
    if (token.value) {
      headers.Authorization = `Bearer ${token.value}`
    }
    try {
      return await $fetch<T>(`${base}${path}`, {
        ...options,
        headers,
      })
    } catch (error: any) {
      const status = error?.statusCode ?? error?.response?.status ?? error?.status
      if (status === 401 || status === 403) {
        await handleAdminAuthError()
        throw new Error('管理员登录状态已过期，请重新登录')
      }
      throw error
    }
  }

  return {
    get: <T>(path: string, options: FetchOptions<'json'> = {}) => request<T>(path, { ...options, method: 'GET' }),
    post: <T>(path: string, body?: any, options: FetchOptions<'json'> = {}) =>
      request<T>(path, { ...options, method: 'POST', body }),
    put: <T>(path: string, body?: any, options: FetchOptions<'json'> = {}) =>
      request<T>(path, { ...options, method: 'PUT', body }),
    del: <T>(path: string, options: FetchOptions<'json'> = {}) => request<T>(path, { ...options, method: 'DELETE' }),
  }
}

