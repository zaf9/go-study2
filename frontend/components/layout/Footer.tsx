'use client';

import { Layout, Typography } from "antd";

const { Footer: AntFooter } = Layout;
const { Text } = Typography;

export default function Footer() {
  return (
    <AntFooter className="text-center bg-white">
      <Text type="secondary">Go Study 2 · 学习与实践并重</Text>
    </AntFooter>
  );
}


