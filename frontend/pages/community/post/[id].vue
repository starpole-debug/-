<template>
  <section v-if="post" class="space-y-6">
    <CommunityPostCard :post="post" @like="like" @favorite="favorite">
      <template #footer>
        <span class="text-muted">👍 {{ reactions.like || 0 }}</span>
        <span class="text-muted">收藏 {{ reactions.favorite || 0 }}</span>
      </template>
    </CommunityPostCard>
    <div v-if="post.attachments?.length" class="grid gap-4 md:grid-cols-3">
      <div v-for="att in post.attachments" :key="att.id || att.file_url" class="overflow-hidden rounded-xl border">
        <img :src="resolveAssetUrl(att.file_url)" class="h-48 w-full object-cover" />
      </div>
    </div>
    <section class="rounded-3xl bg-white p-6 shadow">
      <h2 class="text-lg font-semibold">评论</h2>
      <ul class="mt-3 space-y-3 text-sm text-slate-700">
        <li v-for="comment in comments" :key="comment.id">
          {{ comment.author_name || comment.author_id.slice(0, 6) }}：{{ comment.content }}
        </li>
      </ul>
      <p v-if="errorMessage" class="mt-3 text-sm text-rose-500">{{ errorMessage }}</p>

      <form class="mt-4 flex gap-3" @submit.prevent="submit">
        <input v-model="draft" class="flex-1 rounded-full border px-4 py-2 text-sm" placeholder="说点什么..." />
        <button class="rounded-full bg-primary px-4 py-2 text-sm text-white">发布</button>
      </form>
    </section>
  </section>
  <p v-else class="text-muted">加载中...</p>
</template>

<script setup lang="ts">
import type { CommunityComment } from '@/types'
import { useAssetUrl } from '~/composables/useAssetUrl'
const route = useRoute()
const community = useCommunity()
const auth = useAuth()
const router = useRouter()
const post = ref()
const comments = ref<CommunityComment[]>([])
const reactions = ref<Record<string, number>>({})
const draft = ref('')
const reacting = ref(false)
const errorMessage = ref('')
const { resolveAssetUrl } = useAssetUrl()

const load = async () => {
  const data = await community.fetchPost(route.params.id as string)
  post.value = {
    ...data.post,
    attachments: (data.post?.attachments || []).map((att: any) => ({
      ...att,
      file_url: resolveAssetUrl(att.file_url),
    })),
  }
  comments.value = data.comments
  reactions.value = data.reactions || {}
}

const submit = async () => {
  if (!draft.value.trim()) return
  await community.addComment(route.params.id as string, { content: draft.value })
  draft.value = ''
  await load()
}

const like = async () => {
  await react('like')
}

const favorite = async () => {
  await react('favorite')
}

const react = async (type: string) => {
  if (!auth.isAuthenticated.value) {
    router.push('/login')
    return
  }
  if (reacting.value) return
  reacting.value = true
  errorMessage.value = ''
  try {
    const result = await community.react(route.params.id as string, type)
    if (result?.reactions) {
      reactions.value = result.reactions
    } else {
      await load()
    }
  } catch (error: any) {
    errorMessage.value = error?.data?.error || error?.message || '操作失败，请稍后再试'
  } finally {
    reacting.value = false
  }
}

onMounted(load)
</script>
