# Implementation Plan: Go-Study2 前端UI界面

**Branch**: `009-frontend-ui` | **Date**: 2025-12-10 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/009-frontend-ui/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

为 Go-Study2 提供现代化 Web 界面，支持用户注册/登录、学习内容展示与进度跟踪、测验记录、响应式布局与一体化部署。

**技术方案**：
- **前端**：Next.js 14（React 18 + TypeScript 5）静态导出，Ant Design 5 + Tailwind CSS 响应式布局，Prism.js 代码高亮，SWR 数据缓存，Axios 统一 API 调用
- **后端**：Go 1.24.5 + GoFrame v2.9.5，新增认证（JWT + bcrypt）、学习进度、测验记录等 API，复用现有 `/api/v1/topics` 等学习内容接口
- **数据库**：SQLite3（WAL 模式），文件路径 `backend/data/gostudy.db`，存储用户、进度、测验记录
- **部署**：前端构建产物 `frontend/out/` 由后端静态文件服务托管，API 与静态资源同端口分发

## Technical Context

**Language/Version**: 
- 前端：TypeScript 5.x + React 18 + Next.js 14
- 后端：Go 1.24.5

**Primary Dependencies**: 
- 前端：Ant Design 5、Tailwind CSS、SWR、Axios、Prism.js（按需语言包）
- 后端：GoFrame v2.9.5、GoFrame ORM (gdb)、JWT (github.com/golang-jwt/jwt/v5)、bcrypt (golang.org/x/crypto/bcrypt)、SQLite driver (github.com/mattn/go-sqlite3)

**Storage**: SQLite3（文件路径 `backend/data/gostudy.db`，启用 WAL 模式，自动迁移）

**Testing**: 
- 前端：Jest + React Testing Library（核心组件覆盖率 ≥80%）
- 后端：Go testing + testify（单元测试覆盖率 ≥80%）

**Target Platform**: 
- 前端：现代桌面与移动浏览器（Chrome/Safari/Firefox/Edge 最新两个主版本）
- 后端：Linux/Windows/macOS 服务器，单实例部署

**Project Type**: Web application（frontend + backend）

**Performance Goals**: 
- API 响应时间 p95 < 200ms（不含网络延迟）
- 前端首屏加载 < 2s（3G 网络）
- 静态资源 gzip 压缩后 < 500KB（初始加载）

**Constraints**: 
- 单实例部署，内存占用 < 200MB（后端 + SQLite）
- 前端静态导出，支持离线缓存（Service Worker 可选）
- 数据库文件 < 100MB（预估 1000 用户 + 10000 条进度记录）

**Scale/Scope**: 
- 用户规模：100-1000 并发用户
- 学习内容：4 个主题（lexical_elements/constants/variables/types），约 50 个章节
- 前端页面：约 15 个页面/组件（登录/注册/主题列表/章节详情/测验/历史记录等）
- 后端 API：约 20 个端点（认证 5 个 + 学习内容 4 个已有 + 进度 3 个 + 测验 5 个 + 用户 3 个）

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### 后端原则检查

- **Principle I (代码质量与可维护性):** ✅ 方案采用分层架构（handler/middleware/models），每层职责清晰；函数保持简短（<50 行），易于理解与测试
- **Principle II (显式错误处理):** ✅ 所有 API 返回统一错误码与消息；数据库操作、JWT 验证、密码验证均有明确错误处理，无静默失败
- **Principle III/XXI (全面测试):** ✅ 规划单元测试（handler/models/middleware）与集成测试（auth/learning/quiz flow），目标覆盖率≥80%；每个新增包含 `*_test.go`
- **Principle IV (单一职责):** ✅ handler 仅处理 HTTP 请求/响应，models 仅负责数据库操作，middleware 仅负责请求拦截；JWT/validator 独立为 pkg
- **Principle V/XV (一致文档与中文要求):** ✅ 后端所有注释与错误消息使用中文；每个新增包规划 README 说明用途
- **Principle VI (YAGNI):** ✅ 拒绝 Repository 模式（直接使用 GoFrame ORM）；拒绝复杂权限系统（仅用户认证）；拒绝 Redis（SQLite 足够）
- **Principle VII (安全优先):** ✅ 密码 bcrypt 哈希（cost=10）；JWT secret 从环境变量读取；refresh token HttpOnly Cookie；输入校验（用户名/密码格式）；SQL 注入防护（使用 ORM 参数化查询）
- **Principle VIII/XVIII (可预测结构):** ✅ 遵循标准 Go 布局，仅根目录 `main.go`；新增代码在 `internal/app/` 下，保持现有结构
- **Principle IX (依赖纪律):** ✅ 新增依赖最小：SQLite driver（必需）、JWT 库（必需）、bcrypt（必需）；复用 GoFrame ORM 而非引入新 ORM
- **Principle X (性能优化):** ✅ SQLite 启用 WAL 模式提升并发；关键字段建索引（username/user_id+topic/created_at）；API 响应目标 p95<200ms
- **Principle XI (文档同步):** ✅ 方案包含完成后更新根 README（新增"用户认证"、"学习进度跟踪"、"测验功能"章节）
- **Principle XIV (清晰分层注释):** ✅ 规划 handler/middleware/models 各层中文注释说明职责
- **Principle XVI (浅层逻辑):** ✅ 采用卫语句（early return）避免深层嵌套；复杂逻辑拆分为独立函数（如 `validateUser`、`hashPassword`）
- **Principle XVII (一致开发者体验):** ✅ 数据库自动迁移（无需手动建表）；配置文件清晰（config.yaml）；启动流程简单（构建前端 -> 编译后端 -> 运行）
- **Principle XIX (包级 README):** ✅ 规划 `internal/app/models/README.md`（数据库模型说明）、`internal/pkg/jwt/README.md`（JWT 工具用法）
- **Principle XX (代码质量执行):** ✅ 提交前执行 `go fmt`、`go vet`、`golint`、`go mod tidy`；CI 流程包含质量检查
- **Principle XXII (分层菜单导航):** ✅ 前端使用 AntD Menu 组件支持多级导航；CLI 模式保持现有菜单结构不变
- **Principle XXIII (双学习模式):** ✅ 复用现有 `/api/v1/topics` 等接口，CLI 与 HTTP 共享内容源；新增 API 不影响 CLI 模式
- **Principle XXIV (层次化章节结构):** ✅ 保持现有章节组织（constants/lexical_elements/variables/types），不改变目录结构
- **Principle XXV (HTTP/CLI 一致性):** ✅ 新增 API 遵循现有路由规范（`/api/v1/topic/{topic}/{chapter}`）；响应格式统一（`{code, message, data}`）；错误处理显式（404/40001/50001）

### 前端原则检查

- **Principle XXVI (类型安全优先):** ✅ 所有组件使用 TypeScript；Props/State/API 响应均有类型定义（`types/` 目录）
- **Principle XXVII (静态导出优化):** ✅ Next.js `output: 'export'`；使用 `generateStaticParams` 预生成路由；动态导入重组件（`next/dynamic`）
- **Principle XXVIII (一致 UI/UX):** ✅ 严格使用 Ant Design 组件；通过 ConfigProvider 统一主题；自定义样式基于 AntD 设计 token
- **Principle XXIX (无障碍标准):** ✅ 使用 AntD 语义化组件（自带 ARIA 属性）；表单字段提供 label；错误提示关联表单项
- **Principle XXX (客户端安全):** ✅ access token 仅内存+localStorage（刷新恢复）；refresh token HttpOnly Cookie；敏感信息不存客户端；Axios 拦截器统一错误处理
- **Principle XXXI (组件组织):** ✅ 使用函数组件+Hooks；组件按功能分目录（Auth/Learning/Quiz/Common）；每个组件独立文件
- **Principle XXXII (状态管理):** ✅ 本地状态用 `useState`；跨组件状态用 AuthContext；避免全局状态滥用；Context 按功能域划分
- **Principle XXXIII (API 集成):** ✅ 统一 Axios 实例（`lib/api.ts`）；请求/响应拦截器处理 token 与错误；错误提示用 AntD message
- **Principle XXXIV (样式标准):** ✅ 优先 CSS Modules 隔离样式；全局样式仅 reset+theme；响应式用 AntD Grid；通过 ConfigProvider 配置主题
- **Principle XXXV (静态导出优化):** ✅ 避免 `getServerSideProps`；按需导入 Prism.js 语言包；动态导入重组件；`next/image` 配置 `unoptimized: true`
- **Principle XXXVI (前端测试):** ✅ 核心组件单元测试（LoginForm/ChapterContent/QuizQuestion）；API 层集成测试；目标覆盖率>80%；使用 `data-testid` 定位元素
- **Principle XXXIX (代码质量):** ✅ ESLint+Prettier 统一格式；TypeScript 严格模式；Husky+lint-staged 提交前检查
- **Principle XL (前端文档):** ✅ 规划 `frontend/README.md`（安装/运行/构建指南）；复杂组件添加 JSDoc 注释

