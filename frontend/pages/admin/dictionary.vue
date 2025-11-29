<template>
  <section>
    <h1 class="text-xl font-semibold">配置词典</h1>
    <p v-if="error" class="mt-2 text-sm text-rose-300">{{ error }}</p>
    <p v-else-if="loading" class="mt-2 text-sm text-slate-300">加载中...</p>
    <div class="mt-4 grid gap-4 sm:grid-cols-2">
      <div
        v-for="(item, index) in items"
        :key="item.id || item.label || index"
        class="rounded-2xl border border-white/10 p-4"
      >
        <p class="font-semibold">{{ item.label }}</p>
        <p class="text-xs text-slate-400">{{ item.description }}</p>
      </div>
    </div>
    <p v-if="!items.length && !loading && !error" class="mt-3 text-sm text-slate-400">暂无内容</p>
  </section>
</template>

<script setup lang="ts">
import type { DictionaryItem } from '@/types'

const { loadDictionary } = useAdminConfig()
const items = ref<DictionaryItem[]>([])
const loading = ref(false)
const error = ref('')

definePageMeta({
  layout: 'admin',
  middleware: 'admin',
})

const load = async () => {
  loading.value = true
  error.value = ''
  try {
    const result = await loadDictionary('style')
    items.value = Array.isArray(result) ? result.filter(Boolean) : []
  } catch (err: any) {
    error.value = err?.data?.error || err?.message || '加载失败'
    items.value = []
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>
