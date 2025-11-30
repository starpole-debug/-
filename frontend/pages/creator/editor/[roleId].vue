<template>
  <section class="flex flex-col lg:flex-row gap-6 text-charcoal-900 items-start">
    <!-- Left Sidebar: Basic Info -->
    <aside class="w-full lg:w-80 flex-shrink-0 flex flex-col gap-4 card-soft p-6 rounded-3xl shadow-lg lg:sticky lg:top-24 max-h-[calc(100vh-8rem)] overflow-y-auto custom-scrollbar">
      <header>
        <h1 class="text-xl font-semibold">{{ isNew ? '创建角色' : '编辑角色' }}</h1>
        <p class="text-xs text-charcoal-500 mt-1">配置角色基础信息</p>
      </header>

      <div class="space-y-4">
        <div>
          <label class="text-xs text-charcoal-500">角色名</label>
          <input v-model.trim="form.name" class="mt-1 glass-input text-sm" required placeholder="如：星际航行家" />
        </div>
        <div>
          <label class="text-xs text-charcoal-500">头像 URL</label>
          <input v-model.trim="form.avatar_url" class="mt-1 glass-input text-sm" placeholder="https://example.com/avatar.png" />
          <div class="mt-2 flex items-center gap-3 flex-wrap">
            <input ref="avatarFileRef" type="file" accept="image/*" class="text-xs text-charcoal-500" @change="onAvatarFileChange" />
            <button
              type="button"
              class="px-3 py-1.5 rounded-lg border border-charcoal-200 bg-white text-xs text-charcoal-600 hover:border-accent-yellow hover:text-charcoal-900 transition-colors disabled:opacity-50"
              :disabled="avatarUploading"
              @click="uploadAvatarFile()"
            >
              {{ avatarUploading ? '上传中...' : '上传本地图片' }}
            </button>
          </div>
          <p v-if="avatarUploadError" class="text-[11px] text-rose-400 mt-1">{{ avatarUploadError }}</p>
          <div v-if="form.avatar_url" class="mt-2 flex items-center gap-2 text-xs text-charcoal-500">
            <img :src="avatarPreview" class="h-12 w-12 rounded-lg object-cover border border-charcoal-200" alt="角色头像预览" />
            <span class="truncate">{{ form.avatar_url }}</span>
          </div>
        </div>
        <div>
          <label class="text-xs text-charcoal-500">短描述</label>
          <textarea v-model="form.description" class="mt-1 h-24 glass-input text-sm resize-none" placeholder="一句话介绍角色定位" />
        </div>
        <div>
          <label class="text-xs text-charcoal-500">标签</label>
          <input v-model="form.tags" class="mt-1 glass-input text-sm" placeholder="塔罗, 未来, 温柔" />
        </div>
        <div>
          <label class="text-xs text-charcoal-500">能力/特质</label>
          <input v-model="form.abilities" class="mt-1 glass-input text-sm" placeholder="共情, 讲故事" />
        </div>
      </div>

      <div class="mt-auto space-y-3 pt-4 border-t border-charcoal-100">
        <p v-if="errorMessage" class="text-xs text-rose-500">{{ errorMessage }}</p>
        <p v-if="successMessage" class="text-xs text-emerald-500">{{ successMessage }}</p>
        <button
          @click="submit"
          class="w-full rounded-xl btn-primary py-2.5 text-sm font-medium disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          :disabled="submitting"
        >
          {{ submitting ? '保存中...' : '保存角色' }}
        </button>
        <NuxtLink to="/creator/roles" class="block text-center text-xs text-charcoal-400 hover:text-charcoal-900 transition-colors">返回列表</NuxtLink>
      </div>
    </aside>

    <!-- Main Content: Tabs -->
    <main class="flex-1 flex flex-col card-soft rounded-3xl shadow-lg min-h-[calc(100vh-8rem)]">
      <!-- Tabs Header -->
      <div class="flex items-center gap-1 p-2 border-b border-charcoal-100 overflow-x-auto sticky top-24 bg-white/95 backdrop-blur z-10 rounded-t-3xl">
        <button
          v-for="tab in tabs"
          :key="tab.id"
          @click="activeTab = tab.id"
          class="px-4 py-2 rounded-lg text-sm font-medium transition-colors whitespace-nowrap"
          :class="activeTab === tab.id ? 'bg-bg-cream-200 text-charcoal-900 font-bold' : 'text-charcoal-500 hover:text-charcoal-700 hover:bg-bg-cream-100'"
        >
          {{ tab.label }}
        </button>
      </div>

      <!-- Tabs Content -->
      <div class="p-6">
        <p v-if="loading" class="text-sm text-charcoal-400">加载中...</p>
        <div v-else>
            
          <!-- Tab: Character -->
          <div v-show="activeTab === 'character'" class="space-y-6 max-w-3xl mx-auto">
            <div class="space-y-4">
                <h3 class="text-lg font-medium text-charcoal-900 border-l-4 border-accent-yellow pl-3">角色设定</h3>
                <div>
                    <label class="text-sm text-charcoal-500 block mb-2">Persona / 详细设定</label>
                    <textarea v-model="builder.persona" class="glass-input h-64 text-sm leading-relaxed" placeholder="在此输入角色的详细设定、性格描述、背景故事等..." />
                </div>
                <div>
                    <label class="text-sm text-charcoal-500 block mb-2">First Message / 开场白</label>
                    <textarea v-model="builder.first_message" class="glass-input h-32 text-sm leading-relaxed" placeholder="角色对用户说的第一句话..." />
                </div>
                <div class="grid gap-6 md:grid-cols-2">
                    <div>
                        <label class="text-sm text-charcoal-500 block mb-2">性格特质</label>
                        <input v-model="builder.traits" class="glass-input text-sm" placeholder="傲娇, 腹黑, 治愈" />
                    </div>
                    <div>
                        <label class="text-sm text-charcoal-500 block mb-2">当前场景</label>
                        <input v-model="builder.scenario" class="glass-input text-sm" placeholder="如：暴雨中的咖啡馆" />
                    </div>
                </div>
            </div>
          </div>

          <!-- Tab: World -->
          <div v-show="activeTab === 'world'" class="space-y-6 max-w-3xl mx-auto">
             <div class="flex items-center justify-between">
                <h3 class="text-lg font-medium text-charcoal-900 border-l-4 border-accent-yellow pl-3">世界书 (Worldbook)</h3>
                <button class="text-xs px-3 py-1.5 rounded-lg bg-accent-pink/10 text-accent-pink hover:bg-accent-pink/20 transition-colors" @click="addWorldEntry">+ 添加条目</button>
             </div>
             
             <div class="grid gap-6 md:grid-cols-2">
                <div>
                    <label class="text-sm text-charcoal-500 block mb-2">世界观摘要</label>
                    <textarea v-model="builder.world.summary" class="glass-input h-24 text-sm resize-none" />
                </div>
                <div class="space-y-4">
                    <div>
                        <label class="text-sm text-charcoal-500 block mb-2">场景 / 地点</label>
                        <input v-model="builder.world.scene" class="glass-input text-sm" />
                    </div>
                    <div>
                        <label class="text-sm text-charcoal-500 block mb-2">时间线</label>
                        <input v-model="builder.world.timeline" class="glass-input text-sm" />
                    </div>
                </div>
             </div>
             
             <div>
                <label class="text-sm text-charcoal-500 block mb-2">关键 NPC</label>
                <input v-model="builder.world.npcs" class="glass-input text-sm" placeholder="逗号分隔，如：店长, 神秘人" />
             </div>

             <div class="space-y-3">
                <label class="text-sm text-charcoal-500 block">自定义条目</label>
                <div v-if="builder.world.entries.length === 0" class="text-sm text-charcoal-400 italic py-4 text-center bg-bg-cream-100 rounded-xl">暂无条目</div>
                <div v-for="(entry, idx) in builder.world.entries" :key="idx" class="flex gap-3 items-start bg-bg-cream-100 p-3 rounded-xl border border-charcoal-100">
                  <input v-model="entry.key" class="w-32 rounded-lg border border-charcoal-200 bg-white px-3 py-2 text-sm text-charcoal-900 focus:border-accent-yellow focus:outline-none" placeholder="关键词" />
                  <textarea v-model="entry.value" class="flex-1 h-10 min-h-[2.5rem] rounded-lg border border-charcoal-200 bg-white px-3 py-2 text-sm text-charcoal-900 focus:border-accent-yellow focus:outline-none resize-y" placeholder="详细设定内容" />
                  <button class="text-charcoal-400 hover:text-status-error p-2 transition-colors" @click="builder.world.entries.splice(idx, 1)">
                    <span class="text-lg">×</span>
                  </button>
                </div>
             </div>
          </div>

          <!-- Tab: Examples -->
          <div v-show="activeTab === 'examples'" class="space-y-6 max-w-3xl mx-auto">
            <div class="flex items-center justify-between">
                <h3 class="text-lg font-medium text-charcoal-900 border-l-4 border-accent-yellow pl-3">对话示例</h3>
                <button class="text-xs px-3 py-1.5 rounded-lg bg-accent-pink/10 text-accent-pink hover:bg-accent-pink/20 transition-colors" @click="addExample">+ 添加示例</button>
            </div>
            
            <div v-if="builder.examples.length === 0" class="text-sm text-charcoal-400 italic py-8 text-center bg-bg-cream-100 rounded-xl">
                暂无示例对话，添加示例有助于模型更好地模仿语气。
            </div>
            
            <div class="space-y-4">
                <div v-for="(ex, idx) in builder.examples" :key="idx" class="bg-bg-cream-100 p-4 rounded-xl border border-charcoal-100 space-y-3 relative group">
                    <button class="absolute top-2 right-2 text-charcoal-400 hover:text-status-error opacity-0 group-hover:opacity-100 transition-opacity" @click="builder.examples.splice(idx, 1)">删除</button>
                    <div>
                        <label class="text-xs text-charcoal-500 mb-1 block">用户 (User)</label>
                        <input v-model="ex.user" class="glass-input text-sm" placeholder="用户说的话..." />
                    </div>
                    <div>
                        <label class="text-xs text-accent-pink mb-1 block">角色 (Assistant)</label>
                        <input v-model="ex.assistant" class="glass-input text-sm" placeholder="角色的回复..." />
                    </div>
                </div>
            </div>
          </div>

          <!-- Tab: Preset -->
          <div v-show="activeTab === 'preset'" class="space-y-6 max-w-3xl mx-auto">
            <div class="flex items-center justify-between">
                <h3 class="text-lg font-medium text-charcoal-900 border-l-4 border-accent-yellow pl-3">Prompt 编排</h3>
                <button class="text-xs px-3 py-1.5 rounded-lg bg-accent-pink/10 text-accent-pink hover:bg-accent-pink/20 transition-colors" @click="addBlock">+ 添加 Block</button>
            </div>
            
            <div class="space-y-3">
              <div
                v-for="(block, idx) in builder.preset.blocks"
                :key="idx"
                class="bg-bg-cream-100 p-4 rounded-xl border border-charcoal-100 hover:border-accent-yellow transition-colors"
              >
                <div class="flex items-center justify-between gap-4 mb-3">
                  <div class="flex flex-col gap-1 flex-1">
                    <input v-model="block.name" class="bg-transparent text-sm font-medium text-charcoal-900 border-b border-transparent focus:border-accent-yellow outline-none transition-colors w-full" placeholder="Block Name" />
                  </div>
                  <div class="flex items-center gap-2">
                    <div class="flex bg-white border border-charcoal-200 rounded-lg p-0.5">
                        <button class="text-xs px-2 py-1 text-charcoal-400 hover:text-charcoal-900 hover:bg-charcoal-50 rounded transition-colors" @click="moveBlock(idx, 'up')">↑</button>
                        <button class="text-xs px-2 py-1 text-charcoal-400 hover:text-charcoal-900 hover:bg-charcoal-50 rounded transition-colors" @click="moveBlock(idx, 'down')">↓</button>
                    </div>
                    <label class="flex items-center gap-2 text-xs text-charcoal-400 cursor-pointer select-none">
                        <div class="relative inline-flex items-center cursor-pointer">
                            <input type="checkbox" v-model="block.enabled" class="sr-only peer">
                            <div class="w-9 h-5 bg-charcoal-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-accent-yellow"></div>
                        </div>
                        <span class="w-6">{{ block.enabled ? '启用' : '禁用' }}</span>
                    </label>
                    <button class="text-charcoal-400 hover:text-status-error p-1.5 rounded-lg hover:bg-red-50 transition-colors ml-1" @click="removeBlock(idx)">
                        <span class="text-lg leading-none">×</span>
                    </button>
                  </div>
                </div>
                <textarea
                  v-model="block.content"
                  class="w-full rounded-lg border border-charcoal-200 bg-white px-3 py-2 text-sm text-charcoal-900 focus:border-accent-yellow focus:outline-none transition-colors"
                  rows="3"
                  :placeholder="block.marker ? '此 Block 为 Marker 占位符，运行时将被替换为动态内容' : '在此输入静态 Prompt 内容...'"
                />
              </div>
            </div>
          </div>

          <!-- Tab: Preview -->
          <div v-show="activeTab === 'preview'" class="space-y-6 max-w-3xl mx-auto">
             <h3 class="text-lg font-medium text-charcoal-900 border-l-4 border-accent-yellow pl-3">预览与导出</h3>
             
             <div class="grid gap-4 md:grid-cols-2">
                <div class="space-y-2">
                    <label class="text-sm text-charcoal-500">备份 Prompt (Role Versions)</label>
                    <textarea
                        v-model="form.prompt"
                        class="glass-input h-40 text-sm resize-none"
                        placeholder="如需将当前草稿快照到后端 role_versions，可在此粘贴/生成"
                    />
                    <p class="text-xs text-charcoal-400">保存后可通过 /roles/:id/prompt 访问</p>
                </div>
                <div class="space-y-2">
                    <label class="text-sm text-charcoal-500">导出选项</label>
                    <div class="flex flex-col gap-2">
                        <button class="w-full py-2 rounded-xl border border-charcoal-200 bg-white text-sm text-charcoal-700 hover:bg-bg-cream-100 transition-colors" @click="copy(nebulaJSON, '已复制 Nebula JSON')">
                            复制 Nebula JSON
                        </button>
                        <button class="w-full py-2 rounded-xl border border-charcoal-200 bg-white text-sm text-charcoal-700 hover:bg-bg-cream-100 transition-colors" @click="copy(stJSON, '已复制 ST 友好 JSON')">
                            复制 SillyTavern JSON
                        </button>
                        <button class="w-full py-2 rounded-xl border border-charcoal-200 bg-white text-sm text-charcoal-700 hover:bg-bg-cream-100 transition-colors" @click="copy(presetJSON, '已复制 Preset JSON')">
                            复制 Preset JSON
                        </button>
                    </div>
                    <p v-if="copyMessage" class="text-xs text-status-success text-center mt-2">{{ copyMessage }}</p>
                </div>
             </div>

             <div class="space-y-2">
                <label class="text-sm text-charcoal-500">Prompt 预览 (System)</label>
                <div class="w-full h-64 rounded-xl border border-charcoal-200 bg-white px-4 py-3 text-xs text-charcoal-600 overflow-y-auto whitespace-pre-wrap font-mono custom-scrollbar">
                    {{ promptPreview }}
                </div>
             </div>
          </div>
          
        </div>
      </div>
    </main>
  </section>
