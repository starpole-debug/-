<template>
  <section class="space-y-6">
    <div>
      <h1 class="text-2xl font-semibold">站内信</h1>
      <p class="text-sm text-slate-400 mt-1">支持向全体或指定用户发送通知，可按昵称/邮箱/ID 搜索选择。</p>
    </div>

    <div class="glass-card p-6 space-y-4">
      <div class="grid gap-4 md:grid-cols-2">
        <div class="space-y-2">
          <label class="text-sm text-slate-300">标题</label>
          <input v-model="title" class="w-full bg-white/5 border border-white/10 rounded-lg px-3 py-2 text-sm" placeholder="例如：系统公告" />
        </div>
        <div class="space-y-2">
          <label class="text-sm text-slate-300">搜索用户</label>
          <div class="flex gap-2">
            <input v-model="query" class="flex-1 bg-white/5 border border-white/10 rounded-lg px-3 py-2 text-sm" placeholder="用户名/邮箱/ID" />
            <button class="btn-secondary" @click="searchUsers">搜索</button>
          </div>
          <p class="text-[12px] text-slate-500">结果点击即可添加到接收列表，或勾选“发送给所有人”广播。</p>
        </div>
      </div>

      <div class="space-y-2">
        <label class="text-sm text-slate-300">内容</label>
        <textarea v-model="content" rows="4" class="w-full bg-white/5 border border-white/10 rounded-lg px-3 py-2 text-sm" placeholder="请输入站内信内容" />
      </div>

      <div class="flex items-center gap-3">
        <label class="flex items-center gap-2 text-sm text-slate-300">
          <input v-model="sendAll" type="checkbox" class="accent-indigo-500" />
          发送给所有用户（使用当前搜索结果或全部）
        </label>
      </div>

      <div class="glass-card bg-white/5 border border-white/10 p-4 space-y-3">
        <div class="flex items-center justify-between">
          <h3 class="text-sm font-semibold text-white">已选择用户</h3>
          <button class="text-xs text-rose-300 hover:text-rose-200" @click="selectedIds = []">清空</button>
        </div>
        <div v-if="!selectedIds.length" class="text-sm text-slate-500">暂未选择用户。</div>
        <div v-else class="flex flex-wrap gap-2">
          <span
            v-for="id in selectedIds"
            :key="id"
            class="px-2 py-1 rounded-lg bg-white/5 border border-white/10 text-xs flex items-center gap-2"
          >
            {{ id }}
            <button class="text-rose-300" @click="remove(id)">×</button>
          </span>
        </div>
      </div>

      <div v-if="results.length" class="glass-card bg-white/5 border border-white/10 p-4 space-y-2">
        <div class="text-sm text-slate-300">搜索结果（点击添加）</div>
        <div class="flex flex-wrap gap-2">
          <button
            v-for="u in results"
            :key="u.id"
            class="px-3 py-1 rounded-lg bg-indigo-500/10 text-sm text-indigo-200 border border-indigo-500/20"
            @click="add(u.id)"
          >
            {{ u.username || u.email || u.id }}
          </button>
        </div>
      </div>

      <div class="flex justify-end">
        <button class="btn-primary" :disabled="sending" @click="send">发送</button>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAdminApi } from '~/composables/useAdminApi'

definePageMeta({
  layout: 'admin',
  middleware: 'admin',
})

const api = useAdminApi()
const title = ref('')
const content = ref('')
const query = ref('')
const sendAll = ref(false)
const results = ref<any[]>([])
const selectedIds = ref<string[]>([])
const sending = ref(false)

const searchUsers = async () => {
  const res = await api.get<{ data: any[] }>(`/admin/users?query=${encodeURIComponent(query.value || '')}&limit=20`)
  results.value = res.data || []
}

const add = (id: string) => {
  if (!selectedIds.value.includes(id)) {
    selectedIds.value.push(id)
  }
}

const remove = (id: string) => {
  selectedIds.value = selectedIds.value.filter(i => i !== id)
}

const send = async () => {
  if (!title.value.trim() || !content.value.trim()) {
    alert('请输入标题和内容')
    return
  }
  sending.value = true
  try {
    await api.post('/admin/notifications/broadcast', {
      title: title.value.trim(),
      content: content.value.trim(),
      user_ids: selectedIds.value,
      broadcast: sendAll.value,
      query: sendAll.value ? query.value : '',
      limit: 500,
    })
    alert('已发送')
    selectedIds.value = []
  } finally {
    sending.value = false
  }
}
</script>