### 前后端集成原则检查

- **Principle XLI (API 契约一致):** ✅ 后端响应格式 `{code, message, data}` 与前端 Axios 拦截器匹配；OpenAPI 规范已定义（`contracts/openapi.yaml`）
- **Principle XLII (开发环境分离):** ✅ 后端 `:8080`，前端开发服务器 `:3000`；前端通过 proxy 或 CORS 访问后端
- **Principle XLIII (生产部署集成):** ✅ 前端构建产物 `frontend/out/` 由后端静态文件服务托管；API 路由 `/api/*` 优先于静态文件；SPA 回退到 `index.html`
- **Principle XLIV (共享配置管理):** ✅ API base URL 通过环境变量配置（`NEXT_PUBLIC_API_URL`）；开发/生产环境配置分离（`.env.local` / `.env.production`）

**结论**：✅ 所有宪章原则检查通过，无违规项，可进入实施阶段。

## Project Structure

### Documentation (this feature)

```text
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
backend/
├── configs/
│   ├── config.yaml                    # 新增: auth/database/static 配置段
│   └── certs/                         # HTTPS 证书 (已有)
├── data/                              # 新增: SQLite 数据库文件目录
│   └── gostudy.db                     # SQLite 数据库文件
├── internal/
│   ├── app/
│   │   ├── http_server/
│   │   │   ├── handler/
│   │   │   │   ├── auth.go           # 新增: 认证相关 handler
│   │   │   │   ├── progress.go       # 新增: 学习进度 handler
│   │   │   │   ├── quiz.go           # 新增: 测验相关 handler
│   │   │   │   └── ... (已有 topics/lexical/constants/variables/types)
│   │   │   ├── middleware/
│   │   │   │   ├── auth.go           # 新增: JWT 认证中间件
│   │   │   │   └── ... (已有 logger/format)
│   │   │   ├── router.go             # 更新: 新增认证/进度/测验路由
│   │   │   └── server.go             # 更新: 新增静态文件服务
│   │   └── ... (已有 constants/lexical_elements)
│   ├── config/
│   │   └── config.go                 # 更新: 新增 Auth/Database/Static 配置结构
│   ├── domain/                       # 新增: 领域层
│   │   ├── user/
│   │   │   ├── entity.go             # 用户实体
│   │   │   ├── repository.go         # 用户仓储接口
│   │   │   └── service.go            # 用户服务
│   │   ├── progress/
│   │   │   ├── entity.go             # 学习进度实体
│   │   │   ├── repository.go         # 进度仓储接口
│   │   │   └── service.go            # 进度服务
│   │   └── quiz/
│   │       ├── entity.go             # 测验记录实体
│   │       ├── repository.go         # 测验仓储接口
│   │       └── service.go            # 测验服务
│   ├── infrastructure/               # 新增: 基础设施层
│   │   ├── database/
│   │   │   ├── sqlite.go             # SQLite 连接与初始化
│   │   │   └── migrations.go         # 数据库迁移 (建表 SQL)
│   │   └── repository/               # 仓储实现
│   │       ├── user_repo.go          # 用户仓储实现
│   │       ├── progress_repo.go      # 进度仓储实现
│   │       └── quiz_repo.go          # 测验仓储实现
│   └── pkg/                          # 新增: 共享工具包
│       ├── jwt/
│       │   ├── jwt.go                # JWT 生成与验证
│       │   └── jwt_test.go
│       └── password/
│           ├── password.go           # bcrypt 哈希与验证
│           └── password_test.go
├── tests/
│   ├── contract/                     # 契约测试
│   │   └── auth_api_test.go         # 新增: 认证 API 契约测试
│   ├── integration/                  # 集成测试
│   │   ├── auth_flow_test.go        # 新增: 认证流程测试
│   │   ├── progress_test.go         # 新增: 进度记录测试
│   │   └── quiz_test.go             # 新增: 测验提交测试
│   └── unit/                         # 单元测试
│       ├── jwt_test.go               # 新增: JWT 工具测试
│       ├── password_test.go          # 新增: 密码工具测试
│       └── ... (已有)
├── go.mod                            # 更新: 新增 sqlite3/jwt/crypto 依赖
├── go.sum                            # 自动更新
└── main.go                           # 更新: 初始化数据库连接

frontend/
├── app/                              # Next.js App Router
│   ├── (auth)/                       # 认证路由组
│   │   ├── login/
│   │   │   └── page.tsx              # 登录页
│   │   └── register/
│   │       └── page.tsx              # 注册页
│   ├── (protected)/                  # 受保护路由组
│   │   ├── layout.tsx                # 认证检查 Layout
│   │   ├── topics/
│   │   │   ├── page.tsx              # 主题列表页
│   │   │   └── [topic]/
│   │   │       ├── page.tsx          # 主题详情/章节列表
│   │   │       └── [chapter]/
│   │   │           └── page.tsx      # 章节内容页
│   │   ├── progress/
│   │   │   └── page.tsx              # 学习进度页
│   │   ├── quiz/
│   │   │   ├── [topic]/
│   │   │   │   └── page.tsx          # 测验作答页
│   │   │   └── history/
│   │   │       └── page.tsx          # 测验历史页
│   │   └── profile/
│   │       └── page.tsx              # 用户资料页
│   ├── layout.tsx                    # 根 Layout
│   └── page.tsx                      # 首页 (重定向到 /topics)
├── components/                       # 共享组件
│   ├── auth/
│   │   ├── LoginForm.tsx             # 登录表单
│   │   ├── RegisterForm.tsx          # 注册表单
│   │   └── AuthGuard.tsx             # 认证守卫
│   ├── layout/
│   │   ├── Header.tsx                # 页头
│   │   ├── Footer.tsx                # 页脚
│   │   └── Sidebar.tsx               # 侧边栏
│   ├── learning/
│   │   ├── TopicCard.tsx             # 主题卡片
│   │   ├── ChapterList.tsx           # 章节列表
│   │   ├── ChapterContent.tsx        # 章节内容 (含代码高亮)
│   │   └── ProgressBar.tsx           # 进度条
│   ├── quiz/
│   │   ├── QuizItem.tsx              # 测验题目
│   │   ├── QuizResult.tsx            # 测验结果
│   │   └── QuizHistory.tsx           # 测验历史列表
│   └── common/
│       ├── ErrorBoundary.tsx         # 错误边界
│       ├── Loading.tsx               # 加载状态
│       └── ErrorMessage.tsx          # 错误提示
├── lib/                              # 工具库
│   ├── api.ts                        # Axios 实例与拦截器
│   ├── auth.ts                       # 认证工具 (token 管理)
│   └── constants.ts                  # 常量定义
├── hooks/                            # 自定义 Hooks
│   ├── useAuth.ts                    # 认证状态 Hook
│   ├── useProgress.ts                # 进度数据 Hook
│   └── useQuiz.ts                    # 测验数据 Hook
├── types/                            # TypeScript 类型定义
│   ├── api.ts                        # API 响应类型
│   ├── auth.ts                       # 认证相关类型
│   ├── learning.ts                   # 学习相关类型
│   └── quiz.ts                       # 测验相关类型
├── styles/                           # 全局样式
│   └── globals.css                   # 全局 CSS (含 Tailwind)
├── public/                           # 静态资源
│   ├── favicon.ico
│   └── images/
├── tests/                            # 前端测试
│   ├── components/
│   │   ├── LoginForm.test.tsx
│   │   ├── ChapterContent.test.tsx
│   │   └── QuizItem.test.tsx
│   └── lib/
│       ├── api.test.ts
│       └── auth.test.ts
├── next.config.js                    # Next.js 配置 (output: 'export')
├── tailwind.config.js                # Tailwind 配置
├── tsconfig.json                     # TypeScript 配置
├── package.json                      # 依赖管理
└── README.md                         # 前端开发指南

tests/
├── integration/
│   ├── auth_flow_test.go                # 新增：认证流程集成测试
│   ├── learning_flow_test.go            # 新增：学习流程集成测试
│   └── quiz_flow_test.go                # 新增：测验流程集成测试
└── contract/
    └── api_contract_test.go             # 新增：API 契约测试（验证响应格式）
```

