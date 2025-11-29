<template>
  <div class="flex h-full items-center justify-center bg-[#0b0b12] text-gray-200">
    <div v-if="loading" class="flex flex-col items-center gap-3">
      <div class="i-svg-spinners-90-ring-with-bg text-2xl text-pink-400" />
      <p class="text-sm text-gray-400">正在拉取会话...</p>
    </div>
    <div v-else class="text-center space-y-4">
      <p class="text-lg font-semibold">还没有会话</p>
      <p class="text-sm text-gray-400">去探索页面选择一个角色，然后开始聊天。</p>
      <NuxtLink
        to="/search"
        class="inline-flex items-center gap-2 px-4 py-2 rounded-lg bg-pink-600 hover:bg-pink-500 text-white text-sm transition-colors"
      >
        <div class="i-ph-compass" />
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
