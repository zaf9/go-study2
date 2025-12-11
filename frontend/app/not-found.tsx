"use client";

import { Button, Result, Typography } from "antd";
import { useRouter } from "next/navigation";

const { Paragraph, Text } = Typography;

export default function NotFoundPage() {
  const router = useRouter();

  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-50 p-6">
      <Result
        status="404"
        title="页面未找到"
        subTitle="抱歉，您访问的页面不存在或已被移动。"
        extra={
          <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-center">
            <Button type="primary" onClick={() => router.push("/topics")}>
              返回学习主题
            </Button>
            <Button onClick={() => router.push("/")}>回到首页</Button>
          </div>
        }
      >
        <Paragraph>
          <Text>请检查链接是否正确，或点击上方按钮回到可用页面继续学习。</Text>
        </Paragraph>
      </Result>
    </div>
  );
}
