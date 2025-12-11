import TopicDetailClient from "./TopicDetailClient";
import { buildTopicParams } from "@/lib/static-routes";

export const generateStaticParams = async () => buildTopicParams();

export default function TopicDetailPage({
  params,
}: {
  params: { topic: string };
}) {
  return <TopicDetailClient params={params} />;
}
