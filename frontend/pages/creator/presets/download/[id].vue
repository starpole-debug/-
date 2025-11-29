<template>
  <section class="h-[calc(100vh-6rem)] flex flex-col gap-6 text-slate-100">
    <header class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-semibold">预设详情</h1>
        <p class="text-sm text-slate-400 mt-1">下载或预览该预设，默认只读</p>
      </div>
      <div class="flex gap-3">
        <button
          class="px-4 py-2 rounded-xl bg-slate-800 border border-slate-700 hover:bg-slate-700 transition-colors text-sm"
          @click="exportPreset"
          :disabled="loading"
        >
          下载 JSON
        </button>
        <NuxtLink to="/creator/presets" class="px-4 py-2 rounded-xl bg-slate-800 border border-slate-700 hover:bg-slate-700 transition-colors text-sm">
          返回预设工坊
        </NuxtLink>
      </div>
    </header>

    <div class="flex-1 overflow-y-auto custom-scrollbar bg-slate-900/50 rounded-3xl border border-slate-800 p-6">
      <div v-if="loading" class="text-center text-slate-400 py-12">加载中...</div>
      <div v-else-if="error" class="text-center text-rose-400 py-12">{{ error }}</div>
      <div v-else class="space-y-6">
        <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-3">
          <div>
            <h2 class="text-xl font-semibold text-white">{{ preset.name }}</h2>
            <p class="text-sm text-slate-400 mt-1">{{ preset.description || '暂无描述' }}</p>
          </div>
          <div class="flex items-center gap-3 text-sm text-slate-400">
            <span class="px-3 py-1 rounded-full bg-white/5 border border-white/10">
              模型：{{ preset.model_key || '通用' }}
            </span>
            <span v-if="preset.is_public" class="px-3 py-1 rounded-full bg-emerald-500/10 border border-emerald-500/30 text-emerald-300">
              已发布
            </span>
            <span class="px-3 py-1 rounded-full bg-white/5 border border-white/10">更新于 {{ formatDate(preset.updated_at) }}</span>
          </div>
        </div>

        <div class="grid gap-4 md:grid-cols-2">
          <div class="bg-slate-800/50 border border-slate-700/50 rounded-2xl p-4">
            <h3 class="text-sm font-semibold text-slate-200 mb-2">生成参数</h3>
            <pre class="text-xs text-slate-300 bg-black/30 rounded-xl p-3 overflow-x-auto">{{ formattedGenParams }}</pre>
          </div>
          <div class="bg-slate-800/50 border border-slate-700/50 rounded-2xl p-4">
            <h3 class="text-sm font-semibold text-slate-200 mb-2">元信息</h3>
            <ul class="text-sm text-slate-300 space-y-1">
              <li><span class="text-slate-500">ID：</span>{{ preset.id }}</li>
              <li><span class="text-slate-500">创建者：</span>{{ preset.creator_id || '匿名' }}</li>
              <li><span class="text-slate-500">创建时间：</span>{{ formatDate(preset.created_at) }}</li>
            </ul>
          </div>
        </div>

        <div class="bg-slate-800/50 border border-slate-700/50 rounded-2xl p-4">
          <div class="flex items-center justify-between mb-3">
            <h3 class="text-sm font-semibold text-slate-200">Blocks</h3>
            <span class="text-xs text-slate-500">只读</span>
          </div>
          <div class="space-y-3">
            <div
              v-for="(block, idx) in preset.blocks || []"
              :key="idx"
              class="rounded-xl border border-slate-700 bg-slate-900/60 p-3"
            >
              <div class="flex items-center gap-2 text-sm text-slate-300">
                <span class="font-semibold text-white">{{ block.name || block.id || '未命名 Block' }}</span>
                <span class="px-2 py-0.5 rounded-full bg-white/5 text-xs text-slate-400 border border-white/10">{{ block.role || 'system' }}</span>
                <span
                  v-if="block.marker"
                  class="px-2 py-0.5 rounded-full bg-amber-500/10 text-xs text-amber-300 border border-amber-500/30"
                >
                  Marker
                </span>
                <span
                  class="px-2 py-0.5 rounded-full text-xs border"
                  :class="block.enabled ? 'bg-emerald-500/10 text-emerald-300 border-emerald-500/30' : 'bg-slate-700/40 text-slate-400 border-slate-600/50'"
                >
                  {{ block.enabled ? '启用' : '禁用' }}
                </span>
              </div>
              <p class="mt-2 text-sm text-slate-200 whitespace-pre-wrap break-words">{{ block.content || '（空）' }}</p>
            </div>
          </div>
        </div>

        <div v-if="isOwner" class="flex items-center gap-3">
          <NuxtLink :to="`/creator/presets/${id}`" class="btn-secondary text-sm">
            进入编辑（仅作者）
          </NuxtLink>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useApi } from '@/composables/useApi'
import { useAuth } from '@/composables/useAuth'

definePageMeta({
  layout: 'default',
})

const route = useRoute()
const api = useApi()
const auth = useAuth()
const id = computed(() => route.params.id as string)

const preset = ref<any>({})
const loading = ref(false)
const error = ref('')

const load = async () => {
  loading.value = true
  error.value = ''
  try {
    const res = await api.get<{ data: any }>(`/presets/${id.value}`)
    preset.value = res.data || {}
  } catch (e: any) {
    error.value = e?.message || '加载失败'
  } finally {
    loading.value = false
  }
}

const exportPreset = () => {
  const dataStr = JSON.stringify(preset.value, null, 2)
  const blob = new Blob([dataStr], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `${preset.value.name || 'preset'}.json`
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
}

const formatDate = (str?: string) => (str ? new Date(str).toLocaleDateString() : '')
const formattedGenParams = computed(() => JSON.stringify(preset.value.gen_params || {}, null, 2))
const isOwner = computed(() => preset.value?.creator_id && auth.user.value?.id === preset.value.creator_id)

onMounted(load)
</script>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #334155;
  border-radius: 3px;
}
</style>
