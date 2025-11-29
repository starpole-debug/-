<template>
  <section class="rounded-3xl bg-white p-8 shadow space-y-4">
    <header>
      <h1 class="text-xl font-semibold">账号设置</h1>
      <p class="mt-2 text-sm text-muted">更新昵称与头像信息，保存后立即生效。</p>
    </header>

    <form class="space-y-4" @submit.prevent="submit">
      <div>
        <label class="text-xs text-muted">昵称</label>
        <input v-model.trim="form.nickname" class="mt-1 w-full rounded-xl border px-4 py-2 text-sm" placeholder="展示名称" />
      </div>
      <div>
        <label class="text-xs text-muted">头像</label>
        <div class="mt-2 flex flex-col gap-3 rounded-xl border border-dashed border-slate-300/70 p-4 md:flex-row md:items-center">
          <img :src="avatarPreview" class="h-16 w-16 rounded-full border object-cover" alt="avatar preview" />
          <div class="space-y-2">
            <input ref="avatarInput" type="file" accept="image/*" @change="onAvatarFile" />
            <div class="flex items-center gap-3">
              <input v-model.trim="form.avatar_url" class="flex-1 rounded-xl border px-4 py-2 text-sm" placeholder="https://example.com/avatar.png" />
              <button
                type="button"
                class="rounded-full bg-primary px-4 py-2 text-xs text-white disabled:opacity-50"
                :disabled="uploadingAvatar"
                @click="uploadAvatar"
              >
                {{ uploadingAvatar ? '上传中...' : '上传文件' }}
              </button>
            </div>
          </div>
        </div>
      </div>
      <p v-if="errorMessage" class="text-sm text-rose-500">{{ errorMessage }}</p>
      <p v-if="successMessage" class="text-sm text-emerald-500">{{ successMessage }}</p>
      <button type="submit" class="rounded-full bg-primary px-5 py-2 text-sm text-white disabled:opacity-50 disabled:cursor-not-allowed" :disabled="submitting">
        {{ submitting ? '保存中...' : '保存' }}
      </button>
    </form>
  </section>
</template>

<script setup lang="ts">
import type { User } from '@/types'

definePageMeta({
  middleware: ['auth'],
})

import { useUpload } from '~/composables/useUpload'
import { useAssetUrl } from '~/composables/useAssetUrl'

const api = useApi()
const auth = useAuth()
const uploader = useUpload()
const { resolveAssetUrl } = useAssetUrl()
const defaultAvatar = 'https://placehold.co/80x80?text=Me'
const form = reactive({
  nickname: auth.user.value?.nickname ?? '',
  avatar_url: auth.user.value?.avatarUrl ?? auth.user.value?.avatar_url ?? '',
})
const submitting = ref(false)
const errorMessage = ref('')
const successMessage = ref('')
const avatarInput = ref<HTMLInputElement | null>(null)
const pendingAvatar = ref<File | null>(null)
const uploadingAvatar = ref(false)
const avatarPreview = computed(() => resolveAssetUrl(form.avatar_url || defaultAvatar))

watch(
  () => auth.user.value,
  (user) => {
    if (!user) return
    form.nickname = user.nickname ?? form.nickname
    form.avatar_url = user.avatarUrl ?? user.avatar_url ?? form.avatar_url
  },
)

const submit = async () => {
  submitting.value = true
  errorMessage.value = ''
  successMessage.value = ''
  try {
    const res = await api.put<{ data: User }>('/auth/me', {
      nickname: form.nickname,
      avatar_url: form.avatar_url,
    })
    const current = auth.user.value
    auth.user.value = current
      ? { ...current, nickname: res.data.nickname, avatar_url: res.data.avatar_url, avatarUrl: res.data.avatar_url }
      : { ...res.data, token: '' }
    successMessage.value = '保存成功'
  } catch (error: any) {
    errorMessage.value = error?.data?.error || error?.message || '保存失败，请稍后重试'
  } finally {
    submitting.value = false
  }
}

const onAvatarFile = (e: Event) => {
  const target = e.target as HTMLInputElement
  const file = target.files?.[0]
  if (file) {
    pendingAvatar.value = file
    errorMessage.value = ''
  }
}

const uploadAvatar = async () => {
  if (!pendingAvatar.value) {
    errorMessage.value = '请选择图片后再上传'
    return
  }
  uploadingAvatar.value = true
  errorMessage.value = ''
  successMessage.value = ''
  try {
    const url = await uploader.uploadAvatar(pendingAvatar.value)
    if (url) {
      form.avatar_url = url
      successMessage.value = '头像上传成功，请点击保存以生效'
      pendingAvatar.value = null
      if (avatarInput.value) avatarInput.value.value = ''
    } else {
      errorMessage.value = '上传失败，请稍后再试'
    }
  } catch (error: any) {
    errorMessage.value = error?.data?.error || error?.message || '上传失败，请稍后再试'
  } finally {
    uploadingAvatar.value = false
  }
}
</script>
