"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";

/**
 * 测验历史页面
 * 此页面已重定向到测验中心（/quiz-center）
 * 保留此路由是为了兼容旧链接和书签
 */
export default function QuizHistoryPage() {
  const router = useRouter();

  useEffect(() => {
    // 重定向到测验中心
    router.replace("/quiz-center");
  }, [router]);

  // 返回 null 或加载状态，因为会立即重定向
  return null;
}

