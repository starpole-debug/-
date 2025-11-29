<template>
  <div class="bento-card group h-64 flex flex-col justify-between p-6 cursor-pointer">
    <!-- Visual Area (Abstract Avatar) -->
    <div class="absolute top-0 right-0 w-32 h-32 bg-gradient-to-br from-indigo-500/20 to-purple-500/20 rounded-bl-full blur-2xl transition-all duration-500 group-hover:scale-150 group-hover:opacity-70"></div>
    
    <div class="relative z-10">
      <div class="h-14 w-14 rounded-2xl bg-gradient-to-br from-indigo-500 to-purple-600 text-white flex items-center justify-center text-2xl font-bold shadow-lg shadow-indigo-500/20 group-hover:scale-110 transition-transform duration-300 overflow-hidden">
        <img v-if="avatarUrl" :src="avatarUrl" class="w-full h-full object-cover" alt="avatar" />
        <span v-else>{{ (role.name || '?').charAt(0) }}</span>
      </div>
    </div>

    <div class="relative z-10 space-y-2">
      <div>
        <h3 class="text-xl font-bold text-white group-hover:text-indigo-300 transition-colors">{{ role.name || 'Unknown' }}</h3>
        <p class="text-sm text-slate-400 line-clamp-2 group-hover:text-slate-300 transition-colors">{{ role.description }}</p>
      </div>
      
      <div class="flex flex-wrap gap-2 pt-2">
        <span v-for="tag in (role.tags || []).slice(0, 3)" :key="tag" class="text-[10px] uppercase tracking-wider font-medium text-indigo-300 bg-indigo-500/10 px-2 py-1 rounded-md">
          {{ tag }}
        </span>
      </div>
    </div>

    <!-- Hover Action Overlay -->
    <div class="absolute inset-0 bg-black/60 backdrop-blur-[2px] opacity-0 group-hover:opacity-100 transition-opacity duration-300 flex items-center justify-center gap-4 z-20">
      <NuxtLink :to="`/role/${role.id}`" class="btn-primary text-sm py-2 px-4">查看详情</NuxtLink>
      <button @click.stop="handleChat" class="btn-secondary text-sm py-2 px-4 flex items-center gap-2">
        <div v-if="isLoading" class="i-svg-spinners-90-ring-with-bg" />
        <span v-else>开始对话</span>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Role } from '@/types'
import { useChat } from '@/composables/useChat'
import { useAssetUrl } from '~/composables/useAssetUrl'

const props = defineProps<{ role: Role; sessionId?: string }>()
const { createSession } = useChat()
const router = useRouter()
const isLoading = ref(false)
const { resolveAssetUrl } = useAssetUrl()
const avatarUrl = computed(() => resolveAssetUrl(props.role.avatar_url))

const handleChat = async () => {
  if (props.sessionId) {
    router.push(`/chat/${props.sessionId}`)
    return
  }
  
  isLoading.value = true
  try {
    const session = await createSession(props.role.id)
    router.push(`/chat/${session.id}`)
  } catch (e) {
    console.error('Failed to start chat', e)
  } finally {
    isLoading.value = false
  }
}
</script>
