<template>
  <div class="flex h-full">
    <!-- Sidebar -->
    <aside class="w-80 bg-[#101018] border-r border-white/5 flex flex-col">
      <div class="h-16 px-4 border-b border-white/5 flex items-center justify-between">
        <div class="flex items-center gap-3">
          <div class="w-8 h-8 rounded-full bg-pink-600 text-white text-sm font-semibold flex items-center justify-center shadow-lg shadow-pink-900/40">
            对话
          </div>
          <div class="flex items-center gap-2 text-xs">
            <button
              class="px-3 py-1 rounded-md border border-white/10 transition-colors"
              :class="activeList === 'sessions' ? 'bg-pink-600 text-white shadow-pink-900/30 shadow' : 'text-gray-300 hover:text-white'"
              @click="activeList = 'sessions'"
            >
              会话
            </button>
            <button
              class="px-3 py-1 rounded-md border border-white/10 transition-colors"
              :class="activeList === 'favorites' ? 'bg-pink-600 text-white shadow-pink-900/30 shadow' : 'text-gray-300 hover:text-white'"
              @click="activeList = 'favorites'"
            >
              收藏角色
            </button>
          </div>
        </div>
        <NuxtLink to="/search" class="text-xs text-pink-200 hover:text-white transition-colors">探索</NuxtLink>
      </div>

      <div class="px-4 py-3">
        <div class="relative">
          <div class="i-ph-magnifying-glass text-gray-500 absolute left-3 top-2.5"></div>
          <input
            v-model="search"
            :placeholder="activeList === 'sessions' ? '搜索会话' : '搜索收藏角色'"
            class="w-full bg-white/5 border border-white/10 rounded-xl pl-9 pr-3 py-2 text-sm text-gray-200 placeholder:text-gray-500 focus:outline-none focus:border-pink-500 focus:ring-1 focus:ring-pink-500"
          />
        </div>
      </div>
      <div v-if="presetError" class="px-6 py-1 text-[12px] text-rose-400 bg-rose-500/10 border-t border-b border-rose-500/20">
        {{ presetError }}
      </div>

      <div class="flex-1 overflow-y-auto custom-scroll">
        <div v-if="sessionsLoading && activeList === 'sessions'" class="flex justify-center py-10 text-gray-500">
          <div class="i-svg-spinners-90-ring-with-bg text-xl text-pink-400" />
        </div>
        <div v-else-if="favoritesLoading && activeList === 'favorites'" class="flex justify-center py-10 text-gray-500">
          <div class="i-svg-spinners-90-ring-with-bg text-xl text-pink-400" />
        </div>
        <div v-else-if="activeList === 'sessions' && filteredSessions.length === 0" class="px-4 py-6 text-sm text-gray-500 space-y-2">
          <p>还没有对话记录。</p>
          <NuxtLink to="/search" class="text-pink-300 hover:text-white text-xs inline-flex items-center gap-1">
            <span>去找角色开始聊天</span>
            <div class="i-ph-arrow-right" />
          </NuxtLink>
        </div>
        <div v-else-if="activeList === 'favorites' && filteredFavorites.length === 0" class="px-4 py-6 text-sm text-gray-500 space-y-2">
          <p>还没有收藏的角色。</p>
          <NuxtLink to="/search" class="text-pink-300 hover:text-white text-xs inline-flex items-center gap-1">
            <span>去探索并收藏吧</span>
            <div class="i-ph-arrow-right" />
          </NuxtLink>
        </div>
        <div v-else-if="activeList === 'sessions'" class="px-2 space-y-1">
          <NuxtLink
            v-for="session in filteredSessions"
            :key="session.id"
            :to="`/chat/${session.id}`"
            class="block rounded-xl px-3 py-2 border transition-all duration-150"
            :class="session.id === activeSessionId
              ? 'bg-pink-500/10 border-pink-500/30 shadow-inner shadow-pink-900/30'
              : 'border-transparent hover:bg-white/5 hover:border-white/10'"
          >
            <div class="flex items-center justify-between gap-2">
              <p class="text-sm font-medium text-white truncate">{{ session.title }}</p>
              <span class="text-[10px] text-gray-500 whitespace-nowrap">
                {{ new Date(session.updated_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }) }}
              </span>
            </div>
            <p class="text-xs text-gray-500 truncate mt-1" :title="session.last_message || '暂无消息'">
              {{ session.last_message || '暂无消息' }}
            </p>
          </NuxtLink>
        </div>

        <div v-else class="px-3 space-y-2">
          <div
            v-for="role in filteredFavorites"
            :key="role.id"
            class="flex items-center gap-3 p-3 rounded-xl border border-white/10 hover:border-pink-500/50 hover:bg-white/5 transition-colors"
          >
            <div class="w-10 h-10 rounded-full bg-white/10 overflow-hidden flex items-center justify-center">
              <img v-if="role.avatar_url" :src="resolveRoleAvatar(role.avatar_url)" class="w-full h-full object-cover" />
              <span v-else class="text-sm text-gray-300">{{ role.name?.charAt(0) }}</span>
            </div>
            <div class="flex-1 min-w-0">
              <p class="text-sm text-white truncate">{{ role.name }}</p>
              <p class="text-xs text-gray-500 truncate">{{ role.description }}</p>
            </div>
            <button
              class="px-3 py-1 rounded-lg bg-pink-600 text-white text-xs hover:bg-pink-500 disabled:opacity-50"
              :disabled="isLoading"
              @click="startChatWithRole(role.id)"
            >
              开始对话
            </button>
          </div>
        </div>
      </div>

      <div class="p-4 border-t border-white/5">
        <NuxtLink
          to="/search"
          class="w-full inline-flex justify-center rounded-lg bg-pink-600/90 hover:bg-pink-500 text-white text-sm py-2 transition-colors"
        >
          新建对话
        </NuxtLink>
      </div>
    </aside>

    <!-- Conversation Pane -->
    <section class="flex-1 flex flex-col bg-gradient-to-b from-[#0b0b12] to-[#0e0e18]">
      <div class="h-16 flex items-center justify-between border-b border-white/5 px-6 bg-black/40 backdrop-blur">
        <div class="flex items-center gap-3">
          <button class="p-2 rounded-lg hover:bg-white/5 text-gray-400 hover:text-white transition-colors" @click="router.push('/search')">
            <div class="i-ph-caret-left text-lg" />
          </button>
          <div class="w-11 h-11 rounded-2xl bg-white/5 overflow-hidden flex items-center justify-center text-gray-200 text-sm">
            <img
              v-if="sessionView?.role?.avatar_url"
              :src="resolveRoleAvatar(sessionView.role.avatar_url)"
              class="w-full h-full object-cover"
            />
            <span v-else>{{ sessionView?.role?.name?.charAt(0) || 'AI' }}</span>
          </div>
          <div>
            <div class="text-sm font-semibold text-white leading-tight">{{ sessionView?.role?.name || 'AI 会话' }}</div>
            <div class="text-xs text-gray-500">{{ sessionView?.session?.title || '选择左侧会话开始聊天' }}</div>
          </div>
        </div>

        <div class="flex items-center gap-3">
          <button
            class="px-3 py-1 rounded-lg border border-white/10 bg-white/5 text-xs text-white hover:border-pink-500 hover:text-pink-200 transition-colors"
            @click="showUserSettings = true"
          >
            用户设定
          </button>
          <button
            class="px-3 py-1 rounded-lg border border-white/10 bg-white/5 text-xs text-white hover:border-pink-500 hover:text-pink-200 transition-colors"
            :disabled="isLoading || !sessionView?.session?.role_id"
            @click="startFreshChat"
          >
            开启新聊天
          </button>
          <button
            class="px-3 py-1 rounded-lg border border-white/10 bg-white/5 text-xs text-white hover:border-pink-500 hover:text-pink-200 transition-colors"
            :disabled="!roleHistory.length"
            @click="showHistory = true"
          >
            历史记录
          </button>
          <div class="flex items-center gap-2 bg-white/5 border border-white/10 rounded-lg px-2 py-1 text-xs">
            <select
              v-model="selectedPresetId"
              class="bg-transparent text-gray-100 focus:outline-none"
              :disabled="isLoading"
            >
              <option value="">未使用预设</option>
              <option v-for="p in userPresets" :key="p.id" :value="p.id">
                {{ p.name || '未命名预设' }}
              </option>
            </select>
            <button class="px-2 py-1 rounded-md bg-white/10 hover:bg-white/20" @click="openPresetDetails" :disabled="!selectedPreset">
              查看
            </button>
            <button class="px-2 py-1 rounded-md bg-white/10 hover:bg-white/20 text-rose-300" @click="deleteSelectedPreset" :disabled="!selectedPresetId">
              删除
            </button>
            <button class="px-2 py-1 rounded-md bg-white/10 hover:bg-white/20" @click="triggerPresetImport">
              导入
            </button>
          </div>
          <div class="flex items-center gap-1 bg-white/5 border border-white/10 rounded-lg p-1 text-xs">
            <button
              class="px-3 py-1 rounded-md transition-colors"
              :class="isSFW ? 'bg-pink-600 text-white shadow-pink-900/30 shadow' : 'text-gray-300 hover:text-white'"
              :disabled="isUpdatingSettings || togglingMode || !sessionView"
              @click="handleToggleMode('sfw')"
            >
              SFW
            </button>
            <button
              class="px-3 py-1 rounded-md transition-colors"
              :class="!isSFW ? 'bg-white/10 text-white' : 'text-gray-400 hover:text-white'"
              :disabled="isUpdatingSettings || togglingMode || !sessionView"
              @click="handleToggleMode('nsfw')"
            >
              NSFW
            </button>
          </div>
          <div class="bg-white/5 rounded-lg px-3 py-1.5 flex items-center gap-2 border border-white/10">
            <div class="text-[11px] text-gray-400">模型</div>
            <select
              v-model="selectedModel"
              @change="handleSelectModel(($event.target as HTMLSelectElement).value)"
              class="bg-transparent text-sm text-gray-100 focus:outline-none"
              :disabled="modelChanging || isUpdatingSettings || !sessionView"
            >
              <option value="">默认</option>
              <option v-for="model in models" :key="model.id" :value="model.id">
                {{ model.name }}（{{ formatModelPrice(model) }}）
              </option>
            </select>
            <div v-if="modelChanging" class="i-svg-spinners-3-dots-fade text-gray-400" />
          </div>
          <label class="flex items-center gap-2 text-xs text-gray-300 cursor-pointer">
            <input type="checkbox" v-model="enableStream" class="rounded bg-slate-900 border-slate-700 text-pink-500 focus:ring-pink-500" />
            流式传输
          </label>
        </div>
      </div>

      <div ref="messagesContainer" class="flex-1 overflow-y-auto px-6 py-6 space-y-4 custom-scroll">
        <div v-if="loadingError" class="text-center text-sm text-rose-400 py-10">{{ loadingError }}</div>
        <div v-else-if="isLoading && displayMessages.length === 0" class="flex justify-center py-10">
          <div class="i-svg-spinners-90-ring-with-bg text-2xl text-pink-400" />
        </div>
        <template v-else>
          <div v-if="displayMessages.length === 0" class="w-full py-16 text-center text-gray-500">
            <p class="text-sm">还没有消息，先打个招呼吧。</p>
          </div>
          <div
            v-for="msg in displayMessages"
            :key="msg.id"
            class="group flex gap-3 max-w-5xl"
            :class="msg.role === 'user' ? 'flex-row-reverse ml-auto' : ''"
          >
            <div class="w-10 h-10 rounded-full bg-white/5 flex items-center justify-center text-xs text-gray-300 overflow-hidden">
              <template v-if="msg.role === 'user'">
                我
              </template>
              <template v-else>
                <img
                  v-if="sessionView?.role?.avatar_url"
                  :src="resolveRoleAvatar(sessionView.role.avatar_url)"
                  class="w-full h-full object-cover rounded-full"
                />
                <span v-else>{{ sessionView?.role?.name?.charAt(0) || 'AI' }}</span>
              </template>
            </div>
            <div class="flex flex-col max-w-[75%]" :class="msg.role === 'user' ? 'items-end' : 'items-start'">
              <div
                v-if="editingMessageId === msg.id"
                class="w-full space-y-2"
              >
                <textarea
                  v-model="editingContent"
                  class="w-full rounded-2xl px-4 py-3 text-sm bg-white/10 border border-pink-400 text-white focus:outline-none"
                  rows="3"
                />
                <div class="flex gap-2 text-xs">
                  <button class="px-3 py-1 rounded-md bg-pink-600 text-white" @click="saveEdit(msg)">保存</button>
                  <button class="px-3 py-1 rounded-md bg-white/10 text-slate-200" @click="cancelEdit">取消</button>
                </div>
              </div>
              <div v-else>
                <div
                  class="rounded-2xl px-4 py-3 text-sm leading-relaxed shadow-lg transition"
                  :class="msg.role === 'user'
                    ? 'bg-pink-600 text-white shadow-pink-900/30'
                    : 'bg-white/5 border border-white/10 text-gray-100'"
                  :style="retryingMessageId === msg.id ? 'opacity:0.65; filter:blur(0.2px);' : ''"
                >
                  <div v-html="formatMessageContent(msg.content)"></div>
                  <div v-if="msg.metadata?.reasoning_text" class="mt-2 text-[11px] text-amber-200">
                    <button class="underline hover:text-amber-100" @click="reasoningOpen[msg.id] = !reasoningOpen[msg.id]">
                      {{ reasoningOpen[msg.id] ? '收起思维链' : '展开思维链' }}
                    </button>
                    <div v-if="reasoningOpen[msg.id]" class="mt-1 whitespace-pre-wrap break-words text-amber-100/90">
                      {{ msg.metadata.reasoning_text }}
                    </div>
                  </div>
                  <div v-if="retryingMessageId === msg.id" class="mt-2 flex items-center gap-2 text-[11px] text-pink-200">
                    <div class="i-svg-spinners-3-dots-fade" />
                    重新生成中…
                  </div>
                </div>
                <div class="text-[10px] text-gray-500 mt-1">
                  {{ new Date(msg.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }) }}
                </div>
                <div class="mt-1 flex flex-wrap gap-2 text-[11px] text-slate-400 opacity-0 group-hover:opacity-100 transition">
                  <button v-if="msg.role === 'user' || msg.role === 'assistant'" class="hover:text-white" @click="startEdit(msg)">编辑</button>
                  <button v-if="msg.role === 'assistant'" class="hover:text-white" @click="retryMessage(msg)">重试</button>
                  <button class="hover:text-rose-300" @click="deleteMessage(msg)">撤回</button>
                  <button
                    class="px-2 py-1 rounded-full bg-pink-600 text-white hover:bg-pink-500 disabled:opacity-50 disabled:cursor-not-allowed"
                    :disabled="!isPersistedMessage(msg)"
                    @click="openImageModal(msg)"
                  >生成图片</button>
                </div>
                <div v-if="imageJobs[msg.id]" class="mt-2 rounded-xl border border-white/10 bg-white/5 p-3 text-xs text-slate-200 space-y-2">
                  <div class="flex items-center gap-2">
                    <span class="rounded-full px-2 py-0.5 text-[10px]" :class="imageJobs[msg.id].status === 'succeeded' ? 'bg-emerald-500/20 text-emerald-200' : imageJobs[msg.id].status === 'failed' ? 'bg-rose-500/20 text-rose-200' : 'bg-amber-500/20 text-amber-200'">
                      {{ imageJobs[msg.id].status }}
                    </span>
                  </div>
                  <div v-if="imageJobs[msg.id].error" class="text-rose-300">错误：{{ imageJobs[msg.id].error }}</div>
                  <div v-if="imageJobs[msg.id].url" class="overflow-hidden rounded-lg border border-white/10">
                    <img :src="imageJobs[msg.id].url" class="max-w-xs md:max-w-sm w-full mx-auto object-contain" />
                  </div>
                  <div class="flex justify-end gap-2 text-[11px] text-slate-300">
                    <button class="rounded-full border border-white/20 px-3 py-1 hover:border-white/40 hover:text-white transition" @click="retryImageJob(msg)">重试</button>
                    <button class="rounded-full border border-white/20 px-3 py-1 hover:border-white/40 hover:text-rose-200 transition" @click="clearImageJob(msg.id)">撤回</button>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div v-if="isSendingMessage" class="flex gap-3 items-center text-gray-400">
            <div class="w-10 h-10 rounded-full bg-white/5 flex items-center justify-center">
              <div class="i-svg-spinners-3-dots-fade" />
            </div>
            <div class="text-xs">正在思考...</div>
          </div>
        </template>
      </div>

      <div class="border-t border-white/5 bg-black/40 backdrop-blur px-6 py-4">
        <div class="max-w-5xl mx-auto flex items-end gap-3">
          <textarea
            v-model="input"
            @keydown.enter.exact.prevent="handleSend"
            placeholder="和她聊点什么吧..."
            class="flex-1 bg-white/5 border border-white/10 rounded-2xl px-4 py-3 text-sm text-white focus:outline-none focus:border-pink-500 focus:ring-1 focus:ring-pink-500 min-h-[56px] max-h-32 resize-none"
          />
          <button
            @click="handleSend"
            :disabled="!input.trim() || isSendingMessage || !activeSessionId"
            class="h-12 px-5 rounded-xl bg-pink-600 hover:bg-pink-500 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2 text-sm font-medium shadow-lg shadow-pink-900/30"
          >
            <div class="i-ph-paper-plane-tilt" />
            发送
          </button>
        </div>
        <div class="text-[11px] text-gray-500 text-center mt-2">AI 可能会出错，请谨慎核对关键信息。</div>
      </div>
    </section>

    <!-- User settings modal -->
    <div v-if="showUserSettings" class="modal-backdrop" @click.self="showUserSettings = false">
      <div class="modal-card max-w-md w-full">
        <div class="flex items-center justify-between mb-3">
          <div>
            <h3 class="font-semibold text-white text-lg">用户设定</h3>
            <p class="text-xs text-slate-400">设置 &#123;&#123;user&#125;&#125; 的名字，用于替换占位符并发送给 AI。</p>
          </div>
          <button class="text-gray-400 hover:text-white" @click="showUserSettings = false">✕</button>
        </div>
        <div class="space-y-3">
          <div>
            <label class="text-xs text-slate-400 mb-1 block">用户名 / 昵称</label>
            <input
              v-model.trim="userSettings.name"
              class="w-full rounded-lg border border-white/10 bg-white/5 px-3 py-2 text-sm text-white focus:border-pink-500 focus:outline-none"
              placeholder="如：星极、阿华田"
            />
            <p class="text-[11px] text-slate-500 mt-1">将覆盖账号昵称；&#123;&#123;user&#125;&#125; 占位符会使用此值。</p>
          </div>
          <div>
            <label class="text-xs text-slate-400 mb-1 block">用户设定 / 备注</label>
            <textarea
              v-model.trim="userSettings.bio"
              rows="4"
              class="w-full rounded-lg border border-white/10 bg-white/5 px-3 py-2 text-sm text-white focus:border-pink-500 focus:outline-none"
              placeholder="如：年龄、性格、背景、爱好等，AI 回复时可参考"
            />
            <p class="text-[11px] text-slate-500 mt-1">该设定会一并发送给 AI，帮助更个性化回复。</p>
          </div>
        </div>
        <div class="flex justify-end gap-2 pt-4">
          <button class="px-3 py-1.5 rounded-md bg-white/10 text-sm" @click="showUserSettings = false">取消</button>
          <button class="px-3 py-1.5 rounded-md bg-pink-600 text-sm text-white" @click="saveUserSettings">保存</button>
        </div>
      </div>
    </div>

    <!-- History modal -->
    <div v-if="showHistory" class="modal-backdrop" @click.self="showHistory = false">
      <div class="modal-card">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-400">历史对话</p>
            <h3 class="font-semibold">{{ sessionView?.role?.name || '当前角色' }}</h3>
          </div>
          <button class="text-gray-400 hover:text-white" @click="showHistory = false">✕</button>
        </div>
        <div v-if="roleHistory.length === 0" class="text-sm text-gray-500 py-6 text-center">暂无记录</div>
        <div v-else class="modal-scroll space-y-2 pr-1">
          <div
            v-for="hist in roleHistory"
            :key="hist.id"
            class="p-3 rounded-lg border border-white/10 hover:border-pink-500/40 bg-white/5 transition-colors"
          >
            <div class="flex items-center justify-between text-sm">
              <div class="cursor-pointer" @click="router.push(`/chat/${hist.id}`); showHistory = false">
                <span class="text-white truncate">{{ hist.title }}</span>
                <p class="text-xs text-gray-500 mt-0.5 truncate">{{ hist.model_key || '默认模型' }}</p>
              </div>
              <div class="flex items-center gap-2">
                <span class="text-xs text-gray-400 whitespace-nowrap">{{ new Date(hist.updated_at).toLocaleString() }}</span>
                <button class="text-rose-300 text-xs hover:text-rose-200" @click.stop="deleteHistoryItem(hist.id)">删除</button>
              </div>
            </div>
          </div>
        </div>
        <div class="flex justify-end gap-2 pt-2">
          <button class="px-3 py-1.5 rounded-md bg-white/10 text-sm text-rose-300" @click="deleteAllHistory">全部删除</button>
          <button class="px-3 py-1.5 rounded-md bg-white/10 text-sm" @click="showHistory = false">关闭</button>
        </div>
      </div>
    </div>
    <!-- Preset detail modal -->
    <div v-if="showPresetDetail && presetDetail" class="modal-backdrop" @click.self="showPresetDetail = false">
      <div class="modal-card max-w-3xl w-full">
        <div class="flex items-center justify-between mb-2">
          <div class="flex-1 pr-4">
            <p class="text-sm text-gray-400">预设详情</p>
            <input
              v-model="presetDetail.name"
              class="mt-1 w-full rounded-lg border border-white/10 bg-white/5 px-3 py-2 text-sm text-white focus:border-pink-500 focus:outline-none"
              placeholder="未命名预设"
            />
            <input
              v-model="presetDetail.description"
              class="mt-2 w-full rounded-lg border border-white/10 bg-white/5 px-3 py-2 text-xs text-white focus:border-pink-500 focus:outline-none"
              placeholder="预设描述"
            />
          </div>
          <button class="text-gray-400 hover:text-white" @click="showPresetDetail = false">✕</button>
        </div>
        <div class="text-xs text-slate-400 mb-3">
          模型：{{ presetDetail.model_key || '通用' }} · Blocks：{{ presetDetail.blocks?.length || 0 }}
        </div>
        <div class="modal-scroll space-y-3 pr-1 max-h-[60vh]">
          <div
            v-for="(block, idx) in presetDetail.blocks || []"
            :key="idx"
            class="p-3 rounded-lg border border-white/10 bg-white/5"
          >
            <div class="flex flex-wrap items-center gap-2 text-sm text-white">
              <input
                v-model="block.name"
                class="flex-1 min-w-[140px] rounded-md bg-white/10 border border-white/10 px-2 py-1 text-sm"
                placeholder="Block 名称"
              />
              <input
                v-model="block.id"
                class="w-28 rounded-md bg-white/10 border border-white/10 px-2 py-1 text-xs"
                placeholder="ID"
              />
              <input
                v-model="block.role"
                class="w-24 rounded-md bg-white/10 border border-white/10 px-2 py-1 text-xs"
                placeholder="role"
              />
              <button
                class="px-2 py-0.5 rounded-full text-xs border transition-colors"
                :class="block.enabled ? 'bg-emerald-500/10 text-emerald-300 border-emerald-500/30' : 'bg-slate-700/40 text-slate-400 border-slate-600/50'"
                @click="block.enabled = !block.enabled"
              >
                {{ block.enabled ? '启用' : '禁用' }}
              </button>
              <span v-if="block.marker" class="px-2 py-0.5 rounded-full bg-amber-500/10 text-amber-300 text-xs border border-amber-500/30">
                Marker
              </span>
            </div>
            <textarea
              v-model="block.content"
              class="mt-2 w-full rounded-md border border-white/10 bg-white/5 px-2 py-2 text-sm text-slate-200 focus:border-pink-500 focus:outline-none"
              rows="3"
              placeholder="Block 内容"
            />
          </div>
        </div>
        <div class="flex justify-end gap-2 pt-3">
          <button class="px-3 py-1.5 rounded-md bg-white/10 text-sm" @click="showPresetDetail = false">取消</button>
          <button class="px-3 py-1.5 rounded-md bg-pink-600 text-white text-sm" @click="savePresetDetail">保存</button>
        </div>
      </div>
    </div>
    <!-- Image prompt modal -->
    <div v-if="imageModalVisible" class="modal-backdrop" @click.self="imageModalVisible = false">
      <div class="modal-card max-w-xl w-full">
        <div class="flex items-center justify-between mb-3">
          <h3 class="font-semibold text-white text-lg">生成图片</h3>
          <button class="text-gray-400 hover:text-white" @click="imageModalVisible = false">✕</button>
        </div>
        <p class="text-sm text-slate-400 mb-3">可补充绘图提示词，系统会用预设 + 对话内容自动生成最终 prompt。</p>
        <label class="text-xs text-slate-400 mb-1 block">持久提示词（会自动保存，每次生成都会附加）</label>
        <textarea
          v-model="persistentImagePrompt"
          rows="2"
          class="w-full rounded-xl border border-white/10 bg-white/5 px-3 py-2 text-sm mb-3"
          placeholder="例如：偏好二次元风格，保持角色一致"
        />
        <label class="text-xs text-slate-400 mb-1 block">本次补充提示词</label>
        <textarea
          v-model="imagePrompt"
          rows="4"
          class="w-full rounded-xl border border-white/10 bg-white/5 px-3 py-2 text-sm"
          placeholder="例如：夕阳、雨夜、特写等"
        />
        <div class="mt-4 flex justify-end gap-3">
          <button class="rounded-full border border-white/20 px-4 py-2 text-sm" @click="imageModalVisible = false">取消</button>
          <button
            class="rounded-full bg-primary px-5 py-2 text-sm text-white disabled:opacity-60"
            :disabled="isSubmittingImage"
            @click="submitImageJob"
          >
            {{ isSubmittingImage ? '生成中...' : '生成' }}
          </button>
        </div>
      </div>
    </div>
    <input ref="fileInput" type="file" class="hidden" accept=".json,application/json" @change="handlePresetFile" />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch, nextTick } from 'vue'
