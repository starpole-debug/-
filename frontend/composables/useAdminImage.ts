import type { FetchOptions } from 'ofetch'

type ImageProvider = {
  id?: string
  name: string
  base_url: string
  api_key?: string
  max_concurrency?: number
  weight?: number
  status?: string
  params_json?: string
}

type ImagePreset = {
  id?: string
  name: string
  preset_json: string
  status?: string
}

export const useAdminImage = () => {
  const api = useAdminApi()

  const listProviders = async () => {
    const res = await api.get<{ data?: ImageProvider[] } | ImageProvider[]>('/admin/image-providers')
    return (res as any)?.data ?? res ?? []
  }
  const saveProvider = async (payload: ImageProvider) => {
    if (payload.params_json !== undefined && payload.params_json !== null) {
      try {
        const raw = typeof payload.params_json === 'string' ? payload.params_json.trim() : JSON.stringify(payload.params_json)
        if (raw) {
          const parsed = JSON.parse(raw)
          payload.params_json = JSON.stringify(parsed)
        } else {
          payload.params_json = '{}'
        }
      } catch (e) {
        throw new Error('默认参数需要是合法的 JSON')
      }
    }
    if (payload.id) {
      return api.put<ImageProvider>(`/admin/image-providers/${payload.id}`, payload as any)
    }
    return api.post<ImageProvider>('/admin/image-providers', payload as any)
  }
  const deleteProvider = async (id: string) => api.del(`/admin/image-providers/${id}`)

  const listPresets = async () => {
    const res = await api.get<{ data?: ImagePreset[] } | ImagePreset[]>('/admin/image-presets')
    return (res as any)?.data ?? res ?? []
  }
  const savePreset = async (payload: ImagePreset) => {
    if (payload.id) {
      return api.put<ImagePreset>(`/admin/image-presets/${payload.id}`, payload as any)
    }
    return api.post<ImagePreset>('/admin/image-presets', payload as any)
  }
  const deletePreset = async (id: string) => api.del(`/admin/image-presets/${id}`)

  return {
    listProviders,
    saveProvider,
    deleteProvider,
    listPresets,
    savePreset,
    deletePreset,
  }
}
