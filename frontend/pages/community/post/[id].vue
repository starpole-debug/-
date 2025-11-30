<template>
  <section v-if="post" class="space-y-6">
    <CommunityPostCard :post="post" @like="like" @favorite="favorite">
      <template #footer>
        <div class="flex gap-4">
          <span class="text-charcoal-500 flex items-center gap-1">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-4 h-4 text-status-error">
              <path d="M1 8.25a1.25 1.25 0 112.5 0v7.5a1.25 1.25 0 11-2.5 0v-7.5zM11 3V1.7c0-.268.14-.526.395-.607A2 2 0 0114 3c0 .995-.182 1.948-.514 2.826-.204.54.166 1.174.744 1.174h2.52c1.243 0 2.261 1.01 2.146 2.247a23.864 23.864 0 01-1.341 5.974C17.153 16.323 16.072 17 14.9 17h-3.192a3 3 0 01-1.341-.317l-2.734-1.366A3 3 0 006.292 15H5V8h.963c.685 0 1.258-.483 1.612-1.068a4.011 4.011 0 012.166-1.73c.432-.143.853-.386 1.011-.814.16-.432.448-.779.9-.976.45-.198.937-.312 1.348-.312z" />
            </svg>
            {{ reactions.like || 0 }}
          </span>
          <span class="text-charcoal-500 flex items-center gap-1">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-4 h-4 text-accent-yellow">
              <path fill-rule="evenodd" d="M10 2c-1.716 0-3.408.106-5.07.31C3.806 2.45 3 3.414 3 4.517V17.25a.75.75 0 001.075.676L10 15.082l5.925 2.844A.75.75 0 0017 17.25V4.517c0-1.103-.806-2.068-1.93-2.207A41.403 41.403 0 0010 2z" clip-rule="evenodd" />
            </svg>
            {{ reactions.favorite || 0 }}
          </span>
        </div>
      </template>
    </CommunityPostCard>
    
    <div v-if="post.attachments?.length" class="grid gap-4 md:grid-cols-3">
      <div v-for="att in post.attachments" :key="att.id || att.file_url" class="overflow-hidden rounded-2xl border border-white/50 shadow-sm">
        <img :src="resolveAssetUrl(att.file_url)" class="h-48 w-full object-cover hover:scale-105 transition-transform duration-500" />
      </div>
    </div>

    <section class="card-soft p-6">
      <h2 class="text-lg font-display font-bold text-charcoal-900 mb-4">评论</h2>
      <ul class="space-y-4">
        <li v-for="comment in comments" :key="comment.id" class="flex gap-3 p-3 rounded-xl bg-bg-cream-100/50">
          <div class="w-8 h-8 rounded-full bg-accent-yellow flex items-center justify-center text-xs font-bold text-charcoal-900 shrink-0">
            {{ (comment.author_name || comment.author_id || '?').charAt(0).toUpperCase() }}
          </div>
          <div class="space-y-1">
            <div class="text-xs font-bold text-charcoal-900">
              {{ comment.author_name || comment.author_id.slice(0, 6) }}
            </div>
            <div class="text-sm text-charcoal-700 leading-relaxed">
              {{ comment.content }}
            </div>
          </div>
        </li>
      </ul>
      <p v-if="errorMessage" class="mt-3 text-sm text-status-error">{{ errorMessage }}</p>

      <form class="mt-6 flex gap-3" @submit.prevent="submit">
        <input v-model="draft" class="glass-input rounded-full" placeholder="说点什么..." />
        <button class="btn-primary rounded-full px-6 shadow-lg shadow-charcoal-900/10">发布</button>
      </form>
    </section>
  </section>
  <p v-else class="text-charcoal-500 text-center py-12">加载中...</p>
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