</template>

<script setup lang="ts">
import type { Role } from '@/types'

definePageMeta({
  middleware: ['auth'],
  layout: 'default'
})

const route = useRoute()
const router = useRouter()
const api = useApi()
const { fetchCreatorRole } = useCreatorStats()
const { uploadRoleAvatar, resolveAssetUrl } = useUpload()

const roleId = computed(() => route.params.roleId as string)
const isNew = computed(() => !roleId.value || roleId.value === 'new')
const loading = ref(false)
const submitting = ref(false)
const errorMessage = ref('')
const successMessage = ref('')
const copyMessage = ref('')

const activeTab = ref('character')
const tabs = [
    { id: 'character', label: '角色设定' },
    { id: 'world', label: '世界书' },
    { id: 'examples', label: '对话示例' },
    { id: 'preset', label: 'Prompt 编排' },
    { id: 'preview', label: '预览/导出' }
]

const form = reactive({
  name: '',
  description: '',
  avatar_url: '',
  tags: '',
  abilities: '',
  prompt: '',
})
const avatarFileRef = ref<HTMLInputElement | null>(null)
const avatarUploading = ref(false)
const avatarUploadError = ref('')
const avatarPreview = computed(() => resolveAssetUrl(form.avatar_url))

const builder = reactive({
  persona: '',
  first_message: '',
  traits: '',
  scenario: '',
  world: {
    summary: '',
    scene: '',
    timeline: '',
    npcs: '',
    entries: [] as { key: string; value: string }[],
  },
  examples: [] as { user: string; assistant: string }[],
  model: {
    provider: 'openai-like',
    model: 'gpt-3.5-turbo',
  },
  safety: '',
  nsfw: false,
  preset: {
    name: '默认预设',
    model: {
      provider: 'openai-like',
      key: 'gpt-4o',
      base_url: '',
      api_key: '',
    },
    genParams: {
      temperature: 0.9,
      top_p: 1,
      frequency_penalty: 0,
      presence_penalty: 0,
      max_tokens: 1024,
    },
    blocks: [
      { id: 'safety', name: '安全', role: 'system', content: '', enabled: true, marker: false },
      { id: 'main', name: '主提示', role: 'system', content: "Write {{char}}'s next reply...", enabled: true, marker: false },
      { id: 'world_pre', name: '世界信息(前)', role: 'system', marker: true, enabled: true, content: '' },
      { id: 'char_desc', name: '角色描述', role: 'system', marker: true, enabled: true, content: '' },
      { id: 'personality', name: '人格特质', role: 'system', marker: true, enabled: true, content: '' },
      { id: 'scenario', name: '场景', role: 'system', marker: true, enabled: true, content: '' },
      { id: 'nsfw', name: 'NSFW 扩展', role: 'system', content: '', enabled: false, marker: false },
      { id: 'world_post', name: '世界信息(后)', role: 'system', marker: true, enabled: true, content: '' },
      { id: 'examples', name: '示例对话', role: 'system', marker: true, enabled: true, content: '' },
      { id: 'history', name: '聊天记录', role: 'system', marker: true, enabled: true, content: '' },
    ],
  },
})

