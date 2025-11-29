// composables/useCreatorStats.ts
import { useApi } from '~/composables/useApi'
// 这里保留你原来的类型导入
// import type { DashboardResponse, Role } from '@/types'

export const useCreatorStats = () => {
  const api = useApi()

  // 这个可以继续用 useState，因为它确实是用户级的全局信息，不像 posts 那种列表
  const dashboard = useState<DashboardResponse | null>('creator:dashboard', () => null)
  const roles = useState<Role[]>('creator:roles', () => [])

  const isUnauthorized = (err: any) => {
    const status = err?.statusCode ?? err?.response?.status ?? err?.status
    return status === 401
  }

  const loadDashboard = async () => {
    try {
      const res = await api.get<{ data: DashboardResponse }>('/creator/dashboard')
      dashboard.value = res.data
    } catch (err: any) {
      if (isUnauthorized(err)) {
        // 未登录：你可以在这里选择静默，或者跳登录
        // const router = useRouter()
        // router.push('/login')
        return
      }
      console.error('加载 creator dashboard 失败:', err)
    }
  }

  const loadRoles = async () => {
    try {
      const res = await api.get<{ data: Role[] }>('/creator/roles')
      roles.value = Array.isArray(res.data) ? res.data : []
    } catch (err: any) {
      if (isUnauthorized(err)) {
        return
      }
      console.error('加载 creator roles 失败:', err)
    }
  }

  const fetchCreatorRole = async (id: string) => {
    const res = await api.get<{ data: Role }>(`/creator/roles/${id}`)
    return res.data
  }

  const publishRole = async (id: string) => {
    await api.post(`/roles/${id}/publish`)
    roles.value = roles.value.map((r) => (r.id === id ? { ...r, status: 'published' } : r))
  }

  return {
    dashboard,
    roles,
    loadDashboard,
    loadRoles,
    fetchCreatorRole,
    publishRole,
  }
}
