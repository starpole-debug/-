import type { FetchOptions } from 'ofetch'

const handleUserAuthError = async () => {
  const auth = useAuth()
  const route = useRoute()
  auth.logout()
  if (route.path !== '/login') {
    await navigateTo({
      path: '/login',
      query: route.fullPath ? { redirect: route.fullPath } : undefined,
    })
  }
}

export const useApi = () => {
  const config = useRuntimeConfig()
  const token = useState<string>('auth:token', () => '')
  const tokenCookie = useCookie<string | null>('auth-token', {
    sameSite: 'lax',
    secure: process.env.NODE_ENV === 'production',
    default: () => null,
  })
  const base = config.public.apiBase || '/api'

  const request = async <T>(path: string, options: FetchOptions<'json'> = {}) => {
    const headers: Record<string, string> = {
      Accept: 'application/json',
      ...(options.headers as Record<string, string>),
    }
    const authToken = token.value || tokenCookie.value || ''
    if (authToken) {
      headers.Authorization = `Bearer ${authToken}`
    }
    try {
      return await $fetch<T>(`${base}${path}`, {
        ...options,
        headers,
        credentials: 'include',
      })
    } catch (error: any) {
      const status = error?.statusCode ?? error?.response?.status ?? error?.status
      if (status === 401 || status === 403) {
        await handleUserAuthError()
        throw new Error('登录状态已失效，请重新登录')
      }
      // 把后端返回的 error 字段透传到前端，便于展示「余额不足」等业务错误
      const message =
        error?.data?.error ||
        error?.response?._data?.error ||
        error?.message ||
        '请求失败，请稍后重试'
      throw new Error(message)
    }
  }

  return {
    get: <T>(path: string, options: FetchOptions<'json'> = {}) => request<T>(path, { ...options, method: 'GET' }),
    post: <T>(path: string, body?: any, options: FetchOptions<'json'> = {}) =>
      request<T>(path, { ...options, method: 'POST', body }),
    put: <T>(path: string, body?: any, options: FetchOptions<'json'> = {}) =>
      request<T>(path, { ...options, method: 'PUT', body }),
    patch: <T>(path: string, body?: any, options: FetchOptions<'json'> = {}) =>
      request<T>(path, { ...options, method: 'PATCH', body }),
    del: <T>(path: string, options: FetchOptions<'json'> = {}) => request<T>(path, { ...options, method: 'DELETE' }),
  }
}
