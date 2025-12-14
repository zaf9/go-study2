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

export type ProgressStatus =
  | "not_started"
  | "in_progress"
  | "completed"
  | "tested";

export interface ChapterProgress {
  topic: TopicKey;
  chapter: string;
  status: ProgressStatus;
  // percent: 数值形式的进度（0-100），由后端返回优先使用
  percent?: number;
  scrollProgress?: number;
  readDuration?: number;
  lastPosition?: string;
  lastVisitAt?: string;
  quizScore?: number;
  quizPassed?: boolean;
}

export interface TopicProgressSummary {
  id: TopicKey;
  name: string;
  weight: number;
  progress: number;
  completedChapters: number;
  totalChapters: number;
  lastVisitAt?: string;
}

export interface TopicProgressDetail extends TopicProgressSummary {
  chapters: ChapterProgress[];
}

export interface OverallProgress {
  progress: number;
  completedChapters: number;
  totalChapters: number;
  studyDays: number;
  totalStudyTime: number;
}

export interface NextChapterHint {
  topic: TopicKey;
  chapter: string;
  status: ProgressStatus;
  progress: number;
}

export interface ProgressSnapshot {
  overall: OverallProgress;
  topics: TopicProgressSummary[];
  next?: NextChapterHint | null;
}
