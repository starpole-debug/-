<template>
  <section class="space-y-6 rounded-3xl bg-white p-8 shadow">
    <header class="flex flex-wrap items-center justify-between gap-3">
      <div>
        <h1 class="text-xl font-semibold">收藏夹</h1>
        <p class="text-sm text-muted">展示你标记为「收藏」的社区帖子。</p>
      </div>
      <button class="text-sm text-primary disabled:opacity-50" :disabled="loading" @click="load">
        {{ loading ? '正在刷新...' : '刷新' }}
      </button>
    </header>
    <p v-if="errorMessage" class="text-sm text-rose-500">{{ errorMessage }}</p>
    <div v-if="loading" class="text-sm text-muted">加载中...</div>
    <div v-else-if="!favorites.length" class="text-sm text-muted">还没有收藏的帖子，去社区逛逛吧～</div>
    <div v-else class="space-y-4">
      <CommunityPostCard v-for="post in favorites" :key="post.id" :post="post" />
    </div>
  </section>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: ['auth'],
})

import type { CommunityPost, UserHomePayload } from '@/types'

const api = useApi()
const loading = ref(false)
const errorMessage = ref('')
const favorites = ref<CommunityPost[]>([])

const load = async () => {
  loading.value = true
  errorMessage.value = ''
  try {
    const res = await api.get<{ data: CommunityPost[] }>('/me/favorites')
    favorites.value = Array.isArray(res.data) ? res.data : []
    if (!favorites.value.length) {
      // 后端列表空时用首页聚合接口兜底
      const home = await api.get<{ data: UserHomePayload }>('/me/home')
      favorites.value = Array.isArray(home.data?.favorites) ? home.data.favorites : favorites.value
    }
  } catch (error: any) {
    errorMessage.value = error?.data?.error || error?.message || '无法加载收藏列表，请稍后再试'
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>
