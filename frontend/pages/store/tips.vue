<template>
  <section class="rounded-3xl bg-white p-8 shadow">
    <h1 class="text-xl font-semibold">打赏流程</h1>
    <ol class="mt-4 list-decimal space-y-2 pl-6 text-sm text-muted">
      <li>前端收集金额后调用 /api/store/tips。</li>
      <li>后端 revenue.Service 会写入 revenue_events 并更新 creator_wallets。</li>
      <li>钱包页面会展示实时余额与提现记录。</li>
    </ol>
  </section>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: ['auth'],
})
</script>