**Structure Decision**: 
采用 Web application 结构（frontend + backend 分离开发，生产合并部署）。关键设计决策：

1. **后端分层架构**：
   - **Handler 层** (`http_server/handler/`): 处理 HTTP 请求/响应
   - **Domain 层** (`domain/`): 业务逻辑与实体定义（user/progress/quiz）
   - **Infrastructure 层** (`infrastructure/`): 数据库连接与仓储实现
   - **Pkg 层** (`pkg/`): 共享工具（jwt/password）

2. **前端路由组织**：
   - 使用 Next.js 14 App Router 路由组 `(auth)` 和 `(protected)`
   - 认证路由（login/register）与受保护路由（topics/quiz/profile）分离
   - 受保护路由组使用统一 Layout 进行认证检查

3. **前端独立目录**：`frontend/` 完全独立，构建产物 `out/` 由后端静态文件服务托管

4. **测试分层**：前端使用 Jest，后端使用 Go testing，集成测试覆盖完整用户流程

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

无宪章违规项需要特别说明。本方案遵循 YAGNI 原则，仅引入必要的依赖与模式。

---

## 依赖管理

### 后端新增依赖

需要在 `backend/go.mod` 中添加以下依赖：

```bash
cd backend
go get github.com/mattn/go-sqlite3@latest
go get github.com/golang-jwt/jwt/v5@latest
go get golang.org/x/crypto/bcrypt@latest
go mod tidy
```

**依赖说明**：
- `github.com/mattn/go-sqlite3`: SQLite3 驱动（CGO 依赖，需要 GCC）
- `github.com/golang-jwt/jwt/v5`: JWT 生成与验证
- `golang.org/x/crypto/bcrypt`: 密码哈希

**注意事项**：
- SQLite driver 需要 CGO，Windows 环境需安装 MinGW-w64 或 TDM-GCC
- 交叉编译时需指定 `CGO_ENABLED=1` 和目标平台的 C 编译器

### 前端新增依赖

需要在 `frontend/package.json` 中添加以下依赖：

```bash
cd frontend
npm install --save antd@^5.0.0
npm install --save axios@^1.6.0
npm install --save swr@^2.2.0
npm install --save prismjs@^1.29.0
npm install --save-dev @types/prismjs@^1.26.0
npm install --save-dev jest@^29.0.0
npm install --save-dev @testing-library/react@^14.0.0
npm install --save-dev @testing-library/jest-dom@^6.0.0
```

**依赖说明**：
- `antd`: Ant Design UI 组件库
- `axios`: HTTP 客户端
- `swr`: 数据获取与缓存
- `prismjs`: 代码语法高亮
- `jest` + `@testing-library/react`: 测试框架

---

## Phase 0: 数据库设计与初始化

### 数据库选型与配置

**决策**：使用 SQLite3，启用 WAL（Write-Ahead Logging）模式

**配置参数**：
```yaml
# backend/configs/config.yaml 新增配置
database:
  type: sqlite3
  path: backend/data/gostudy.db
  maxOpenConns: 10
  maxIdleConns: 5
  connMaxLifetime: 3600  # 秒
  pragmas:
    - journal_mode=WAL
    - busy_timeout=5000
    - synchronous=NORMAL
    - cache_size=-64000    # 64MB
    - foreign_keys=ON

jwt:
  secret: ${JWT_SECRET}    # 环境变量，生产环境必须设置
  accessTokenExpiry: 604800  # 7天（秒）
  refreshTokenExpiry: 604800 # 7天（秒）
  issuer: go-study2
```

### 数据库表结构设计

#### 1. users 表（用户）

```sql
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_users_username ON users(username);
```

**字段说明**：
- `id`: 自增主键
- `username`: 用户名，唯一索引，长度 3-50，仅字母数字下划线
- `password_hash`: bcrypt 哈希后的密码，cost=10
- `created_at`: 创建时间
- `updated_at`: 更新时间

**验证规则**：
- 用户名正则：`^[A-Za-z0-9_]{3,50}$`
- 密码正则：`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d).{8,}$`（至少 8 位，包含大小写字母与数字）

#### 2. learning_progress 表（学习进度）

```sql
CREATE TABLE IF NOT EXISTS learning_progress (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    topic TEXT NOT NULL,
    chapter TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'not_started',
    last_visit DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_position TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, topic, chapter)
);

CREATE INDEX idx_progress_user_topic ON learning_progress(user_id, topic);
CREATE INDEX idx_progress_last_visit ON learning_progress(user_id, last_visit DESC);
```

**字段说明**：
- `id`: 自增主键
- `user_id`: 用户 ID，外键关联 users.id
- `topic`: 主题标识（lexical_elements/constants/variables/types）
- `chapter`: 章节标识（如 "boolean", "integer" 等）
- `status`: 学习状态（not_started/in_progress/done）
- `last_visit`: 最近访问时间
- `last_position`: 最近位置（JSON 字符串，如 `{"scroll": 1200}` 或 `{"anchor": "section-3"}`）
- `created_at`: 创建时间
- `updated_at`: 更新时间

**唯一约束**：`(user_id, topic, chapter)` 组合唯一，更新为幂等覆盖

**索引策略**：
- `idx_progress_user_topic`: 支持按用户+主题查询进度列表
- `idx_progress_last_visit`: 支持"继续上次学习"功能（按最近访问时间倒序）

#### 3. quiz_records 表（测验记录）

```sql
CREATE TABLE IF NOT EXISTS quiz_records (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    topic TEXT NOT NULL,
    chapter TEXT,
    score INTEGER NOT NULL,
    total INTEGER NOT NULL,
    duration_ms INTEGER NOT NULL,
    answers TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_quiz_user_topic ON quiz_records(user_id, topic);
CREATE INDEX idx_quiz_created_at ON quiz_records(user_id, created_at DESC);
```

**字段说明**：
- `id`: 自增主键
- `user_id`: 用户 ID，外键关联 users.id
- `topic`: 主题标识
- `chapter`: 章节标识（nullable，综合测验时为空）
- `score`: 得分
- `total`: 总分
- `duration_ms`: 答题耗时（毫秒）
- `answers`: 答题详情（JSON 字符串，格式见下方）
- `created_at`: 提交时间

**answers 字段 JSON 格式**：
```json
[
  {
    "id": "q1",
    "choices": ["A"],
    "correct": true,
    "correctAnswer": ["A"]
  },
  {
    "id": "q2",
    "choices": ["B", "C"],
    "correct": false,
    "correctAnswer": ["A", "C"]
  }
]
```

**索引策略**：
- `idx_quiz_user_topic`: 支持按用户+主题筛选历史记录
- `idx_quiz_created_at`: 支持按时间倒序查询历史记录

#### 4. refresh_tokens 表（刷新令牌）

```sql
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    token_hash TEXT NOT NULL UNIQUE,
    expires_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_refresh_tokens_user ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_expires ON refresh_tokens(expires_at);
```

**字段说明**：
- `id`: 自增主键
- `user_id`: 用户 ID
- `token_hash`: refresh token 的 SHA256 哈希（避免明文存储）
- `expires_at`: 过期时间
- `created_at`: 创建时间

**清理策略**：定时任务（或登录/刷新时触发）删除 `expires_at < NOW()` 的记录

### 数据库初始化流程

**文件**：`backend/internal/infrastructure/database/sqlite.go`

**职责**：
1. 读取配置文件中的数据库路径与参数
2. 初始化 GoFrame gdb 连接（使用 SQLite driver）
3. 执行表结构迁移（CREATE TABLE IF NOT EXISTS）
4. 设置 PRAGMA 参数（WAL、busy_timeout 等）
5. 验证连接可用性（Ping）

