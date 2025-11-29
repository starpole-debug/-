<template>
  <section>
    <h1 class="text-xl font-semibold">角色审核</h1>
    <div class="mt-4 space-y-3">
      <div v-for="role in roles" :key="role.id" class="rounded-2xl border border-white/10 p-4">
        <p class="font-semibold">{{ role.name }}<span class="ml-2 text-xs text-slate-400">{{ role.status }}</span></p>
        <p class="text-xs text-slate-400">{{ role.description }}</p>
      </div>
    </div>
    </section>
</template>

<script setup lang="ts">
const { roles, loadRoles } = useAdminConfig()

definePageMeta({
  layout: 'admin',
  middleware: 'admin',
})


onMounted(loadRoles)
</script>
