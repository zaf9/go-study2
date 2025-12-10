'use client';

import type React from "react";
import { useState } from "react";
import { Layout } from "antd";
import AuthGuard from "@/components/auth/AuthGuard";
import Header from "@/components/layout/Header";
import Sidebar from "@/components/layout/Sidebar";
import Footer from "@/components/layout/Footer";

const { Content } = Layout;

export default function ProtectedLayout({ children }: { children: React.ReactNode }) {
  const [collapsed, setCollapsed] = useState(false);

  return (
    <AuthGuard>
      <Layout style={{ minHeight: "100vh" }}>
        <Sidebar collapsed={collapsed} onCollapse={setCollapsed} />
        <Layout>
          <Header onToggleSidebar={() => setCollapsed((prev) => !prev)} />
          <Content className="p-6 bg-gray-50">{children}</Content>
          <Footer />
        </Layout>
      </Layout>
    </AuthGuard>
  );
}


