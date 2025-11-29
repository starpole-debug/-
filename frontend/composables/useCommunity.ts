import { ref } from 'vue'
import { useApi } from '~/composables/useApi'
import { useAdminApi } from '~/composables/useAdminApi'

export const useCommunity = () => {
  const api = useApi()
  const adminApi = useAdminApi()
  const posts = ref<CommunityPost[]>([])
  const dictionary = ref<Record<string, DictionaryItem[]>>({})
  const favorites = ref<CommunityPost[]>([])

  const fetchFeed = async (params: { sort?: string; filter?: string; search?: string } = {}) => {
    const search = new URLSearchParams()
    if (params.sort) search.set('sort', params.sort)
    if (params.filter) search.set('filter', params.filter)
    if (params.search) search.set('search', params.search)

    const data = await api.get<{ data: CommunityPost[] }>(`/community?${search.toString()}`)
    posts.value = Array.isArray(data.data) ? data.data : []
  }

  const fetchPost = async (id: string) => {
    const data = await api.get<{
      data: {
        post: CommunityPost
        comments: CommunityComment[]
        reactions: Record<string, number>
      }
    }>(`/community/${id}`)

    return data.data
  }

  const fetchAdminPosts = async (
    params: { query?: string; visibility?: string; limit?: number; offset?: number } = {}
  ) => {
    const search = new URLSearchParams()
    if (params.query) search.set('query', params.query)
    if (params.visibility && params.visibility !== 'all') search.set('visibility', params.visibility)
    if (params.limit) search.set('limit', String(params.limit))
    if (params.offset) search.set('offset', String(params.offset))

    const res = await adminApi.get(`/admin/posts${search.toString() ? `?${search}` : ''}`)

    const raw =
      Array.isArray(res) ? res :
        Array.isArray(res?.data) ? res.data :
          Array.isArray(res?.posts) ? res.posts :
            Array.isArray(res?.items) ? res.items :
              []

    posts.value = raw.filter(x => x && typeof x === 'object')
    return posts.value
  }

  const createPost = async (payload: { title: string; content: string; attachments?: string[]; link_url?: string; link_type?: string }) => {
    return await api.post('/community', payload)
  }

  const addComment = async (id: string, payload: { content: string }) => {
    return await api.post(`/community/${id}/comments`, payload)
  }

  const react = async (id: string, reaction: string) => {
    const res = await api.post<{ data: { status: string; reactions?: Record<string, number>; activated?: boolean } }>(
      `/community/${id}/reactions`,
      { type: reaction },
    )
    return {
      reactions: res.data?.reactions,
      activated: Boolean(res.data?.activated),
    }
  }

  const loadDictionary = async () => {
    const data = await api.get<{ data: Record<string, DictionaryItem[]> }>('/community/dictionary')
    dictionary.value = data.data || {}
  }

  const fetchUserProfile = async (id: string) => {
    const data = await api.get<{ data: User }>(`/community/users/${id}`)
    return data.data
  }

  const fetchUserPosts = async (id: string) => {
    const data = await api.get<{ data: CommunityPost[] }>(`/community/users/${id}/posts`)
    return data.data || []
  }

  const followUser = async (id: string) => {
    const res = await api.post<{ data: { following: boolean; follower_count: number; following_count: number } }>(
      `/community/users/${id}/follow`,
      {},
    )
    return res.data
  }

  const fetchFavorites = async () => {
    try {
      const data = await api.get<{ data: CommunityPost[] }>('/me/favorites')
      favorites.value = Array.isArray(data.data) ? data.data : []
      // 防止后端返回空列表但概览接口里有数据，做一次兜底
      if (!favorites.value.length) {
        const home = await api.get<{ data: UserHomePayload }>('/me/home')
        favorites.value = Array.isArray(home.data?.favorites) ? home.data.favorites : favorites.value
      }
    } catch (error) {
      // /me/favorites 失败时兜底尝试从个人主页获取收藏
      try {
        const home = await api.get<{ data: UserHomePayload }>('/me/home')
        favorites.value = Array.isArray(home.data?.favorites) ? home.data.favorites : []
      } catch {
        favorites.value = []
        throw error
      }
    }
    return favorites.value
  }

  return {
    posts,
    dictionary,
    favorites,
    fetchFeed,
    fetchPost,
    fetchAdminPosts,
    createPost,
    addComment,
    react,
    loadDictionary,
    fetchUserProfile,
    fetchUserPosts,
    followUser,
    fetchFavorites
  }
}
