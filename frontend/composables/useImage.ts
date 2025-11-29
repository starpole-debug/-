type ImageJob = {
  id: string
  status: string
  result_url?: string
  error?: string
  final_prompt?: string
  negative_prompt?: string
}

export const useImage = () => {
  const api = useApi()

  const createJob = async (sessionId: string, messageId?: string, prompt?: string) => {
    const res = await api.post<{ data: ImageJob }>('/chat/images', {
      session_id: sessionId,
      message_id: messageId,
      prompt,
    })
    return res.data
  }

  const getJob = async (id: string) => {
    const res = await api.get<{ data: ImageJob }>(`/chat/images/${id}`)
    return res.data
  }

  return { createJob, getJob }
}