const normalizeList = (input: string) =>
  input
    .split(/[,，]/)
    .map((item) => item.trim())
    .filter(Boolean)

const load = async () => {
  if (isNew.value) {
    return
  }
  loading.value = true
  errorMessage.value = ''
  try {
    const role = await fetchCreatorRole(roleId.value)
    if (!role) return
    form.name = role.name
    form.description = role.description
    form.avatar_url = role.avatar_url || ''
    form.tags = (role.tags || []).join(', ')
    form.abilities = (role.abilities || []).join(', ')
    if (role.data) {
        Object.assign(builder, role.data)
    } else {
        builder.persona = role.description || ''
    }
  } catch (error: any) {
    errorMessage.value = error?.data?.error || error?.message || '加载角色失败'
  } finally {
    loading.value = false
  }
}

const uploadAvatarFile = async (file?: File | null) => {
  if (avatarUploading.value) return
  const targetFile = file || avatarFileRef.value?.files?.[0] || null
  if (!targetFile) {
    avatarUploadError.value = '请选择图片后再上传'
    return
  }
  avatarUploadError.value = ''
  avatarUploading.value = true
  try {
    const url = await uploadRoleAvatar(targetFile)
    if (url) {
      form.avatar_url = url
      if (avatarFileRef.value) avatarFileRef.value.value = ''
    } else {
      avatarUploadError.value = '上传失败，请稍后重试'
    }
  } catch (error: any) {
    avatarUploadError.value = error?.data?.error || error?.message || '上传失败，请稍后重试'
  } finally {
    avatarUploading.value = false
  }
}

