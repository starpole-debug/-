<template>
  <section class="space-y-4">
    <header class="flex items-center justify-between">
      <h1 class="text-xl font-semibold">角色列表</h1>
      <NuxtLink to="/creator/editor/new" class="rounded-full bg-primary px-4 py-2 text-sm text-white">新建</NuxtLink>
    </header>

    <div class="space-y-3">
      <h2 class="text-sm text-slate-300">已发布</h2>
      <div v-if="published.length === 0" class="text-sm text-slate-500">暂无已发布角色</div>
      <CreatorRoleRow v-for="role in published" :key="role.id" :role="role" />
    </div>

    <div class="space-y-3">
      <h2 class="text-sm text-slate-300">草稿</h2>
      <div v-if="drafts.length === 0" class="text-sm text-slate-500">暂无草稿</div>
      <CreatorRoleRow v-for="role in drafts" :key="role.id" :role="role" @publish="handlePublish" />
    </div>
  </section>
</template>

<script setup lang="ts">
const { roles, loadRoles, publishRole } = useCreatorStats()
const published = computed(() => roles.value.filter((r) => r.status === 'published'))
const drafts = computed(() => roles.value.filter((r) => r.status !== 'published'))

const handlePublish = async (id: string) => {
  await publishRole(id)
}

onMounted(loadRoles)
</script>