import type { SessionView } from '@/composables/useChat'
import { useRole } from '@/composables/useRole'
import { useApi } from '@/composables/useApi'
import { useAssetUrl } from '~/composables/useAssetUrl'
import { useImage } from '@/composables/useImage'
import { useAuth } from '@/composables/useAuth'

definePageMeta({
  layout: 'chat',
  middleware: ['auth'],
})

const route = useRoute()
const router = useRouter()
const api = useApi()
const auth = useAuth()
const { createJob, getJob } = useImage()
const {
  sessions,
  models,
  messages,
  isLoading,
  isUpdatingSettings,
  isSendingMessage,
  createSession,
  listSessions,
  fetchModels,
  fetchSession,
  sendMessage,
  updateSettings,
  retryAssistantMessage,
} = useChat()
const { favorites, fetchFavorites } = useRole()
const { resolveAssetUrl } = useAssetUrl()
const resolveRoleAvatar = (url?: string) => resolveAssetUrl(url || '')
const userPresets = ref<any[]>([])
const selectedPresetId = ref('')
const presetError = ref('')
const fileInput = ref<HTMLInputElement | null>(null)
const selectedPreset = computed(() => userPresets.value.find(p => p.id === selectedPresetId.value))
const cachedPresetId = ref('')
const formatModelPrice = (model: any) => {
  const coins = Number(model?.price_coins || 0)
  if (coins > 0) return `${coins} 币/次`
  if (model?.price_hint) return model.price_hint
  return '免费'
}

