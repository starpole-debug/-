type Credentials = {
  identifier?: string
  username?: string
  email?: string
  code?: string
  password: string
  nickname?: string
}

type AuthUser = {
  id: string
  username: string
  token: string
  nickname?: string
  isAdmin?: boolean
}

export const useAuth = () => {
  const storedUser = useState<AuthUser | null>('auth:user', () => null)
  const tokenState = useState<string>('auth:token', () => '')
  const status = useState<'idle' | 'authenticating'>('auth:status', () => 'idle')
  const authError = useState<string | null>('auth:error', () => null)

  // 持久化的 token（存在 cookie 里）
  const tokenCookie = useCookie<string | null>('auth-token', {
    sameSite: 'lax',
    secure: process.env.NODE_ENV === 'production',
    default: () => null,
  })

  // 运行时配置里建议加一个 apiBase，比如 http://localhost:8080/api
  // nuxt.config.ts:
  // runtimeConfig: {
  //   public: {
  //     apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8080/api',
  //   },
  // },
  const config = useRuntimeConfig()
  const apiBase = config.public.apiBase || '/api'

  // 如果 cookie 里有 token，但内存里没有，先同步一下
  if (!tokenState.value && tokenCookie.value) {
    tokenState.value = tokenCookie.value
  }

  const isAuthenticated = computed(() => Boolean(storedUser.value?.token))
  const hasToken = computed(() => Boolean(tokenState.value || tokenCookie.value))
  const isAuthenticating = computed(() => status.value === 'authenticating')

  const setToken = (token: string | null) => {
    tokenState.value = token ?? ''
    tokenCookie.value = token
  }

  const apiGet = async <T>(path: string): Promise<T> => {
    const result = await $fetch<{ data: T }>(`${apiBase}${path}`, {
      method: 'GET',
      headers: tokenState.value
        ? {
            Authorization: `Bearer ${tokenState.value}`,
          }
        : undefined,
      credentials: 'include',
    })
    return result.data
  }

  const apiPost = async <T>(path: string, body: any): Promise<T> => {
    const result = await $fetch<{ data: T }>(`${apiBase}${path}`, {
      method: 'POST',
      body,
      headers: tokenState.value
        ? {
            Authorization: `Bearer ${tokenState.value}`,
          }
        : undefined,
      credentials: 'include',
    })
    return result.data
  }

  /**
   * 登录
   * 期望后端返回：
   * {
   *   token: string;
   *   user: { id: string; username: string }
   * }
   */
  const buildUser = (payload: { id: string; username: string; nickname?: string; is_admin?: boolean }, token: string): AuthUser => ({
    id: payload.id,
    username: payload.username,
    nickname: payload.nickname,
    isAdmin: Boolean(payload.is_admin),
    token,
  })

  const login = async (credentials: Credentials) => {
    status.value = 'authenticating'
    authError.value = null
    try {
      const identifier = credentials.identifier || credentials.username || credentials.email || ''
      const result = await apiPost<{ token: string; user: { id: string; username: string; nickname?: string; is_admin?: boolean } }>(
        '/auth/login',
        { identifier, password: credentials.password },
      )
      const fullUser = buildUser(result.user, result.token)

      storedUser.value = fullUser
      setToken(result.token)

      return fullUser
    } catch (error: any) {
      const message = error?.data?.error || error?.message || '登录失败，请检查账号密码'
      authError.value = message
      throw new Error(message)
    } finally {
      status.value = 'idle'
    }
  }

  /**
   * 注册
   * 这里先简单按和登录一样的结构来：
   * POST /auth/register -> { token, user }
   * 如果你后端定义不一样，把 body / 返回值类型改一下就行
   */
  const register = async (credentials: Credentials) => {
    status.value = 'authenticating'
    authError.value = null
    try {
      const payload = {
        username: credentials.username,
        email: credentials.email,
        code: credentials.code,
        password: credentials.password,
        nickname: credentials.nickname?.trim() || credentials.username,
      }
      const result = await apiPost<{ token: string; user: { id: string; username: string; nickname?: string; is_admin?: boolean } }>(
        '/auth/register',
        payload,
      )
      const fullUser = buildUser(result.user, result.token)

      storedUser.value = fullUser
      setToken(result.token)

      return fullUser
    } catch (error: any) {
      const message = error?.data?.error || error?.message || '注册失败，请稍后重试'
      authError.value = message
      throw new Error(message)
    } finally {
      status.value = 'idle'
    }
  }

  /**
   * 从后端拉取当前登录用户信息
   * layout 里可以 onMounted(() => auth.fetchProfile())
   * 解决你之前的 auth.fetchProfile is not a function + 500 问题
   */
  const fetchProfile = async () => {
    // 先确保 tokenState 里有东西
    if (!tokenState.value && tokenCookie.value) {
      tokenState.value = tokenCookie.value
    }

    if (!tokenState.value) {
      // 没有 token，当没登录处理
      return null
    }

    status.value = 'authenticating'
    authError.value = null

    try {
      const result = await apiGet<{ id: string; username: string; nickname?: string; is_admin?: boolean }>('/auth/me')

      const fullUser: AuthUser = {
        id: result.id,
        username: result.username,
        nickname: result.nickname,
        isAdmin: Boolean(result.is_admin),
        token: tokenState.value!,
      }

      storedUser.value = fullUser
      return fullUser
    } catch (error: any) {
      // token 失效 / 被撤销：当成登出处理
      console.warn('fetchProfile failed, clearing auth state', error)
      logout()
      return null
    } finally {
      status.value = 'idle'
    }
  }

  const sendCode = async (email: string, purpose: 'signup' | 'reset') => {
    await apiPost('/auth/send-code', { email, purpose })
  }

  const resetPassword = async (email: string, code: string, password: string) => {
    status.value = 'authenticating'
    authError.value = null
    try {
      await apiPost('/auth/password/reset', { email, code, password })
    } catch (error: any) {
      const message = error?.data?.error || error?.message || '重置密码失败'
      authError.value = message
      throw new Error(message)
    } finally {
      status.value = 'idle'
    }
  }

  const logout = () => {
    storedUser.value = null
    setToken(null)
    authError.value = null
  }

  return {
    user: storedUser,
    isAuthenticated,
    hasToken,
    isAuthenticating,
    authError,
    login,
    register,
    sendCode,
    resetPassword,
    logout,
    fetchProfile,
  }
}
