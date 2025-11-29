<template>
  <form class="flex items-center gap-3 rounded-2xl border border-slate-200 p-3" @submit.prevent="onSubmit">
    <textarea
      v-model="draft"
      class="h-16 flex-1 resize-none bg-transparent text-sm outline-none disabled:opacity-60"
      placeholder="输入新的问题或指令..."
      :disabled="props.disabled"
    ></textarea>
    <button
      type="submit"
      class="rounded-full bg-primary px-4 py-2 text-sm text-white disabled:cursor-not-allowed disabled:bg-slate-400"
      :disabled="props.disabled"
    >
      {{ props.disabled ? '发送中...' : '发送' }}
    </button>
  </form>
</template>

<script setup lang="ts">
const emit = defineEmits<{ (e: 'send', value: string): void }>()
const props = defineProps<{ disabled?: boolean }>()
const draft = ref('')

const onSubmit = () => {
  if (props.disabled) return
  if (!draft.value.trim()) return
  emit('send', draft.value)
  draft.value = ''
}
</script>
