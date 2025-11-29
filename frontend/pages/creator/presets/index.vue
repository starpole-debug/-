<template>
  <div class="h-[calc(100vh-6rem)] flex flex-col gap-6 text-slate-100">
    <header class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-semibold">预设工坊 (Preset Studio)</h1>
        <p class="text-sm text-slate-400 mt-1">管理您的 Prompt 预设，可独立于角色使用</p>
      </div>
      <div class="flex gap-3">
        <button class="px-4 py-2 rounded-xl bg-slate-800 border border-slate-700 hover:bg-slate-700 transition-colors text-sm" @click="importPreset">
          导入预设
        </button>
        <NuxtLink to="/creator/presets/new" class="px-4 py-2 rounded-xl bg-primary hover:bg-primary/90 transition-colors text-sm font-medium text-white">
          新建预设
        </NuxtLink>
      </div>
    </header>

    <div class="flex-1 overflow-y-auto custom-scrollbar bg-slate-900/50 rounded-3xl border border-slate-800 p-6">
      <div v-if="loading" class="text-center py-12 text-slate-400">加载中...</div>
      <div v-else-if="presets.length === 0" class="text-center py-12 text-slate-500">
        <p>暂无预设</p>
        <p class="text-xs mt-2">点击右上角“新建预设”开始创作</p>
      </div>
      <div v-else class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        <div
          v-for="preset in presets"
          :key="preset.id"
          class="bg-slate-800/50 border border-slate-700/50 rounded-xl p-4 hover:border-primary/50 transition-colors group cursor-pointer relative"
          @click="openPreset(preset.id)"
        >
          <div class="flex justify-between items-start mb-2">
            <h3 class="font-medium text-lg truncate">{{ preset.name }}</h3>
            <div class="flex gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
              <button class="p-1.5 text-slate-400 hover:text-white bg-slate-700/50 rounded-lg" title="导出" @click.stop="exportPreset(preset)">
                <span class="text-xs">⬇️</span>
              </button>
              <button
                class="p-1.5 text-slate-400 hover:text-emerald-400 bg-slate-700/50 rounded-lg"
                :title="preset.is_public ? '取消发布' : '发布到市场'"
                @click.stop="togglePublish(preset)"
              >
                <span class="text-xs">{{ preset.is_public ? '📤' : '🚀' }}</span>
              </button>
              <button class="p-1.5 text-slate-400 hover:text-rose-400 bg-slate-700/50 rounded-lg" title="删除" @click.stop="deletePreset(preset.id)">
                <span class="text-xs">🗑️</span>
              </button>
            </div>
          </div>
          <p class="text-sm text-slate-400 line-clamp-2 h-10 mb-3">{{ preset.description || '暂无描述' }}</p>
          <div class="flex items-center justify-between text-xs text-slate-500">
            <span class="flex items-center gap-2">
              <span>{{ preset.model_key || '通用' }}</span>
              <span v-if="preset.is_public" class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full bg-emerald-500/10 text-emerald-300 border border-emerald-500/30">
                已发布
              </span>
            </span>
            <span>{{ formatDate(preset.updated_at) }}</span>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Hidden file input for import -->
    <input type="file" ref="fileInput" class="hidden" accept=".json" @change="handleFileImport" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useApi } from '@/composables/useApi'

definePageMeta({
  layout: 'default',
  middleware: ['auth']
})

const api = useApi()
const router = useRouter()
const loading = ref(true)
const presets = ref<any[]>([])
const fileInput = ref<HTMLInputElement | null>(null)

const loadPresets = async () => {
  loading.value = true
  try {
    const res = await api.get<{ data: any[] }>('/presets')
    presets.value = res.data || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString()
}

const deletePreset = async (id: string) => {
  if (!confirm('确定要删除这个预设吗？')) return
  try {
    await api.del(`/presets/${id}`)
    await loadPresets()
  } catch (e) {
    alert('删除失败')
  }
}

const exportPreset = (preset: any) => {
  const dataStr = JSON.stringify(preset, null, 2)
  const blob = new Blob([dataStr], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `${preset.name || 'preset'}.json`
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
}

const importPreset = () => {
  fileInput.value?.click()
}

const handleFileImport = async (event: Event) => {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (!file) return
  
  const reader = new FileReader()
  reader.onload = async (e) => {
    try {
      const content = e.target?.result as string
      const data = JSON.parse(content)
      // Clean up ID to create new
      delete data.id
      delete data.created_at
      delete data.updated_at
      delete data.creator_id
      
      await api.post('/presets', data)
      await loadPresets()
      alert('导入成功')
    } catch (err) {
      console.error(err)
      alert('导入失败：格式错误')
    }
  }
  reader.readAsText(file)
  // Reset input
  if (fileInput.value) fileInput.value.value = ''
}

const togglePublish = async (preset: any) => {
  try {
    await api.post(`/presets/${preset.id}/publish`, { is_public: !preset.is_public })
    await loadPresets()
  } catch (e) {
    alert('发布操作失败')
  }
}

const openPreset = (id: string) => {
  router.push(`/creator/presets/download/${id}`)
}

onMounted(loadPresets)
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
