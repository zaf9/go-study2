"use client";

import { Layout, Menu } from "antd";
import {
    BookOutlined,
    HomeOutlined,
    LineChartOutlined,
    QuestionCircleOutlined,
} from "@ant-design/icons";
import { usePathname, useRouter } from "next/navigation";
import { useMemo } from "react";

const { Sider } = Layout;

interface SidebarProps {
    collapsed: boolean;
    onCollapse: (next: boolean) => void;
}

export default function Sidebar({ collapsed, onCollapse }: SidebarProps) {
    const router = useRouter();
    const pathname = usePathname();

    const selectedKeys = useMemo(() => {
        if (!pathname) return [];
        if (pathname.startsWith("/dashboard")) return ["/dashboard"];
        if (pathname.startsWith("/topics")) return ["/topics"];
        if (pathname.startsWith("/progress")) return ["/progress"];
        if (pathname.startsWith("/quiz")) return ["/quiz"]; // 匹配 /quiz 和 /quiz-center
        return [pathname];
    }, [pathname]);

    const items = [
        {
            key: "/dashboard",
            icon: <HomeOutlined />,
            label: "首页",
            onClick: () => router.push("/dashboard"),
        },
        {
            key: "/topics",
            icon: <BookOutlined />,
            label: "学习主题",
            onClick: () => router.push("/topics"),
        },
        {
            key: "/progress",
            icon: <LineChartOutlined />,
            label: "学习进度",
            onClick: () => router.push("/progress"),
        },
        {
            key: "/quiz",
            icon: <QuestionCircleOutlined />,
            label: "章节测验",
            // 导航到测验中心页面
            onClick: () => router.push("/quiz-center"),
        },
    ];

    return (
        <Sider
            collapsible
            collapsed={collapsed}
            onCollapse={onCollapse}
            trigger={null}
            breakpoint="lg"
            width={220}
            style={{
                borderRight: "1px solid rgba(0, 0, 0, 0.06)",
                zIndex: 100,
            }}
        >
            <div className={`flex h-16 items-center px-6 transition-all duration-300 ${collapsed ? "justify-center px-0" : "justify-start"}`}>
                <div className="flex items-center gap-3">
                    <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-blue-600 font-bold text-white shadow-lg shadow-blue-500/20">
                        G
                    </div>
                    {!collapsed && (
                        <span className="text-base font-bold tracking-tight text-white opacity-90">
                            Go Study 2
                        </span>
                    )}
                </div>
            </div>
            <Menu
                theme="dark"
                mode="inline"
                selectedKeys={selectedKeys}
                items={items}
                style={{ borderRight: 0 }}
            />
        </Sider>
    );
}
