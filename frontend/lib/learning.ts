import api from "./api";
import { API_PATHS } from "./constants";
import { ChapterContent, ChapterSummary, TopicSummary } from "@/types/learning";

export async function fetchTopics(): Promise<TopicSummary[]> {
  const { data } = await api.get<{
    topics: Array<{ id: string; title: string; description?: string }>;
  }>(API_PATHS.topics);
  const topics = data?.topics ?? [];
  return topics.map<TopicSummary>((item) => ({
    key: item.id as TopicSummary["key"],
    title: item.title,
    summary: item.description ?? "",
    chapterCount: 0,
  }));
}

export async function fetchChapters(topic: string): Promise<ChapterSummary[]> {
  const { data } = await api.get<{
    items: Array<{ id: number; title: string; name: string }>;
  }>(API_PATHS.topicMenu(topic));
  const items = data?.items ?? [];
  return items.map<ChapterSummary>((item) => ({
    id: item.name,
    topicKey: topic as ChapterSummary["topicKey"],
    title: item.title,
    order: item.id,
  }));
}

export async function fetchChapterContent(
  topic: string,
  chapter: string,
): Promise<ChapterContent> {
  const { data } = await api.get<{ title: string; content: string }>(
    API_PATHS.chapterContent(topic, chapter),
  );
  return {
    id: chapter,
    topicKey: topic as ChapterContent["topicKey"],
    title: data?.title ?? chapter,
    markdown: data?.content ?? "",
  };
}
