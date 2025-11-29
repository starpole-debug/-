<template>
  <section class="space-y-8">
    <header class="relative overflow-hidden rounded-3xl border border-white/10 bg-gradient-to-r from-indigo-900/60 via-slate-900/70 to-purple-900/50 p-8 shadow-xl">
      <div class="pointer-events-none absolute inset-0 bg-[radial-gradient(circle_at_20%_20%,rgba(99,102,241,0.25),transparent_40%),radial-gradient(circle_at_80%_0%,rgba(236,72,153,0.2),transparent_35%)]" />
      <div class="relative flex flex-col gap-3 md:flex-row md:items-center md:justify-between text-white">
        <div>
          <p class="text-sm text-indigo-200/80">创作者钱包</p>
          <h1 class="text-2xl font-bold">收益总览</h1>
          <p class="text-xs text-slate-200/70">1 元 = 1000 平台币，金额以平台币计价，提现时按汇率折算。</p>
          <p v-if="errorMessage" class="mt-2 text-sm text-rose-300">
            {{ errorMessage }}
          </p>
        </div>
        <div class="text-sm text-slate-200/80">
          收益来源：角色/预设分成、其它收入。
        </div>
      </div>
      <div class="relative mt-6 grid gap-4 sm:grid-cols-3">
        <div class="stat-card">
          <p class="label">可用余额</p>
          <p class="value">{{ coinDisplay(wallet?.available_balance) }}</p>
          <p class="hint">≈ ¥{{ rmbDisplay(wallet?.available_balance) }}</p>
        </div>
        <div class="stat-card">
          <p class="label">冻结金额</p>
          <p class="value">{{ coinDisplay(wallet?.frozen_balance) }}</p>
          <p class="hint">提现审核中</p>
        </div>
        <div class="stat-card">
          <p class="label">累计收入</p>
          <p class="value">{{ coinDisplay(wallet?.total_earned) }}</p>
          <p class="hint">≈ ¥{{ rmbDisplay(wallet?.total_earned) }}</p>
        </div>
      </div>
    </header>

    <section class="rounded-3xl border border-white/10 bg-white/5 p-6 shadow-lg text-white">
      <div class="flex flex-col gap-2 md:flex-row md:items-center md:justify-between">
        <div>
          <h2 class="text-lg font-semibold">提现申请</h2>
          <p class="text-xs text-slate-300">输入提现平台币数量，按 1000 币 ≈ ¥1 折算。</p>
        </div>
        <p class="text-xs text-slate-300">当前可用：{{ coinDisplay(wallet?.available_balance) }}（≈ ¥{{ rmbDisplay(wallet?.available_balance) }}）</p>
      </div>
      <form class="mt-4 flex flex-col gap-3 md:flex-row" @submit.prevent="submit">
        <input v-model.number="amount" type="number" min="0" placeholder="提现币数，例如 10000" class="input-glass flex-1" />
        <input v-model="channel" placeholder="收款渠道（如支付宝）" class="input-glass flex-1" />
        <button
          class="btn-primary px-5 py-2 text-sm"
          :disabled="submitting || amount <= 0"
        >
          {{ submitting ? '提交中...' : '提交' }}
        </button>
      </form>
      <p class="mt-2 text-xs text-slate-300">预计到账：¥{{ rmbDisplay(amount) }}</p>

      <div class="mt-6">
        <h3 class="text-sm font-semibold text-white/80">近期收入</h3>
        <ul class="mt-3 divide-y divide-white/10 text-sm text-slate-200">
          <li v-if="!events.length" class="py-3 text-slate-400">暂无收入记录</li>
          <li v-for="event in events" :key="event.id" class="py-3 flex items-center justify-between">
            <div>
              <p class="font-medium text-white">{{ event.event_type || '收益' }}</p>
              <p class="text-xs text-slate-400">{{ formatDate(event.created_at) }}</p>
            </div>
            <div class="text-right">
              <p class="text-emerald-300 font-semibold">+ {{ coinDisplay(event.amount) }}</p>
              <p class="text-[11px] text-slate-400">≈ ¥{{ rmbDisplay(event.amount) }}</p>
            </div>
          </li>
        </ul>
      </div>
    </section>
  </section>
</template>

<script setup lang="ts">
import type { CreatorWallet, RevenueEvent } from '@/types'

definePageMeta({
  middleware: ['auth'],
})

const RATE = 1000 // 1000 币 = 1 元
const api = useApi()
const wallet = ref<CreatorWallet | null>(null)
const events = ref<RevenueEvent[]>([])
const amount = ref(10000)
const channel = ref('支付宝')
const errorMessage = ref('')
const submitting = ref(false)

const coinDisplay = (v?: number | null) => Number(v || 0).toLocaleString()
const rmbDisplay = (v?: number | null) => (Number(v || 0) / RATE).toFixed(2)
const formatDate = (ts?: string) => (ts ? new Date(ts).toLocaleString() : '')

const load = async () => {
  errorMessage.value = ''
  try {
    const data = await api.get<{ data: { wallet: CreatorWallet; events: RevenueEvent[]; payouts: any[] } }>('/me/wallet')
    wallet.value = data.data.wallet
    events.value = data.data.events || []
  } catch (error: any) {
    errorMessage.value = error?.data?.error || error?.message || '无法加载钱包信息，请稍后再试'
  }
}

const submit = async () => {
  submitting.value = true
  errorMessage.value = ''
  try {
    await api.post('/me/payouts', { amount: amount.value, channel: channel.value })
    await load()
  } catch (error: any) {
    errorMessage.value = error?.data?.error || error?.message || '提交失败，请稍后再试'
  } finally {
    submitting.value = false
  }
}

onMounted(load)
</script>

<style scoped>
.stat-card {
  @apply rounded-2xl bg-white/10 border border-white/10 p-4 backdrop-blur text-white;
}
.stat-card .label {
  @apply text-xs text-indigo-100/80;
}
.stat-card .value {
  @apply text-3xl font-semibold mt-2;
}
.stat-card .hint {
  @apply text-xs text-slate-200/80 mt-1;
}
.input-glass {
  @apply rounded-xl border border-white/15 bg-white/10 px-4 py-2 text-sm text-white placeholder:text-slate-400 focus:outline-none focus:border-indigo-400;
}
.btn-primary {
  @apply inline-flex items-center justify-center rounded-xl bg-primary text-white font-semibold shadow-lg shadow-indigo-500/20 hover:opacity-90 transition disabled:opacity-60 disabled:cursor-not-allowed;
}
</style>
