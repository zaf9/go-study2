# Go-Study2 前端（Next.js 14）

基于 Next.js 14（App Router）+ TypeScript + Ant Design，静态导出由后端托管。

## 快速开始

```bash
cd frontend
npm install
npm run dev   # http://localhost:3000
```

环境变量：

- `NEXT_PUBLIC_API_URL`：后端 API 基址（默认 `/api/v1`，生产建议设为完整域名）。
- `NEXT_PUBLIC_API_URL`：后端 API 基址。建议在本地开发时指向后端服务（例如 `http://localhost:8080`）。

开发提示：
- 如果未设置 `NEXT_PUBLIC_API_URL`，Next.js 开发服务器会尝试将以 `/api/v1` 开头的请求发送到自身，可能导致 404。为避免此类问题：
	1. 在 `frontend/.env.local` 中设置 `NEXT_PUBLIC_API_URL=http://localhost:8080`（或你的后端地址），然后重启 dev server；
	2. 项目已在 `next.config.ts` 中为 dev 模式添加了 `rewrites`，会将 `/api/v1/:path*` 代理到 `NEXT_PUBLIC_API_URL`（如果未设置则默认 `http://localhost:8080`）。

示例 `.env.local`：
```
NEXT_PUBLIC_API_URL=http://localhost:8080
```

重启 Dev Server：
```bash
cd frontend
npm run dev
```

## 构建与导出

```bash
npm run build   # next.config.ts 已配置 output: 'export'，产物输出到 frontend/out
```

构建完成后，`frontend/out/` 可直接由后端静态托管（参见 `backend/internal/app/http_server/server.go`）。

## 测试

```bash
npm test -- --coverage
```

核心覆盖：

- `lib/api` 拦截器与鉴权流程
- `contexts/AuthContext` 登录/注销逻辑
- 组件与页面集成测试（见 `tests/`）

## 代码规范

- ESLint + TypeScript 严格模式
- SWR 全局缓存：`revalidateOnFocus=false`，`dedupingInterval=60s`（见 `app/(protected)/layout.tsx`）
- Prism.js 语言包按需动态加载，减少首包体积