const editingMessageId = ref('')
const editingContent = ref('')

const startEdit = (msg: any) => {
  editingMessageId.value = msg.id
  editingContent.value = msg.content
}
const cancelEdit = () => {
  editingMessageId.value = ''
  editingContent.value = ''
}
const saveEdit = async (msg: any) => {
  if (!editingContent.value.trim()) return
  try {
    await api.patch(`/chat/messages/${msg.id}`, { content: editingContent.value.trim() })
    editingMessageId.value = ''
    editingContent.value = ''
    await fetchSession(activeSessionId.value)
  } catch (e: any) {
    loadingError.value = e?.message || '修改失败'
  }
}

const deleteMessage = async (msg: any) => {
  if (!activeSessionId.value) return
  if (!msg?.id || typeof msg.id !== 'string' || !msg.id.includes('-')) {
    loadingError.value = '消息还在发送中，稍后再试'
    return
  }
  try {
    await api.del(`/chat/messages/${msg.id}`)
    await fetchSession(activeSessionId.value)
  } catch (e: any) {
    loadingError.value = e?.message || '撤回失败'
  }
}

const retryMessage = async (msg: any) => {
  if (!activeSessionId.value) return
  // 仅对已保存的 AI 消息重试，避免临时 ID 触发 UUID 错误
  if (!msg?.id || typeof msg.id !== 'string' || !msg.id.includes('-')) {
    loadingError.value = '消息还在发送中，请稍后重试'
    return
  }
  if (msg.role !== 'assistant') {
    loadingError.value = '只能重试 AI 回复'
    return
  }
  try {
    retryingMessageId.value = msg.id
    await retryAssistantMessage(msg.id)
    await fetchSession(activeSessionId.value)
  } catch (e: any) {
    loadingError.value = e?.message || '重试失败'
  } finally {
    retryingMessageId.value = ''
  }
}

