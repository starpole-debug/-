<template>
  <section class="space-y-8">
    <header class="flex flex-col gap-2">
      <h1 class="text-xl font-semibold">绘图管理</h1>
      <p class="text-sm text-slate-300">配置 NovelAI 等绘图 API，并管理用于生成绘图提示词的预设。</p>
    </header>

    <!-- Providers -->
    <div class="space-y-4 rounded-2xl border border-white/10 p-6">
      <div class="flex items-center justify-between">
        <div>
          <h2 class="text-lg font-semibold">绘图 API</h2>
          <p class="text-xs text-slate-400">支持配置多个 API，随机选择并发最少的可用项。</p>
        </div>
        <div class="flex items-center gap-3">
          <button class="rounded-full border border-white/20 px-4 py-2 text-sm text-rose-200" @click="clearProviders">清空</button>
          <button class="rounded-full bg-white/10 px-4 py-2 text-sm text-white" @click="startCreateProvider">
            {{ editingProviderId ? '取消编辑' : '新增 API' }}
          </button>
        </div>
      </div>

      <form v-if="providerFormVisible" class="grid gap-4 md:grid-cols-2" @submit.prevent="saveProviderForm">
        <label class="text-sm">
          名称
          <input v-model.trim="providerForm.name" class="mt-1 w-full rounded-xl border border-white/10 bg-white/5 px-4 py-2 text-sm" required />
        </label>
        <label class="text-sm">
          Base URL（如 https://api.novelai.net/v1）
          <input v-model.trim="providerForm.base_url" class="mt-1 w-full rounded-xl border border-white/10 bg-white/5 px-4 py-2 text-sm" required />
        </label>
        <label class="text-sm">
          API Key
          <input v-model.trim="providerForm.api_key" type="password" autocomplete="off" class="mt-1 w-full rounded-xl border border-white/10 bg-white/5 px-4 py-2 text-sm" />
        </label>
        <label class="text-sm">
          模型（如 nai-diffusion-3）
          <input v-model.trim="providerForm.selected_model" class="mt-1 w-full rounded-xl border border-white/10 bg-white/5 px-4 py-2 text-sm" placeholder="如 nai-diffusion-3" />
        </label>
        <label class="text-sm">
          最大并发
          <input v-model.number="providerForm.max_concurrency" type="number" min="1" class="mt-1 w-full rounded-xl border border-white/10 bg-white/5 px-4 py-2 text-sm" />
        </label>
        <label class="text-sm">
          权重
          <input v-model.number="providerForm.weight" type="number" min="1" class="mt-1 w-full rounded-xl border border-white/10 bg-white/5 px-4 py-2 text-sm" />
        </label>
        <label class="text-sm md:col-span-2">
          默认参数（JSON，宽高/步数/CFG/采样器等）
          <textarea
            v-model="providerForm.params_json"
            rows="3"
            class="mt-1 w-full rounded-xl border border-white/10 bg-white/5 px-4 py-2 text-sm font-mono"
            placeholder='{"width":512,"height":768,"steps":28,"sampler":"k_euler_ancestral","scale":11}'
          />
        </label>
        <label class="text-sm">
          状态
          <select v-model="providerForm.status" class="mt-1 w-full rounded-xl border border-white/10 bg-slate-900 text-white px-4 py-2 text-sm">
            <option value="active">Active</option>
            <option value="inactive">Inactive</option>
          </select>
        </label>
        <div class="md:col-span-2 flex gap-3">
          <button type="submit" class="rounded-full bg-primary px-5 py-2 text-sm text-white" :disabled="savingProvider">
            {{ savingProvider ? '保存中...' : (editingProviderId ? '保存修改' : '创建 API') }}
          </button>
          <button type="button" class="rounded-full border border-white/20 px-5 py-2 text-sm" @click="cancelProviderForm">取消</button>
        </div>
      </form>

      <div class="overflow-auto rounded-xl border border-white/10">
        <table class="min-w-full text-sm">
          <thead class="bg-white/5 text-xs uppercase text-slate-400">
            <tr>
              <th class="px-4 py-3">名称</th>
              <th>Base URL</th>
              <th>模型</th>
              <th>并发</th>
              <th>权重</th>
              <th>状态</th>
              <th>Key</th>
              <th></th>
            </tr>
          </thead>
          <tbody v-if="providers.length">
            <tr v-for="p in providers" :key="p.id" class="border-t border-white/5">
              <td class="px-4 py-3 font-medium">{{ p.name }}</td>
              <td class="text-xs text-slate-300">{{ p.base_url }}</td>
              <td class="text-xs text-slate-300">{{ p.selected_model || '-' }}</td>
              <td>{{ p.max_concurrency || 0 }}</td>
              <td>{{ p.weight || 1 }}</td>
              <td>
                <span :class="pill(p.status === 'active')" class="rounded-full px-3 py-1 text-xs">
                  {{ p.status }}
                </span>
              </td>
              <td>{{ p.api_key ? '已配置' : '未配置' }}</td>
              <td class="space-x-3 text-xs">
                <button class="text-primary" @click="editProvider(p)">编辑</button>
                <button class="text-rose-300" @click="removeProvider(p.id!)">删除</button>
              </td>
            </tr>
          </tbody>
          <tbody v-else>
            <tr><td class="p-4 text-slate-400" colspan="6">暂无配置</td></tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Presets -->
    <div class="space-y-4 rounded-2xl border border-white/10 p-6">
      <div class="flex items-center justify-between">
        <div>
          <h2 class="text-lg font-semibold">绘图 Prompt 预设</h2>
          <p class="text-xs text-slate-400">为生成绘画用的 LLM 提供规范；Active 的预设将被使用。</p>
        </div>
        <button class="rounded-full bg-white/10 px-4 py-2 text-sm text-white" @click="startCreatePreset">
          {{ editingPresetId ? '取消编辑' : '新增预设' }}
        </button>
      </div>

      <form v-if="presetFormVisible" class="space-y-3" @submit.prevent="savePresetForm">
        <div class="grid gap-4 md:grid-cols-2">
          <label class="text-sm">
            名称
            <input v-model.trim="presetForm.name" class="mt-1 w-full rounded-xl border border-white/10 bg-white/5 px-4 py-2 text-sm" required />
          </label>
          <label class="text-sm">
            Prompt 模型（用于生成绘画提示词）
            <select v-model="presetForm.prompt_model_key" class="mt-1 w-full rounded-xl border border-white/10 bg-slate-900 text-white px-4 py-2 text-sm">
              <option value="">默认（mock）</option>
              <option v-for="m in modelOptions" :key="m.id" :value="m.id">
                {{ m.name }} ({{ m.model_name }})
              </option>
            </select>
          </label>
          <label class="text-sm">
            状态
            <select v-model="presetForm.status" class="mt-1 w-full rounded-xl border border-white/10 bg-slate-900 text-white px-4 py-2 text-sm">
              <option value="active">Active</option>
              <option value="inactive">Inactive</option>
            </select>
          </label>
        </div>
        <label class="text-sm block">
          预设 JSON（instruction/style/negative 等）
          <textarea
            v-model="presetForm.preset_json"
            class="mt-1 w-full rounded-xl border border-white/10 bg-white/5 px-4 py-2 text-sm font-mono"
            rows="8"
            placeholder='{"instruction":"...","style":"...","negative":"..."}'
            required
          />
        </label>
        <div class="flex gap-3">
          <button type="submit" class="rounded-full bg-primary px-5 py-2 text-sm text-white" :disabled="savingPreset">
            {{ savingPreset ? '保存中...' : (editingPresetId ? '保存修改' : '创建预设') }}
          </button>
          <button type="button" class="rounded-full border border-white/20 px-5 py-2 text-sm" @click="cancelPresetForm">取消</button>
        </div>
      </form>

      <div class="overflow-auto rounded-xl border border-white/10">
        <table class="min-w-full text-sm">
          <thead class="bg-white/5 text-xs uppercase text-slate-400">
            <tr>
              <th class="px-4 py-3">名称</th>
              <th>状态</th>
              <th>预设</th>
              <th></th>
            </tr>
          </thead>
          <tbody v-if="presets.length">
            <tr v-for="p in presets" :key="p.id" class="border-t border-white/5 align-top">
              <td class="px-4 py-3 font-medium">{{ p.name }}</td>
              <td>
                <span :class="pill(p.status === 'active')" class="rounded-full px-3 py-1 text-xs">
                  {{ p.status }}
                </span>
              </td>
              <td class="text-xs text-slate-300 whitespace-pre-wrap max-w-xl">
                {{ p.preset_json }}
              </td>
              <td class="space-x-3 text-xs">
                <button class="text-primary" @click="editPreset(p)">编辑</button>
                <button class="text-rose-300" @click="removePreset(p.id!)">删除</button>
              </td>
            </tr>
          </tbody>
          <tbody v-else>
            <tr><td class="p-4 text-slate-400" colspan="4">暂无预设</td></tr>
          </tbody>
        </table>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { useAdminImage } from '@/composables/useAdminImage'

