import type { DictionaryItem, ModelConfig, Role } from '@/types'
import { useAdminApi } from '~/composables/useAdminApi'

type ModelInput = Partial<ModelConfig> & { api_key?: string }

export const useAdminConfig = () => {
  const api = useAdminApi()
  const models = useState<ModelConfig[]>('admin:models', () => [])
  const roles = useState<Role[]>('admin:roles', () => [])

  const loadModels = async () => {
    const data = await api.get<{ data?: ModelConfig[] } | { data: ModelConfig[] } | ModelConfig[]>('/admin/models')
    const raw = (data as any)?.data ?? data
    models.value = Array.isArray(raw) ? raw.filter(Boolean) : []
  }

  const saveModel = async (payload: ModelInput) => {
    if (payload.id) {
      await api.put(`/admin/models/${payload.id}`, payload)
    } else {
      await api.post('/admin/models', payload)
    }
    await loadModels()
  }

  const deleteModel = async (id: string) => {
    await api.del(`/admin/models/${id}`)
    await loadModels()
  }

  const loadRoles = async () => {
    const data = await api.get<{ data?: Role[] } | { data: Role[] } | Role[]>('/admin/roles')
    const raw = (data as any)?.data ?? data
    roles.value = Array.isArray(raw) ? raw.filter(Boolean) : []
  }

  const loadDictionary = async (group?: string) => {
    const query = group ? `?group=${group}` : ''
    const data = await api.get<{ data?: DictionaryItem[] } | { data: DictionaryItem[] } | DictionaryItem[]>(
      `/admin/dictionary${query}`,
    )
    const raw = (data as any)?.data ?? data
    return Array.isArray(raw) ? raw.filter(Boolean) : []
  }

  const saveDictionary = async (payload: Partial<DictionaryItem>) => {
    await api.post('/admin/dictionary', payload)
  }

  const deleteDictionary = async (id: string) => {
    await api.del(`/admin/dictionary/${id}`)
  }

  return { models, roles, loadModels, saveModel, deleteModel, loadRoles, loadDictionary, saveDictionary, deleteDictionary }
}
