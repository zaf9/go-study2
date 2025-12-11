import QuizPageClient from "../QuizPageClient";
import { buildTopicChapterParams } from "@/lib/static-routes";

export const generateStaticParams = async () => buildTopicChapterParams();

export default function ChapterQuizPage({
  params,
}: {
  params: { topic: string; chapter: string };
}) {
  return <QuizPageClient params={params} />;
}