const onAvatarFileChange = async (e: Event) => {
  const target = e.target as HTMLInputElement
  const file = target.files?.[0] || null
  if (file) {
    await uploadAvatarFile(file)
  }
}

const submit = async () => {
  submitting.value = true
  errorMessage.value = ''
  successMessage.value = ''
  try {
    const payload = {
      name: form.name,
      description: form.description,
      avatar_url: form.avatar_url,
      tags: normalizeList(form.tags),
      abilities: normalizeList(form.abilities),
      data: builder,
    }
    let savedRole: Role | null = null
    if (isNew.value) {
      const res = await api.post<{ data: Role }>('/roles', payload)
      savedRole = res.data
    } else {
      const res = await api.put<{ data: Role }>(`/roles/${roleId.value}`, payload)
      savedRole = res.data
    }
    if (!savedRole) {
      throw new Error('保存失败')
    }
    const targetId = savedRole.id
    if (form.prompt.trim()) {
      await api.post(`/roles/${targetId}/prompt`, { prompt: form.prompt })
    }
    successMessage.value = '保存成功'
    if (isNew.value) {
      await router.replace(`/creator/editor/${targetId}`)
    }
  } catch (error: any) {
    errorMessage.value = error?.data?.error || error?.message || '保存失败，请稍后重试'
  } finally {
    submitting.value = false
  }
}

