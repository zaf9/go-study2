import {
  getProgress,
  getTopicProgress,
  updateProgress,
  UpdateProgressPayload,
} from "@/services/progressService";
import {
  ChapterProgress,
  TopicProgressSummary,
} from "@/types/learning";

export type SaveProgressRequest = UpdateProgressPayload;

export async function fetchAllProgress(): Promise<TopicProgressSummary[]> {
  const snapshot = await getProgress();
  return snapshot.topics;
}

export async function fetchProgressByTopic(
  topic: string,
): Promise<ChapterProgress[]> {
  const detail = await getTopicProgress(topic);
  return detail.chapters;
}

export async function saveProgress(payload: SaveProgressRequest): Promise<void> {
  await updateProgress(payload);
}
