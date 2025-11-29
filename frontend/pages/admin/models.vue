<template>
  <section class="space-y-6">
    <header class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
      <div>
        <h1 class="text-xl font-semibold">模型配置</h1>
        <p class="text-sm text-slate-300">配置不同 provider、API Base、温度、最大 token 等参数。</p>
      </div>
      <button class="rounded-full bg-white/10 px-4 py-2 text-sm text-white" @click="startCreate">
        {{ showForm ? '取消' : '新增模型' }}
      </button>
    </header>

    <form v-if="showForm" class="space-y-4 rounded-2xl border border-white/10 p-6" @submit.prevent="submit">
      <div class="grid gap-4 md:grid-cols-2">
        <label class="text-sm">
          名称
          <input v-model.trim="form.name" class="mt-1 w-full rounded-xl border border-white/10 bg-white/5 px-4 py-2 text-sm" required />
        </label>
        <label class="text-sm">
          Provider
          <select v-model="form.provider" class="mt-1 w-full rounded-xl border border-white/10 bg-slate-900 text-white px-4 py-2 text-sm">
            <option value="openai">OpenAI 兼容</option>
            <option value="deepseek">DeepSeek</option>
            <option value="custom">Custom</option>
            <option value="mock">Mock</option>
          </select>
        </label>
        <label class="text-sm md:col-span-2">
          描述
          <textarea v-model.trim="form.description" class="mt-1 w-full rounded-xl border border-white/10 bg-white/5 px-4 py-2 text-sm" rows="2" />
        </label>
        <label class="text-sm md:col-span-2">
          Base URL
          <input v-model.trim="form.base_url" class="mt-1 w-full rounded-xl border border-white/10 bg-white/5 px-4 py-2 text-sm" />
        </label>
        <label class="text-sm">
          模型名称
          <input v-model.trim="form.model_name" class="mt-1 w-full rounded-xl border border-white/10 bg-white/5 px-4 py-2 text-sm" required />
        </label>
        <label class="text-sm">
          API Key（留空表示沿用原值）
          <input v-model="form.api_key" type="password" autocomplete="off" class="mt-1 w-full rounded-xl border border-white/10 bg-white/5 px-4 py-2 text-sm" />
        </label>
        <label class="text-sm">
          温度
          <input v-model.number="form.temperature" type="number" step="0.1" min="0" max="2" class="mt-1 w-full rounded-xl border border-white/10 bg-white/5 px-4 py-2 text-sm" />
        </label>
        <label class="text-sm">
          最大 Tokens
          <input v-model.number="form.max_tokens" type="number" min="0" class="mt-1 w-full rounded-xl border border-white/10 bg-white/5 px-4 py-2 text-sm" />
        </label>
        <label class="text-sm">
          每次调用价格（平台币）
          <input v-model.number="form.price_coins" type="number" min="0" class="mt-1 w-full rounded-xl border border-white/10 bg-white/5 px-4 py-2 text-sm" />
        </label>
        <label class="text-sm">
          创作者分成（角色）
          <input v-model.number="form.share_role_pct" type="number" min="0" max="1" step="0.01" class="mt-1 w-full rounded-xl border border-white/10 bg-white/5 px-4 py-2 text-sm" placeholder="0.2 表示 20%" />
        </label>
        <label class="text-sm">
          预设作者分成
          <input v-model.number="form.share_preset_pct" type="number" min="0" max="1" step="0.01" class="mt-1 w-full rounded-xl border border-white/10 bg-white/5 px-4 py-2 text-sm" placeholder="0.1 表示 10%" />
        </label>
        <label class="text-sm">
          价格备注
          <input v-model.trim="form.price_hint" class="mt-1 w-full rounded-xl border border-white/10 bg-white/5 px-4 py-2 text-sm" placeholder="如：按次/按千 tokens 描述" />
        </label>
        <label class="text-sm">
          状态
          <select v-model="form.status" class="mt-1 w-full rounded-xl border border-white/10 bg-slate-900 text-white px-4 py-2 text-sm">
            <option value="active">Active</option>
            <option value="inactive">Inactive</option>
          </select>
        </label>
        <label class="flex items-center gap-2 text-sm">
          <input v-model="form.is_default" type="checkbox" class="h-4 w-4" />
          设为默认模型
        </label>
      </div>
      <div class="flex gap-3">
        <button type="submit" class="rounded-full bg-primary px-5 py-2 text-sm text-white" :disabled="saving">
          {{ saving ? '保存中...' : editingId ? '保存修改' : '创建模型' }}
        </button>
        <button type="button" class="rounded-full border border-white/20 px-5 py-2 text-sm" @click="cancelForm">取消</button>
      </div>
    </form>

    <p v-if="error" class="text-sm text-rose-300">{{ error }}</p>
    <p v-else-if="loading" class="text-sm text-slate-300">加载中...</p>
    <div class="overflow-auto rounded-2xl border border-white/10">
      <table class="min-w-full text-left text-sm">
        <thead class="bg-white/5 text-xs uppercase text-slate-400">
          <tr>
            <th class="px-4 py-3">名称</th>
            <th>Provider</th>
            <th>模型名</th>
            <th>单次价格(币)</th>
            <th>分成(角色/预设)</th>
            <th>状态</th>
            <th>默认</th>
            <th>API Key</th>
            <th></th>
          </tr>
        </thead>
        <tbody v-if="models && models.length">
          <tr v-for="model in models" :key="model.id" class="border-t border-white/10">
            <td class="px-4 py-3">
              <p class="font-medium">{{ model.name }}</p>
              <p class="text-xs text-slate-400">{{ model.description }}</p>
            </td>
            <td>{{ model.provider }}</td>
            <td>{{ model.model_name }}</td>
            <td>{{ model.price_coins || 0 }}</td>
            <td>{{ (model.share_role_pct || 0) }} / {{ (model.share_preset_pct || 0) }}</td>
            <td>
              <span
                class="rounded-full px-3 py-1 text-xs"
                :class="model.status === 'active' ? 'bg-emerald-500/20 text-emerald-200' : 'bg-rose-500/20 text-rose-200'"
              >
                {{ model.status }}
              </span>
            </td>
            <td>{{ model.is_default ? 'Yes' : 'No' }}</td>
            <td>{{ model.has_api_key ? '已配置' : '未配置' }}</td>
            <td class="space-x-3">
              <button class="text-xs text-primary" @click="startEdit(model)">编辑</button>
              <button class="text-xs text-rose-400" @click="() => remove(model.id)">删除</button>
            </td>
          </tr>
        </tbody>
        <tbody v-else>
          <tr>
            <td class="p-4 text-sm text-slate-400" colspan="8">暂无配置，请先创建一个模型。</td>
          </tr>
        </tbody>
      </table>
    </div>
    </section>
