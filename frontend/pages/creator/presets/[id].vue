<template>
  <section class="h-[calc(100vh-6rem)] flex flex-col lg:flex-row gap-6 text-slate-100 overflow-hidden">
    <!-- Left Sidebar: Basic Info -->
    <aside class="w-full lg:w-80 flex-shrink-0 flex flex-col gap-4 bg-slate-900/80 p-6 rounded-3xl shadow-lg overflow-y-auto custom-scrollbar">
      <header>
        <h1 class="text-xl font-semibold">{{ isNew ? '新建预设' : '编辑预设' }}</h1>
        <p class="text-xs text-slate-400 mt-1">配置 Prompt 编排与模型参数</p>
      </header>

      <div class="space-y-4">
        <div>
          <label class="text-xs text-slate-400">预设名称</label>
          <input v-model.trim="form.name" class="mt-1 w-full rounded-xl border border-slate-700 bg-slate-900 px-4 py-2 text-sm text-white focus:border-primary focus:outline-none" required placeholder="如：通用越狱 V1" />
        </div>
        <div>
          <label class="text-xs text-slate-400">描述</label>
          <textarea v-model="form.description" class="mt-1 h-24 w-full rounded-xl border border-slate-700 bg-slate-900 px-4 py-2 text-sm text-white focus:border-primary focus:outline-none resize-none" placeholder="简要描述此预设的用途..." />
        </div>
        <div>
          <label class="text-xs text-slate-400">绑定模型 (可选)</label>
          <input v-model.trim="form.model_key" class="mt-1 w-full rounded-xl border border-slate-700 bg-slate-900 px-4 py-2 text-sm text-white focus:border-primary focus:outline-none" placeholder="如：gpt-4o" />
        </div>
        <div class="flex items-center gap-2">
          <input type="checkbox" v-model="form.is_public" id="is_public" class="rounded bg-slate-900 border-slate-700 text-primary focus:ring-primary" />
          <label for="is_public" class="text-sm text-slate-300">公开此预设</label>
        </div>
      </div>

      <div class="mt-auto space-y-3 pt-4 border-t border-slate-800">
        <p v-if="errorMessage" class="text-xs text-rose-500">{{ errorMessage }}</p>
        <p v-if="successMessage" class="text-xs text-emerald-500">{{ successMessage }}</p>
        <button
          @click="exportJSON"
          class="w-full rounded-xl bg-slate-800 border border-slate-700 py-2.5 text-sm font-medium text-white hover:bg-slate-700 transition-colors"
        >
          导出 JSON
        </button>
        <button
          @click="submit"
          class="w-full rounded-xl bg-primary py-2.5 text-sm font-medium text-white hover:bg-primary/90 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          :disabled="submitting"
        >
          {{ submitting ? '保存中...' : '保存预设' }}
        </button>
        <NuxtLink to="/creator/presets" class="block text-center text-xs text-slate-400 hover:text-white transition-colors">返回列表</NuxtLink>
      </div>
    </aside>

    <!-- Main Content: Block Editor -->
    <main class="flex-1 flex flex-col bg-slate-900/80 rounded-3xl shadow-lg overflow-hidden">
      <div class="flex items-center justify-between p-4 border-b border-slate-800 bg-slate-900/95 backdrop-blur z-10">
        <h3 class="text-lg font-medium text-white border-l-4 border-primary pl-3">Prompt 编排</h3>
        <button class="text-xs px-3 py-1.5 rounded-lg bg-indigo-500/20 text-indigo-300 hover:bg-indigo-500/30 transition-colors" @click="addBlock">+ 添加 Block</button>
      </div>

      <div class="flex-1 overflow-y-auto custom-scrollbar p-6 space-y-3">
        <div class="rounded-2xl border border-indigo-800/50 bg-indigo-900/30 text-slate-200 p-4 space-y-2">
          <div class="flex items-center gap-2 text-sm font-medium text-indigo-200">
            <span class="i-ph-magic-wand-duotone" />
            正则折叠内部标签（防止隐藏内容泄漏）
          </div>
          <p class="text-xs text-slate-300 leading-relaxed">
            在角色卡或 Block 中，如需隐藏内部备注/推理/系统信息，可用自定义标签包裹：
            例如 <code>&lt;think&gt;内部推理&lt;/think&gt;</code>、<code>&lt;note reason=\"cot\"&gt;作者备注&lt;/note&gt;</code>。
            系统会使用通用正则自动折叠这些内容，避免输出到用户端。
          </p>
          <div class="rounded-xl border border-slate-700 bg-slate-800/60 p-3 text-xs font-mono text-slate-200 overflow-x-auto">
            /&lt;([a-zA-Z0-9_-]+)(?:\s+[^&gt;]*)?&gt;[\s\S]*?&lt;\/\1&gt;/g
          </div>
          <p class="text-xs text-slate-400">
            支持任意标签名与属性，跨多行；不处理嵌套/缺失闭合标签。折叠后内容会显示为 <code>[folded]</code>。
          </p>
        </div>
        <div v-if="loading" class="text-center text-slate-400">加载中...</div>
        <div
          v-else
          v-for="(block, idx) in form.blocks"
          :key="idx"
          class="bg-slate-900/50 p-4 rounded-xl border border-slate-800 hover:border-slate-700 transition-colors"
        >
          <div class="flex items-center justify-between gap-4 mb-3">
            <div class="flex flex-col gap-1 flex-1">
              <input v-model="block.name" class="bg-transparent text-sm font-medium text-white border-b border-transparent focus:border-indigo-500 outline-none transition-colors w-full" placeholder="Block Name" />
              <div class="flex items-center gap-3">
                <div class="flex items-center gap-1 text-xs text-slate-500">
                  <span>ID:</span>
                  <input v-model="block.id" class="bg-transparent border-b border-transparent focus:border-indigo-500 outline-none w-24 transition-colors" placeholder="ID" />
                </div>
                <label class="flex items-center gap-1.5 text-xs text-slate-400 cursor-pointer hover:text-slate-300 transition-colors">
                  <input type="checkbox" v-model="block.marker" class="rounded bg-slate-800 border-slate-600 text-indigo-500 focus:ring-indigo-500/50" />
                  Marker
                </label>
              </div>
            </div>
            <div class="flex items-center gap-2">
              <div class="flex bg-slate-800 rounded-lg p-0.5">
                <button class="text-xs px-2 py-1 text-slate-400 hover:text-white hover:bg-slate-700 rounded transition-colors" @click="moveBlock(idx, 'up')">↑</button>
                <button class="text-xs px-2 py-1 text-slate-400 hover:text-white hover:bg-slate-700 rounded transition-colors" @click="moveBlock(idx, 'down')">↓</button>
              </div>
              <label class="flex items-center gap-2 text-xs text-slate-400 cursor-pointer select-none">
                <div class="relative inline-flex items-center cursor-pointer">
                  <input type="checkbox" v-model="block.enabled" class="sr-only peer">
                  <div class="w-9 h-5 bg-slate-700 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-indigo-600"></div>
                </div>
                <span class="w-6">{{ block.enabled ? '启用' : '禁用' }}</span>
              </label>
              <button class="text-slate-500 hover:text-rose-400 p-1.5 rounded-lg hover:bg-rose-500/10 transition-colors ml-1" @click="removeBlock(idx)">
                <span class="text-lg leading-none">×</span>
              </button>
            </div>
          </div>
          <textarea
            v-model="block.content"
            class="w-full rounded-lg border border-slate-700 bg-slate-900 px-3 py-2 text-sm text-white focus:border-primary focus:outline-none transition-colors"
            rows="3"
            :placeholder="block.marker ? '此 Block 为 Marker 占位符，运行时将被替换为动态内容' : '在此输入静态 Prompt 内容...'"
          />
        </div>
      </div>
    </main>
  </section>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useApi } from '@/composables/useApi'