**代码框架**：
```go
// backend/internal/infrastructure/database/sqlite.go
package database

import (
    "context"
    "database/sql"
    "fmt"
    "os"
    "path/filepath"

    "github.com/gogf/gf/v2/database/gdb"
    "github.com/gogf/gf/v2/frame/g"
    _ "github.com/mattn/go-sqlite3"
)

// InitDB 初始化数据库连接与表结构
func InitDB(ctx context.Context, dbPath string) error {
    // 1. 确保数据目录存在
    dir := filepath.Dir(dbPath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return fmt.Errorf("创建数据目录失败: %w", err)
    }

    // 2. 配置 GoFrame gdb
    gdb.SetConfig(gdb.Config{
        "default": gdb.ConfigGroup{
            gdb.ConfigNode{
                Type: "sqlite",
                Link: dbPath,
            },
        },
    })

    // 3. 获取数据库实例
    db := g.DB()

    // 4. 设置 PRAGMA
    pragmas := []string{
        "PRAGMA journal_mode=WAL;",
        "PRAGMA busy_timeout=5000;",
        "PRAGMA synchronous=NORMAL;",
        "PRAGMA cache_size=-64000;",
        "PRAGMA foreign_keys=ON;",
    }
    for _, pragma := range pragmas {
        if _, err := db.Exec(ctx, pragma); err != nil {
            return fmt.Errorf("设置 PRAGMA 失败: %w", err)
        }
    }

    // 5. 执行表结构迁移
    if err := Migrate(ctx, db); err != nil {
        return fmt.Errorf("数据库迁移失败: %w", err)
    }

    // 6. 验证连接
    if err := db.PingMaster(); err != nil {
        return fmt.Errorf("数据库连接验证失败: %w", err)
    }

    return nil
}
```

```go
// backend/internal/infrastructure/database/migrations.go
package database

import (
    "context"
    "github.com/gogf/gf/v2/database/gdb"
)

// Migrate 执行数据库迁移
func Migrate(ctx context.Context, db gdb.DB) error {
    // 执行 CREATE TABLE IF NOT EXISTS 语句
    // （具体 SQL 见上方表结构设计）
    migrations := []string{
        createUsersTableSQL,
        createLearningProgressTableSQL,
        createQuizRecordsTableSQL,
        createRefreshTokensTableSQL,
    }

    for _, sql := range migrations {
        if _, err := db.Exec(ctx, sql); err != nil {
            return err
        }
    }

    return nil
}
```

**架构说明**：
- **领域层** (`internal/domain/`): 定义实体 (entity.go)、仓储接口 (repository.go)、业务服务 (service.go)
- **基础设施层** (`internal/infrastructure/`): 实现数据库连接 (database/)、仓储实现 (repository/)
- **应用层** (`internal/app/http_server/handler/`): HTTP handler 调用领域服务，不直接操作数据库

**调用位置**：`backend/main.go` 的 `main()` 函数中，在启动 HTTP 服务器之前调用 `database.InitDB()`

---

## Phase 1: 后端 API 实现

### 1. 认证相关 API

**文件**：`backend/internal/app/http_server/handler/auth.go`

#### 1.1 POST /api/v1/auth/register（用户注册）

**请求体**：
```json
{
  "username": "testuser",
  "password": "Test1234"
}
```

**响应**：
```json
{
  "code": 20000,
  "message": "注册成功",
  "data": {
    "accessToken": "eyJhbGc...",
    "expiresIn": 604800
  }
}
```

**错误码**：
- `40004`: 参数错误（用户名/密码格式不符）
- `40009`: 用户名已存在
- `50001`: 服务器错误

**实现要点**（分层架构）：
1. **Handler 层** (`handler/auth.go`): 
   - 解析请求参数
   - 调用 `user.Service.Register()` 进行注册
   - 设置 HttpOnly Cookie（refresh token）
   - 返回 access token
2. **Service 层** (`domain/user/service.go`):
   - 使用 `validator` 包验证用户名与密码格式
   - 调用 `user.Repository.FindByUsername()` 检查唯一性
   - 使用 `password.Hash()` 哈希密码（bcrypt cost=10）
   - 调用 `user.Repository.Create()` 创建用户
   - 调用 `jwt.GenerateTokenPair()` 生成 token
3. **Repository 层** (`infrastructure/repository/user_repo.go`):
   - 实现 `FindByUsername()`: 查询 users 表
   - 实现 `Create()`: 插入 users 表
   - 实现 `SaveRefreshToken()`: 插入 refresh_tokens 表（SHA256 哈希）

#### 1.2 POST /api/v1/auth/login（用户登录）

**请求体**：
```json
{
  "username": "testuser",
  "password": "Test1234"
}
```

**响应**：同注册接口

**错误码**：
- `40001`: 用户名或密码错误
- `40004`: 参数错误

**实现要点**（分层架构）：
1. **Handler 层**: 调用 `user.Service.Login()`，设置 Cookie，返回 token
2. **Service 层**: 
   - 调用 `user.Repository.FindByUsername()` 获取用户
   - 使用 `password.Verify()` 验证密码（bcrypt）
   - 调用 `jwt.GenerateTokenPair()` 生成 token
   - 调用 `user.Repository.SaveRefreshToken()` 保存 refresh token
3. **Repository 层**: 实现数据库查询与更新

#### 1.3 POST /api/v1/auth/logout（退出登录）

**请求头**：`Authorization: Bearer <access_token>`

**响应**：
```json
{
  "code": 20000,
  "message": "退出成功",
  "data": null
}
```

**实现要点**：
1. 从 JWT 中提取 `user_id`
2. 删除 `refresh_tokens` 表中对应记录
3. 清除 Cookie（设置 `Max-Age=-1`）

#### 1.4 POST /api/v1/auth/refresh（刷新 access token）

**请求**：从 Cookie 中读取 `refresh_token`

**响应**：
```json
{
  "code": 20000,
  "message": "刷新成功",
  "data": {
    "accessToken": "eyJhbGc...",
    "expiresIn": 604800
  }
}
```

**错误码**：
- `40002`: refresh token 过期或无效
- `40001`: 认证失败

**实现要点**：
1. 从 Cookie 中获取 `refresh_token`
2. 计算 SHA256 哈希
3. 查询 `refresh_tokens` 表验证有效性与过期时间
4. 生成新的 access token
5. 可选：生成新的 refresh token 并更新数据库（rotation 策略）

#### 1.5 GET /api/v1/auth/profile（获取用户信息）

**请求头**：`Authorization: Bearer <access_token>`

**响应**：
```json
{
  "code": 20000,
  "message": "success",
  "data": {
    "id": 1,
    "username": "testuser"
  }
}
```

**实现要点**：
1. 使用认证中间件验证 JWT
2. 从 JWT 中提取 `user_id`
3. 查询 `users` 表返回用户信息（不包含 `password_hash`）

### 2. 学习进度 API

**文件**：`backend/internal/app/http_server/handler/progress.go`

#### 2.1 GET /api/v1/progress（获取所有进度）

**请求头**：`Authorization: Bearer <access_token>`

**响应**：
```json
{
  "code": 20000,
  "message": "success",
  "data": [
    {
      "topic": "constants",
      "chapter": "boolean",
      "status": "done",
      "lastVisit": "2025-12-10T10:30:00Z",
      "position": "{\"scroll\":1200}"
    }
  ]
}
```

#### 2.2 GET /api/v1/progress/{topic}（获取指定主题进度）

**路径参数**：`topic` (lexical_elements/constants/variables/types)

**响应**：同上，仅返回指定主题的进度记录

#### 2.3 POST /api/v1/progress（记录学习进度）

**请求体**：
```json
{
  "topic": "constants",
  "chapter": "boolean",
  "status": "in_progress",
  "position": "{\"scroll\":1200}"
}
```

**响应**：
```json
{
  "code": 20000,
  "message": "进度已保存",
  "data": null
}
```

**实现要点**：
1. 使用 `INSERT OR REPLACE` 或 `ON CONFLICT UPDATE` 实现幂等更新
2. 自动更新 `last_visit` 为当前时间
3. 验证 `topic` 与 `status` 枚举值

