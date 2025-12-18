const RAW_API_BASE = process.env.NEXT_PUBLIC_API_URL?.trim() ?? "";

// 统一 API 基址，确保附加 /api/v1（或已有版本号则尊重）
function normalizeApiBase(raw: string): string {
  if (!raw) {
    return "/api/v1";
  }
  const cleaned = raw.replace(/\/+$/, "");
  if (/\/api\/v\d+$/.test(cleaned)) {
    return cleaned;
  }
  if (cleaned.endsWith("/api")) {
    return `${cleaned}/v1`;
  }
  return `${cleaned}/api/v1`;
}

export const API_BASE_URL = normalizeApiBase(RAW_API_BASE);

export const ACCESS_TOKEN_KEY = "go-study2.access_token";
export const REMEMBER_ME_KEY = "go-study2.remember_me";
export const REQUEST_TIMEOUT = 10000;

export const API_PATHS = {
  login: "/auth/login",
  register: "/auth/register",
  refresh: "/auth/refresh",
  profile: "/auth/profile",
  changePassword: "/auth/change-password",
  logout: "/auth/logout",
  topics: "/topics",
  topicMenu: (topic: string) => `/topic/${topic}`,
  chapterContent: (topic: string, chapter: string) =>
    `/topic/${topic}/${chapter}`,
  progress: "/progress",
  progressByTopic: (topic: string) => `/progress/${topic}`,
  quiz: (topic: string, chapter: string) => `/quiz/${topic}/${chapter}`,
  quizSubmit: "/quiz/submit",
  quizHistory: "/quiz/history",
  quizHistoryByTopic: (topic: string) => `/quiz/history/${topic}`,
  quizReview: (sessionId: string) => `/quiz/history/${sessionId}`,
  quizStats: (topic: string, chapter: string) => `/quiz/${topic}/${chapter}/stats`,
};
