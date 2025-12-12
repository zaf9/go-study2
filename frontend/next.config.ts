import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: "export",
  images: {
    unoptimized: true,
  },
  trailingSlash: true,
  // 在开发环境下将相对的 /api/v1 请求代理到后端服务，便于本地调试。
  // 优先使用 NEXT_PUBLIC_API_URL 环境变量，否则默认到 http://localhost:8080
  async rewrites() {
    const backend = (process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080").replace(/\/+$/, "");
    return [
      {
        source: "/api/v1/:path*",
        destination: `${backend}/api/v1/:path*`,
      },
    ];
  },
};

export default nextConfig;
