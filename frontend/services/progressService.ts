import useSWR, { mutate as mutateCache } from "swr";
import api from "@/lib/api";
import { API_BASE_URL, API_PATHS } from "@/lib/constants";
import {
  ChapterProgress,
  NextChapterHint,
  ProgressSnapshot,
  TopicProgressDetail,
} from "@/types/learning";

export interface UpdateProgressPayload {
  topic: string;
  chapter: string;
  readDuration?: number;
  scrollProgress?: number;
  lastPosition?: string;
  quizScore?: number;
  quizPassed?: boolean;
  estimatedSeconds?: number;
  forceSync?: boolean;
}

export const progressKeys = {
  overview: "progress/overview",
  topic: (topic: string) => ["progress/topic", topic],
};

const MAX_RETRY = 5;
const BASE_DELAY_MS = 400;
let latestProgress: UpdateProgressPayload | null = null;
let beforeUnloadRegistered = false;

function jitterDelay(attempt: number): number {
  const base = BASE_DELAY_MS * Math.pow(2, attempt);
  const jitter = Math.floor(Math.random() * 150);
  return base + jitter;
}

function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => {
    setTimeout(resolve, ms);
  });
}

async function withRetry<T>(task: () => Promise<T>): Promise<T> {
  let lastError: unknown;
  for (let i = 0; i < MAX_RETRY; i += 1) {
    try {
      return await task();
    } catch (err) {
      lastError = err;
      if (i === MAX_RETRY - 1) {
        break;
      }
      await sleep(jitterDelay(i));
    }
  }
  throw lastError;
}

function normalizePayload(
  payload: UpdateProgressPayload,
): Record<string, unknown> {
  return {
    topic: payload.topic,
    chapter: payload.chapter,
    read_duration: payload.readDuration ?? 0,
    scroll_progress: payload.scrollProgress ?? 0,
    last_position: payload.lastPosition ?? "",
    quiz_score: payload.quizScore ?? 0,
    quiz_passed: payload.quizPassed ?? false,
    estimated_seconds: payload.estimatedSeconds ?? 0,
    force_sync: payload.forceSync ?? false,
  };
}

async function flushPendingProgress(): Promise<void> {
  if (!latestProgress || typeof window === "undefined") {
    return;
  }
  const body = JSON.stringify(normalizePayload({ ...latestProgress, forceSync: true }));
  try {
    if (navigator.sendBeacon) {
      const url = `${API_BASE_URL}${API_PATHS.progress}`;
      const blob = new Blob([body], { type: "application/json" });
      navigator.sendBeacon(url, blob);
      return;
    }
    await api.post(API_PATHS.progress, JSON.parse(body));
  } catch {
    // beforeunload 兜底忽略错误
  }
}

function ensureBeforeUnload() {
  if (beforeUnloadRegistered || typeof window === "undefined") {
    return;
  }
  beforeUnloadRegistered = true;
  window.addEventListener("beforeunload", () => {
    void flushPendingProgress();
  });
}

export async function updateProgress(payload: UpdateProgressPayload) {
  ensureBeforeUnload();
  latestProgress = { ...payload };
  const response = await withRetry(() =>
    api.post(API_PATHS.progress, normalizePayload(payload)),
  );
  void mutateCache(progressKeys.overview);
  void mutateCache(progressKeys.topic(payload.topic));
  return response as unknown as {
    status: string;
    overall?: ProgressSnapshot["overall"];
    topic?: TopicProgressDetail;
  };
}

export async function getProgress(): Promise<ProgressSnapshot> {
  const data = (await api.get<ProgressSnapshot>(
    API_PATHS.progress,
  )) as unknown as ProgressSnapshot;
  return {
    overall: data?.overall ?? {
      progress: 0,
      completedChapters: 0,
      totalChapters: 0,
      studyDays: 0,
      totalStudyTime: 0,
    },
    topics: data?.topics ?? [],
    next: data?.next ?? null,
  };
}

export async function getTopicProgress(
  topic: string,
): Promise<TopicProgressDetail> {
  const data = (await api.get<{
    topic: TopicProgressDetail;
    chapters?: ChapterProgress[];
  }>(API_PATHS.progressByTopic(topic))) as unknown as {
    topic?: TopicProgressDetail;
    chapters?: ChapterProgress[];
  };
  const summary = data?.topic ?? {
    id: topic as TopicProgressDetail["id"],
    name: topic,
    weight: 0,
    progress: 0,
    completedChapters: 0,
    totalChapters: 0,
  };
  return {
    ...summary,
    id: summary.id as TopicProgressDetail["id"],
    chapters: data?.chapters ?? [],
  };
}

export function useProgressOverview() {
  return useSWR<ProgressSnapshot>(progressKeys.overview, getProgress, {
    revalidateOnFocus: false,
  });
}

export function useTopicProgress(topic?: string) {
  const key = topic ? progressKeys.topic(topic) : null;
  return useSWR<TopicProgressDetail | null>(
    key,
    () => (topic ? getTopicProgress(topic) : Promise.resolve(null)),
    { revalidateOnFocus: false },
  );
}

export const __test__ = {
  jitterDelay,
  sleep,
  withRetry,
};