definePageMeta({
  layout: 'admin',
  middleware: 'admin',
})

const { listProviders, saveProvider, deleteProvider, listPresets, savePreset, deletePreset } = useAdminImage()

const providers = ref<any[]>([])
const presets = ref<any[]>([])
const modelOptions = ref<any[]>([])

const providerFormVisible = ref(false)
const presetFormVisible = ref(false)
const editingProviderId = ref<string | null>(null)
const editingPresetId = ref<string | null>(null)
const savingProvider = ref(false)
const savingPreset = ref(false)

const providerForm = reactive({
  id: '',
  name: '',
  base_url: '',
  api_key: '',
  max_concurrency: 5,
  weight: 1,
  status: 'active',
  params_json: '',
  selected_model: '',
})
const presetForm = reactive({
  id: '',
  name: '',
  preset_json: '',
  prompt_model_key: '',
  status: 'active',
})

const pill = (ok: boolean) => (ok ? 'bg-emerald-500/20 text-emerald-200' : 'bg-rose-500/20 text-rose-200')

const loadAll = async () => {
  const ps = await listProviders()
  providers.value = Array.isArray(ps) ? ps.filter(p => p && (p.id || p.name)) : []
  const pr = await listPresets()
  presets.value = Array.isArray(pr) ? pr.filter(p => p && (p.id || p.name)) : []
  await loadModels()
}

