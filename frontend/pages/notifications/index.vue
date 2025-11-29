<template>
  <section class="rounded-3xl bg-white p-8 shadow">
    <header class="flex items-center justify-between">
      <h1 class="text-xl font-semibold">通知</h1>
      <button class="text-sm text-primary" :disabled="loading" @click="markAll">
        {{ loading ? '处理中...' : '全部已读' }}
      </button>
    </header>

    <p v-if="error" class="mt-3 text-sm text-rose-500">
      {{ error }}
    </p>

    <ul class="mt-4 space-y-3 text-sm text-slate-700">
      <li v-for="item in notifications" :key="item.id" class="rounded-xl border border-slate-200 p-4">
        <p class="font-medium">{{ item.title }}</p>
        <p class="text-xs text-muted">{{ item.content }}</p>
      </li>
      <li v-if="!notifications.length && !loading" class="text-xs text-muted">
        暂无通知
      </li>
    </ul>
  </section>
</template>

<script setup lang="ts">
import type { Notification } from '@/types'

const api = useApi()

const notifications = ref<Notification[]>([])
const loading = ref(false)
const error = ref('')

let destroyed = false
onBeforeUnmount(() => {
  destroyed = true
})

const isUnauthorized = (err: any) => {
  const status = err?.statusCode ?? err?.response?.status ?? err?.status
  return status === 401
}

const load = async () => {
  loading.value = true
  error.value = ''
  try {
    const data = await api.get<{ data: Notification[] }>('/notifications')
    if (destroyed) return
    notifications.value = Array.isArray(data.data) ? data.data : []
  } catch (err: any) {
    if (isUnauthorized(err)) {
      // 未登录：这里选择静默，不抛错，不打断路由
      return
    }
    console.error('加载通知失败:', err)
    error.value = err?.data?.error || err?.message || '加载通知失败'
  } finally {
    if (!destroyed) {
      loading.value = false
    }
  }
}

const markAll = async () => {
  loading.value = true
  error.value = ''
  try {
    await api.post('/notifications/all/read', {})
    await load()
  } catch (err: any) {
    if (isUnauthorized(err)) {
      return
    }
    console.error('全部已读操作失败:', err)
    error.value = err?.data?.error || err?.message || '操作失败'
  } finally {
    if (!destroyed) {
      loading.value = false
    }
  }
}

onMounted(load)
</script>
