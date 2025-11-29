<template>
  <div class="rounded-xl border border-slate-200 bg-white p-4 shadow-sm">
    <p class="text-sm text-muted">{{ title }}</p>
    <p class="mt-2 text-2xl font-semibold">{{ value }}</p>
    <p v-if="hint" class="text-xs text-muted">{{ hint }}</p>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{ title: string; value: string | number; hint?: string }>()
</script>
