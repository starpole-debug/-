# 项目待办事项清单 (Project Remaining Tasks)

根据对当前代码库 (`d:\studio`) 的分析，以下是按模块划分的未完成功能和建议的开发任务。

## 1. 社区 (Community)

目前状态：已实现基础的信息流、详情页、发帖、评论和点赞。
**待完成功能：**

- [ ] **动态筛选 (Filter Tabs)**
    - [ ] **前端**: 实现 "最新" (Latest)、"热门" (Hot)、"关注" (Following) 标签的点击切换逻辑。
    - [ ] **后端**: 更新 `/community` 接口，支持 `sort` (latest/hot) 和 `filter` (following) 参数。
- [ ] **用户个人主页 (User Profile)**
    - [ ] **前端**: 创建 `/community/user/:id` 页面，展示用户的头像、简介、发布的帖子列表。
    - [ ] **后端**: 提供获取用户公开信息和帖子列表的接口。
- [ ] **搜索功能 (Search)**
    - [ ] **前端**: 在社区顶部添加搜索栏。
    - [ ] **后端**: 实现帖子内容的关键词搜索接口。
- [ ] **通知系统 (Notifications)**
    - [ ] **前端**: 完善通知中心 UI，展示评论、点赞、关注等消息。
    - [ ] **后端**: 完善通知触发逻辑（当有人评论/点赞时写入通知表）及读取接口。

## 2. 创作者生态 (Creator Ecosystem)

目前状态：已实现基础数据看板 (Dashboard) 和角色列表。
**待完成功能：**

- [ ] **角色编辑器 (Role Editor)**
    - [ ] **前端**: 完善 `/creator/editor/[roleId].vue`。目前仅为静态表单，需对接 `/api/roles/:id` 的 GET (回显) 和 PUT (更新) 接口。
    - [ ] **功能**: 支持上传头像、设置背景图、配置详细 Prompt、设置开场白、配置语音合成 (TTS) 参数。
- [ ] **发布流程 (Publishing Flow)**
    - [ ] **前端**: 添加 "保存草稿" 和 "发布上架" 的区分。
    - [ ] **后端**: 确保 `role_versions` 表支持状态管理 (draft/published)。
- [ ] **创作者收益设置 (Monetization Settings)**
    - [ ] **前端**: 允许创作者为角色设置 "付费聊天" 或 "订阅查看" 的价格。
    - [ ] **后端**: 扩展 `roles` 表以存储定价策略。

## 3. 对话聊天 (Chat)

目前状态：已实现会话创建、历史记录和基础消息发送。
**待完成功能：**

- [ ] **高级交互 (Advanced Interaction)**
    - [ ] **流式响应 (Streaming)**: 目前看代码似乎是等待完整响应，建议改为 Server-Sent Events (SSE) 或 WebSocket 以提升体验。
    - [ ] **多模态支持**: 支持发送图片或文件。
- [ ] **角色/模型选择 (Role/Model Selection)**
    - [ ] **前端**: 在创建会话前提供更丰富的角色选择 UI (卡片墙、分类筛选)。
- [ ] **记忆与上下文 (Memory & Context)**
    - [ ] **后端**: 优化 Context 管理，支持长期记忆 (Long-term Memory) 或向量数据库检索 (RAG)。

## 4. 商店与支付 (Store & Payments)

目前状态：仅有基础的打赏 UI (`/store`)，且逻辑较为简单。
**待完成功能：**

- [ ] **支付网关集成 (Payment Gateway)**
    - [ ] **后端**: 对接真实支付渠道 (如 Stripe, WeChat Pay, Alipay)。目前 `/store/tips` 似乎仅记录数据，未处理真实扣款。
- [ ] **商品系统 (Product System)**
    - [ ] **前端**: 展示可购买的虚拟商品（如：会员、积分包、角色解锁卡）。
    - [ ] **后端**: 完善订单系统 (Orders) 和库存/权限管理。

## 5. 系统与通用 (System & General)

- [ ] **管理后台 (Admin Panel)**
    - [ ] 检查 `/admin` 路由下的功能是否完善（用户管理、内容审核、系统配置）。
- [ ] **移动端适配 (Mobile Responsiveness)**
    - [ ] 检查所有新页面的响应式布局，确保在手机端体验良好。