const loadModels = async () => {
  // 复用 admin 模型接口，给绘图 prompt 选择用
  const adminConfig = useAdminConfig()
  if (adminConfig.models.value.length === 0) {
    await adminConfig.loadModels()
  }
  modelOptions.value = adminConfig.models.value
}

onMounted(loadAll)

const startCreateProvider = () => {
  editingProviderId.value = null
  Object.assign(providerForm, { id: '', name: '', base_url: '', api_key: '', max_concurrency: 5, weight: 1, status: 'active', params_json: '' })
  providerFormVisible.value = true
}
const editProvider = (p: any) => {
  editingProviderId.value = p.id
  const paramsJson = typeof p.params_json === 'string' ? p.params_json : JSON.stringify(p.params_json || {}, null, 2)
  Object.assign(providerForm, { ...p, api_key: '', params_json: paramsJson })
  providerFormVisible.value = true
}
const cancelProviderForm = () => {
  providerFormVisible.value = false
  editingProviderId.value = null
}
const saveProviderForm = async () => {
  savingProvider.value = true
  try {
    const payload = { ...providerForm }
    if (!editingProviderId.value) delete payload.id
    if (payload.params_json) {
      try {
        const parsed = JSON.parse(payload.params_json as any)
        payload.params_json = JSON.stringify(parsed)
      } catch (e) {
        alert('默认参数需要合法的 JSON')
        savingProvider.value = false
        return
      }
    }
    await saveProvider(payload as any)
    await loadAll()
    providerFormVisible.value = false
    editingProviderId.value = null
  } catch (e: any) {
    alert(e?.message || '保存失败')
  } finally {
    savingProvider.value = false
  }
}
const removeProvider = async (id: string) => {
  if (!confirm('确认删除该 API 配置？')) return
  await deleteProvider(id)
  await loadAll()
}
const clearProviders = async () => {
  if (!providers.value.length) return
  if (!confirm('确认清空所有绘图 API 配置？')) return
  const ids = providers.value.map(p => p.id).filter(Boolean)
  for (const id of ids) {
    await deleteProvider(id as string)
  }
  providers.value = []
}

const startCreatePreset = () => {
  editingPresetId.value = null
  Object.assign(presetForm, { id: '', name: '', preset_json: '', status: 'active' })
  presetFormVisible.value = true
}
const editPreset = (p: any) => {
  editingPresetId.value = p.id
  Object.assign(presetForm, { ...p })
  presetFormVisible.value = true
}
const cancelPresetForm = () => {
  presetFormVisible.value = false
  editingPresetId.value = null
}
const savePresetForm = async () => {
  savingPreset.value = true
  try {
    const payload = { ...presetForm }
    if (!editingPresetId.value) delete payload.id
    await savePreset(payload as any)
    await loadAll()
    presetFormVisible.value = false
    editingPresetId.value = null
  } catch (e: any) {
    alert(e?.message || '保存失败')
  } finally {
    savingPreset.value = false
  }
}
const removePreset = async (id: string) => {
  if (!confirm('确认删除该预设？')) return
  await deletePreset(id)
  await loadAll()
}
</script>
