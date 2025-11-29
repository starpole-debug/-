<template>
  <section v-if="role" class="space-y-6">
    <header class="rounded-3xl bg-gradient-to-br from-slate-900 to-indigo-900 text-white p-8 shadow-xl border border-white/10">
      <div class="flex flex-col lg:flex-row gap-8 items-start">
        <div class="w-28 h-28 rounded-2xl bg-white/5 border border-white/10 shadow-lg overflow-hidden flex items-center justify-center text-2xl font-bold">
          <img v-if="avatarUrl" :src="avatarUrl" class="w-full h-full object-cover" alt="角色头像" />
          <span v-else>{{ (role.name || '?').charAt(0) }}</span>
        </div>
        <div class="flex-1 space-y-3">
          <p class="text-sm uppercase tracking-wide text-indigo-200/80">AI Persona</p>
          <h1 class="text-4xl font-bold leading-tight">{{ role.name }}</h1>
          <p class="text-slate-200/80 text-base max-w-3xl">{{ role.description || '创作者正在补充描述。' }}</p>
          <div class="flex flex-wrap gap-2">
            <span v-for="tag in role.tags || []" :key="tag" class="text-xs rounded-full bg-white/10 border border-white/10 px-3 py-1 text-indigo-100 uppercase tracking-wide">
              #{{ tag }}
            </span>
          </div>
          <div class="flex flex-col sm:flex-row gap-3 pt-2">
            <button class="rounded-full bg-primary px-6 py-3 text-white text-sm shadow-lg shadow-indigo-500/30 hover:translate-y-[-2px] transition-transform" @click="startChat">立即对话</button>
            <button
              class="rounded-full border border-white/30 px-6 py-3 text-sm text-white flex items-center gap-2 hover:border-primary transition-colors"
              @click="toggleFavorite"
            >
              <span v-if="isFavorited">💖 已收藏</span>
              <span v-else>🤍 收藏角色</span>
            </button>
            <NuxtLink to="/community/new" class="rounded-full border border-white/30 px-6 py-3 text-sm text-white hover:border-primary transition-colors">写一篇互动故事</NuxtLink>
          </div>
          <p v-if="errorMessage" class="text-sm text-rose-300">
            {{ errorMessage }}
          </p>
        </div>
      </div>
    </header>
    <section class="rounded-3xl bg-white/5 border border-white/10 p-6 shadow text-slate-100">
      <h2 class="text-xl font-semibold text-white">行为设定</h2>
      <p class="mt-2 text-sm text-slate-300">{{ abilitiesText }}</p>
    </section>
  </section>
  <p v-else class="text-center text-muted">载入中...</p>
</template>

<script setup lang="ts">
import { useAssetUrl } from '~/composables/useAssetUrl'

const route = useRoute()
const router = useRouter()
const { fetchRole, favoriteRole, unfavoriteRole, fetchFavorites } = useRole()
const chat = useChat()
const role = ref()
const errorMessage = ref('')
const isFavorited = computed(() => !!role.value?.is_favorited)
const { resolveAssetUrl } = useAssetUrl()
const avatarUrl = computed(() => resolveAssetUrl(role.value?.avatar_url))
const abilitiesText = computed(() => {
  const list = role.value?.abilities || []
  if (Array.isArray(list) && list.length) return list.join(' / ')
  return '创作者正在补充更详细的 Prompt 片段。'
})

const load = async () => {
  role.value = await fetchRole(route.params.id as string)
}

const startChat = async () => {
  if (!role.value) return
  errorMessage.value = ''
  try {
    const session = await chat.createSession(role.value.id)
    router.push(`/chat/${session.id}`)
  } catch (error: any) {
    errorMessage.value = error?.message || '无法创建会话，请稍后重试'
  }
}

const toggleFavorite = async () => {
  if (!role.value) return
  try {
    if (isFavorited.value) {
      await unfavoriteRole(role.value.id)
      role.value = { ...role.value, is_favorited: false }
    } else {
      await favoriteRole(role.value.id, role.value)
      role.value = { ...role.value, is_favorited: true }
    }
    await fetchFavorites()
  } catch (error: any) {
    errorMessage.value = error?.message || '收藏失败，请稍后重试'
  }
}

onMounted(load)
</script>