const nebulaJSON = computed(() => {
  const entriesObj: Record<string, string[]> = {}
  builder.world.entries.forEach((e) => {
    if (!e.key.trim()) return
    entriesObj[e.key.trim()] = normalizeList(e.value)
  })
  const obj = {
    version: 'nebula-1',
    name: form.name || '未命名角色',
    persona: builder.persona,
    first_message: builder.first_message,
    traits: normalizeList(builder.traits),
    scenario: builder.scenario,
    world: {
      summary: builder.world.summary,
      scene: builder.world.scene,
      timeline: builder.world.timeline,
      npcs: normalizeList(builder.world.npcs),
      entries: entriesObj,
    },
    examples: builder.examples.filter((e) => e.user || e.assistant),
    model_hints: builder.model,
    safety: builder.safety,
    tags: normalizeList(form.tags),
    nsfw: builder.nsfw,
  }
  return JSON.stringify(obj, null, 2)
})

const stJSON = computed(() => {
  const obj = {
    name: form.name || '未命名角色',
    description: builder.persona || form.description,
    first_mes: builder.first_message,
    personality: builder.traits,
    scenario: builder.scenario,
    character_book: {
        entries: builder.world.entries.map(e => ({
            keys: normalizeList(e.key),
            content: e.value,
            enabled: true,
            insertion_order: 100,
        }))
    },
    mes_example: builder.examples
      .filter((e) => e.user || e.assistant)
      .map((e) => `${e.user}\n${e.assistant}`)
      .join('\n---\n'),
    creator_notes: builder.safety,
    tags: normalizeList(form.tags),
  }
  return JSON.stringify(obj, null, 2)
})