### 3. 测验 API

**文件**：`backend/internal/app/http_server/handler/quiz.go`

#### 3.1 GET /api/v1/quiz/{topic}/{chapter}（获取测验题目）

**路径参数**：
- `topic`: 主题标识
- `chapter`: 章节标识

**响应**：
```json
{
  "code": 20000,
  "message": "success",
  "data": [
    {
      "id": "q1",
      "stem": "以下哪个是 Go 的关键字？",
      "options": [
        {"id": "A", "label": "const"},
        {"id": "B", "label": "final"},
        {"id": "C", "label": "let"},
        {"id": "D", "label": "var"}
      ],
      "multi": true,
      "answer": ["A", "D"],
      "explanation": "Go 的关键字包括 const 和 var"
    }
  ]
}
```

**实现要点**：
1. 复用现有学习内容模块的 `GetQuiz()` 函数
2. 如果章节不存在测验，返回空数组（不报错）
3. 前端需要时可支持随机抽题（后续扩展）

#### 3.2 POST /api/v1/quiz/submit（提交测验）

**请求体**：
```json
{
  "topic": "constants",
  "chapter": "boolean",
  "answers": [
    {"id": "q1", "choices": ["A", "D"]},
    {"id": "q2", "choices": ["B"]}
  ]
}
```

**响应**：
```json
{
  "code": 20000,
  "message": "提交成功",
  "data": {
    "score": 8,
    "total": 10,
    "correctIds": ["q1"],
    "wrongIds": ["q2"],
    "submittedAt": "2025-12-10T11:00:00Z",
    "durationMs": 120000
  }
}
```

**实现要点**：
1. 获取标准答案（调用学习内容模块）
2. 逐题对比用户答案与标准答案（单选/多选均需完全匹配）
3. 计算得分（每题 1 分，全对计分）
4. 构造 `answers` JSON 字符串
5. 插入 `quiz_records` 表
6. 返回评分结果

#### 3.3 GET /api/v1/quiz/history（获取测验历史）

**请求头**：`Authorization: Bearer <access_token>`

**查询参数**：
- `from`: 起始时间（可选，ISO 8601 格式）
- `to`: 结束时间（可选）

**响应**：
```json
{
  "code": 20000,
  "message": "success",
  "data": [
    {
      "id": 1,
      "topic": "constants",
      "chapter": "boolean",
      "score": 8,
      "total": 10,
      "durationMs": 120000,
      "createdAt": "2025-12-10T11:00:00Z"
    }
  ]
}
```

#### 3.4 GET /api/v1/quiz/history/{topic}（按主题获取历史）

**路径参数**：`topic`

**查询参数**：同上

**响应**：同上，仅返回指定主题的记录

### 4. 中间件实现

**文件**：`backend/internal/app/http_server/middleware/auth.go`

#### JWT 认证中间件

**职责**：
1. 从请求头 `Authorization: Bearer <token>` 中提取 JWT
2. 验证 JWT 签名与过期时间
3. 提取 `user_id` 并注入到请求上下文（`r.SetCtxVar("user_id", userID)`）
4. 如果验证失败，返回 401 错误

**代码框架**：
```go
func Auth(r *ghttp.Request) {
    token := r.Header.Get("Authorization")
    if token == "" || !strings.HasPrefix(token, "Bearer ") {
        r.Response.WriteJson(g.Map{
            "code":    40001,
            "message": "未提供认证令牌",
            "data":    nil,
        })
        r.ExitAll()
        return
    }

    tokenString := strings.TrimPrefix(token, "Bearer ")
    claims, err := jwt.VerifyToken(tokenString)
    if err != nil {
        r.Response.WriteJson(g.Map{
            "code":    40002,
            "message": "令牌无效或已过期",
            "data":    nil,
        })
        r.ExitAll()
        return
    }

    r.SetCtxVar("user_id", claims.UserID)
    r.Middleware.Next()
}
```

#### CORS 中间件（开发环境）

**文件**：`backend/internal/app/http_server/middleware/cors.go`

**职责**：
1. 处理 OPTIONS 预检请求
2. 设置 CORS 响应头（`Access-Control-Allow-Origin`、`Access-Control-Allow-Credentials` 等）
3. 仅在开发环境启用（生产环境前后端同域，无需 CORS）

### 5. 路由注册更新

**文件**：`backend/internal/app/http_server/router.go`

**新增路由**：
```go
func RegisterRoutes(s *ghttp.Server) {
    h := handler.New()

    // API v1 路由组
    s.Group("/api/v1", func(group *ghttp.RouterGroup) {
        group.Middleware(middleware.Format)

        // 认证路由（无需 JWT 验证）
        group.POST("/auth/register", h.Register)
        group.POST("/auth/login", h.Login)
        group.POST("/auth/refresh", h.RefreshToken)

        // 需要认证的路由
        group.Group("/", func(authGroup *ghttp.RouterGroup) {
            authGroup.Middleware(middleware.Auth)

            // 用户信息
            authGroup.GET("/auth/profile", h.GetProfile)
            authGroup.POST("/auth/logout", h.Logout)

            // 学习进度
            authGroup.GET("/progress", h.GetAllProgress)
            authGroup.GET("/progress/:topic", h.GetTopicProgress)
            authGroup.POST("/progress", h.SaveProgress)

            // 测验
            authGroup.GET("/quiz/:topic/:chapter", h.GetQuiz)
            authGroup.POST("/quiz/submit", h.SubmitQuiz)
            authGroup.GET("/quiz/history", h.GetQuizHistory)
            authGroup.GET("/quiz/history/:topic", h.GetTopicQuizHistory)
        })

        // 学习内容路由（已有，无需认证）
        group.ALL("/topics", h.GetTopics)
        group.ALL("/topic/lexical_elements", h.GetLexicalMenu)
        group.ALL("/topic/lexical_elements/:chapter", h.GetLexicalContent)
        // ... 其他已有路由
    })

    // 静态文件服务（托管前端构建产物）
    s.SetServerRoot("frontend/out")
    s.AddStaticPath("/", "frontend/out")
    s.SetRewrite("/", "/index.html")  // SPA 回退
}
```

### 6. 工具包实现

#### JWT 工具

**文件**：`backend/internal/pkg/jwt/jwt.go`

**职责**：
- `GenerateToken(userID int64) (string, error)`: 生成 access token
- `GenerateRefreshToken(userID int64) (string, error)`: 生成 refresh token
- `VerifyToken(tokenString string) (*Claims, error)`: 验证并解析 token

**依赖**：`github.com/golang-jwt/jwt/v5`

#### 密码工具

**文件**：`backend/internal/pkg/password/password.go`

**职责**：
- `Hash(password string) (string, error)`: bcrypt 哈希密码（cost=10）
- `Verify(hashedPassword, password string) error`: 验证密码

#### 输入验证

**位置**：集成在 Service 层或使用 GoFrame 内置验证

**职责**：
- 验证用户名格式：`^[A-Za-z0-9_]{3,50}$`
- 验证密码强度：`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d).{8,}$`
- 验证主题枚举值：lexical_elements/constants/variables/types
- 验证学习状态枚举值：not_started/in_progress/done

---

## Phase 2: 前端实现

### 1. 项目初始化

**步骤**：
1. 在 `frontend/` 目录执行 `npx create-next-app@latest . --typescript --tailwind --app --no-src-dir`
2. 安装依赖：`npm install antd axios swr prismjs @types/prismjs`
3. 配置 `next.config.js`：
```javascript
module.exports = {
  output: 'export',
  images: {
    unoptimized: true,
  },
  trailingSlash: true,
};
```
4. 配置 `tailwind.config.js`：启用 AntD 前缀避免冲突
5. 配置 `tsconfig.json`：启用路径别名（`@/*` -> `./app/*`）

### 2. 核心功能实现

#### 2.1 认证流程

**组件**：
- `app/login/page.tsx`: 登录页
- `app/register/page.tsx`: 注册页
- `components/Auth/LoginForm.tsx`: 登录表单（AntD Form + 验证规则）
- `components/Auth/RegisterForm.tsx`: 注册表单
- `components/Auth/ProtectedRoute.tsx`: 路由守卫（未登录重定向）

