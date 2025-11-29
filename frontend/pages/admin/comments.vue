<template>
  <section>
    <header class="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
      <div>
        <h1 class="text-xl font-semibold">评论管理</h1>
        <p class="text-sm text-slate-300">按帖子 ID 过滤，批量隐藏或删除评论。</p>
      </div>
      <div class="flex gap-2">
        <input v-model="postId" placeholder="帖子 ID" class="rounded-full border border-white/10 bg-white/5 px-4 py-2 text-sm" />
        <select v-model="visibility" class="rounded-full border border-white/10 bg-white/5 px-3 py-2 text-sm">
          <option value="all">全部</option>
          <option value="public">公开</option>
          <option value="hidden">隐藏</option>
          <option value="deleted">删除</option>
        </select>
        <button class="rounded-full bg-white/10 px-4 py-2 text-sm text-white" :disabled="loading" @click="loadComments">
          {{ loading ? '加载中...' : '刷新' }}
        </button>
      </div>
    </header>
    <p v-if="error" class="mt-3 text-sm text-rose-300">{{ error }}</p>

    <div class="mt-6 space-y-3">
      <article
        v-for="comment in comments"
        :key="comment.id"
        class="rounded-2xl border border-white/10 bg-white/5 p-4 text-sm text-slate-100"
      >
        <header class="flex flex-wrap items-center justify-between gap-2 text-xs text-slate-400">
          <span>#{{ comment.id.slice(0, 8) }} · 帖子 {{ comment.post_id.slice(0, 8) }}</span>
          <span>作者 {{ comment.author_id.slice(0, 6) }}</span>
          <span>{{ formatTime(comment.created_at) }}</span>
          <span class="rounded-full border border-white/10 px-2 py-0.5 text-[10px] uppercase">{{ comment.visibility }}</span>
        </header>
        <p class="mt-3 text-slate-100">{{ comment.content }}</p>
        <div class="mt-3 flex gap-3 text-xs">
          <button class="text-primary" :disabled="actionLoading === comment.id" @click="hide(comment.id)">
            {{ actionLoading === comment.id ? '隐藏中...' : '隐藏' }}
          </button>
          <button class="text-rose-400" :disabled="actionLoading === comment.id" @click="remove(comment.id)">
            {{ actionLoading === comment.id ? '删除中...' : '删除' }}
          </button>
        </div>
      </article>
      <p v-if="!comments.length && !loading" class="text-sm text-slate-400">暂无评论。</p>
    </div>
    </section>
</template>

<script setup lang="ts">
import dayjs from 'dayjs'
import type { CommunityComment } from '@/types'

definePageMeta({
  layout: 'admin',
  middleware: 'admin',
})


const api = useAdminApi()
const postId = ref('')
const visibility = ref('all')
const comments = ref<CommunityComment[]>([])
const loading = ref(false)
const error = ref('')
const actionLoading = ref('')

const loadComments = async () => {
  loading.value = true
  error.value = ''
  try {
    const search = new URLSearchParams()
    if (postId.value) search.set('post_id', postId.value)
    if (visibility.value && visibility.value !== 'all') search.set('visibility', visibility.value)
    const query = search.toString()
    const data = await api.get<{ data?: CommunityComment[] } | { data: CommunityComment[] } | CommunityComment[]>(
      `/admin/comments${query ? `?${query}` : ''}`,
    )
    const raw = (data as any)?.data ?? data
    comments.value = Array.isArray(raw) ? raw.filter(Boolean) : []
  } catch (err: any) {
    error.value = err?.data?.error || err?.message || '加载失败'
  } finally {
    loading.value = false
  }
}

const hide = async (id: string) => {
  actionLoading.value = id
  error.value = ''
  try {
    await api.post(`/admin/comments/${id}/hide`)
    await loadComments()
  } catch (err: any) {
    error.value = err?.data?.error || err?.message || '隐藏失败'
  } finally {
    actionLoading.value = ''
  }
}

const remove = async (id: string) => {
  if (!confirm('确认删除该评论？')) return
  actionLoading.value = id
  error.value = ''
  try {
    await api.del(`/admin/comments/${id}`)
    await loadComments()
  } catch (err: any) {
    error.value = err?.data?.error || err?.message || '删除失败'
  } finally {
    actionLoading.value = ''
  }
}

const formatTime = (value?: string) => {
  if (!value) return ''
  return dayjs(value).format('YYYY-MM-DD HH:mm')
}

onMounted(loadComments)
</script>
