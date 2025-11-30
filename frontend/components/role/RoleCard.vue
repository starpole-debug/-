<template>
  <NuxtLink :to="`/role/${role.id}`" class="card-soft card-soft-hover group block h-full overflow-hidden">
    <div class="relative aspect-[16/9] overflow-hidden">
      <img 
        :src="avatarUrl" 
        :alt="role.name"
        class="h-full w-full object-cover transition-transform duration-700 group-hover:scale-110"
        loading="lazy"
      />
      <div class="absolute inset-0 bg-gradient-to-t from-black/20 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
    </div>
    
    <div class="p-4 space-y-2">
      <div class="flex items-center justify-between mb-1.5">
        <h3 class="text-lg font-display font-semibold text-charcoal-900 line-clamp-2">{{ role.name }}</h3>
        <div class="flex items-center gap-1 text-[11px] font-medium text-charcoal-500 bg-bg-cream-200 px-2 py-1 rounded-full">
          <span class="i-ph-heart-fill text-accent-pink-soft text-sm"></span>
          {{ role.favorite_count || 0 }}
        </div>
      </div>
      <p class="text-sm text-charcoal-500 line-clamp-2 leading-relaxed">
        {{ role.description }}
      </p>
    </div>
  </NuxtLink>
</template>

<script setup lang="ts">
import type { Role } from '../../types'
import { useAssetUrl } from '~/composables/useAssetUrl'

const props = defineProps<{
  role: Role
  sessionId?: string
}>()

const router = useRouter()
const { resolveAssetUrl } = useAssetUrl()
const avatarUrl = computed(() => resolveAssetUrl(props.role.avatar_url))
</script>