**状态管理**：
- `contexts/AuthContext.tsx`: 提供 `user`、`login()`、`logout()`、`isAuthenticated`
- `hooks/useAuth.ts`: 封装认证逻辑（调用 API + 更新 Context）

**token 管理**：
- `lib/auth.ts`: 
  - `setAccessToken(token)`: 存入内存 + localStorage
  - `getAccessToken()`: 优先内存，回退 localStorage
  - `clearTokens()`: 清空内存与 localStorage
  - `refreshAccessToken()`: 调用 `/api/v1/auth/refresh`

**Axios 拦截器**：
```typescript
// lib/api.ts
axios.interceptors.request.use((config) => {
  const token = getAccessToken();
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

axios.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (error.response?.status === 401) {
      // 尝试刷新 token
      try {
        const newToken = await refreshAccessToken();
        setAccessToken(newToken);
        // 重试原请求
        error.config.headers.Authorization = `Bearer ${newToken}`;
        return axios.request(error.config);
      } catch {
        // 刷新失败，清空状态并重定向登录
        clearTokens();
        window.location.href = '/login';
      }
    }
    return Promise.reject(error);
  }
);
```

#### 2.2 学习内容展示

**页面**：
- `app/topics/page.tsx`: 主题列表（调用 `/api/v1/topics`）
- `app/topics/[topic]/page.tsx`: 章节列表（调用 `/api/v1/topic/{topic}`）
- `app/topics/[topic]/[chapter]/page.tsx`: 章节详情（调用 `/api/v1/topic/{topic}/{chapter}`）

**组件**：
- `components/Learning/TopicCard.tsx`: 主题卡片（显示标题、简介、进度百分比）
- `components/Learning/ChapterList.tsx`: 章节列表（显示标题、状态图标）
- `components/Learning/ChapterContent.tsx`: 章节内容（Markdown 渲染 + Prism.js 代码高亮）

**代码高亮**：
```typescript
// components/Learning/ChapterContent.tsx
import Prism from 'prismjs';
import 'prismjs/components/prism-go';
import 'prismjs/components/prism-typescript';
// ... 按需引入其他语言

useEffect(() => {
  Prism.highlightAll();
}, [content]);
```

#### 2.3 学习进度跟踪

**功能**：
1. 进入章节页时，调用 `POST /api/v1/progress` 记录 `status=in_progress`
2. 滚动到底部时，更新 `status=done`
3. 离开页面时，保存 `last_position`（滚动位置）
4. 主题列表页显示进度百分比（已完成章节数 / 总章节数）
5. 首页提供"继续上次学习"按钮（查询 `last_visit` 最新的记录）

**Hooks**：
- `hooks/useProgress.ts`: 封装进度查询与保存逻辑
- `hooks/useScrollPosition.ts`: 监听滚动位置并防抖保存

#### 2.4 测验功能

**页面**：
- `app/quiz/[topic]/page.tsx`: 测验作答页
- `app/quiz/history/page.tsx`: 测验历史记录页

**组件**：
- `components/Quiz/QuizQuestion.tsx`: 单题组件（单选/多选 Radio/Checkbox）
- `components/Quiz/QuizResult.tsx`: 结果展示（得分、正确/错误题目列表）
- `components/Quiz/QuizHistory.tsx`: 历史记录表格（AntD Table + 时间筛选）

**状态管理**：
- `hooks/useQuiz.ts`: 封装题目获取、答案提交、历史查询

**防重复提交**：
```typescript
const [submitting, setSubmitting] = useState(false);

const handleSubmit = async () => {
  if (submitting) return;
  setSubmitting(true);
  try {
    await submitQuiz(answers);
  } finally {
    setSubmitting(false);
  }
};
```

### 3. 响应式设计

**断点定义**（Tailwind）：
```javascript
// tailwind.config.js
module.exports = {
  theme: {
    screens: {
      'sm': '640px',   // Mobile
      'md': '768px',   // Tablet
      'lg': '1024px',  // Desktop
      'xl': '1280px',
    },
  },
};
```

**布局策略**：
- Mobile (<768px): 单列布局，导航栏折叠为汉堡菜单
- Tablet (768-1024px): 2 列卡片布局
- Desktop (>1024px): 3-4 列卡片布局

**AntD 响应式**：
```tsx
<Row gutter={[16, 16]}>
  <Col xs={24} sm={12} md={8} lg={6}>
    <TopicCard />
  </Col>
</Row>
```

### 4. 错误处理

**全局错误边界**：
```tsx
// components/Common/ErrorBoundary.tsx
class ErrorBoundary extends React.Component {
  componentDidCatch(error, errorInfo) {
    console.error(error, errorInfo);
    message.error('页面加载失败，请刷新重试');
  }
  render() {
    return this.props.children;
  }
}
```

**API 错误处理**：
```typescript
// lib/api.ts
axios.interceptors.response.use(
  (response) => {
    if (response.data.code !== 20000) {
      message.error(response.data.message || '操作失败');
      return Promise.reject(new Error(response.data.message));
    }
    return response.data.data;
  },
  (error) => {
    const msg = error.response?.data?.message || '网络错误，请重试';
    message.error(msg);
    return Promise.reject(error);
  }
);
```

### 5. 性能优化

**代码分割**：
```tsx
import dynamic from 'next/dynamic';

const ChapterContent = dynamic(() => import('@/components/Learning/ChapterContent'), {
  loading: () => <Spin />,
  ssr: false,
});
```

**SWR 缓存策略**：
```typescript
const { data, error } = useSWR('/api/v1/topics', fetcher, {
  revalidateOnFocus: false,
  dedupingInterval: 60000, // 1 分钟内不重复请求
});
```

**图片优化**：
- 使用 WebP 格式
- 配置 `next/image` 的 `unoptimized: true`（静态导出限制）

---

## Phase 3: 集成测试

### 1. 后端集成测试

**文件**：`tests/integration/auth_flow_test.go`

**测试场景**：
1. 注册 -> 登录 -> 获取 profile -> 退出 -> 再次登录（验证 refresh token）
2. 注册重复用户名（验证唯一性约束）
3. 错误密码登录（验证 bcrypt 验证）
4. 过期 token 访问受保护路由（验证中间件）

**文件**：`tests/integration/learning_flow_test.go`

**测试场景**：
1. 获取主题列表 -> 获取章节列表 -> 获取章节内容
2. 记录学习进度 -> 查询进度（验证幂等更新）
3. 多次访问同一章节（验证 `last_visit` 更新）

**文件**：`tests/integration/quiz_flow_test.go`

**测试场景**：
1. 获取测验题目 -> 提交答案 -> 查询历史记录
2. 提交全对答案（验证评分）
3. 提交全错答案（验证评分）
4. 按主题筛选历史记录

### 2. 前端集成测试

**文件**：`frontend/__tests__/auth.test.tsx`

**测试场景**：
1. 渲染登录表单 -> 输入用户名密码 -> 提交 -> 验证跳转
2. 输入错误格式用户名 -> 验证错误提示
3. 登录失败 -> 验证错误消息显示

**工具**：Jest + React Testing Library + MSW（Mock Service Worker）

---

## Phase 4: 部署配置

### 1. 前端构建

**命令**：
```bash
cd frontend
npm run build
```

**产物**：`frontend/out/` 目录（包含 HTML/CSS/JS/assets）

### 2. 后端配置

**静态文件服务**：
```go
// backend/internal/app/http_server/server.go
func (s *Server) Start() error {
    // 配置静态文件服务
    s.server.SetServerRoot("frontend/out")
    s.server.AddStaticPath("/", "frontend/out")
    
    // SPA 回退（所有非 API 路径返回 index.html）
    s.server.BindHandler("/*", func(r *ghttp.Request) {
        if !strings.HasPrefix(r.URL.Path, "/api/") {
            r.Response.ServeFile("frontend/out/index.html")
        }
    })
    
    return s.server.Run()
}
```

### 3. 环境变量

**生产环境必须设置**：
```bash
export JWT_SECRET="your-secret-key-min-32-chars"
export DB_PATH="backend/data/gostudy.db"
```

### 4. 启动流程

```bash
# 1. 构建前端
cd frontend && npm run build && cd ..

# 2. 编译后端
cd backend && go build -o ../bin/go-study2 . && cd ..

# 3. 运行
./bin/go-study2
```

