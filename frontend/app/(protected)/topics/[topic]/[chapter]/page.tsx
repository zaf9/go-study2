import ChapterPageClient from "./ChapterPageClient";
import { buildTopicChapterParams } from "@/lib/static-routes";

export const generateStaticParams = async () => buildTopicChapterParams();

export default function ChapterPage({
  params,
}: {
  params: { topic: string; chapter: string };
}) {
  return <ChapterPageClient params={params} />;
}
