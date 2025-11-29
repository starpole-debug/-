<template>
  <section>
    <header class="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
      <div>
        <h1 class="text-xl font-semibold">社区守望</h1>
        <p class="text-sm text-slate-300">按可见性筛选帖子并快速下发隐藏或限流操作。</p>
      </div>
      <div class="flex gap-2">
        <input v-model="query" placeholder="搜索标题/内容/ID" class="rounded-full border border-white/10 bg-white/5 px-4 py-2 text-sm" />
        <select v-model="visibility" class="rounded-full border border-white/10 bg-white/5 px-3 py-2 text-sm">
          <option value="all">全部</option>
          <option value="public">公开</option>
          <option value="hidden">隐藏</option>
          <option value="limited">限流</option>
          <option value="deleted">删除</option>
        </select>
        <button class="rounded-full bg-white/10 px-4 py-2 text-sm text-white" :disabled="loading" @click="loadPosts">
          {{ loading ? '加载中...' : '刷新' }}
        </button>
      </div>
    </header>
    <p v-if="error" class="mt-3 text-sm text-rose-300">{{ error }}</p>
    <div class="mt-6 space-y-4">
      <div
        v-for="post in community.posts"
        :key="post.id || post._id || post.uuid || 'no-id'"
        class="rounded-2xl border border-white/10 bg-white/5 p-5 text-sm text-slate-100"
      >
        <div class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <p class="text-base font-semibold text-white">{{ post.title }}</p>
            <p class="text-xs text-slate-400">
             #{{ safe(post.id, 8) }} · 作者 {{ safe(post.author_id, 6) }} · {{ formatTime(post.created_at) }}
            </p>
          </div>
          <span class="rounded-full border border-white/10 px-3 py-1 text-xs uppercase text-slate-300">
            {{ post.visibility || 'public' }}
          </span>
        </div>
        <p class="mt-3 max-h-24 overflow-hidden text-ellipsis text-slate-200">{{ post.content }}</p>
        <div class="mt-4 flex flex-wrap gap-2 text-xs text-slate-300">
          <NuxtLink :to="`/community/post/${post.id}`" class="rounded-full border border-white/10 px-3 py-1 text-primary">查看详情</NuxtLink>
          <button
            class="rounded-full border border-white/10 px-3 py-1"
            :disabled="isModerating(post.id, 'hide')"
            @click="moderate(post.id, 'hide')"
          >
            {{ isModerating(post.id, 'hide') ? '隐藏中...' : '隐藏' }}
          </button>
          <button
            class="rounded-full border border-white/10 px-3 py-1"
            :disabled="isModerating(post.id, 'limit')"
            @click="moderate(post.id, 'limit')"
          >
            {{ isModerating(post.id, 'limit') ? '限流中...' : '限流' }}
          </button>
          <button
            class="rounded-full border border-white/10 px-3 py-1"
            :disabled="isModerating(post.id, 'restore')"
            @click="moderate(post.id, 'restore')"
          >
            {{ isModerating(post.id, 'restore') ? '恢复中...' : '恢复公开' }}
          </button>
        </div>
      </div>
      <p v-if="!community.posts.length && !loading" class="text-sm text-slate-400">暂无帖子。</p>
    </div>
    </section>
</template>

<script setup lang="ts">
import dayjs from 'dayjs'

definePageMeta({
  layout: 'admin',
  middleware: 'admin',
})

const safe = (v: any, len = 8) => {
  if (!v) return '—'
  return String(v).slice(0, len)
}


const community = useCommunity()
const api = useAdminApi()
const query = ref('')
const visibility = ref('all')
const loading = ref(false)
const error = ref('')
const moderating = ref('')

const loadPosts = async () => {
  loading.value = true
  error.value = ''
  try {
    await community.fetchAdminPosts({ query: query.value, visibility: visibility.value })
  } catch (err: any) {
    error.value = err?.data?.error || err?.message || '加载失败'
  } finally {
    loading.value = false
  }
}

const moderate = async (postId: string, action: 'hide' | 'limit' | 'restore') => {
  moderating.value = `${postId}:${action}`
  try {
    if (action === 'restore') {
      await api.post(`/admin/posts/${postId}/unhide`)
      await api.post(`/admin/posts/${postId}/unlimit`)
    } else {
      const endpoint = action === 'hide' ? 'hide' : 'limit'
      await api.post(`/admin/posts/${postId}/${endpoint}`)
    }
    await loadPosts()
  } catch (err: any) {
    error.value = err?.data?.error || err?.message || '操作失败'
  } finally {
    moderating.value = ''
  }
}

const isModerating = (postId: string, action: string) => moderating.value === `${postId}:${action}`

const formatTime = (value?: string) => {
  if (!value) return ''
  return dayjs(value).format('YYYY-MM-DD HH:mm')
}

onMounted(loadPosts)
</script>
