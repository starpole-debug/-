<template>
  <section>
    <h1 class="text-2xl font-semibold">后台总览</h1>
    <p class="mt-2 text-sm text-slate-300">使用左侧导航切换模块。</p>
    </section>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'admin',
  middleware: 'admin',
})

</script>
