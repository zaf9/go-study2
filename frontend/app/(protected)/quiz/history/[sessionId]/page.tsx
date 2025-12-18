import QuizReviewPageClient from "./QuizReviewPageClient";

// 为静态导出提供 generateStaticParams
// 由于 sessionId 是动态的，返回一个占位符路径以满足静态导出的要求
// 实际运行时，客户端会通过 useParams 获取真实的 sessionId
export const generateStaticParams = async () => {
  // 返回一个占位符路径，满足静态导出的要求
  // 实际的路由参数会在客户端运行时通过 useParams 获取
  return [{ sessionId: "placeholder" }];
};

export default function QuizReviewPage() {
  return <QuizReviewPageClient />;
}
