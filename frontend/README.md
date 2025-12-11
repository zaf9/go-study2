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
