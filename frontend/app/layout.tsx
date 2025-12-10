import type { Metadata } from "next";
import localFont from "next/font/local";
import { ConfigProvider, App as AntdApp } from "antd";
import { StyleProvider } from "@ant-design/cssinjs";
import zhCN from "antd/locale/zh_CN";
import "./globals.css";

const geistSans = localFont({
  src: "./fonts/GeistVF.woff",
  variable: "--font-geist-sans",
  weight: "100 900",
});
const geistMono = localFont({
  src: "./fonts/GeistMonoVF.woff",
  variable: "--font-geist-mono",
  weight: "100 900",
});

export const metadata: Metadata = {
  title: "Go Study 2",
  description: "Go 学习平台前端",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="zh-CN">
      <body className={`${geistSans.variable} ${geistMono.variable} antialiased`}>
        <StyleProvider hashPriority="high">
          <ConfigProvider locale={zhCN} theme={{ token: { colorPrimary: "#1677ff" } }}>
            <AntdApp>{children}</AntdApp>
          </ConfigProvider>
        </StyleProvider>
      </body>
    </html>
  );
}
