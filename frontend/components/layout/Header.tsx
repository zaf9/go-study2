"use client";

import { Button, Dropdown, Layout, Space, Typography, message } from "antd";
import { MenuOutlined, UserOutlined } from "@ant-design/icons";
import { useRouter, usePathname } from "next/navigation";
import useAuth from "@/hooks/useAuth";
import { useMemo } from "react";

interface HeaderProps {
    onToggleSidebar?: () => void;
}

const { Header: AntHeader } = Layout;
const { Text } = Typography;

export default function Header({ onToggleSidebar }: HeaderProps) {
    const { user, logout } = useAuth();
    const router = useRouter();
    const pathname = usePathname();

    const pageTitle = useMemo(() => {
        if (!pathname) return "";
        if (pathname.startsWith("/dashboard")) return "首页";
        if (pathname.startsWith("/topics/")) return "课程详情";
        if (pathname.startsWith("/topics")) return "学习主题";
        if (pathname.startsWith("/progress")) return "学习进度";
        if (pathname.startsWith("/quiz-center")) return "测验中心";
        if (pathname.startsWith("/quiz")) return "章节测验";
        return "";
    }, [pathname]);

    const handleLogout = async () => {
        await logout();
        message.success("已退出登录");
        router.replace("/login");
    };

    const menuItems = [
        {
            key: "profile",
            label: (
                <div className="flex flex-col">
                    <Text strong>{user?.username}</Text>
                    <Text type="secondary">已登录</Text>
                </div>
            ),
            disabled: true,
        },
        { type: "divider" as const },
        {
            key: "logout",
            label: "退出登录",
            onClick: handleLogout,
        },
    ];

    return (
        <AntHeader
            className="sticky top-0 z-50 w-full px-0"
            style={{
                height: "64px",
                background: "rgba(255, 255, 255, 0.72)",
                backdropFilter: "blur(20px)",
                WebkitBackdropFilter: "blur(20px)",
                borderBottom: "1px solid rgba(0, 0, 0, 0.06)",
                lineHeight: "64px",
            }}
        >
            <div className="flex h-full w-full items-center justify-between px-6 lg:px-8">
                <Space size={12} align="center">
                    <Button
                        icon={<MenuOutlined />}
                        type="text"
                        className="flex items-center justify-center hover:bg-black/5"
                        onClick={onToggleSidebar}
                    />
                    <div className="h-4 w-[1px] bg-black/10 mx-1" />
                    <Text className="text-base font-medium tracking-tight text-[#1d1d1f]">
                        {pageTitle}
                    </Text>
                </Space>
                <Dropdown menu={{ items: menuItems }} trigger={["click"]}>
                    <Button
                        type="default"
                        shape="round"
                        icon={<UserOutlined />}
                        className="flex items-center border-none bg-black/5 font-medium text-[#1d1d1f] shadow-none hover:bg-black/10"
                    >
                        <span className="ml-2">
                            {user?.username || "未登录"}
                        </span>
                    </Button>
                </Dropdown>
            </div>
        </AntHeader>
    );
}
