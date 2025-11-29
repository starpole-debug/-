import type { User } from '@/types'
import { useAdminApi } from '~/composables/useAdminApi'

interface FetchParams {
  query?: string
  limit?: number
  offset?: number
}

interface CreateUserPayload {
  username: string
  email: string
  password: string
  nickname: string
  is_admin?: boolean
}

export const useAdminUsers = () => {
  const api = useAdminApi()
  const users = useState<User[]>('admin:users', () => [])
  const loading = useState<boolean>('admin:users:loading', () => false)

  const buildQuery = (params: FetchParams = {}) => {
    const search = new URLSearchParams()
    if (params.query) search.set('query', params.query)
    if (params.limit) search.set('limit', String(params.limit))
    if (params.offset) search.set('offset', String(params.offset))
    const query = search.toString()
    return query ? `?${query}` : ''
  }

  const fetchUsers = async (params: FetchParams = {}) => {
    loading.value = true
    try {
      const data = await api.get<{ data: User[] }>(`/admin/users${buildQuery(params)}`)
      users.value = data.data
      return users.value
    } finally {
      loading.value = false
    }
  }

  const createUser = async (payload: CreateUserPayload) => {
    await api.post('/admin/users', payload)
  }

  const banUser = async (id: string) => {
    await api.post(`/admin/users/${id}/ban`)
  }

  const unbanUser = async (id: string) => {
    await api.post(`/admin/users/${id}/unban`)
  }

  const deleteUser = async (id: string) => {
    await api.del(`/admin/users/${id}`)
  }

  return { users, loading, fetchUsers, createUser, banUser, unbanUser, deleteUser }
}