</template>

<script setup lang="ts">
import type { ModelConfig } from '@/types'

const { models, loadModels, saveModel, deleteModel } = useAdminConfig()
const loading = ref(false)
const error = ref('')

definePageMeta({
  layout: 'admin',
  middleware: 'admin',
})

const defaultForm = () => ({
  name: '',
  description: '',
  provider: 'openai',
  base_url: 'https://api.openai.com/v1',
  model_name: '',
  api_key: '',
  temperature: 0.8,
  max_tokens: 1024,
  price_coins: 0,
  share_role_pct: 0.5,
  share_preset_pct: 0.5,
  price_hint: '',
  status: 'active',
  is_default: false,
})

const showForm = ref(false)
const editingId = ref<string | null>(null)
const saving = ref(false)
const form = reactive(defaultForm())

const resetForm = () => {
  Object.assign(form, defaultForm())
}

const startCreate = () => {
  editingId.value = null
  resetForm()
  showForm.value = true
}

const startEdit = (model: ModelConfig) => {
  editingId.value = model.id
  Object.assign(form, {
    name: model.name,
    description: model.description || '',
    provider: model.provider,
    base_url: model.base_url,
    model_name: model.model_name,
    api_key: '',
    temperature: model.temperature ?? 0.8,
    max_tokens: model.max_tokens ?? 0,
    price_coins: model.price_coins ?? 0,
    share_role_pct: model.share_role_pct ?? 0,
    share_preset_pct: model.share_preset_pct ?? 0,
    price_hint: model.price_hint || '',
    status: model.status,
    is_default: model.is_default ?? false,
  })
  showForm.value = true
}

const cancelForm = () => {
  showForm.value = false
  editingId.value = null
  resetForm()
}

const submit = async () => {
  saving.value = true
  try {
    await saveModel({ ...form, id: editingId.value || undefined })
    cancelForm()
  } catch (error) {
    console.error('save model failed', error)
  } finally {
    saving.value = false
  }
}

const remove = async (id: string) => {
  if (!confirm('确定删除该模型吗？')) return
  await deleteModel(id)
}

const load = async () => {
  loading.value = true
  error.value = ''
  try {
    await loadModels()
  } catch (err: any) {
    error.value = err?.data?.error || err?.message || '加载失败'
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>