**访问地址**：`http://localhost:8080`（API 与前端同端口）

---

## Phase 5: 实施检查清单

### 后端实施顺序

**注意**: 详细任务列表见 `tasks.md`，此处为概览。

1. **数据库与基础设施层** (优先级: P0)
   - [ ] 添加 SQLite driver 依赖到 `go.mod`
   - [ ] 实现 `internal/infrastructure/database/sqlite.go`（数据库初始化）
   - [ ] 实现 `internal/infrastructure/database/migrations.go`（表结构迁移）
   - [ ] 实现领域实体（`internal/domain/user/entity.go` 等）
   - [ ] 实现仓储接口（`internal/domain/user/repository.go` 等）
   - [ ] 实现仓储实现（`internal/infrastructure/repository/user_repo.go` 等）
   - [ ] 编写单元测试（覆盖率≥80%）

2. **工具包** (优先级: P0)
   - [ ] 实现 `internal/pkg/jwt/jwt.go`（JWT 生成与验证）
   - [ ] 实现 `internal/pkg/password/password.go`（密码哈希与验证）
   - [ ] 编写单元测试

3. **中间件** (优先级: P0)
   - [ ] 实现 `internal/app/http_server/middleware/auth.go`（JWT 认证中间件）
   - [ ] 实现 `internal/app/http_server/middleware/cors.go`（CORS 中间件，开发环境）
   - [ ] 编写单元测试

4. **认证 API** (优先级: P1)
   - [ ] 实现 `internal/domain/user/service.go`（用户业务逻辑）
   - [ ] 实现 `internal/app/http_server/handler/auth.go`（注册/登录/登出/刷新/profile）
   - [ ] 更新 `internal/app/http_server/router.go` 注册认证路由
   - [ ] 编写单元测试与集成测试

5. **学习进度 API** (优先级: P2)
   - [ ] 实现 `internal/domain/progress/service.go`（进度业务逻辑）
   - [ ] 实现 `internal/app/http_server/handler/progress.go`（获取/记录进度）
   - [ ] 更新 `internal/app/http_server/router.go` 注册进度路由
   - [ ] 编写单元测试与集成测试

6. **测验 API** (优先级: P3)
   - [ ] 实现 `internal/domain/quiz/service.go`（测验业务逻辑）
   - [ ] 实现 `internal/app/http_server/handler/quiz.go`（获取题目/提交/历史记录）
   - [ ] 更新 `internal/app/http_server/router.go` 注册测验路由
   - [ ] 编写单元测试与集成测试

7. **静态文件服务** (优先级: P4)
   - [ ] 更新 `internal/app/http_server/server.go` 配置静态文件托管
   - [ ] 实现 SPA 回退逻辑
   - [ ] 测试 API 与静态资源路由优先级

8. **配置与初始化** (优先级: P0)
   - [ ] 更新 `backend/configs/config.yaml` 添加数据库与 JWT 配置
   - [ ] 更新 `backend/main.go` 初始化数据库连接
   - [ ] 编写环境变量文档

### 前端实施顺序

1. **项目初始化** (优先级: P0)
   - [ ] 创建 Next.js 14 项目（TypeScript + Tailwind + App Router）
   - [ ] 安装依赖（Ant Design + Axios + SWR + Prism.js）
   - [ ] 配置 `next.config.js`（静态导出）
   - [ ] 配置 `tailwind.config.js` 与 `tsconfig.json`

2. **基础设施** (优先级: P0)
   - [ ] 实现 `lib/api.ts`（Axios 实例与拦截器）
   - [ ] 实现 `lib/auth.ts`（token 管理）
   - [ ] 实现 `lib/constants.ts`（常量定义）
   - [ ] 定义 `types/` 下所有类型

3. **认证功能** (优先级: P1)
   - [ ] 实现 `contexts/AuthContext.tsx`（认证上下文）
   - [ ] 实现 `hooks/useAuth.ts`（认证 Hook）
   - [ ] 实现 `components/Auth/LoginForm.tsx`（登录表单）
   - [ ] 实现 `components/Auth/RegisterForm.tsx`（注册表单）
   - [ ] 实现 `components/Auth/ProtectedRoute.tsx`（路由守卫）
   - [ ] 实现 `app/login/page.tsx` 与 `app/register/page.tsx`
   - [ ] 编写组件测试

4. **布局与导航** (优先级: P1)
   - [ ] 实现 `app/layout.tsx`（根布局 + AntD ConfigProvider）
   - [ ] 实现 `components/Layout/Header.tsx`（顶部导航栏）
   - [ ] 实现 `components/Layout/Footer.tsx`（页脚）
   - [ ] 配置 `styles/globals.css` 与 `styles/theme.ts`

5. **学习内容展示** (优先级: P2)
   - [ ] 实现 `components/Learning/TopicCard.tsx`（主题卡片）
   - [ ] 实现 `components/Learning/ChapterList.tsx`（章节列表）
   - [ ] 实现 `components/Learning/ChapterContent.tsx`（章节内容 + 代码高亮）
   - [ ] 实现 `app/topics/page.tsx`（主题列表页）
   - [ ] 实现 `app/topics/[topic]/page.tsx`（章节列表页）
   - [ ] 实现 `app/topics/[topic]/[chapter]/page.tsx`（章节详情页）
   - [ ] 编写组件测试

6. **学习进度跟踪** (优先级: P2)
   - [ ] 实现 `hooks/useProgress.ts`（进度管理 Hook）
   - [ ] 实现 `hooks/useScrollPosition.ts`（滚动位置监听）
   - [ ] 实现 `components/Learning/ProgressIndicator.tsx`（进度指示器）
   - [ ] 集成进度记录到章节页
   - [ ] 实现"继续上次学习"功能

7. **测验功能** (优先级: P3)
   - [ ] 实现 `hooks/useQuiz.ts`（测验管理 Hook）
   - [ ] 实现 `components/Quiz/QuizQuestion.tsx`（题目组件）
   - [ ] 实现 `components/Quiz/QuizResult.tsx`（结果组件）
   - [ ] 实现 `components/Quiz/QuizHistory.tsx`（历史记录组件）
   - [ ] 实现 `app/quiz/[topic]/page.tsx`（测验作答页）
   - [ ] 实现 `app/quiz/history/page.tsx`（历史记录页）
   - [ ] 编写组件测试

8. **错误处理与优化** (优先级: P4)
   - [ ] 实现 `components/Common/ErrorBoundary.tsx`（错误边界）
   - [ ] 实现 `components/Common/Loading.tsx`（加载状态）
   - [ ] 实现 `components/Common/ErrorMessage.tsx`（错误提示）
   - [ ] 配置 SWR 缓存策略
   - [ ] 实现代码分割（动态导入重组件）
   - [ ] 优化图片与静态资源

### 集成与部署

1. **集成测试** (优先级: P4)
   - [ ] 后端集成测试（auth/learning/quiz flow）
   - [ ] 前端集成测试（端到端流程）
   - [ ] API 契约测试（验证响应格式）

2. **部署准备** (优先级: P4)
   - [ ] 前端构建测试（`npm run build`）
   - [ ] 后端编译测试（`go build`）
   - [ ] 静态文件服务测试（验证路由优先级）
   - [ ] 环境变量配置文档
   - [ ] 更新根 README.md

3. **文档更新** (优先级: P4)
   - [ ] 更新 `backend/README.md`（新增 API 说明）
   - [ ] 创建 `frontend/README.md`（安装/运行/构建指南）
   - [ ] 更新根 `README.md`（新增功能章节）
   - [ ] 创建 `docs/API.md`（API 文档）
   - [ ] 创建 `docs/DEPLOYMENT.md`（部署指南）

---

## 风险评估与缓解策略

### 技术风险

