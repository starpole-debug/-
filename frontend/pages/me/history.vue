<template>
  <section class="min-h-screen px-4 py-10 space-y-6">
    <div class="max-w-6xl mx-auto flex items-center justify-between">
      <div class="space-y-1">
        <p class="text-xs uppercase tracking-[0.2em] text-charcoal-400">History</p>
        <h1 class="text-2xl font-display font-bold text-charcoal-900">浏览历史</h1>
        <p class="text-sm text-charcoal-500">回看你最近浏览过的社区帖子和角色。</p>
      </div>
      <NuxtLink to="/me" class="text-xs font-medium text-charcoal-500 hover:text-charcoal-900">返回个人中心 →</NuxtLink>
    </div>

    <div class="max-w-6xl mx-auto card-soft p-6">
      <div v-if="loading" class="py-12 text-center text-charcoal-400 flex flex-col items-center gap-3">
        <div class="i-svg-spinners-90-ring-with-bg text-2xl" />
        <p class="text-sm">加载浏览历史...</p>
      </div>

      <div v-else-if="errorMessage" class="py-12 text-center text-rose-500 text-sm">
        {{ errorMessage }}
      </div>

      <div v-else>
        <div v-if="recentViews.length === 0" class="py-12 text-center text-charcoal-400 text-sm">
          暂无浏览记录，去社区逛逛吧。
        </div>
        <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
          <CommunityPostCard v-for="post in recentViews" :key="post.id" :post="post" />
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import type { UserHomePayload, CommunityPost } from '@/types'

definePageMeta({
  middleware: ['auth'],
})

const api = useApi()
const loading = ref(false)
const errorMessage = ref('')
const recentViews = ref<CommunityPost[]>([])

const load = async () => {
  loading.value = true
  errorMessage.value = ''
  try {
    const res = await api.get<{ data: UserHomePayload }>('/me/home')
    recentViews.value = res.data?.recent_views || []
  } catch (error: any) {
    errorMessage.value = error?.data?.error || error?.message || '加载浏览历史失败'
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<style scoped>
.i-svg-spinners-90-ring-with-bg {
  display: inline-block;
}
</style>
