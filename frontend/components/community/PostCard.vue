<template>
  <article class="card-soft card-soft-hover group overflow-hidden cursor-pointer w-full max-w-[32rem] text-sm">
    <!-- Visual Header (Placeholder for Post Image) -->
    <div
      class="relative h-48 w-full overflow-hidden bg-bg-cream-200 bg-center bg-cover"
      :style="coverUrl ? { backgroundImage: `url(${coverUrl})` } : {}"
    >
      <div class="absolute inset-0 bg-gradient-to-br from-bg-cream-100/25 to-bg-pink-100/25 mix-blend-overlay"></div>
      <div class="absolute inset-0 bg-charcoal-900/5"></div>
      
      <!-- Hover Overlay -->
      <div class="absolute inset-0 bg-charcoal-900/10 opacity-0 group-hover:opacity-100 transition-opacity duration-300 flex items-center justify-center gap-4 backdrop-blur-[2px]">
        <button 
          @click.stop.prevent="handleLike"
          class="p-2 rounded-full backdrop-blur-md transition-all duration-300 active:scale-90"
          :class="isLiked ? 'bg-status-error text-white shadow-lg shadow-status-error/30' : 'bg-white/80 hover:bg-white text-charcoal-700'"
        >
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-5 h-5" :class="{ 'animate-bounce-short': isLiked }">
            <path d="M1 8.25a1.25 1.25 0 112.5 0v7.5a1.25 1.25 0 11-2.5 0v-7.5zM11 3V1.7c0-.268.14-.526.395-.607A2 2 0 0114 3c0 .995-.182 1.948-.514 2.826-.204.54.166 1.174.744 1.174h2.52c1.243 0 2.261 1.01 2.146 2.247a23.864 23.864 0 01-1.341 5.974C17.153 16.323 16.072 17 14.9 17h-3.192a3 3 0 01-1.341-.317l-2.734-1.366A3 3 0 006.292 15H5V8h.963c.685 0 1.258-.483 1.612-1.068a4.011 4.011 0 012.166-1.73c.432-.143.853-.386 1.011-.814.16-.432.448-.779.9-.976.45-.198.937-.312 1.348-.312z" />
          </svg>
        </button>
        <button 
          @click.stop.prevent="handleFavorite"
          class="p-2 rounded-full backdrop-blur-md transition-all duration-300 active:scale-90"
          :class="isFavorited ? 'bg-accent-yellow text-charcoal-900 shadow-lg shadow-accent-yellow/30' : 'bg-white/80 hover:bg-white text-charcoal-700'"
        >
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-5 h-5" :class="{ 'animate-bounce-short': isFavorited }">
            <path fill-rule="evenodd" d="M10 2c-1.716 0-3.408.106-5.07.31C3.806 2.45 3 3.414 3 4.517V17.25a.75.75 0 001.075.676L10 15.082l5.925 2.844A.75.75 0 0017 17.25V4.517c0-1.103-.806-2.068-1.93-2.207A41.403 41.403 0 0010 2z" clip-rule="evenodd" />
          </svg>
        </button>
      </div>
    </div>

    <div class="p-4 space-y-2.5">
      <header class="flex items-start justify-between gap-4">
        <h3 class="text-base font-display font-semibold text-charcoal-900 leading-snug group-hover:text-charcoal-700 transition-colors line-clamp-2">
          {{ post.title }}
        </h3>
      </header>

      <p class="text-[13px] text-charcoal-500 line-clamp-3 leading-relaxed">
        {{ post.content }}
      </p>

      <div v-if="post.link_url" class="flex items-center gap-2 text-xs">
        <button
          class="inline-flex items-center gap-1 px-2 py-1 rounded-lg bg-bg-cream-200 text-charcoal-700 border border-accent-yellow/30 hover:bg-accent-yellow/20 transition-colors"
          @click.stop.prevent="openLink"
        >
          <span>{{ linkLabel }}</span>
          <span class="text-[10px]">↗</span>
        </button>
      </div>

      <div class="flex flex-wrap gap-2 pt-1">
        <span v-for="topic in post.topic_ids?.slice(0, 2) || []" :key="topic" class="text-[10px] font-medium px-2 py-0.5 rounded bg-bg-cream-200 text-charcoal-500 border border-charcoal-200/50">
          #{{ topic }}
        </span>
      </div>

      <footer class="pt-2.5 mt-2.5 border-t border-charcoal-100 flex items-center justify-between text-[11px] text-charcoal-500">
        <div class="flex items-center gap-2">
          <img 
            v-if="post.author_avatar" 
            :src="authorAvatarUrl" 
            class="w-5 h-5 rounded-full object-cover border border-white shadow-sm"
            alt="Avatar"
          />
          <div v-else class="w-5 h-5 rounded-full bg-accent-yellow flex items-center justify-center text-[10px] text-charcoal-900 font-bold border border-white shadow-sm">
            {{ (post.author_name || post.author_id || '?').charAt(0).toUpperCase() }}
          </div>
          <NuxtLink
            v-if="post.author_id"
            :to="`/community/user/${post.author_id}`"
            class="hover:text-charcoal-900 transition font-medium"
          >
            {{ post.author_name || post.author_id.slice(0, 8) }}
          </NuxtLink>
          <span v-else>匿名用户</span>
        </div>
        <slot name="footer" />
      </footer>
    </div>
  </article>
</template>

<script setup lang="ts">
import type { CommunityPost, CommunityComment } from '@/types'
import { useAssetUrl } from '~/composables/useAssetUrl'
const router = useRouter()

const props = defineProps<{ post: CommunityPost }>()
const emit = defineEmits(['like', 'favorite'])
const { resolveAssetUrl } = useAssetUrl()

const isLiked = ref(false)
const isFavorited = ref(false)
const coverUrl = computed(() => resolveAssetUrl(props.post.attachments?.[0]?.file_url))
const authorAvatarUrl = computed(() => resolveAssetUrl(props.post.author_avatar))
const isExternal = computed(() => Boolean(props.post.link_url?.startsWith('http://') || props.post.link_url?.startsWith('https://')))
const linkLabel = computed(() => {
  if (props.post.link_type === 'role') return '查看角色'
  if (props.post.link_type === 'preset') return '查看预设'
  if (isExternal.value) return '打开外链'
  return '打开链接'
})

const handleLike = () => {
  isLiked.value = !isLiked.value
  emit('like')
}

const handleFavorite = () => {
  isFavorited.value = !isFavorited.value
  emit('favorite')
}

const openLink = () => {
  const href = props.post.link_url || ''
  if (!href) return
  if (isExternal.value) {
    window.open(href, '_blank', 'noopener')
    return
  }
  router.push(href)
}
</script>
