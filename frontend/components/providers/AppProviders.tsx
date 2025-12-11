"use client";

import { ConfigProvider, App as AntdApp } from "antd";
import { StyleProvider } from "@ant-design/cssinjs";
import zhCN from "antd/locale/zh_CN";
import ErrorBoundary from "@/components/common/ErrorBoundary";
import { AuthProvider } from "@/contexts/AuthContext";

export default function AppProviders({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <StyleProvider hashPriority="high">
      <ConfigProvider
        locale={zhCN}
        theme={{ token: { colorPrimary: "#1677ff" } }}
      >
        <AntdApp>
          <AuthProvider>
            <ErrorBoundary>{children}</ErrorBoundary>
          </AuthProvider>
        </AntdApp>
      </ConfigProvider>
    </StyleProvider>
  );
}