| 风险项 | 影响 | 概率 | 缓解策略 |
|--------|------|------|----------|
| SQLite CGO 依赖导致交叉编译困难 | 中 | 高 | 提供 Docker 构建环境；文档说明 MinGW 安装步骤；考虑纯 Go SQLite 实现（modernc.org/sqlite）作为备选 |
| JWT secret 泄露导致安全风险 | 高 | 低 | 强制从环境变量读取；启动时验证 secret 长度≥32；生产环境使用密钥管理服务 |
| SQLite 并发写入锁冲突 | 中 | 中 | 启用 WAL 模式；设置 busy_timeout；限制单实例部署；监控锁等待时间 |
| 前端静态导出限制动态路由 | 中 | 中 | 使用 `generateStaticParams` 预生成已知路由；客户端获取动态数据；文档说明新增章节需重新构建 |
| refresh token 被盗用 | 高 | 低 | HttpOnly Cookie + SameSite=Lax；实现 token rotation；记录刷新日志；异常检测（IP/UA 变化） |

### 实施风险

| 风险项 | 影响 | 概率 | 缓解策略 |
|--------|------|------|----------|
| 前后端接口契约不一致 | 高 | 中 | 使用 OpenAPI 规范作为单一事实来源；编写契约测试；前后端并行开发时先 mock API |
| 测试覆盖率不达标（<80%） | 中 | 中 | 优先编写核心路径测试；使用覆盖率工具监控；代码审查检查测试完整性 |
| 响应式布局在移动端显示异常 | 中 | 中 | 使用 AntD 响应式组件；在多设备测试；使用浏览器开发工具模拟移动设备 |
| 数据库迁移失败导致启动失败 | 高 | 低 | 迁移脚本幂等设计（CREATE IF NOT EXISTS）；启动前备份数据库；提供迁移回滚机制 |
| 代码高亮库体积过大影响加载速度 | 低 | 中 | 按需引入语言包；使用动态导入；监控构建产物体积；设置体积预算（<500KB） |

### 业务风险

| 风险项 | 影响 | 概率 | 缓解策略 |
|--------|------|------|----------|
| 用户数据丢失（数据库损坏） | 高 | 低 | 定期备份数据库文件；启用 SQLite WAL 模式；提供数据导出功能；文档说明备份策略 |
| 密码重置功能缺失导致用户锁定 | 中 | 中 | 当前版本不实现密码重置（YAGNI）；文档说明管理员可直接修改数据库；后续版本可添加邮箱验证 |
| 学习进度数据不一致 | 中 | 低 | 使用唯一约束保证幂等更新；前端防抖避免频繁请求；后端事务保证原子性 |
| 测验作弊（查看源码获取答案） | 低 | 高 | 当前版本接受此风险（学习工具非考试系统）；后续可实现服务端评分；文档说明设计意图 |

---

## 实施里程碑

### Milestone 1: 后端基础设施 (Week 1-2)
- 数据库初始化与迁移
- JWT 工具与认证中间件
- 用户模型与认证 API
- 单元测试覆盖率≥80%

**验收标准**：
- ✅ 可通过 Postman 完成注册/登录/刷新/登出流程
- ✅ 数据库表结构正确创建并包含索引
- ✅ JWT token 正确生成与验证
- ✅ 所有单元测试通过

### Milestone 2: 学习进度与测验 API (Week 2-3)
- 学习进度模型与 API
- 测验记录模型与 API
- 集成测试覆盖主要流程

**验收标准**：
- ✅ 可通过 API 记录与查询学习进度
- ✅ 可通过 API 提交测验并查询历史记录
- ✅ 集成测试覆盖 auth/learning/quiz 完整流程
- ✅ API 响应时间 p95 < 200ms

### Milestone 3: 前端认证与布局 (Week 3-4)
- 项目初始化与依赖安装
- 认证功能（登录/注册/路由守卫）
- 布局组件（Header/Footer）
- 响应式设计基础

**验收标准**：
- ✅ 用户可完成注册/登录流程
- ✅ 未登录访问受保护页面自动重定向
- ✅ 布局在 Mobile/Tablet/Desktop 正确显示
- ✅ 核心组件测试覆盖率≥80%

### Milestone 4: 学习内容展示 (Week 4-5)
- 主题列表与章节列表页
- 章节详情页（含代码高亮）
- 学习进度跟踪
- "继续上次学习"功能

**验收标准**：
- ✅ 用户可浏览主题与章节列表
- ✅ 章节内容正确渲染，代码高亮正常
- ✅ 学习进度正确记录与显示
- ✅ "继续学习"按钮跳转到正确章节

### Milestone 5: 测验功能 (Week 5-6)
- 测验作答页
- 测验结果展示
- 测验历史记录页
- 防重复提交

**验收标准**：
- ✅ 用户可完成测验作答与提交
- ✅ 评分结果正确显示
- ✅ 历史记录可按主题与时间筛选
- ✅ 提交按钮正确防抖

### Milestone 6: 集成与部署 (Week 6-7)
- 前端构建与静态文件服务
- 端到端集成测试
- 文档更新
- 部署验证

**验收标准**：
- ✅ 前端构建产物正确托管
- ✅ API 与静态资源路由优先级正确
- ✅ 所有集成测试通过
- ✅ README 与 API 文档更新完整
- ✅ 生产环境部署成功

---

## 后续扩展计划（不在当前范围）

以下功能遵循 YAGNI 原则，当前版本**不实现**，仅作为后续迭代参考：

1. **密码重置功能**：需要邮箱验证或管理员重置，当前版本用户可联系管理员直接修改数据库
2. **用户头像与个人资料编辑**：当前仅支持用户名，不支持昵称/头像/简介等扩展字段
3. **学习统计与数据可视化**：如学习时长、完成率趋势图等，当前仅提供基础进度百分比
4. **社交功能**：如学习笔记分享、讨论区等，当前为单用户学习工具
5. **多语言支持**：当前仅中文，国际化需要前后端全面改造
6. **离线模式**：当前需要网络连接，Service Worker 离线缓存可作为后续优化
7. **测验题库管理后台**：当前测验题目硬编码在学习内容模块，管理后台需要独立开发
8. **学习路径推荐**：基于进度与测验成绩的智能推荐，需要算法支持
9. **证书与徽章系统**：完成学习后颁发证书，需要设计激励机制
10. **移动端原生应用**：当前为响应式 Web 应用，原生 App 需要独立开发

---

## 总结

本实施计划详细设计了 Go-Study2 前端 UI 界面的完整技术方案，涵盖：

### 后端扩展
- **数据库**：SQLite3（WAL 模式）+ 4 张表（users/learning_progress/quiz_records/refresh_tokens）
- **新增 API**：认证 5 个 + 学习进度 3 个 + 测验 5 个，共 13 个端点
- **安全机制**：JWT（access + refresh token）+ bcrypt 密码哈希 + HttpOnly Cookie + 输入校验
- **复用现有**：学习内容接口（`/api/v1/topics` 等）保持不变，CLI 模式不受影响

### 前端实现
- **技术栈**：Next.js 14（App Router + 静态导出）+ TypeScript + Ant Design 5 + Tailwind CSS
- **核心功能**：用户认证 + 学习内容展示 + 进度跟踪 + 测验作答与历史记录
- **响应式设计**：Mobile/Tablet/Desktop 三档断点，AntD Grid 布局
- **性能优化**：代码分割 + SWR 缓存 + 按需导入 + 静态预生成

### 部署方案
- **一体化部署**：前端构建产物由后端静态文件服务托管，API 与前端同端口（`:8080`）
- **开发分离**：前端开发服务器 `:3000`，后端 `:8080`，通过 CORS 或 proxy 通信
- **环境变量**：JWT secret、数据库路径等敏感配置通过环境变量管理

### 质量保证
- **测试覆盖率**：后端单元测试≥80%，前端核心组件≥80%，集成测试覆盖主要流程
- **宪章合规**：所有 44 条宪章原则检查通过，无违规项
- **风险缓解**：识别 15 项技术/实施/业务风险，并提供具体缓解策略

### 实施路径
- **6 个里程碑**：从后端基础设施到前端功能再到集成部署，预计 6-7 周完成
- **优先级清晰**：P0（基础设施）-> P1（认证）-> P2（学习内容）-> P3（测验）-> P4（优化与部署）
- **可追溯性**：每个功能点对应 spec.md 中的需求编号，确保不遗漏

**下一步行动**：执行 `/speckit.tasks` 命令将本计划分解为可执行任务列表。

---

## Complexity Tracking

无宪章违规项需要特别说明。本方案遵循 YAGNI 原则，仅引入必要的依赖与模式。
