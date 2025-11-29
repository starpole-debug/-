import { ref } from 'vue'
import { useApi } from './useApi'

export interface ChatSession {
  id: string
  user_id: string
  role_id: string
  model_key: string
  title: string
  last_message?: string
  mode: string
  status: string
  settings: ChatSettings
  created_at: string
  updated_at: string
}

export interface ChatSettings {
  temperature: number
  max_tokens: number
  narrative_focus: string
  action_richness: string
  sfw_mode: boolean
  immersive: boolean
}

export interface ChatMessage {
  id: string
  session_id: string
  role: 'user' | 'assistant' | 'system'
  content: string
  is_important: boolean
  created_at: string
  metadata?: any
}

export interface SessionView {
  session: ChatSession
  role: any // Role type
  world: any // WorldSummary type
  messages: ChatMessage[]
}

export interface ChatModel {
  id: string
  name: string
  description?: string
  provider?: string
  price_hint?: string
  price_coins?: number
}

export interface CreateSessionPayload {
  role_id: string
  model_key?: string
  title?: string
}

export const useChat = () => {
  const api = useApi()
  const config = useRuntimeConfig()
  const tokenState = useState<string>('auth:token', () => '')
  const tokenCookie = useCookie<string | null>('auth-token', {
    sameSite: 'lax',
    secure: process.env.NODE_ENV === 'production',
    default: () => null,
  })
  const sessions = ref<ChatSession[]>([])
  const currentSession = ref<ChatSession | null>(null)
  const messages = ref<ChatMessage[]>([])
  const models = ref<ChatModel[]>([])
  const isLoading = ref(false)
  const isUpdatingSettings = ref(false)
  const isSendingMessage = ref(false)
  const error = ref<string | null>(null)

  const normalizeCreatePayload = (payload: string | CreateSessionPayload, modelKey?: string): CreateSessionPayload => {
    if (typeof payload === 'string') {
      return {
        role_id: payload,
        model_key: modelKey,
      }
    }
    return {
      ...payload,
      model_key: payload.model_key ?? modelKey,
    }
  }

  const listSessions = async () => {
    isLoading.value = true
    error.value = null
    try {
      const res = await api.get<{ data: ChatSession[] }>('/chat/sessions')
      sessions.value = res.data || []
      return sessions.value
    } catch (e: any) {
      error.value = e.message
      throw e
    } finally {
      isLoading.value = false
    }
  }

  const fetchModels = async () => {
    try {
      const res = await api.get<{ data: ChatModel[] }>('/chat/models')
      models.value = res.data || []
      return models.value
    } catch (e: any) {
      // Do not break chat if models are not configured
      models.value = []
      error.value = e.message
      return []
    }
  }

  const createSession = async (payload: string | CreateSessionPayload, modelKey?: string) => {
    isLoading.value = true
    error.value = null
    try {
      const body = normalizeCreatePayload(payload, modelKey)
      if (!body.role_id) {
        throw new Error('缺少角色 ID')
      }
      const res = await api.post<{ data: ChatSession }>('/chat/sessions', body)
      // Prepend to local session list for sidebar freshness
      sessions.value = [res.data, ...sessions.value.filter(s => s.id !== res.data.id)]
      return res.data
    } catch (e: any) {
      error.value = e.message
      throw e
    } finally {
      isLoading.value = false
    }
  }

  const fetchSession = async (sessionId: string) => {
    isLoading.value = true
    error.value = null
    try {
      const res = await api.get<{ data: SessionView }>(`/chat/sessions/${sessionId}`)
      messages.value = Array.isArray(res.data.messages) ? res.data.messages : []
      const lastMsg = messages.value.length ? messages.value[messages.value.length - 1].content : ''
      currentSession.value = { ...res.data.session, last_message: lastMsg }
      // keep sidebar list updated
      sessions.value = [
        ...sessions.value.filter(s => s.id !== res.data.session.id),
        { ...res.data.session, last_message: lastMsg },
      ].sort((a, b) => new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime())
      return res.data
    } catch (e: any) {
      error.value = e.message
      throw e
    } finally {
      isLoading.value = false
    }
  }

  const updateSettings = async (sessionId: string, payload: Partial<ChatSettings> & { mode?: string; model_key?: string }) => {
    isUpdatingSettings.value = true
    try {
      const res = await api.patch<{ data: ChatSession }>(`/chat/sessions/${sessionId}/settings`, payload)
      currentSession.value = res.data
      sessions.value = sessions.value.map(s => (s.id === sessionId ? res.data : s))
      return res.data
    } finally {
      isUpdatingSettings.value = false
    }
  }

  const retryAssistantMessage = async (messageId: string) => {
    const res = await api.post<{ data: ChatMessage[] }>(`/chat/messages/${messageId}/retry`)
    messages.value = Array.isArray(res.data) ? res.data : []
    return messages.value
  }

  const sendMessage = async (sessionId: string, content: string, options?: { preset?: any; stream?: boolean }) => {
    // Optimistic update for user message
    const tempUserId = `user-${Date.now()}`
    const userMsg: ChatMessage = {
      id: tempUserId,
      session_id: sessionId,
      role: 'user',
      content,
      is_important: false,
      created_at: new Date().toISOString(),
    }
    messages.value.push(userMsg)

    // Prepare assistant placeholder for stream
    let tempAssistantId: string | null = null
    if (options?.stream) {
      tempAssistantId = `assistant-${Date.now()}`
      messages.value.push({
        id: tempAssistantId,
        session_id: sessionId,
        role: 'assistant',
        content: '',
        is_important: false,
        created_at: new Date().toISOString(),
        metadata: {},
      })
    }

    isSendingMessage.value = true
    try {
      if (options?.stream) {
        const base = config.public.apiBase || '/api'
        const authToken = tokenState.value || tokenCookie.value || ''
        const headers: Record<string, string> = { 'Content-Type': 'application/json' }
        if (authToken) headers.Authorization = `Bearer ${authToken}`

        const resp = await fetch(`${base}/chat/sessions/${sessionId}/messages`, {
          method: 'POST',
          headers,
          body: JSON.stringify({ content, preset: options?.preset, stream: true }),
        })
        if (!resp.ok || !resp.body) {
          throw new Error(`请求失败 ${resp.status}`)
        }
        const reader = resp.body.getReader()
        const decoder = new TextDecoder()
        let buffer = ''
        let contentAccum = ''
        let reasoningAccum = ''
        while (true) {
          const { done, value } = await reader.read()
          if (done) break
          buffer += decoder.decode(value, { stream: true })
          const parts = buffer.split('\n')
          buffer = parts.pop() || ''
          for (const part of parts) {
            const line = part.trim()
            if (!line) continue
            try {
              const obj = JSON.parse(line)
              if (obj.done) continue
              if (obj.content) contentAccum += obj.content
              if (obj.reasoning) reasoningAccum += obj.reasoning
              if (tempAssistantId) {
                const target = messages.value.find(m => m.id === tempAssistantId)
                if (target) {
                  target.content = contentAccum
                  target.metadata = { ...(target.metadata || {}), reasoning_text: reasoningAccum }
                }
              }
            } catch (err) {
              console.warn('stream parse error', err, line)
            }
          }
        }
        // 最终刷新会话，替换临时消息
        await fetchSession(sessionId)
      } else {
        const res = await api.post<{ data: ChatMessage[] }>(`/chat/sessions/${sessionId}/messages`, {
          content,
          preset: options?.preset,
          stream: false,
        })
        messages.value = Array.isArray(res.data) ? res.data : []
      }

      // bump updated_at locally so sidebar ordering stays fresh
      const now = new Date().toISOString()
      const lastMsg = messages.value.length ? messages.value[messages.value.length - 1].content : content
      if (currentSession.value?.id === sessionId) {
        currentSession.value = { ...currentSession.value, updated_at: now, last_message: lastMsg }
      }
      const existing = sessions.value.find(s => s.id === sessionId)
      const updatedSession = currentSession.value?.id === sessionId
        ? currentSession.value
        : existing
          ? { ...existing, updated_at: now, last_message: lastMsg }
          : null
      if (updatedSession) {
        sessions.value = [
          updatedSession,
          ...sessions.value.filter(s => s.id !== sessionId),
        ]
      }
    } catch (e: any) {
      error.value = e.message
      messages.value = messages.value.filter(m => m.id !== tempUserId && m.id !== tempAssistantId)
      throw e
    } finally {
      isSendingMessage.value = false
    }
  }

  return {
    sessions,
    currentSession,
    messages,
    models,
    isLoading,
    isUpdatingSettings,
    isSendingMessage,
    error,
    retryAssistantMessage,
    listSessions,
    fetchModels,
    createSession,
    fetchSession,
    updateSettings,
    sendMessage,
  }
}
