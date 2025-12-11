export type TopicKey = "lexical_elements" | "constants" | "variables" | "types";

export interface TopicSummary {
  key: TopicKey;
  title: string;
  summary: string;
  chapterCount: number;
  order?: number;
}

export interface ChapterSummary {
  id: string;
  topicKey: TopicKey;
  title: string;
  summary?: string;
  order?: number;
}

export interface ChapterContent {
  id: string;
  topicKey: TopicKey;
  title: string;
  markdown: string;
  lastVisit?: string;
  lastPosition?: string;
}

export type ProgressStatus = "not_started" | "in_progress" | "done";

export interface LearningProgress {
  topic: TopicKey;
  chapter: string;
  status: ProgressStatus;
  lastVisit: string;
  lastPosition?: string;
}
