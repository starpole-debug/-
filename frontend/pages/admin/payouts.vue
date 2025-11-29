<template>
  <section class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-semibold">提现申请</h1>
        <p class="text-sm text-slate-400 mt-1">创作者提交的提现请求，支持筛选和审批。</p>
      </div>
      <div class="flex items-center gap-3">
        <select v-model="status" class="bg-white/5 border border-white/10 rounded-lg px-3 py-2 text-sm">
          <option value="">全部状态</option>
          <option value="requested">待审核</option>
          <option value="approved">已通过</option>
          <option value="rejected">已驳回</option>
        </select>
        <button class="btn-secondary" @click="load">刷新</button>
      </div>
    </div>

    <div class="glass-card p-4">
      <div v-if="isLoading" class="text-sm text-slate-400 py-6 text-center">加载中...</div>
      <div v-else-if="!payouts.length" class="text-sm text-slate-500 py-6 text-center">暂无记录</div>
      <div v-else class="overflow-x-auto">
        <table class="min-w-full text-sm">
          <thead class="text-slate-400">
            <tr class="border-b border-white/5">
              <th class="py-2 text-left">ID</th>
              <th class="py-2 text-left">创作者</th>
              <th class="py-2 text-left">金额</th>
              <th class="py-2 text-left">渠道</th>
              <th class="py-2 text-left">状态</th>
              <th class="py-2 text-left">创建时间</th>
              <th class="py-2 text-right">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="p in payouts" :key="p.id" class="border-b border-white/5">
              <td class="py-2">{{ p.id.slice(0, 8) }}</td>
              <td class="py-2">{{ p.creator_id }}</td>
              <td class="py-2">+{{ p.amount }}</td>
              <td class="py-2">{{ p.channel || '-' }}</td>
              <td class="py-2">
                <span
                  class="px-2 py-1 rounded text-xs"
                  :class="p.status === 'approved' ? 'bg-green-500/20 text-green-300' : p.status === 'rejected' ? 'bg-rose-500/20 text-rose-200' : 'bg-amber-500/20 text-amber-200'"
                >
                  {{ p.status }}
                </span>
              </td>
              <td class="py-2">{{ new Date(p.created_at).toLocaleString() }}</td>
              <td class="py-2 text-right space-x-2">
                <button
                  class="px-3 py-1 rounded-lg bg-green-500/20 text-green-200 border border-green-500/30 disabled:opacity-40"
                  :disabled="p.status !== 'requested' || actionLoading"
                  @click="updateStatus(p.id, 'approve')"
                >
                  通过
                </button>
                <button
                  class="px-3 py-1 rounded-lg bg-rose-500/20 text-rose-200 border border-rose-500/30 disabled:opacity-40"
                  :disabled="p.status !== 'requested' || actionLoading"
                  @click="updateStatus(p.id, 'reject')"
                >
                  拒绝
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAdminApi } from '~/composables/useAdminApi'

definePageMeta({
  layout: 'admin',
  middleware: 'admin',
})

const api = useAdminApi()
const payouts = ref<any[]>([])
const status = ref('')
const isLoading = ref(false)
const actionLoading = ref(false)

const load = async () => {
  isLoading.value = true
  try {
    const res = await api.get<{ data: any[] }>(`/admin/payouts?status=${status.value || ''}`)
    payouts.value = res.data || []
  } finally {
    isLoading.value = false
  }
}

const updateStatus = async (id: string, action: 'approve' | 'reject') => {
  actionLoading.value = true
  try {
    await api.post(`/admin/payouts/${id}/${action}`)
    await load()
  } finally {
    actionLoading.value = false
  }
}

onMounted(load)
</script>
