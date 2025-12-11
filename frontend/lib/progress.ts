import api from "./api";
import { API_PATHS } from "./constants";
import { LearningProgress, ProgressStatus } from "@/types/learning";

export interface SaveProgressRequest {
  topic: string;
  chapter: string;
  status: ProgressStatus;
  position?: string;
}

export async function fetchAllProgress(): Promise<LearningProgress[]> {
  const { data } = await api.get<LearningProgress[]>(API_PATHS.progress);
  return data || [];
}

export async function fetchProgressByTopic(
  topic: string,
): Promise<LearningProgress[]> {
  const { data } = await api.get<LearningProgress[]>(
    API_PATHS.progressByTopic(topic),
  );
  return data || [];
}

export async function saveProgress(
  payload: SaveProgressRequest,
): Promise<void> {
  await api.post(API_PATHS.progress, payload);
}
