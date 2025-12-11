"use client";

import { useCallback } from "react";
import {
  UpdateProgressPayload,
  updateProgress,
  useProgressOverview,
  useTopicProgress,
} from "@/services/progressService";

export default function useProgress(topic?: string) {
  const {
    data: overview,
    error: overviewError,
    isLoading: overviewLoading,
    mutate: mutateOverview,
  } = useProgressOverview();
  const {
    data: topicDetail,
    error: topicError,
    isLoading: topicLoading,
    mutate: mutateTopic,
  } = useTopicProgress(topic);

  const recordProgress = useCallback(
    async (payload: UpdateProgressPayload) => {
      await updateProgress(payload);
      await Promise.all([
        mutateOverview(),
        topic ? mutateTopic() : Promise.resolve(),
      ]);
    },
    [mutateOverview, mutateTopic, topic],
  );

  return {
    overview,
    topicDetail,
    chapters: topicDetail?.chapters ?? [],
    next: overview?.next ?? null,
    error: overviewError || topicError,
    isLoading: overviewLoading || (!!topic && topicLoading),
    recordProgress,
    refresh: () => Promise.all([mutateOverview(), mutateTopic()]),
  };
}
