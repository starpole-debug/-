<template>
  <section class="rounded-3xl bg-white p-8 shadow">
    <h1 class="text-xl font-semibold">发布帖子</h1>
    <form class="mt-4 space-y-3" @submit.prevent="submit">
      <div>
        <label class="text-xs text-muted">标题</label>
        <input v-model="payload.title" class="mt-1 w-full rounded-xl border px-4 py-2 text-sm" />
      </div>
      <div>
        <label class="text-xs text-muted">内容</label>
        <textarea v-model="payload.content" class="mt-1 h-48 w-full rounded-xl border px-4 py-2 text-sm"></textarea>
      </div>
      <div class="grid gap-3 md:grid-cols-2">
        <div>
          <label class="text-xs text-muted">跳转链接（可选）</label>
          <input
            v-model="payload.link_url"
            class="mt-1 w-full rounded-xl border px-4 py-2 text-sm"
            placeholder="如 /role/<id> 或 /creator/presets/download/<id> 或 https://example.com"
          />
        </div>
        <div>
          <label class="text-xs text-muted">链接类型</label>
          <select v-model="payload.link_type" class="mt-1 w-full rounded-xl border px-4 py-2 text-sm">
            <option value="">自动识别</option>
            <option value="role">角色卡</option>
            <option value="preset">预设</option>
            <option value="external">外部链接</option>
            <option value="internal">站内链接</option>
          </select>
        </div>
      </div>
      <div>
        <label class="text-xs text-muted">上传图片</label>
        <div class="mt-2 flex flex-col gap-2 rounded-xl border border-dashed border-slate-300/70 p-4">
          <div class="flex items-center gap-3">
            <input ref="fileRef" type="file" accept="image/*" @change="onFileChange" />
            <button
              class="rounded-full bg-primary px-4 py-1 text-xs text-white disabled:opacity-50"
              type="button"
              :disabled="uploading"
              @click="triggerUpload"
            >
              {{ uploading ? '上传中...' : '上传图片' }}
            </button>
          </div>
          <p v-if="uploadError" class="text-xs text-rose-500">{{ uploadError }}</p>
          <div v-if="attachments.length" class="flex flex-wrap gap-3">
            <div v-for="item in attachments" :key="item" class="relative h-20 w-28 overflow-hidden rounded-lg border">
              <img :src="resolveAssetUrl(item)" class="h-full w-full object-cover" />
              <button type="button" class="absolute right-1 top-1 rounded-full bg-black/60 px-2 text-[10px] text-white" @click="removeAttachment(item)">
                ×
              </button>
            </div>
          </div>
        </div>
      </div>
      <p v-if="errorMessage" class="text-sm text-rose-500">{{ errorMessage }}</p>
      <button
        class="rounded-full bg-primary px-6 py-2 text-sm text-white disabled:opacity-50 disabled:cursor-not-allowed"
        :disabled="submitting"
      >
        {{ submitting ? '提交中...' : '提交' }}
      </button>
    </form>
  </section>
</template>

<script setup lang="ts">
import { useUpload } from '~/composables/useUpload'
import { useAssetUrl } from '~/composables/useAssetUrl'
definePageMeta({
  middleware: ['auth'],
})

const community = useCommunity()
const router = useRouter()
const payload = reactive({ title: '', content: '', link_url: '', link_type: '' })
const { resolveAssetUrl } = useAssetUrl()
const errorMessage = ref('')
const submitting = ref(false)
const uploader = useUpload()
const fileRef = ref<HTMLInputElement | null>(null)
const pendingFile = ref<File | null>(null)
const attachments = ref<string[]>([])
const uploading = ref(false)
const uploadError = ref('')

const submit = async () => {
  errorMessage.value = ''
  // 如果还有未上传的文件，先执行上传，避免遗漏
  if (pendingFile.value) {
    await triggerUpload()
    if (uploadError.value) {
      return
    }
  }
  if (uploading.value) {
    errorMessage.value = '图片正在上传，请稍候再提交'
    return
  }
  submitting.value = true
  try {
    await community.createPost({ ...payload, attachments: attachments.value })
    router.push('/community')
  } catch (error: any) {
    errorMessage.value = error?.data?.error || error?.message || '发布失败，请稍后再试'
  } finally {
    submitting.value = false
  }
}

const onFileChange = async (e: Event) => {
  const target = e.target as HTMLInputElement
  const file = target.files?.[0]
  if (file) {
    pendingFile.value = file
    uploadError.value = ''
    // 选中文件后立即尝试上传，避免忘记点击按钮
    await triggerUpload()
  }
}

const triggerUpload = async () => {
  if (uploading.value) return
  if (!pendingFile.value) {
    uploadError.value = '请选择图片后再上传'
    return
  }
  uploadError.value = ''
  uploading.value = true
  try {
    const url = await uploader.uploadPostImage(pendingFile.value)
    if (url) {
      attachments.value.push(url)
      pendingFile.value = null
      if (fileRef.value) fileRef.value.value = ''
    } else {
      uploadError.value = '上传失败，请稍后重试'
    }
  } catch (error: any) {
    uploadError.value = error?.data?.error || error?.message || '上传失败，请稍后重试'
  } finally {
    uploading.value = false
  }
}

const removeAttachment = (url: string) => {
  attachments.value = attachments.value.filter((item) => item !== url)
}
</script>
