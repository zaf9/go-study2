import QuizPageClient from "./QuizPageClient";
import { buildTopicParams } from "@/lib/static-routes";

export const generateStaticParams = async () => buildTopicParams();

export default function QuizPage({ params }: { params: { topic: string } }) {
  return <QuizPageClient params={params} />;
}

