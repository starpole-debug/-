type AdminCredentials = {
  username: string
  password: string
  adminSecret: string
  accessKey?: string
}

type AdminUser = {
  id: string
  username: string
  nickname?: string
  is_admin?: boolean
}

export const useAdminAuth = () => {
  const adminUser = useState<AdminUser | null>('admin:user', () => null)
  const tokenState = useState<string>('admin:token', () => '')
  const status = useState<'idle' | 'authenticating'>('admin:status', () => 'idle')
  const authError = useState<string | null>('admin:error', () => null)

  const tokenCookie = useCookie<string | null>('admin-token', {
    sameSite: 'lax',
    secure: process.env.NODE_ENV === 'production',
    default: () => null,
  })
  const accessKeyCookie = useCookie<string | null>('admin-access', {
    sameSite: 'lax',
    secure: process.env.NODE_ENV === 'production',
    default: () => null,
  })
  const decodedAccessKey = computed(() => {
    try {
      return accessKeyCookie.value ? decodeURIComponent(accessKeyCookie.value) : ''
    } catch {
      return accessKeyCookie.value || ''
    }
  })

  if (!tokenState.value && tokenCookie.value) {
    tokenState.value = tokenCookie.value
  }

  const isAuthenticated = computed(() => Boolean(tokenState.value))
  const isAuthenticating = computed(() => status.value === 'authenticating')

  const setToken = (token: string | null) => {
    tokenState.value = token ?? ''
    tokenCookie.value = token
  }

  const clearAccessKey = () => {
    accessKeyCookie.value = null
  }

  const login = async (credentials: AdminCredentials) => {
    status.value = 'authenticating'
    authError.value = null
    const config = useRuntimeConfig()
    const apiBase = config.public.apiBase || '/api'
    try {
      const payload = {
        username: credentials.username,
        password: credentials.password,
        admin_secret: credentials.adminSecret,
        access_key: credentials.accessKey || decodedAccessKey.value || undefined,
      }
      const result = await $fetch<{ data: { token: string; user: AdminUser } }>(`${apiBase}/admin/login`, {
        method: 'POST',
        body: payload,
        credentials: 'include',
        // 避免中文放在 Header，统一走请求体
      })
      adminUser.value = result.data.user
      setToken(result.data.token)
      return result.data.user
    } catch (error: any) {
      const message = error?.data?.error || error?.message || '管理员登录失败，请重试'
      authError.value = message
      if (message.toLowerCase().includes('access denied')) {
        clearAccessKey()
      }
      throw new Error(message)
    } finally {
      status.value = 'idle'
    }
  }

  const logout = () => {
    adminUser.value = null
    authError.value = null
    setToken(null)
  }

  return {
    user: adminUser,
    token: tokenState,
    authError,
    isAuthenticated,
    isAuthenticating,
    login,
    logout,
    clearAccessKey,
  }
}