const openImageModal = (msg: any) => {
  if (!isPersistedMessage(msg)) {
    loadingError.value = '消息尚未保存，无法生成图片'
    return
  }
  imageTargetMsg.value = msg
  imagePrompt.value = ''
  imageModalVisible.value = true
}

const clearImageJob = (messageId: string) => {
  if (!messageId) return
  delete imageJobs[messageId]
}

const submitImageJob = async () => {
  if (!sessionView.value?.session?.id || !imageTargetMsg.value || isSubmittingImage.value) return
  if (!isPersistedMessage(imageTargetMsg.value)) {
    loadingError.value = '消息尚未保存，无法生成图片'
    return
  }
  const messageId = imageTargetMsg.value.id
  const combinedPrompt = [persistentImagePrompt.value, imagePrompt.value].filter(p => p && p.trim()).join('\n')
  imageJobs[messageId] = { status: 'loading', promptUsed: combinedPrompt }
  isSubmittingImage.value = true
  imageModalVisible.value = false
  try {
    const job = await createJob(sessionView.value.session.id, messageId, combinedPrompt)
    imageJobs[messageId] = {
      ...imageJobs[messageId],
      status: job.status,
      url: (job as any).result_url,
      error: job.error || '',
      jobId: job.id,
      finalPrompt: (job as any).final_prompt,
    }
    if (job.status !== 'succeeded' && job.id) {
      for (let i = 0; i < 10; i++) {
        await new Promise(resolve => setTimeout(resolve, 2000))
        const latest = await getJob(job.id)
        imageJobs[messageId] = {
          ...imageJobs[messageId],
          status: latest.status,
          url: (latest as any).result_url,
          error: latest.error || '',
          jobId: latest.id,
          finalPrompt: (latest as any).final_prompt,
        }
        if (latest.status === 'succeeded' || latest.status === 'failed') break
      }
    }
  } catch (e: any) {
    imageJobs[messageId] = { ...imageJobs[messageId], status: 'failed', error: e?.message || '生成失败' }
  } finally {
    isSubmittingImage.value = false
  }
}

