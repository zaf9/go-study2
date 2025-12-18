import api from "./api";
import { API_PATHS } from "./constants";
import {
  QuizHistoryItem,
  QuizItem,
  QuizReviewDetail,
  QuizStats,
  QuizSubmitRequest,
  QuizSubmitResult,
} from "@/types/quiz";

export async function fetchQuizQuestions(
  topic: string,
  chapter: string,
): Promise<QuizItem[]> {
  const resp = (await api.get<{ questions?: QuizItem[] }>(
    API_PATHS.quiz(topic, chapter),
  )) as unknown as { questions?: QuizItem[] } | QuizItem[];
  if (Array.isArray(resp)) {
    return resp;
  }
  return Array.isArray(resp?.questions) ? resp.questions : [];
}

export async function submitQuiz(
  request: QuizSubmitRequest & { durationMs?: number },
): Promise<QuizSubmitResult> {
  const resp = (await api.post<QuizSubmitResult>(
    API_PATHS.quizSubmit,
    request,
  )) as unknown as QuizSubmitResult;
  return resp;
}

export async function fetchQuizHistory(
  topic?: string,
): Promise<QuizHistoryItem[]> {
  const path = topic
    ? API_PATHS.quizHistoryByTopic(topic)
    : API_PATHS.quizHistory;
  const resp = (await api.get<QuizHistoryItem[]>(path)) as unknown as
    | QuizHistoryItem[]
    | { items?: QuizHistoryItem[] };
  if (Array.isArray(resp)) {
    return resp;
  }
  return Array.isArray((resp as { items?: QuizHistoryItem[] })?.items)
    ? (resp as { items?: QuizHistoryItem[] }).items || []
    : [];
}

export async function fetchQuizReview(
  sessionId: string,
): Promise<QuizReviewDetail> {
  const resp = (await api.get<QuizReviewDetail>(
    API_PATHS.quizReview(sessionId),
  )) as unknown as QuizReviewDetail;
  return resp;
}

export async function fetchQuizStats(
  topic: string,
  chapter: string,
): Promise<QuizStats> {
  const resp = (await api.get<QuizStats>(
    API_PATHS.quizStats(topic, chapter),
  )) as unknown as QuizStats;
  return resp;
}
