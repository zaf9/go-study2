"use client";

import type React from "react";
import { useState } from "react";
import { SWRConfig } from "swr";
import { Layout } from "antd";
import AuthGuard from "@/components/auth/AuthGuard";
import Header from "@/components/layout/Header";
import Sidebar from "@/components/layout/Sidebar";
import Footer from "@/components/layout/Footer";

const { Content } = Layout;

export default function ProtectedLayout({
    children,
}: {
    children: React.ReactNode;
}) {
    const [collapsed, setCollapsed] = useState(false);

    return (
        <SWRConfig
            value={{
                revalidateOnFocus: false,
                dedupingInterval: 60000,
                errorRetryInterval: 2000,
            }}
        >
            <AuthGuard>
                <Layout style={{ height: "100vh", overflow: "hidden" }}>
                    <Sidebar collapsed={collapsed} onCollapse={setCollapsed} />
                    <Layout style={{ display: "flex", flexDirection: "column" }}>
                        <Header onToggleSidebar={() => setCollapsed((prev) => !prev)} />
                        <Content
                            className="bg-[#f5f5f7]"
                            style={{
                                flex: 1,
                                overflowY: "auto",
                                overflowX: "hidden",
                                position: "relative",
                            }}
                        >
                            {children}
                            <Footer />
                        </Content>
                    </Layout>
                </Layout>
            </AuthGuard>
        </SWRConfig>
    );
}
