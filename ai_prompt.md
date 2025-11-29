你是一个资深全栈工程师，正在维护一个现有项目，技术栈是：

- backend：Go + Gin，代码在 backend/ 目录
- frontend：Nuxt 3 + Vue 3 + TypeScript，代码在 frontend/ 目录

这个项目是一个「AI 头像 / 角色聊天」网站，目前已经有很多半成品能力，但以下几个核心功能没有完全实现，请你在**不破坏现有结构的前提下**，补全并打通它们：

============================
【总体目标】
============================

1. 实现一个真正可用的「用户注册 / 登录 / JWT 鉴权」系统，替换掉现在前端的 mock 登录逻辑。
2. 完善「模型管理（Admin 添加模型）」功能，让管理员可以配置：
   - provider（字符串，例如 "openai" / "deepseek" / "custom"）
   - base_url（模型 API 的根地址）
   - model_name（例如 "gpt-4o-mini"）
   - api_key（安全存储，不在前端明文暴露）
   - 温度、最大 token 等参数（可选）
3. 完成「AI 聊天 + 角色预设」：
   - 用户在前端选择一个角色（avatar/role），打开聊天页面；
   - 后端根据角色预设，把角色信息拼成系统 prompt；
   - 调用 LLM API（通过你实现的 llm.Client）生成回复；
   - 把对话历史保存到数据库，前端展示完整聊天记录。

下面是详细要求。

============================
【一、用户注册 / 登录 / 鉴权】
============================

当前状况（需要你先阅读代码确认）：

- frontend 里有一个 composable：`frontend/composables/useAuth.ts`，目前是 mock 登录（密码写死为 123456，生成 mock token）。
- backend 里有 auth 相关 handler/service/repository，但实际上没有完备的：
  - `POST /api/auth/register`
  - `POST /api/auth/login`
  - `GET  /api/auth/me` 仅部分实现或依赖未完成的 token 逻辑。

请你做以下修改：

1. 在 backend 实现用户表和模型（如果已经有 User 模型，请沿用），字段至少包括：
   - id（主键）
   - username 或 email（作为登录用账号）
   - password_hash（密码 hash）
   - created_at / updated_at
   - is_banned / deleted_at 如果项目里已有这些字段，请保持兼容。

2. 在 `backend/internal/service/auth` 和 `backend/internal/handler/auth` 里实现：
   - `POST /api/auth/register`
     - 请求体：`{ username: string, password: string }`
     - 校验用户名是否已存在
     - 使用 bcrypt（或等价算法）对密码进行 hash
     - 写入数据库
     - 自动登录：注册成功后生成 JWT token，返回用户信息 + token
   - `POST /api/auth/login`
     - 请求体：`{ username: string, password: string }`
     - 根据 username 查用户
     - 校验密码 hash
     - 如果用户被封禁 / 已删除，则返回错误（例如 403 / 400 with message "account disabled"）
     - 登录成功时生成 JWT token，返回用户信息 + token
   - `GET /api/auth/me`
     - 使用 JWT 中间件解析 Authorization: Bearer <token>
     - 把 user_id 写入 context
     - 返回当前用户的基础信息（id + username）

3. JWT 中间件：
   - 如果项目已有类似 `middleware.Authenticator`，请补全它，使其：
     - 从 Authorization header 解析 Bearer token
     - 校验签名、过期时间
     - 在 gin.Context 中设置 userID，后续 handler 可以读取
   - JWT 的密钥从环境变量读取，例如：`AUTH_JWT_SECRET`

4. 前端 `frontend/composables/useAuth.ts`：
   - 替换掉 `mockLoginRequest`，改为真正调用后端：
     - `POST /api/auth/login`
     - `POST /api/auth/register`
     - `GET  /api/auth/me`
   - 使用 cookie 保存 token，例如 cookie 名：`auth-token`
   - 提供以下方法和状态：
     - `login(credentials)`
     - `register(credentials)`
     - `logout()`
     - `fetchProfile()`：如果 cookie 里有 token，则请求 `/auth/me`，恢复用户状态
     - `user`、`isAuthenticated`、`isAuthenticating`、`authError`
   - 保持与现有布局 `layouts/default.vue` 中对 `useAuth()` 的调用兼容（例如 `onMounted(() => auth.fetchProfile())`）。

5. 所有需要登录的接口（例如聊天、发帖等，视项目情况）应使用 JWT 中间件保护。

============================
【二、模型管理（Admin 添加模型）】
============================

当前 Admin 页面可以「添加模型」，但只是简单地在本地列表里新增一条记录，没有配置模型的 API 信息，后端对应的模型配置也不完整。

请你实现一个完整的模型管理功能，后端负责保存配置并为 LLM 调用提供 ModelConfig：

1. 在 backend 的 model / repository / service 层中，扩展 Model 或 AIModel 结构，使其包含至少这些字段：
   - id
   - name（展示名，如 "GPT-4o 小模型"）
   - description
   - provider（字符串，如 "openai" / "deepseek" / "custom"）
   - base_url（API 的根地址，例如 "https://api.openai.com/v1"）
   - model_name（用于请求的模型名，例如 "gpt-4o-mini"）
   - api_key（建议存数据库时加密 / 或存环境变量 key 名）
   - temperature（float，可选）
   - max_tokens（int，可选）
   - status（active / inactive）

2. 在后台（admin）接口中实现：
   - `GET    /api/admin/models`：列出所有模型
   - `POST   /api/admin/models`：创建模型配置
   - `PUT    /api/admin/models/:id`：更新模型配置
   - `DELETE /api/admin/models/:id`：删除或软删除模型配置
   - 尽量使用已有的 admin service / handler 结构，保持风格一致。