const retryImageJob = async (msg: any) => {
  if (!sessionView.value?.session?.id || !msg?.id || isSubmittingImage.value) return
  if (!isPersistedMessage(msg)) {
    loadingError.value = '消息尚未保存，无法重试生成'
    return
  }
  const existing = imageJobs[msg.id] || {}
  const prompt = existing?.promptUsed || existing?.finalPrompt || ''
  // If we don't have a stored prompt, open modal to let user supply one instead of retrying the whole chat.
  if (!prompt) {
    imageTargetMsg.value = msg
    imageModalVisible.value = true
    return
  }
  imageJobs[msg.id] = { ...(existing as any), status: 'loading', error: '' }
  isSubmittingImage.value = true
  try {
    const job = await createJob(sessionView.value.session.id, msg.id, prompt)
    imageJobs[msg.id] = {
      ...imageJobs[msg.id],
      status: job.status,
      url: (job as any).result_url,
      error: job.error || '',
      jobId: job.id,
      finalPrompt: (job as any).final_prompt,
    }
    if (job.status !== 'succeeded' && job.id) {
      for (let i = 0; i < 10; i++) {
        await new Promise(resolve => setTimeout(resolve, 2000))
        const latest = await getJob(job.id)
        imageJobs[msg.id] = {
          ...imageJobs[msg.id],
          status: latest.status,
          url: (latest as any).result_url,
          error: latest.error || '',
          jobId: latest.id,
          finalPrompt: (latest as any).final_prompt,
        }
        if (latest.status === 'succeeded' || latest.status === 'failed') break
      }
    }
  } catch (e: any) {
    imageJobs[msg.id] = { ...imageJobs[msg.id], status: 'failed', error: e?.message || '生成失败' }
  } finally {
    isSubmittingImage.value = false
  }
}

