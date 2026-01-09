/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'export',
  images: {
    unoptimized: true,
  },
  // 移除 rewrites，因为与静态导出不兼容
};

module.exports = nextConfig;
