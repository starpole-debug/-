<template>
  <section class="space-y-6">
    <header class="hero-card relative overflow-hidden rounded-3xl border border-white/10 p-6 md:p-8">
      <div class="pointer-events-none absolute inset-y-0 right-0 w-1/2 opacity-60 blur-3xl" :style="heroStyle"></div>
      <div class="relative flex items-center justify-between gap-4 flex-wrap">
        <div>
          <h1 class="text-2xl font-bold text-white">充值记录</h1>
          <p class="text-sm text-slate-200/70">展示所有已支付的充值订单（含支付宝/微信）。</p>
        </div>
        <NuxtLink class="btn-primary" to="/store">前往充值</NuxtLink>
      </div>
    </header>

    <section class="glass-card">
      <div v-if="loading && !payments.length" class="text-sm text-slate-200/70">加载中...</div>
      <div v-else-if="!payments.length" class="text-sm text-slate-200/70">暂无充值记录。</div>
      <ul v-else class="divide-y divide-white/5">
        <li v-for="order in payments" :key="order.id" class="py-4 flex items-center justify-between text-sm text-slate-200/80">
          <div class="space-y-1">
            <p class="text-white font-medium break-all">订单 {{ order.out_trade_no }}</p>
            <p class="text-xs text-slate-400">方式：{{ order.pay_type || '未知' }} | 状态：{{ order.status }}</p>
          </div>
          <div class="text-right">
            <p class="text-lg font-semibold text-white">¥{{ (order.money_cents / 100).toFixed(2) }}</p>
            <p class="text-xs text-emerald-300">+ {{ order.coins }} 币</p>
            <p class="text-[11px] text-slate-400">{{ formatDate(order.created_at) }}</p>
          </div>
        </li>
      </ul>
      <div v-if="payments.length" class="mt-4 flex justify-center">
        <button
          class="btn-secondary px-4 py-2 text-sm"
          :disabled="paymentLoading || !paymentHasMore"
          @click="loadPayments(false)"
        >
          {{ paymentHasMore ? (paymentLoading ? '加载中...' : '加载更多') : '已加载全部' }}
        </button>
      </div>
    </section>
  </section>
</template>

<script setup lang="ts">
import type { PaymentOrder } from '@/types'

definePageMeta({
  middleware: ['auth'],
})

const api = useApi()
const payments = ref<PaymentOrder[]>([])
const paymentLoading = ref(false)
const paymentHasMore = ref(true)
const paymentPageSize = 20
const paymentOffset = ref(0)
const loading = ref(false)

const accentPool = [
  ['#a855f7', '#6366f1'],
  ['#4ade80', '#22d3ee'],
  ['#f97316', '#f43f5e'],
]
const heroStyle = computed(() => {
  const palette = accentPool[new Date().getHours() % accentPool.length]
  return `background: radial-gradient(circle at 30% 20%, ${palette[0]}55, transparent), radial-gradient(circle at 80% 0%, ${palette[1]}88, transparent)`
})

const loadPayments = async (reset = false) => {
  if (paymentLoading.value) return
  paymentLoading.value = true
  if (reset) {
    payments.value = []
    paymentOffset.value = 0
    paymentHasMore.value = true
    loading.value = true
  }
  if (!paymentHasMore.value) {
    paymentLoading.value = false
    loading.value = false
    return
  }
  try {
    const res = await api.get<{ data: PaymentOrder[] }>(
      `/store/payments?status=paid&limit=${paymentPageSize}&offset=${paymentOffset.value}`,
    )
    const list = Array.isArray(res.data) ? res.data : []
    if (list.length < paymentPageSize) paymentHasMore.value = false
    payments.value = [...payments.value, ...list]
    paymentOffset.value += list.length
  } finally {
    paymentLoading.value = false
    loading.value = false
  }
}

const formatDate = (ts?: string) => {
  if (!ts) return ''
  const d = new Date(ts)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')} ${String(
    d.getHours(),
  ).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

onMounted(() => {
  loadPayments(true)
})
</script>

<style scoped>
.btn-primary {
  @apply inline-flex items-center gap-2 rounded-full bg-primary px-5 py-2 text-sm font-semibold text-white shadow hover:opacity-90 transition-all;
}
.btn-secondary {
  @apply inline-flex items-center rounded-full border border-white/30 px-5 py-2 text-sm text-white/80 hover:border-white/60 transition;
}
.hero-card {
  background: rgba(15, 23, 42, 0.6);
  box-shadow: 0 25px 80px rgba(15, 23, 42, 0.45);
}
.glass-card {
  @apply relative rounded-3xl border border-white/10 p-6;
  background: rgba(15, 23, 42, 0.55);
  box-shadow: 0 10px 40px rgba(15, 23, 42, 0.3);
}
</style>
