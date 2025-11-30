<template>
  <section class="space-y-8">
    <header class="card-soft p-8 animate-fade-up relative overflow-hidden">
      <div class="absolute top-0 right-0 w-32 h-32 bg-accent-yellow/20 rounded-full blur-3xl"></div>
      <div class="relative z-10">
        <h1 class="text-2xl font-bold text-charcoal-900">平台充值</h1>
        <p class="text-sm text-charcoal-500 mt-2 max-w-xl">
          您在这里支付的金额会折算为平台内虚拟货币，用于后续与 AI 角色对话等功能消耗。
        </p>
      </div>
    </header>

    <div class="card-soft p-8 animate-fade-up delay-100 space-y-6">
      <div>
        <h2 class="text-lg font-bold text-charcoal-900 mb-4 flex items-center gap-2">
          <span class="w-1.5 h-6 rounded-full bg-accent-pink"></span>
          选择支持金额
        </h2>
        <StoreTipOptions
          :amounts="tipOptions?.amounts || []"
          @select="pay"
          :class="{ 'pointer-events-none opacity-50': !tipOptions?.amounts?.length }"
        />
        <p v-if="tipOptionsLoading" class="mt-2 text-xs text-charcoal-400">加载可用打赏档位...</p>
        <div class="mt-4 flex flex-wrap gap-2 items-center">
          <span class="text-sm text-charcoal-500">支付方式：</span>
          <button
            v-for="option in payTypes"
            :key="option.value"
            type="button"
            class="px-3 py-1.5 rounded-full border text-sm transition-colors"
            :class="payType === option.value ? 'border-accent-pink bg-accent-pink/10 text-accent-pink' : 'border-charcoal-200 text-charcoal-500 hover:border-accent-pink/50 hover:text-charcoal-900'"
            @click="payType = option.value"
          >
            {{ option.label }}
          </button>
        </div>
        <div class="mt-4 flex flex-col sm:flex-row gap-3 items-start sm:items-center">
          <label class="text-sm text-charcoal-500">自定义金额（元）</label>
          <input
            v-model.number="customAmount"
            type="number"
            min="1"
            step="0.01"
            class="glass-input w-full sm:w-48 text-sm"
            placeholder="输入金额"
          />
          <button
            class="btn-primary text-sm px-4 py-2 rounded-full disabled:opacity-50 disabled:cursor-not-allowed shadow-lg shadow-charcoal-900/10"
            :disabled="!customAmount || customAmount <= 0 || tipping"
            @click="pay(customAmount)"
          >
            立即支付
          </button>
        </div>
      </div>

      <p class="text-xs text-charcoal-400">当前汇率：1 元 = 1000 平台币</p>
      <p v-if="statusMessage" class="text-sm text-status-success">{{ statusMessage }}</p>
      <p v-if="errorMessage" class="text-sm text-status-error">{{ errorMessage }}</p>
    </div>
  </section>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: ['auth'],
})

const api = useApi()
const router = useRouter()
const tipping = ref(false)
const tipOptions = ref<{ amounts: number[]; descriptions: string[] } | null>(null)
const tipOptionsLoading = ref(false)
const statusMessage = ref('')
const errorMessage = ref('')
const polling = ref<ReturnType<typeof setInterval> | null>(null)
const customAmount = ref<number | null>(null)
const payType = ref<'alipay' | 'wxpay'>('alipay')
const payTypes = [
  { value: 'alipay', label: '支付宝' },
  { value: 'wxpay', label: '微信' },
]

const loadTipOptions = async () => {
  tipOptionsLoading.value = true
  try {
    const res = await api.get<{ data: { amounts: number[]; descriptions: string[] } }>('/store/options')
    tipOptions.value = res.data
  } catch (error: any) {
    console.error('loadTipOptions error', error)
  } finally {
    tipOptionsLoading.value = false
  }
}

const pay = async (amount: number) => {
  if (tipping.value) return
  if (!amount || amount <= 0) {
    errorMessage.value = '请输入正确的金额'
    return
  }
  statusMessage.value = ''
  errorMessage.value = ''
  tipping.value = true
  try {
    const res = await api.post<{ data: { pay_url: string; out_trade_no: string; coins: number; amount: number } }>(
      '/store/payments',
      { amount, pay_type: payType.value },
    )
    const payUrl = res.data?.pay_url
    const outTradeNo = res.data?.out_trade_no
    if (payUrl) {
      window.open(payUrl, '_blank')
    }
    statusMessage.value = `已创建订单，请在新页面完成支付，金额 ¥${amount.toFixed(2)}，可得 ${Math.round(amount * 1000)} 平台币`
    // 开始轮询订单状态
    if (polling.value) clearInterval(polling.value)
    if (outTradeNo) {
      polling.value = setInterval(async () => {
        try {
          const order = await api.get<{ data: any }>(`/store/payments/${outTradeNo}`)
          if (order.data?.status === 'paid') {
            statusMessage.value = `支付成功，到账 ${order.data?.coins || Math.round(amount * 1000)} 平台币`
            if (polling.value) {
              clearInterval(polling.value)
              polling.value = null
            }
            // 稍等片刻后跳转个人主页
            setTimeout(() => router.push('/me'), 800)
          }
        } catch {
          // ignore transient errors
        }
      }, 3500)
    }
  } catch (error: any) {
    errorMessage.value = error?.data?.error || error?.message || '打赏失败，请稍后再试'
  } finally {
    tipping.value = false
  }
}

onMounted(loadTipOptions)
onBeforeUnmount(() => {
  if (polling.value) clearInterval(polling.value)
})
</script>
