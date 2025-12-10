'use client';

import { Layout, Menu } from "antd";
import { BookOutlined, HomeOutlined } from "@ant-design/icons";
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
    if (pathname.startsWith("/topics")) return ["/topics"];
    return [pathname];
  }, [pathname]);

  const items = [
    {
      key: "/topics",
      icon: <BookOutlined />,
      label: "学习主题",
      onClick: () => router.push("/topics"),
    },
    {
      key: "/",
      icon: <HomeOutlined />,
      label: "首页",
      onClick: () => router.push("/topics"),
    },
  ];

  return (
    <Sider
      collapsible
      collapsed={collapsed}
      onCollapse={onCollapse}
      breakpoint="lg"
      width={220}
      className="min-h-screen"
    >
      <div className="h-12 text-center text-white flex items-center justify-center font-semibold text-base">
        Go Study 2
      </div>
      <Menu theme="dark" mode="inline" selectedKeys={selectedKeys} items={items} />
    </Sider>
  );
}


