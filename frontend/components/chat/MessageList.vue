<template>
  <div class="flex h-full flex-col gap-4 overflow-y-auto rounded-2xl bg-slate-50 p-4">
    <div v-for="message in messages" :key="message.id" class="flex" :class="message.sender_type === 'user' ? 'justify-end' : 'justify-start'">
      <div :class="[message.sender_type === 'user' ? 'bg-primary text-white' : 'bg-white text-slate-800', 'max-w-xl rounded-2xl px-4 py-2 shadow']">
        <p class="text-xs text-white/70" v-if="message.sender_type === 'user'">我</p>
        <p class="text-xs text-primary/70" v-else>AI</p>
        <p class="mt-1 text-sm leading-relaxed">{{ message.content }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { ChatMessage } from '@/types'

const props = defineProps<{ messages: ChatMessage[] }>()
</script>
