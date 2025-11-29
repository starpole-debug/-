export interface User {
  id: string
  username: string
  email?: string
  nickname: string
  avatar_url?: string
  avatarUrl?: string
  is_admin?: boolean
  is_banned?: boolean
  created_at?: string
  deleted_at?: string
  follower_count?: number
  following_count?: number
  is_following?: boolean
}

export interface Role {
  id: string
  name: string
  description: string
  avatar_url?: string
  tags?: string[]
  abilities?: string[]
  status?: string
  creator_id?: string
  is_favorited?: boolean
  favorite_count?: number
}

export interface ChatSession {
  id: string
  title: string
  role_id: string
  model_id?: string
  mode?: string
  settings: ChatSessionSettings
  created_at?: string
}

export interface ChatMessage {
  id: string
  session_id: string
  sender_type: 'user' | 'assistant'
  content: string
  created_at?: string
}

export interface ChatSessionSettings {
  temperature: number
  max_tokens: number
  narrative_focus: string
  action_richness: string
  sfw_mode: boolean
  immersive: boolean
}

export interface WorldSummary {
  summary?: string
  scene?: string
  timeline?: string
  npcs?: string[]
}

export interface ChatSessionView {
  label: string
  description?: string
  enabled?: boolean
}

export interface CommunityPost {
  id: string
  title: string
  content: string
  author_id: string
  author_name?: string
  author_avatar?: string
  link_url?: string
  link_type?: string
  style_id?: string
  topic_ids?: string[]
  visibility?: string
  created_at?: string
  attachments?: CommunityAttachment[]
}

export interface CommunityComment {
  id: string
  post_id: string
  author_id: string
  author_name?: string
  author_avatar?: string
  content: string
  visibility?: string
  created_at?: string
}

export interface ModelConfig {
  id: string
  name: string
  description?: string
  provider: string
  base_url: string
  model_name: string
  temperature?: number
  max_tokens?: number
  status: string
  is_default?: boolean
  is_enabled?: boolean
  has_api_key?: boolean
  price_coins?: number
  price_hint?: string
  share_role_pct?: number
  share_preset_pct?: number
}

export interface Notification {
  id: string
  title: string
  content: string
  type: string
  is_read?: boolean
  created_at?: string
}

export interface CommunityAttachment {
  id: string
  post_id: string
  file_url: string
  file_type?: string
  file_name?: string
  file_size?: number
  created_at?: string
}

export interface RevenueEvent {
  id: string
  event_type: string
  amount: number
  status: string
  created_at?: string
}

export interface CreatorWallet {
  available_balance: number
  frozen_balance: number
  total_earned: number
}

export interface UserAsset {
  user_id: string
  balance: number
  monthly_tickets: number
  updated_at?: string
}

export interface UserHomeStats {
  session_total: number
  favorite_total: number
  recent_view_total: number
  post_total: number
}

export interface CreatorHomePrompt {
  roles_total: number
  published_roles: number
  wallet_balance: number
  has_creator_access: boolean
}

export interface UserHomePayload {
  user: User
  assets: UserAsset
  stats: UserHomeStats
  favorites: CommunityPost[]
  recent_views: CommunityPost[]
  my_posts: CommunityPost[]
  creator: CreatorHomePrompt
  payments?: PaymentOrder[]
}

export interface PaymentOrder {
  id: string
  user_id: string
  out_trade_no: string
  provider_trade_no?: string
  pay_type?: string
  status: string
  money_cents: number
  coins: number
  created_at?: string
  updated_at?: string
}
