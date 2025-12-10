'use client';

import useSWR from "swr";
import { useCallback } from "react";
import { fetchAllProgress, fetchProgressByTopic, saveProgress, SaveProgressRequest } from "@/lib/progress";
import { LearningProgress } from "@/types/learning";

export default function useProgress(topic?: string) {
  const key = topic ? ["progress", topic] : ["progress", "all"];
  const { data, error, isLoading, mutate } = useSWR<LearningProgress[]>(key, () =>
    topic ? fetchProgressByTopic(topic) : fetchAllProgress()
  );

  const recordProgress = useCallback(
    async (payload: SaveProgressRequest) => {
      await saveProgress(payload);
      await mutate();
    },
    [mutate]
  );

  const latest = (data ?? []).reduce<LearningProgress | null>((acc, item) => {
    if (!acc) return item;
    return new Date(item.lastVisit) > new Date(acc.lastVisit) ? item : acc;
  }, null);

  return {
    progress: data ?? [],
    latest,
    error,
    isLoading,
    recordProgress,
  };
}

