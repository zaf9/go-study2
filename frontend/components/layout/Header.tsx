'use client';

import { Button, Dropdown, Layout, Space, Typography, message } from "antd";
import { MenuOutlined, UserOutlined } from "@ant-design/icons";
import { useRouter } from "next/navigation";
import useAuth from "@/hooks/useAuth";

interface HeaderProps {
  onToggleSidebar?: () => void;
}

const { Header: AntHeader } = Layout;
const { Text } = Typography;

export default function Header({ onToggleSidebar }: HeaderProps) {
  const { user, logout } = useAuth();
  const router = useRouter();

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
    <AntHeader className="flex items-center justify-between bg-white px-4 shadow-sm">
      <Space size="large" align="center">
        <Button icon={<MenuOutlined />} type="text" onClick={onToggleSidebar} />
        <Text className="text-lg font-semibold">Go Study 2</Text>
      </Space>
      <Dropdown menu={{ items: menuItems }} trigger={["click"]}>
        <Button type="text" icon={<UserOutlined />} className="flex items-center">
          <span className="ml-2">{user?.username || "未登录"}</span>
        </Button>
      </Dropdown>
    </AntHeader>
  );
}


