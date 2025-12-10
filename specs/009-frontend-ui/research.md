# research.md - 009-frontend-ui

## 决策与依据

### 1. 代码高亮库选择
- **Decision**: 采用 Prism.js，按需引入语言包（Go、TypeScript、JavaScript、JSON、bash、markdown），仅在章节内容组件中动态加载样式。
- **Rationale**: Prism 支持 tree-shaking 与按需语言，体积更小；生态成熟，SSR/SSG 下纯前端渲染无副作用；便于与 AntD/Tailwind 并存。
- **Alternatives considered**: highlight.js（全量包更大，按需配置复杂且默认样式覆盖 AntD 需额外工作）；react-syntax-highlighter（体积较大且依赖运行时样式注入，不利于静态导出）。

### 2. JWT 存储与刷新策略
- **Decision**: access token 保存在内存并同步 localStorage（便于刷新后恢复）；refresh token 放在 HttpOnly、SameSite=Lax 的 Cookie，新增 `/api/v1/auth/refresh` 端点用于续期；Axios 请求拦截器自动附加 Authorization，401/过期时静默刷新一次，失败则清空状态并重定向登录。
- **Rationale**: 避免将长期凭证暴露给 JS（Cookie HttpOnly）；access token 内存优先减少 XSS 窃取风险，localStorage 仅用于页面刷新恢复；显式 refresh 端点满足“token 即将过期时自动刷新”需求。
- **Alternatives considered**: 仅用 localStorage 保存 access token（XSS 风险高）；仅用 Cookie（需额外 CSRF 防护且与静态导出分发复杂）；手动重新登录（违背自动刷新要求）。

### 3. 静态导出与动态路由策略
- **Decision**: Next.js `output: 'export'`，通过内置清单生成 `generateStaticParams`（topics: lexical_elements/constants/variables/types；章节列表来自本地 manifest，与后端章节保持同步）；页面骨架静态化，实际章节与测验数据在客户端通过 API 获取并使用 SWR 缓存。
- **Rationale**: 满足“构建输出静态 HTML/CSS/JS”；预生成全部已知路径避免 export 限制；客户端获取可复用既有 API 并保证数据实时性；SWR 提供缓存与重试，契合离线回退需求。
- **Alternatives considered**: 纯客户端动态路由不导出（违背静态构建要求）；SSR/ISR（与 `output: 'export'` 冲突）；仅在构建期从后端抓取章节内容（耦合构建与后端运行，易脆弱）。

### 4. Ant Design + Tailwind 协同
- **Decision**: 以 AntD 组件为主，使用 ConfigProvider 统一主题；Tailwind 仅用于布局/间距/响应式工具类；局部样式用 CSS Modules，避免直接覆盖 AntD 内部类；开启 AntD CSS-in-JS 服务器侧抽取（构建时）确保静态样式落地。
- **Rationale**: 兼顾设计一致性与快速布局；避免 Tailwind 全局样式污染组件；CSS Modules 提供隔离并便于主题扩展。
- **Alternatives considered**: 仅 Tailwind（需大量自定义组件，违背“遵循 AntD”）；直接修改 AntD less 变量（静态导出构建链复杂）；无主题统一（多处重复样式）。

### 5. SQLite 并发与持久化
- **Decision**: 启用 WAL 模式与 busy_timeout，放置数据库文件在 `backend/data/gostudy.db`；对 `users.username`、`learning_progress.user_id/topic/chapter`、`quiz_records.user_id/topic` 建索引（已有定义保持）；初始化时自动迁移并校验 WAL。
- **Rationale**: WAL 提升读并发且适配单实例；busy_timeout 减少锁冲突；索引契合查询路径（按用户/主题/章节、按时间倒序）。
- **Alternatives considered**: DELETE 日志模式（写锁影响并发）；额外引入外部数据库（违背资源约束）。

### 6. API 调用与错误体验
- **Decision**: 统一 Axios 实例，超时 10s，响应拦截处理业务错误码；错误提示统一用 AntD message/notification，网络失败提供重试按钮，401 触发登录重定向。
- **Rationale**: 满足显式错误处理要求；减少重复代码；提升故障可见性与可恢复性。
- **Alternatives considered**: 原生 fetch（重复配置，缺少全局拦截）；逐接口 try-catch（重复、易遗漏统一错误提示）。