3. 在 `internal/pkg/llm` 中：
   - 定义 `ModelConfig` 结构，承接上述字段。
   - 在聊天时，从数据库中加载对应模型的配置（例如默认使用某个 active 模型，或者会话中记录使用哪个模型）。

4. 前端 Admin 模型管理页面：
   - 增加表单字段：
     - provider（下拉框）
     - base_url
     - model_name
     - api_key（type="password"，不要在列表中回显）
     - temperature（可选）
     - max_tokens（可选）
   - 提交时调用后端的 admin models 接口。
   - 列表页显示 name / provider / model_name / status 等，api_key 不要展示明文。

============================
【三、AI 聊天 & 角色预设（Prompt 注入）】
============================

目标：实现一个真正可用的「角色聊天」，根据前端选择的角色，把角色特质 / 设定拼接成系统 prompt，再加上对话历史和用户输入，一起发给 LLM。

当前情况（你需要先阅读现有 chat 相关代码再改）：

- backend 中有：
  - `internal/handler/chat/handler.go`
  - `internal/service/chat/service.go`
  - `internal/repository/chat_repository.go`
  - `internal/pkg/llm/client.go` 和 `mock.go`
- 目前 `llm.Client` 的实现只是 mock：简单 Echo 用户输入，没有真正调用任何外部 AI。
- 角色（role/avatar）相关的模型和接口已经存在（用于前台展示 featured 角色），但没有把角色设定参与到聊天。

请你做：

1. 设计并实现一个真实的 LLM Client：
   - 在 `internal/pkg/llm` 下新增一个 `openai_client.go`（或者 generic client），实现 `Client` 接口：
     - `Generate(ctx context.Context, prompt string, model *ModelConfig, history []model.ChatMessage) (string, error)`
   - 使用 HTTP 调用类似 OpenAI Chat Completions 风格的 API（可按以下约定）：
     - POST `${base_url}/chat/completions`
     - headers:
       - `Authorization: Bearer <api_key>`
       - `Content-Type: application/json`
     - body 包含：
       - `model: model.ModelName`
       - `messages: []{role, content}`，其中第一条为 system（角色预设 + 全局指令），后面是历史对话和当前用户消息。
   - 将当前的 mock 实现保留为开发模式用的一个 provider（例如 provider = "mock" 时使用 mock client）。

2. 在 chat service 中，将角色预设注入到系统 prompt：
   - 当用户对某个角色发起聊天时（创建 chat session 时应该已经能拿到角色 id 或信息），从数据库中读取该角色：
     - 角色名称
     - 角色 persona / description / 特质设定（如果没有单独字段，可以复用 description，并按需要调整格式）
   - 构造一个 system prompt，例如：
     - 「你现在扮演一个虚拟角色：{角色名}。角色设定如下：{角色描述}。与你对话的用户是平台的终端用户，请始终以该角色的口吻、风格来回答问题。」
   - 在调用 `llm.Client.Generate` 时，把这个 system prompt 放在 messages 的第一条 system 消息中。

3. 聊天历史：
   - 保持 `chat_messages` 表现有结构，如果没有，可以包含：
     - id
     - session_id
     - role（"user" / "assistant" / "system"）
     - content
     - created_at
   - 每次用户发送一条消息时：
     - 在数据库中插入一条 role="user" 的记录
     - 调用 LLM 生成回复
     - 把生成的回复保存为一条 role="assistant" 的记录
     - 返回当前会话的消息列表给前端（或者只返回最新两条也可以，只要前端能重构历史）

4. 前端聊天页面：
   - 使用现有聊天 UI（如果已经有），把发送消息逻辑接到：
     - `POST /api/chat/sessions/:id/messages`
   - 发送时携带当前选中的角色 id 或在会话创建时绑定好角色。
   - 展示时区分左右（用户/AI）。

============================
【四、工程与安全要求】
============================

1. 所有新加的配置（如 JWT 密钥、默认模型等）从环境变量读取，并在 README 或注释中说明：
   - `AUTH_JWT_SECRET`
   - `DEFAULT_MODEL_ID`（或其它你需要的变量）
2. 不要在前端暴露任何敏感信息（API Key 等），由后端持有并调用外部 LLM。
3. 尽可能复用当前项目已有的：
   - router 注册方式
   - error/response 格式（看 `internal/pkg/response`）
   - middleware 架构
   - admin service / handler pattern
4. 为关键逻辑（尤其是 auth / chat / llm）添加必要的错误处理和日志，避免 panic。

============================
【最终交付预期】
============================

完成之后，应满足以下行为：

1. 用户可以在前端完成：
   - 注册账号（用户名 + 密码）
   - 登录账号，登录状态持久化在 cookie 中
   - 刷新页面后仍然保持登录（通过 `/auth/me` 恢复）
2. Admin 后台可以：
   - 添加 / 编辑 / 删除 模型配置
   - 至少配置 provider / base_url / model_name / api_key / status
3. 普通用户可以：
   - 选择一个角色 → 进入聊天页
   - 输入内容 → 后端调用真实 LLM，生成有角色风格的回复
   - 聊天记录可持续（刷新页面后还能看到历史）
4. 整个过程里，mock LLM 实现仍然可以作为 fallback，用于本地无网络或调试。

请你在理解现有代码结构后，按上述目标修改后端和前端代码，保持代码风格一致，并确保项目在本地可以正常启动和运行。
