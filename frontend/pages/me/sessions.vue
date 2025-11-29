<template>
  <section>
    <h1 class="text-xl font-semibold">历史会话</h1>
    <p v-if="errorMessage" class="mt-2 text-sm text-rose-500">
      {{ errorMessage }}
    </p>
    <div class="mt-4 space-y-3">
      <div
        v-for="session in chat.sessions"
        :key="session.id"
        class="flex items-center justify-between rounded-xl border border-slate-200 p-4"
      >
        <div>
          <p class="font-medium">{{ session.title }}</p>
          <p class="text-xs text-muted">{{ session.created_at }}</p>
        </div>
        <NuxtLink :to="`/chat/${session.id}`" class="text-sm text-primary">继续</NuxtLink>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: ['auth'],
})

const chat = useChat()
const errorMessage = ref('')

const loadSessions = async () => {
  errorMessage.value = ''
  try {
    await chat.refreshSessions()
  } catch (error: any) {
    errorMessage.value = error?.message || '无法加载历史会话，请稍后重试'
  }
}

onMounted(loadSessions)
</script>