definePageMeta({
  layout: 'default',
  middleware: ['auth']
})

const route = useRoute()
const router = useRouter()
const api = useApi()

const id = computed(() => route.params.id as string)
const isNew = computed(() => !id.value || id.value === 'new')

const loading = ref(false)
const submitting = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

const form = reactive({
  name: '',
  description: '',
  model_key: '',
  is_public: false,
  blocks: [] as any[],
  gen_params: {}
})

const load = async () => {
  if (isNew.value) {
    // 不预置 Block 名称，留给创作者自定义
    form.blocks = [{ id: '', name: '', role: 'system', content: '', enabled: true, marker: false }]
    return
  }
  loading.value = true
  try {
    const res = await api.get<{ data: any }>(`/presets/${id.value}`)
    const p = res.data
    form.name = p.name
    form.description = p.description
    form.model_key = p.model_key
    form.is_public = p.is_public
    form.blocks = p.blocks || []
    form.gen_params = p.gen_params || {}
  } catch (e: any) {
    errorMessage.value = e.message || '加载失败'
  } finally {
    loading.value = false
  }
}

const exportJSON = () => {
  const data = {
    version: 'nebula-preset-1',
    meta: {
      name: form.name || '未命名预设',
      description: form.description,
      author_id: 'user', // Placeholder
    },
    model: {
      key: form.model_key,
      params: form.gen_params,
    },
    blocks: form.blocks.map(b => ({
      id: b.id,
      name: b.name,
      role: b.role,
      content: b.content,
      enabled: b.enabled,
      marker: b.marker
    })),
    order: form.blocks.map(b => b.id),
  }
  const dataStr = JSON.stringify(data, null, 2)
  const blob = new Blob([dataStr], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `${form.name || 'preset'}.json`
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
}

const submit = async () => {
  submitting.value = true
  errorMessage.value = ''
  successMessage.value = ''
  try {
    const payload = { ...form }
    if (isNew.value) {
      const res = await api.post<{ data: any }>('/presets', payload)
      successMessage.value = '创建成功'
      router.replace(`/creator/presets/${res.data.id}`)
    } else {
      await api.put(`/presets/${id.value}`, payload)
      successMessage.value = '保存成功'
    }
  } catch (e: any) {
    errorMessage.value = e.message || '保存失败'
  } finally {
    submitting.value = false
  }
}

const addBlock = () => {
  form.blocks.push({
    id: '',
    name: '',
    role: 'system',
    content: '',
    enabled: true,
    marker: false
  })
}

const removeBlock = (idx: number) => {
  form.blocks.splice(idx, 1)
}

const moveBlock = (idx: number, direction: 'up' | 'down') => {
  const target = direction === 'up' ? idx - 1 : idx + 1
  if (target < 0 || target >= form.blocks.length) return
  const newBlocks = [...form.blocks]
  ;[newBlocks[idx], newBlocks[target]] = [newBlocks[target], newBlocks[idx]]
  form.blocks = newBlocks
}

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
