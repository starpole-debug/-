<template>
  <div class="grid gap-4 sm:grid-cols-3">
    <button
      v-for="amount in amounts"
      :key="amount"
      type="button"
      class="rounded-2xl border border-slate-200 p-4 text-center hover:border-primary"
      @click="() => emit('select', amount)"
    >
      <p class="text-2xl font-semibold">¥{{ amount }}</p>
      <p class="text-xs text-muted">一次性支持</p>
    </button>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{ amounts: number[] }>()
const emit = defineEmits<{ (e: 'select', amount: number): void }>()
</script>
