# 对话系统实施计划（依据 `task` 规划书）

## 0. 目标
- 交付一个覆盖用户端会话体验、创作者调试、后台配置的完整落地路径。
- 先产出用户聊天页 MVP，再逐步扩展创作者/后台能力。

## 1. 用户端会话体验

### 1.1 入口流程
1. 在角色卡列表页注入 CTA → 角色详情页。
2. 详情页展示：封面、简介、标签、创作者、点赞/收藏/使用量、`开始对话` 按钮。
3. 点击 `开始对话` 触发：
   - 若已登录：创建 Session 记录 → 导航到 `/chat/:sessionId`。
   - 未登录：弹窗提示登录或开启临时会话（限制轮数）。

### 1.2 对话页面布局
- 左侧：角色信息 + 会话列表 + 新建/重命名/删除/收藏。
- 中间：聊天主区域，支持消息类型、状态、长回复 streaming、中断按钮。
- 右侧：角色/世界信息、参数调节（温度、token、风格、SFW/NSFW、沉浸模式）。
- 底部输入区：多行输入、Enter 发送、快捷指令、预制按钮、发送/停止。

### 1.3 消息交互细节
| 功能 | 描述 |
| --- | --- |
| 气泡类型 | 用户 / AI / 系统 / 事件 |
| 状态 | 发送中、失败重试、撤回（可选） |
| 操作 | 复制、引用、👍/👎 反馈、标记关键剧情 |
| Streaming | 渐进展示 + 「中断生成」 |

### 1.4 输入栏增强
- 快捷插入：“继续场景 / 新场景 / 描述表情动作”等。
- 预制指令：角色语气调节、总结剧情。

### 1.5 右侧信息面板
- 角色信息：名/称呼/性格标签/背景摘要/创作者备注，可展开更多。
- 世界书摘要：世界观/场景/时间线/NPC，支持刷新（调用 summarization）。
- 参数控制：
  - Temperature 滑条
  - 回复长度快捷项
  - 叙事 vs 对话、动作描述丰富度
  - SFW/NSFW 开关 + 流程（确认、协议）
  - 沉浸模式开关
- 参数变更对当前 Session config 生效，无需刷新。

### 1.6 会话管理
- 列表展示当前角色或全平台会话。
- 操作：新建、重命名、删除、收藏/重要。
- 导出：文本/Markdown/JSON。
- 分享：生成只读链接（后续版本）。

### 1.7 状态提示 & 错误反馈
- 顶部展示：模型名称、SFW/NSFW、当前场景。
- 错误情况：模型超时、额度不足、拒答等提示与 retry 流程。

## 2. 创作者 / 管理端

### 2.1 角色卡调试页
- 基于用户聊天页扩展：增加右侧调试 Tab。
- 额外显示：模型名称、Session 参数、预设、Prompt 结构、Token 估算、世界书片段引用。

### 2.2 模型与预设管理
1. 模型列表：`model_key`、名称、状态、类型(SFW/NSFW)、接口、key、限速。
2. 预设管理：
   - 绑定模型 & 生效范围(全局/角色类型)。
   - system/safety/style prompt。
   - 版本号 & 修改人，支持版本回滚。
3. SFW/NSFW 策略：
   - 整体开关、未成年限制、角色白/黑名单、地区限制。

### 2.3 对话监控
- 查看匿名化样本，按角色/模型/时间搜索。
- 支持禁用角色卡或用户。

## 3. 后端 / 服务端结构

### 3.1 模块划分
| 模块 | 职责 |
| --- | --- |
| Session | CRUD、状态维护（模型、模式、参数） |
| Chat Engine | 接收消息、拼 prompt、调用模型、写消息 |
| Memory Engine | 短/中/长期记忆、摘要 |
| Model Adapter | 不同模型 API 适配、重试、超时 |
| Preset & Config | 平台预设、模型预设 |
| 角色/世界管理 | 读取卡片与世界书、筛选片段 |

### 3.2 数据结构（Postgres）
- `users(id, username, email, is_creator, is_admin, nsfw_allowed, ...)`
- `character_cards(id, creator_id, name, description, tags, cover_image, config_json, ...)`
- `worldbooks(id, character_id, data_json)`
- `model_presets(id, model_key, type, preset_text, version, updated_by)`
- `sessions(id, user_id, character_id, model_key, mode, settings_json, status)`
- `messages(id, session_id, role, content, created_at, is_important)`
- `session_memory(id, session_id, type, summary_text, created_at)`

### 3.3 对话流程
1. 前端发送 session_id + message + 参数变化。
2. ChatController：
   - 权限校验（NSFW、封禁）。
   - 读取 session info、平台预设、模型预设、角色卡、世界书摘要、记忆。
   - 调用 Memory Engine 生成剧情摘要。
   - 组装 prompt（System 规则 / 模式 / 角色 / 世界 / 摘要 / 历史 / 当前消息）。
3. 调用 Model Adapter → 获取回复。
4. 写入 `messages`（user+assistant），触发 `session_memory` 条件 summarization。
5. 响应前端（文本 + metadata）。

### 3.4 记忆系统 MVP
- `messages`: 全量存历史。
- `session_memory`:
  - type = `mtm`，存摘要。
  - 每 30 条消息触发 summarizer → 存摘要 + 旧消息打 `is_archived`。
  - Prompt 时只取最近 10 条 + 所有 `mtm`。

## 4. 分阶段执行路线
1. **阶段 A：用户聊天页 MVP**
   - 完成入口 → Session 创建 → 基本聊天 → 右侧参数（精简版）。
   - 实现后端 Session + Message + Prompt 组装 + 模型调用。
2. **阶段 B：增强体验**
   - Streaming、消息操作、收藏/引用、输入栏快捷指令。
   - 世界书摘要 & 参数控制拓展。
3. **阶段 C：记忆系统 & 会话管理**
   - 引入 session_memory、导出、分享、收藏。
4. **阶段 D：创作者调试页**
   - 调试 Tab、Prompt/Token 可视化。
5. **阶段 E：后台模型/预设**
   - 模型/预设管理、SFW/NSFW 策略。
6. **阶段 F：对话监控 & 安全**
   - 样本审查、角色/用户封禁。

---
以上计划为实施的基准。下一步：依阶段展开任务、逐一交付。 

