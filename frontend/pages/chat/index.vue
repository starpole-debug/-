<template>
  <div class="flex h-full w-full items-center justify-center bg-bg-cream-100 text-charcoal-900">
    <div v-if="loading" class="flex flex-col items-center gap-3">
      <div class="w-8 h-8 border-2 border-accent-yellow/30 border-t-accent-yellow rounded-full animate-spin" />
      <p class="text-sm text-charcoal-500 font-medium">正在拉取会话...</p>
    </div>
    <div v-else class="text-center space-y-6 max-w-md px-6">
      <div class="w-24 h-24 mx-auto bg-white rounded-full shadow-lg shadow-accent-yellow/20 flex items-center justify-center">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" class="w-12 h-12 text-accent-yellow">
          <path d="M4.913 2.658c2.075-.27 4.19-.408 6.337-.408 2.147 0 4.262.139 6.337.408 1.922.25 3.291 1.861 3.405 3.727a4.403 4.403 0 001.032 2.211c1.017 1.2 1.017 2.936 0 4.136a4.402 4.402 0 00-1.032 2.211c-.114 1.866-1.483 3.477-3.405 3.727a68.766 68.766 0 01-12.674 0c-1.922-.25-3.291-1.861-3.405-3.727a4.402 4.402 0 00-1.032-2.211c-1.017-1.2-1.017-2.936 0-4.136a4.402 4.402 0 001.032-2.211c.114-1.866 1.483-3.477 3.405-3.727z" />
        </svg>
      </div>
      <div class="space-y-2">
        <h2 class="text-2xl font-display font-bold text-charcoal-900">还没有会话</h2>
        <p class="text-charcoal-500">去探索页面选择一个角色，然后开始聊天。</p>
      </div>
      <NuxtLink
        to="/search"
        class="btn-primary inline-flex items-center gap-2 shadow-lg shadow-charcoal-900/10"
      >
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-5 h-5">
          <path fill-rule="evenodd" d="M9 3.5a5.5 5.5 0 100 11 5.5 5.5 0 000-11zM2 9a7 7 0 1112.452 4.391l3.328 3.329a.75.75 0 11-1.06 1.06l-3.329-3.328A7 7 0 012 9z" clip-rule="evenodd" />
        </svg>
        去探索角色
      </NuxtLink>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'chat',
  middleware: ['auth'],
})

const router = useRouter()
const loading = ref(true)
const { sessions, listSessions } = useChat()

onMounted(async () => {
  try {
    await listSessions()
    if (sessions.value.length > 0) {
      router.replace(`/chat/${sessions.value[0].id}`)
      return
    }
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
})
</script>