const sessionView = ref<SessionView | null>(null)
const input = ref('')
const search = ref('')
const sessionsLoading = ref(true)
const favoritesLoading = ref(true)
const activeList = ref<'sessions' | 'favorites'>('sessions')
const modelChanging = ref(false)
const selectedModel = ref('')
const loadingError = ref('')
const retryingMessageId = ref('')
const messagesContainer = ref<HTMLElement | null>(null)
const togglingMode = ref(false)
const showHistory = ref(false)
const showPresetDetail = ref(false)
const presetDetail = ref<any | null>(null)
const imageJobs = reactive<Record<string, { status: string; url?: string; error?: string; jobId?: string; finalPrompt?: string; promptUsed?: string }>>({})
const imageModalVisible = ref(false)
const imagePrompt = ref('')
const imageTargetMsg = ref<any>(null)
const persistentImagePrompt = ref('')
const isSubmittingImage = ref(false)
const roleDisplayName = computed(() => sessionView.value?.role?.name || '')
const userSettings = reactive<{ name: string; bio: string }>({ name: '', bio: '' })
const showUserSettings = ref(false)
const enableStream = ref(false)
const reasoningOpen = reactive<Record<string, boolean>>({})
const userDisplayName = computed(() =>
  userSettings.name.trim() ||
  auth.user.value?.nickname ||
  auth.user.value?.username ||
  '用户',
)

const escapeHtml = (str: string) =>
  str
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;')

const foldHiddenTags = (text: string) =>
  (text || '')
    // 隐藏 user-profile 等用户附加标签
    .replace(/<user-profile[^>]*>[\s\S]*?<\/user-profile>/gi, '')
    // 压缩多余空行
    .replace(/\n{2,}/g, '\n')
    .replace(
    /<([a-zA-Z0-9_-]+)(?:\s+[^>]*)?>[\s\S]*?<\/\1>/g,
    () => '',
  )

const applyPlaceholders = (text: string) => {
  const roleName = roleDisplayName.value || '角色'
  const username = userDisplayName.value || '用户'
  let out = text || ''
  out = out.replace(/\{\{\s*char\s*\}\}/gi, roleName)
  out = out.replace(/\{\{\s*user\s*\}\}/gi, username)
  return out
}

const normalizeMessageContent = (raw: string) => {
  let base = (raw || '').trim()
  base = applyPlaceholders(base)
  base = foldHiddenTags(base)
  base = base.replace(/\n{2,}/g, '\n').trim()
  return base
}

const formatMessageContent = (content: string) => {
  const folded = normalizeMessageContent(content)
  if (!folded) return ''

  const codeBlocks: string[] = []
  const codeToken = (i: number) => `__CODE_BLOCK_${i}__`
  let textWithTokens = folded.replace(/```([\s\S]*?)```/g, (_m, code) => {
    const escaped = escapeHtml((code || '').trim())
    const token = codeToken(codeBlocks.length)
    codeBlocks.push(`<pre class="msg-code-block"><code>${escaped}</code></pre>`)
    return token
  })

  let out = escapeHtml(textWithTokens)

  const esc = (s: string) => s.replace(/[-/\\^$*+?.()|[\]{}]/g, '\\$&')
  const doublePatterns = [
    { open: '"', close: '"' },
    { open: '“', close: '”' },
    { open: '「', close: '」' },
    { open: '『', close: '』' },
  ]
  const singlePatterns = [
    { open: "'", close: "'" },
    { open: '‘', close: '’' },
  ]

  doublePatterns.forEach((p) => {
    const regex = new RegExp(`${esc(p.open)}([^${esc(p.close)}]+?)${esc(p.close)}`, 'g')
    out = out.replace(regex, `${p.open}<span class="dialogue-double">$1</span>${p.close}`)
  })
  singlePatterns.forEach((p) => {
    const regex = new RegExp(`${esc(p.open)}([^${esc(p.close)}]+?)${esc(p.close)}`, 'g')
    out = out.replace(regex, `${p.open}<span class="dialogue-single">$1</span>${p.close}`)
  })

  // 处理已转义的英文引号 (&quot; / &#39;)
  out = out.replace(/&quot;(.+?)&quot;/g, '&quot;<span class="dialogue-double">$1</span>&quot;')
  out = out.replace(/&#39;(.+?)&#39;/g, '&#39;<span class="dialogue-single">$1</span>&#39;')

  out = out.replace(/\n/g, '<br />')
  codeBlocks.forEach((block, idx) => {
    out = out.replace(codeToken(idx), block)
  })
  return out
}

const activeSessionId = computed(() => route.params.sessionId as string)
const isSFW = computed(() => (sessionView.value?.session?.mode || 'sfw') === 'sfw')
const firstMessage = computed(() => {
  const data = sessionView.value?.role?.data as Record<string, any> | undefined
  if (!data) return ''
  return data.first_message || data.first_mes || ''
})
const displayMessages = computed(() => {
  const cleaned = (messages.value || []).map(m => ({
    ...m,
    content: normalizeMessageContent(m.content || ''),
  }))
  const list = cleaned.filter(m => (m.content || '').length > 0 || m.metadata?.type === 'first-message')
  if (firstMessage.value) {
    const hasExisting = list.some(m => m.metadata?.type === 'first-message')
    if (!hasExisting) {
      list.unshift({
        id: 'role-first-message',
        session_id: activeSessionId.value || '',
        role: 'assistant',
        content: normalizeMessageContent(firstMessage.value),
        is_important: false,
        created_at: sessionView.value?.session?.created_at || new Date().toISOString(),
        metadata: { type: 'first-message' },
      })
    }
  }
  return list
})
const roleHistory = computed(() => {
  if (!sessionView.value?.session?.role_id) return []
  return sessions.value
    .filter(s => s.role_id === sessionView.value?.session?.role_id)
    .sort((a, b) => new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime())
})

const latestSessionsByRole = computed(() => {
  const map = new Map<string, typeof sessions.value[number]>()
  const sorted = [...sessions.value].sort((a, b) => new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime())
  for (const s of sorted) {
    if (!map.has(s.role_id)) {
      map.set(s.role_id, s)
    }
  }
  return Array.from(map.values())
})

const filteredSessions = computed(() => {
  const term = search.value.trim().toLowerCase()
  const list = latestSessionsByRole.value
  if (!term) return list
  return list.filter(item =>
    item.title.toLowerCase().includes(term) ||
    item.model_key?.toLowerCase().includes(term),
  )
})

const isPersistedMessage = (msg: any) => {
  if (!msg?.id || typeof msg.id !== 'string') return false
  return /^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i.test(msg.id)
}

const handleAfterDeletion = async (deletedIds: string[]) => {
  try {
    await listSessions()
  } catch (e: any) {
    sessions.value = sessions.value.filter(s => !deletedIds.includes(s.id))
    loadingError.value = e?.message || '会话列表刷新失败'
  }
  if (deletedIds.includes(activeSessionId.value)) {
    messages.value = []
    sessionView.value = null
    const fallback = sessions.value.find(s => !deletedIds.includes(s.id))
    if (fallback) {
      await router.push(`/chat/${fallback.id}`)
    } else {
      await router.push('/search')
    }
    return
  }
  if (activeSessionId.value && sessions.value.some(s => s.id === activeSessionId.value)) {
    try {
      await fetchSession(activeSessionId.value)
    } catch (e: any) {
      loadingError.value = e?.message || '无法加载会话'
    }
  }
}

const deleteHistoryItem = async (sessionId: string) => {
  if (!confirm('删除该对话并清除所有消息？')) return
  loadingError.value = ''
  try {
    await api.del(`/chat/sessions/${sessionId}`)
    await handleAfterDeletion([sessionId])
  } catch (e: any) {
    loadingError.value = e?.message || '删除历史失败'
  }
}

const deleteAllHistory = async () => {
  if (!confirm('删除该角色的全部历史记录？此操作不可恢复。')) return
  if (!roleHistory.value.length) return
  loadingError.value = ''
  try {
    const toDelete = roleHistory.value.map(s => s.id)
    await Promise.all(toDelete.map(id => api.del(`/chat/sessions/${id}`)))
    await handleAfterDeletion(toDelete)
  } catch (e: any) {
    loadingError.value = e?.message || '删除历史失败'
  }
}

const filteredFavorites = computed(() => {
  const term = search.value.trim().toLowerCase()
  const list = [...favorites.value]
  if (!term) return list
  return list.filter(item =>
    item.name.toLowerCase().includes(term) ||
    item.description?.toLowerCase().includes(term) ||
    item.tags?.some(tag => tag.toLowerCase().includes(term)),
  )
})

const scrollToBottom = () => {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTo({ top: messagesContainer.value.scrollHeight, behavior: 'smooth' })
    }
  })
}

