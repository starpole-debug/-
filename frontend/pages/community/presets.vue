<template>
  <section class="space-y-8">
    <header class="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
      <div>
        <p class="text-sm text-indigo-300 uppercase tracking-wide">Preset Market</p>
        <h1 class="text-3xl font-bold text-white">预设市场</h1>
        <p class="text-slate-400 text-sm mt-1">浏览社区发布的 Prompt 预设，直接使用或再编辑</p>
      </div>
      <NuxtLink to="/creator/presets" class="btn-primary text-sm font-medium">我的预设</NuxtLink>
    </header>

    <div class="bg-slate-900/60 border border-slate-800 rounded-3xl p-6 min-h-[60vh]">
      <div v-if="loading" class="text-center text-slate-400 py-10">加载中...</div>
      <div v-else-if="presets.length === 0" class="text-center text-slate-500 py-10">
        <p>暂无已发布的预设</p>
        <p class="text-xs mt-2">成为第一个发布的人吧</p>
      </div>
      <div v-else class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        <div
          v-for="preset in presets"
          :key="preset.id"
          class="group relative bg-slate-800/60 border border-slate-700/50 rounded-2xl p-4 hover:border-primary/50 transition-colors"
        >
          <div class="flex items-start justify-between gap-3 mb-2">
            <div>
              <h3 class="text-lg font-semibold text-white truncate">{{ preset.name }}</h3>
              <p class="text-xs text-slate-500">{{ preset.model_key || '通用模型' }}</p>
            </div>
            <span class="text-[11px] text-slate-500">{{ formatDate(preset.updated_at) }}</span>
          </div>
          <p class="text-sm text-slate-300 line-clamp-3 min-h-[60px]">{{ preset.description || '暂无描述' }}</p>
          <div class="mt-4 flex items-center justify-between text-xs text-slate-500">
            <span class="inline-flex items-center gap-1 px-2 py-1 rounded-full bg-indigo-500/10 text-indigo-300 border border-indigo-500/30">
              社区发布
            </span>
            <NuxtLink :to="`/creator/presets/download/${preset.id}`" class="text-indigo-300 hover:text-indigo-100 text-sm font-medium flex items-center gap-1">
              下载 / 详情
              <span class="text-xs">→</span>
            </NuxtLink>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useApi } from '@/composables/useApi'

definePageMeta({
  layout: 'default',
})

const api = useApi()
const presets = ref<any[]>([])
const loading = ref(false)

const load = async () => {
  loading.value = true
  try {
    const res = await api.get<{ data: any[] }>('/presets/public')
    presets.value = res.data || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const formatDate = (dateStr: string) => new Date(dateStr).toLocaleDateString()

onMounted(load)
</script>
