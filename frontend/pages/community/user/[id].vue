<template>
  <section class="space-y-8">
    <header class="rounded-3xl bg-white p-6 shadow flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
      <div class="flex items-center gap-4">
        <img
          :src="profile?.avatar_url || defaultAvatar"
          class="h-16 w-16 rounded-2xl border border-slate-200 object-cover"
          alt="avatar"
        />
        <div>
          <p class="text-xs text-muted">用户 ID：{{ profile?.id || route.params.id }}</p>
          <h1 class="text-2xl font-semibold text-slate-900">
            {{ profile?.nickname || profile?.username || '社区用户' }}
          </h1>
          <p class="text-sm text-slate-500">
            粉丝 {{ profile?.follower_count ?? 0 }} · 关注 {{ profile?.following_count ?? 0 }}
          </p>
        </div>
      </div>
      <div class="flex items-center gap-3">
        <button
          v-if="!isSelf"
          class="rounded-full px-5 py-2 text-sm font-medium text-white transition"
          :class="profile?.is_following ? 'bg-slate-700' : 'bg-primary hover:brightness-110'"
          :disabled="followPending"
          @click="toggleFollow"
        >
          {{ followPending ? '处理中...' : profile?.is_following ? '已关注' : '关注' }}
        </button>
        <NuxtLink to="/community" class="rounded-full border border-slate-200 px-4 py-2 text-sm text-slate-700 hover:border-primary">
          返回社区
        </NuxtLink>
      </div>
    </header>

    <section class="space-y-4">
      <header class="flex items-center justify-between">
        <h2 class="text-lg font-semibold text-slate-900">TA 的帖子</h2>
        <span class="text-xs text-muted">公开帖子按时间倒序排列</span>
      </header>
      <p v-if="loading" class="text-sm text-muted">加载中...</p>
      <p v-else-if="errorMessage" class="text-sm text-rose-500">{{ errorMessage }}</p>
      <p v-else-if="!posts.length" class="text-sm text-muted">还没有公开帖子。</p>
      <div v-else class="columns-1 md:columns-2 lg:columns-3 space-y-4">
        <CommunityPostCard v-for="post in posts" :key="post.id" :post="post" class="break-inside-avoid" />
      </div>
    </section>
  </section>
</template>

<script setup lang="ts">
import type { CommunityPost, User } from '@/types'

const route = useRoute()
const community = useCommunity()
const auth = useAuth()

const profile = ref<User | null>(null)
const posts = ref<CommunityPost[]>([])
const loading = ref(false)
const followPending = ref(false)
const errorMessage = ref('')
const defaultAvatar = 'https://placehold.co/80x80?text=User'

const isSelf = computed(() => profile.value?.id && auth.user.value?.id === profile.value.id)

const load = async () => {
  loading.value = true
  errorMessage.value = ''
  try {
    const id = route.params.id as string
    profile.value = await community.fetchUserProfile(id)
    posts.value = await community.fetchUserPosts(id)
  } catch (error: any) {
    errorMessage.value = error?.data?.error || error?.message || '加载失败，请稍后再试'
  } finally {
    loading.value = false
  }
}

const toggleFollow = async () => {
  if (!profile.value || isSelf.value) return
  followPending.value = true
  errorMessage.value = ''
  try {
    const res = await community.followUser(profile.value.id)
    profile.value = {
      ...profile.value,
      is_following: res.following,
      follower_count: res.follower_count,
      following_count: res.following_count,
    }
  } catch (error: any) {
    errorMessage.value = error?.data?.error || error?.message || '操作失败，请稍后再试'
  } finally {
    followPending.value = false
  }
}

watch(
  () => route.params.id,
  () => load(),
  { immediate: true },
)
</script>