const loadSession = async (id: string) => {
  if (!id) return
  loadingError.value = ''
  try {
    sessionView.value = await fetchSession(id)
    selectedModel.value = sessionView.value?.session?.model_key || ''
    scrollToBottom()
  } catch (error: any) {
    loadingError.value = error?.message || '无法加载会话'
  }
}

const bootstrap = async () => {
  if (process.client) {
    const cachedPresets = localStorage.getItem('chat:user-presets')
    if (cachedPresets) {
      try {
        userPresets.value = JSON.parse(cachedPresets)
      } catch (_) {}
    }
    const cachedSelected = localStorage.getItem('chat:selected-preset-id')
    if (cachedSelected) {
      cachedPresetId.value = cachedSelected
      selectedPresetId.value = cachedSelected
    }
    const cachedStream = localStorage.getItem('chat:stream-enabled')
    if (cachedStream) {
      enableStream.value = cachedStream === '1'
    }
  }
  try {
    await Promise.all([listSessions(), fetchModels(), fetchFavorites(), loadUserPresets()])
  } catch (error) {
    console.error(error)
  } finally {
    sessionsLoading.value = false
    favoritesLoading.value = false
  }
  if (activeSessionId.value) {
    await loadSession(activeSessionId.value)
  }

  if (process.client) {
    const saved = localStorage.getItem('chat:persistent-image-prompt')
    if (saved) {
      persistentImagePrompt.value = saved
    }
    const savedUser = localStorage.getItem('chat:user-profile')
    if (savedUser) {
      try {
        const parsed = JSON.parse(savedUser)
        userSettings.name = parsed.name || ''
        userSettings.bio = parsed.bio || ''
      } catch (_) {}
    }
  }
}

onMounted(bootstrap)

watch(
  () => route.params.sessionId,
  async (newId) => {
    input.value = ''
    if (typeof newId === 'string' && newId) {
      await loadSession(newId)
    }
  },
)

watch(
  () => displayMessages.value.length,
  () => {
    scrollToBottom()
  },
)

watch(
  () => sessionView.value?.session?.model_key,
  (val) => {
    selectedModel.value = val || ''
  },
)

const saveUserSettings = () => {
  if (process.client) {
    localStorage.setItem('chat:user-profile', JSON.stringify({ name: userSettings.name || '', bio: userSettings.bio || '' }))
  }
  showUserSettings.value = false
}

watch(
  () => selectedPresetId.value,
  (val) => {
    cachedPresetId.value = val || ''
    if (process.client) {
      if (val) {
        localStorage.setItem('chat:selected-preset-id', val)
      } else {
        localStorage.removeItem('chat:selected-preset-id')
      }
    }
  },
)

watch(
  () => enableStream.value,
  (val) => {
    if (process.client) {
      localStorage.setItem('chat:stream-enabled', val ? '1' : '0')
    }
  },
)

watch(
  () => persistentImagePrompt.value,
  (val) => {
    if (process.client) {
      localStorage.setItem('chat:persistent-image-prompt', val || '')
    }
  },
)

const handleSend = async () => {
  const content = input.value.trim()
  if (!content || isSendingMessage.value || !activeSessionId.value) return

  input.value = ''
  try {
    const profileParts = []
    if (userDisplayName.value) profileParts.push(`<name>${userDisplayName.value}</name>`)
    if (userSettings.bio) profileParts.push(`<bio>${userSettings.bio}</bio>`)
    const userProfileTag = profileParts.length ? `<user_input>\n${profileParts.join('\n')}\n</user_input>\n` : ''
    await sendMessage(activeSessionId.value, `${userProfileTag}${content}`, { preset: selectedPreset.value, stream: enableStream.value })
    scrollToBottom()
  } catch (error) {
    console.error(error)
    input.value = content
    loadingError.value = error instanceof Error ? error.message : ''
  }
}

const handleSelectModel = async (value: string) => {
  if (!sessionView.value) return
  modelChanging.value = true
  try {
    const updated = await updateSettings(sessionView.value.session.id, { model_key: value || undefined })
    if (sessionView.value) {
      sessionView.value = { ...sessionView.value, session: updated }
    }
    selectedModel.value = value
  } catch (error) {
    console.error(error)
  } finally {
    modelChanging.value = false
  }
}

const handleToggleMode = async (mode: 'sfw' | 'nsfw') => {
  if (!sessionView.value || sessionView.value.session.mode === mode) return
  togglingMode.value = true
  try {
    const updated = await updateSettings(sessionView.value.session.id, { mode, sfw_mode: mode === 'sfw' })
    if (sessionView.value) {
      sessionView.value = { ...sessionView.value, session: updated }
    }
  } catch (error) {
    console.error(error)
    loadingError.value = error instanceof Error ? error.message : ''
  } finally {
    togglingMode.value = false
  }
}

