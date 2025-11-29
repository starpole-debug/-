<template>
  <section>
    <header class="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
      <div>
        <h1 class="text-xl font-semibold">用户管理</h1>
        <p class="text-sm text-slate-300">搜索、封禁或删除用户账号。</p>
      </div>
      <div class="flex gap-3">
        <input v-model="query" placeholder="用户名 / 邮箱 / 昵称" class="rounded-full border border-white/10 bg-white/5 px-4 py-2 text-sm" />
        <button class="rounded-full bg-white/10 px-4 py-2 text-sm text-white" :disabled="loading" @click="search">
          {{ loading ? '搜索中...' : '搜索' }}
        </button>
        <button class="rounded-full bg-primary px-4 py-2 text-sm text-white" @click="showForm = !showForm">{{ showForm ? '收起' : '创建用户' }}</button>
      </div>
    </header>

    <form v-if="showForm" class="mt-6 grid gap-3 rounded-2xl border border-white/10 p-5 sm:grid-cols-2" @submit.prevent="create">
      <label class="text-sm">
        用户名
        <input v-model="form.username" class="mt-1 w-full rounded-full border border-white/10 bg-white/5 px-4 py-2 text-sm" required />
      </label>
      <label class="text-sm">
        邮箱
        <input v-model="form.email" class="mt-1 w-full rounded-full border border-white/10 bg-white/5 px-4 py-2 text-sm" />
      </label>
      <label class="text-sm">
        昵称
        <input v-model="form.nickname" class="mt-1 w-full rounded-full border border-white/10 bg-white/5 px-4 py-2 text-sm" />
      </label>
      <label class="text-sm">
        密码
        <input v-model="form.password" type="password" class="mt-1 w-full rounded-full border border-white/10 bg-white/5 px-4 py-2 text-sm" />
      </label>
      <label class="flex items-center gap-2 text-sm">
        <input v-model="form.is_admin" type="checkbox" class="h-4 w-4" />
        管理员账号
      </label>
      <div class="sm:col-span-2">
        <button class="rounded-full bg-white/10 px-5 py-2 text-sm text-white" :disabled="submitting">
          {{ submitting ? '创建中...' : '创建' }}
        </button>
      </div>
    </form>
    <p v-if="error" class="mt-3 text-sm text-rose-300">{{ error }}</p>

    <div class="mt-6 overflow-auto rounded-2xl border border-white/10">
      <table class="min-w-full text-left text-sm">
        <thead class="bg-white/5 text-xs uppercase text-slate-400">
          <tr>
            <th class="px-4 py-3">ID</th>
            <th>用户名</th>
            <th>Email</th>
            <th>昵称</th>
            <th>Admin</th>
            <th>封禁</th>
            <th>创建时间</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="user in users" :key="user.id" class="border-t border-white/10">
            <td class="px-4 py-3 text-xs text-slate-400">{{ user.id }}</td>
            <td>{{ user.username }}</td>
            <td>{{ user.email || '-' }}</td>
            <td>{{ user.nickname }}</td>
            <td>{{ user.is_admin ? '是' : '否' }}</td>
            <td>
              <span
                class="rounded-full px-3 py-1 text-xs"
                :class="user.is_banned ? 'bg-rose-500/20 text-rose-200' : 'bg-emerald-500/20 text-emerald-200'"
              >
                {{ user.is_banned ? '已封禁' : '正常' }}
              </span>
            </td>
            <td>{{ formatTime(user.created_at) }}</td>
            <td class="space-x-2">
              <button class="text-xs text-primary" :disabled="actionLoading === user.id" @click="toggleBan(user)">
                {{ actionLoading === user.id ? '处理中...' : user.is_banned ? '解封' : '封禁' }}
              </button>
              <button class="text-xs text-rose-400" :disabled="actionLoading === user.id" @click="remove(user.id)">删除</button>
            </td>
          </tr>
        </tbody>
      </table>
      <p v-if="!users.length && !loading" class="p-4 text-sm text-slate-400">暂无用户。</p>
    </div>
    </section>
</template>

<script setup lang="ts">
import dayjs from 'dayjs'
import type { User } from '@/types'

definePageMeta({
  layout: 'admin',
  middleware: 'admin'
})


const { users, loading, fetchUsers, createUser, banUser, unbanUser, deleteUser } = useAdminUsers()
const query = ref('')
const showForm = ref(false)
const submitting = ref(false)
const actionLoading = ref('')
const error = ref('')
const form = reactive({
  username: '',
  email: '',
  password: '',
  nickname: '',
  is_admin: false
})

const search = async () => {
  error.value = ''
  try {
    await fetchUsers({ query: query.value })
  } catch (err: any) {
    error.value = err?.data?.error || err?.message || '加载失败'
  }
}

const create = async () => {
  submitting.value = true
  error.value = ''
  try {
    await createUser(form)
    Object.assign(form, { username: '', email: '', password: '', nickname: '', is_admin: false })
    await search()
    showForm.value = false
  } catch (err: any) {
    error.value = err?.data?.error || err?.message || '创建失败'
  } finally {
    submitting.value = false
  }
}

const toggleBan = async (user: User) => {
  actionLoading.value = user.id
  try {
    if (user.is_banned) {
      await unbanUser(user.id)
    } else {
      await banUser(user.id)
    }
    await search()
  } catch (err: any) {
    error.value = err?.data?.error || err?.message || '操作失败'
  } finally {
    actionLoading.value = ''
  }
}

const remove = async (id: string) => {
  if (!confirm('确定删除该用户吗？')) return
  actionLoading.value = id
  try {
    await deleteUser(id)
    await search()
  } catch (err: any) {
    error.value = err?.data?.error || err?.message || '删除失败'
  } finally {
    actionLoading.value = ''
  }
}

const formatTime = (value?: string) => {
  if (!value) return ''
  return dayjs(value).format('YYYY-MM-DD')
}

onMounted(search)
</script>
