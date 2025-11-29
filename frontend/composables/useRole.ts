import type { Role } from '@/types'
import { useAssetUrl } from '~/composables/useAssetUrl'

export const useRole = () => {
  const api = useApi()
  const roles = useState<Role[]>('roles:list', () => [])
  const featured = useState<Role[]>('roles:featured', () => [])
  const favorites = useState<Role[]>('roles:favorites', () => [])
  const isLoading = useState<boolean>('roles:loading', () => false)
  const hasLoaded = useState<boolean>('roles:loaded', () => false)
  const errorMessage = useState<string | null>('roles:error', () => null)
  const { resolveAssetUrl } = useAssetUrl()

  const normalizeRole = (role: any): Role => {
    const tags =
      Array.isArray(role?.tags)
        ? role.tags
        : typeof role?.tags === 'string'
          ? role.tags.split(/[,，]/).map((t: string) => t.trim()).filter(Boolean)
          : []
    const avatar = role?.avatar_url || role?.avatarUrl
    return {
      ...role,
      avatar_url: resolveAssetUrl(avatar),
      tags,
    }
  }

  const normalizeList = (list: any) => (Array.isArray(list) ? list.map(normalizeRole) : [])

  const fetchRoles = async (force = false) => {
    if (isLoading.value) return
    if (!force && hasLoaded.value) return
    isLoading.value = true
    errorMessage.value = null
    try {
      const data = await api.get<{ data: Role[] }>('/roles')
      roles.value = normalizeList(data.data)
      hasLoaded.value = true
    } catch (error) {
      errorMessage.value = error instanceof Error ? error.message : '加载角色失败'
      throw error
    } finally {
      isLoading.value = false
    }
  }

  const fetchFeatured = async () => {
    const data = await api.get<{ data: Role[] }>('/roles/featured')
    featured.value = normalizeList(data.data)
  }

  const fetchRole = async (id: string) => {
    const data = await api.get<{ data: Role }>(`/roles/${id}`)
    return normalizeRole(data.data)
  }

  const fetchFavorites = async () => {
    try {
      const res = await api.get<{ data: Role[] }>('/roles/favorites')
      favorites.value = normalizeList(res.data)
      return favorites.value
    } catch (error) {
      favorites.value = []
      return favorites.value
    }
  }

  const favoriteRole = async (id: string, role?: Role) => {
    await api.post(`/roles/${id}/favorite`)
    const candidate = normalizeRole(role || roles.value.find(r => r.id === id) || featured.value.find(r => r.id === id))
    if (candidate && !favorites.value.some(r => r.id === id)) {
      favorites.value = [candidate, ...favorites.value]
    } else if (!candidate) {
      await fetchFavorites()
    }
  }

  const unfavoriteRole = async (id: string) => {
    await api.del(`/roles/${id}/favorite`)
    favorites.value = favorites.value.filter(r => r.id !== id)
  }

  return {
    roles,
    featured,
    favorites,
    isLoading,
    hasLoaded,
    errorMessage,
    fetchRoles,
    fetchFeatured,
    fetchRole,
    fetchFavorites,
    favoriteRole,
    unfavoriteRole,
  }
}