const startChatWithRole = async (roleId: string) => {
  if (!roleId || isLoading.value) return
  try {
    const session = await createSession(roleId)
    if (session?.id) {
      router.push(`/chat/${session.id}`)
      activeList.value = 'sessions'
    }
  } catch (error) {
    console.error(error)
    loadingError.value = error instanceof Error ? error.message : ''
  }
}

const startFreshChat = async () => {
  const roleId = sessionView.value?.session?.role_id
  if (!roleId) return
  try {
    const session = await createSession(roleId, sessionView.value?.session?.model_key)
    if (session?.id) {
      router.push(`/chat/${session.id}`)
      activeList.value = 'sessions'
    }
  } catch (error) {
    console.error(error)
    loadingError.value = error instanceof Error ? error.message : ''
  }
}

const triggerPresetImport = () => {
  presetError.value = ''
  fileInput.value?.click()
}

const handlePresetFile = async (event: Event) => {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (!file) return
  try {
    const text = await file.text()
    const data = JSON.parse(text)
    const preset = normalizePreset(data)
    if (!preset.blocks?.length) {
      throw new Error('预设中未找到 blocks')
    }
    // 保存到后端
    const res = await api.post<{ data: any }>('/presets', preset)
    await loadUserPresets()
    selectedPresetId.value = res.data?.id || ''
    cachedPresetId.value = selectedPresetId.value
    if (process.client) {
      localStorage.setItem('chat:selected-preset-id', selectedPresetId.value)
    }
    presetError.value = ''
  } catch (err: any) {
    presetError.value = err?.message || '预设导入失败'
    selectedPresetId.value = ''
  } finally {
    if (fileInput.value) fileInput.value.value = ''
  }
}

const normalizePreset = (data: any) => {
  const presetData = data.preset || data
  const model = presetData.model || {}
  const meta = presetData.meta || {}
  const rawBlocks = presetData.blocks || []
  const blocks = rawBlocks.map((b: any) => ({
    id: b.id,
    name: b.name,
    role: b.role || 'system',
    content: b.content || '',
    enabled: Boolean(b.enabled), // 默认关闭，作者设置的 enabled 保留
    marker: Boolean(b.marker),
  }))
  return {
    name: presetData.name || meta.name,
    description: presetData.description || meta.description,
    model_key: presetData.model_key || model.key || presetData.modelKey,
    gen_params: presetData.gen_params || model.params || presetData.genParams || {},
    blocks,
  }
}

const loadUserPresets = async () => {
  try {
    const res = await api.get<{ data: any[] }>('/presets')
    userPresets.value = res.data || []
    if (process.client) {
      localStorage.setItem('chat:user-presets', JSON.stringify(userPresets.value))
      if (cachedPresetId.value && userPresets.value.some(p => p.id === cachedPresetId.value)) {
        selectedPresetId.value = cachedPresetId.value
      } else if (!userPresets.value.some(p => p.id === selectedPresetId.value)) {
        selectedPresetId.value = ''
      }
    }
  } catch (e) {
    console.error(e)
    if (process.client) {
      const cached = localStorage.getItem('chat:user-presets')
      if (cached) {
        try {
          userPresets.value = JSON.parse(cached)
        } catch (_) {}
      }
      if (cachedPresetId.value && userPresets.value.some(p => p.id === cachedPresetId.value)) {
        selectedPresetId.value = cachedPresetId.value
      }
    }
  }
}

const deleteSelectedPreset = async () => {
  if (!selectedPresetId.value) return
  if (!confirm('确定要删除这个预设吗？')) return
  try {
    await api.del(`/presets/${selectedPresetId.value}`)
    await loadUserPresets()
    selectedPresetId.value = ''
    cachedPresetId.value = ''
    if (process.client) {
      localStorage.removeItem('chat:selected-preset-id')
    }
  } catch (e) {
    alert('删除失败')
  }
}

const openPresetDetails = () => {
  if (!selectedPreset.value) return
  // 深拷贝一份用于编辑，避免直接改列表
  const clone = JSON.parse(JSON.stringify(selectedPreset.value || {}))
  if (Array.isArray(clone.blocks)) {
    clone.blocks = clone.blocks.map((b: any, idx: number) => ({
      ...b,
      name: b?.name || b?.id || `Block ${idx + 1}`,
      id: b?.id || '',
      role: b?.role || 'system',
    }))
  }
  presetDetail.value = clone
  showPresetDetail.value = true
}

const savePresetDetail = async () => {
  if (!presetDetail.value || !presetDetail.value.id) return
  try {
    await api.put(`/presets/${presetDetail.value.id}`, presetDetail.value)
    await loadUserPresets()
    showPresetDetail.value = false
    // 刷新选中项
    if (presetDetail.value.id === selectedPresetId.value) {
      const updated = userPresets.value.find(p => p.id === selectedPresetId.value)
      if (updated) {
        presetDetail.value = JSON.parse(JSON.stringify(updated))
      }
    }
    if (process.client) {
      localStorage.setItem('chat:user-presets', JSON.stringify(userPresets.value))
      if (selectedPresetId.value) {
        localStorage.setItem('chat:selected-preset-id', selectedPresetId.value)
      }
    }
  } catch (e) {
    alert('保存预设失败')
  }
}
</script>

<style scoped>
.custom-scroll::-webkit-scrollbar {
  width: 8px;
}
.custom-scroll::-webkit-scrollbar-thumb {
  background-color: rgba(255, 255, 255, 0.08);
  border-radius: 999px;
}
.custom-scroll::-webkit-scrollbar-track {
  background: transparent;
}

.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 50;
}
.modal-card {
  width: 480px;
  max-height: 70vh;
  background: #0f111a;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 16px;
  padding: 16px;
  color: white;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.modal-scroll {
  overflow-y: auto;
}

:deep(.dialogue-double) {
  color: #9ec5ff;
  font-weight: 800;
}
:deep(.dialogue-single) {
  color: #7c8497;
  font-style: italic;
}
:deep(.msg-code-block) {
  position: relative;
  background: #0b1220;
  border: 1px solid rgba(255, 255, 255, 0.12);
  border-radius: 12px;
  padding: 12px 12px 12px 14px;
  overflow-x: auto;
  font-family: 'JetBrains Mono', 'Fira Code', Menlo, monospace;
  color: #d8dee9;
  line-height: 1.45;
}
:deep(.msg-code-block::before) {
  content: 'CODE';
  position: absolute;
  top: -10px;
  left: 12px;
  background: #111827;
  color: #7dd3fc;
  font-size: 10px;
  font-weight: 700;
  padding: 2px 6px;
  border-radius: 6px;
  border: 1px solid rgba(125, 211, 252, 0.3);
}
</style>