const promptPreview = computed(() => {
  const lines: string[] = []
  if (builder.safety) lines.push('[SAFETY]', builder.safety)
  lines.push('[ROLE]', builder.persona || form.description || '')
  const traits = normalizeList(builder.traits)
  if (traits.length) lines.push('Traits: ' + traits.join(', '))
  if (builder.scenario) lines.push('Scenario: ' + builder.scenario)
  if (builder.world.summary) lines.push('[WORLD] ' + builder.world.summary)
  if (builder.world.scene || builder.world.timeline) {
    lines.push(`Scene: ${builder.world.scene || ''} Timeline: ${builder.world.timeline || ''}`)
  }
  const npcs = normalizeList(builder.world.npcs)
  if (npcs.length) lines.push('NPCs: ' + npcs.join(', '))
  if (builder.examples.length) {
    lines.push('[EXAMPLES]')
    builder.examples.forEach((ex, i) => {
      lines.push(`${i + 1}. User: ${ex.user} | Assistant: ${ex.assistant}`)
    })
  }
  return lines.filter(Boolean).join('\n')
})

const presetJSON = computed(() => {
  const blocks = builder.preset.blocks.map((b) => ({
    id: b.id,
    name: b.name,
    role: b.role,
    content: b.content,
    enabled: b.enabled,
    marker: b.marker || false,
  }))
  const obj = {
    version: 'nebula-preset-1',
    meta: {
      name: builder.preset.name || form.name || '未命名预设',
    },
    model: builder.preset.model,
    gen_params: builder.preset.genParams,
    blocks,
    order: builder.preset.blocks.map(b => b.id),
    format: {
      wi: '{0}',
      scenario: '{{scenario}}',
      personality: '{{traits}}',
    },
  }
  return JSON.stringify(obj, null, 2)
})

const copy = async (text: string, toast?: string) => {
  try {
    if (navigator.clipboard?.writeText) {
      await navigator.clipboard.writeText(text)
    } else {
      const ta = document.createElement('textarea')
      ta.value = text
      document.body.appendChild(ta)
      ta.select()
      document.execCommand('copy')
      document.body.removeChild(ta)
    }
    if (toast) {
      copyMessage.value = toast
      setTimeout(() => (copyMessage.value = ''), 1500)
    }
  } catch (error) {
    console.error('copy failed', error)
    copyMessage.value = '复制失败，请手动全选复制'
  }
}

const addWorldEntry = () => builder.world.entries.push({ key: '', value: '' })
const addExample = () => builder.examples.push({ user: '', assistant: '' })

const moveBlock = (idx: number, direction: 'up' | 'down') => {
  const target = direction === 'up' ? idx - 1 : idx + 1
  if (target < 0 || target >= builder.preset.blocks.length) return
  const newBlocks = [...builder.preset.blocks]
  ;[newBlocks[idx], newBlocks[target]] = [newBlocks[target], newBlocks[idx]]
  builder.preset.blocks = newBlocks
}

const addBlock = () => {
  builder.preset.blocks.push({
    id: `block_${Date.now()}`,
    name: '新区块',
    role: 'system',
    content: '',
    enabled: true,
    marker: false
  })
}

const removeBlock = (idx: number) => {
  builder.preset.blocks.splice(idx, 1)
}

onMounted(load)
</script>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 3px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: #94a3b8;
}
</style>
