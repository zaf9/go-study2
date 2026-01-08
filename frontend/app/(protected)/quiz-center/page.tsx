/**
 * T118: 测验中心页面
 * 提供测验历史记录的统一入口，支持卡片和表格视图切换
 */

"use client";

import { Typography } from "antd";
import { TrophyOutlined } from "@ant-design/icons";
import { useRouter } from "next/navigation";
import QuizCenter from "@/components/quiz/QuizCenter";
import useAuth from "@/hooks/useAuth";
import { useEffect } from "react";
import PageContainer from "@/components/layout/PageContainer";

const { Title } = Typography;

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
    <PageContainer className="space-y-6" data-testid="quiz-center-page">
      {/* 页面标题 */}
      <div>
        <Title level={2} className="flex items-center gap-2">
          <TrophyOutlined className="text-yellow-500" />
          测验中心
        </Title>
      </div>

      {/* 测验中心组件 */}
      <QuizCenter />
    </PageContainer>
  );
}
