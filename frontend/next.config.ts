/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'export',
  images: {
    unoptimized: true,
  },
  // 注意：rewrites 与静态导出（output: 'export'）不兼容
  // 仅在开发环境启用，生产环境注释掉
  ...(process.env.NODE_ENV !== 'production' && {
    async rewrites() {
      const backend = (process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080").replace(/\/+$/, "");
      return [
        {
          source: "/api/v1/:path*",
          destination: `${backend}/api/v1/:path*`,
        },
      ];
    },
  }),
};

export default nextConfig;
