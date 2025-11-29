<template>
  <div class="flex items-center justify-between rounded-xl border border-slate-200/60 bg-white/5 px-4 py-3">
    <div class="space-y-1">
      <p class="font-medium text-white">{{ role.name }}</p>
      <div class="flex items-center gap-2 text-xs">
        <span
          class="rounded-full px-2 py-0.5 border"
          :class="role.status === 'published' ? 'border-emerald-400 text-emerald-200' : 'border-amber-300 text-amber-100'"
        >
          {{ role.status === 'published' ? '已发布' : '草稿' }}
        </span>
        <span class="text-slate-400">状态：{{ role.status }}</span>
      </div>
    </div>
    <div class="flex items-center gap-3">
      <NuxtLink :to="`/creator/editor/${role.id}`" class="text-sm text-primary">编辑</NuxtLink>
      <button
        v-if="role.status !== 'published'"
        class="rounded-full bg-primary/80 px-3 py-1 text-xs text-white hover:bg-primary disabled:opacity-50"
        @click="$emit('publish', role.id)"
      >
        发布
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Role } from '@/types'
defineProps<{ role: Role }>()
defineEmits<{
  (e: 'publish', id: string): void
}>()
</script>
