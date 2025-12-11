# Go-Study2 前端（Next.js 14）

基于 Next.js 14（App Router）+ TypeScript + Ant Design，采用静态导出托管于后端。

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
npm run build         # 生成 .next
npm run export        # 若需要单独导出，可生成 out/ 目录
```

本仓库已启用 `generateStaticParams` 预生成 topics/quiz 路由，产物在 `frontend/out/`，后端静态托管。

## 测试

```bash
npm test -- --coverage
```

核心覆盖：
- `lib/api` 拦截器与鉴权流程
- `contexts/AuthContext` 登录/注销逻辑
- 组件与页面集成测试（见 `tests/`）

## 代码规范

- ESLint + TypeScript 严格模式。
- SWR 全局缓存：`revalidateOnFocus=false`，`dedupingInterval=60s`，见 `app/(protected)/layout.tsx`。
- Prism.js 语言包按需动态加载，减少首包体积。*** End Patch
This is a [Next.js](https://nextjs.org) project bootstrapped with [`create-next-app`](https://nextjs.org/docs/app/api-reference/cli/create-next-app).

## Getting Started

First, run the development server:

```bash
npm run dev
# or
yarn dev
# or
pnpm dev
# or
bun dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

You can start editing the page by modifying `app/page.tsx`. The page auto-updates as you edit the file.

This project uses [`next/font`](https://nextjs.org/docs/app/building-your-application/optimizing/fonts) to automatically optimize and load [Geist](https://vercel.com/font), a new font family for Vercel.

## Learn More

To learn more about Next.js, take a look at the following resources:

- [Next.js Documentation](https://nextjs.org/docs) - learn about Next.js features and API.
- [Learn Next.js](https://nextjs.org/learn) - an interactive Next.js tutorial.

You can check out [the Next.js GitHub repository](https://github.com/vercel/next.js) - your feedback and contributions are welcome!

## Deploy on Vercel

The easiest way to deploy your Next.js app is to use the [Vercel Platform](https://vercel.com/new?utm_medium=default-template&filter=next.js&utm_source=create-next-app&utm_campaign=create-next-app-readme) from the creators of Next.js.

Check out our [Next.js deployment documentation](https://nextjs.org/docs/app/building-your-application/deploying) for more details.
