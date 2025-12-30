/**
 * T118: 测验中心页面
 * 提供测验历史记录的统一入口
 */

"use client";

import { Breadcrumb, Button, Card, Typography, Space } from "antd";
import { HomeOutlined, TrophyOutlined } from "@ant-design/icons";
import { useRouter } from "next/navigation";
import QuizCenter from "@/components/quiz/QuizCenter";
import useAuth from "@/hooks/useAuth";
import { useEffect } from "react";

const { Title, Text } = Typography;

export default function QuizCenterPage() {
  const router = useRouter();
  const { user } = useAuth();

  // 未认证用户重定向到登录页
  useEffect(() => {
    if (!user) {
      router.replace("/login");
    }
  }, [user, router]);

  if (!user) {
    return null;
  }

  return (
    <div className="quiz-center-page container mx-auto px-4 py-8" data-testid="quiz-center-page">
      {/* 面包屑导航 */}
      <Breadcrumb
        className="mb-6"
        items={[
          {
            href: "/dashboard",
            title: (
              <>
                <HomeOutlined />
                <span>主页</span>
              </>
            ),
          },
          {
            title: "测验中心",
          },
        ]}
      />

      {/* 页面标题 */}
      <div className="mb-8">
        <Space direction="vertical" size="small">
          <Title level={2}>
            <TrophyOutlined className="mr-2" />
            测验中心
          </Title>
          <Text type="secondary">查看和管理你的测验记录</Text>
        </Space>
      </div>

      {/* 统计摘要区域 */}
      <Card className="mb-6" data-testid="stats-summary">
        <Text>统计摘要将在后续版本中添加</Text>
      </Card>

      {/* 测验中心组件 */}
      <QuizCenter />

      {/* 返回按钮 */}
      <div className="mt-8">
        <Button onClick={() => router.push("/dashboard")}>返回主页</Button>
      </div>
    </div>
  );
}
