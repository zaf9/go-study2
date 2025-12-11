import api from "./api";
import { API_PATHS } from "./constants";
import {
  QuizHistoryItem,
  QuizItem,
  QuizResult,
  QuizSubmitRequest,
} from "@/types/quiz";

export async function fetchQuizQuestions(
  topic: string,
  chapter: string,
): Promise<QuizItem[]> {
  const { data } = await api.get<QuizItem[]>(API_PATHS.quiz(topic, chapter));
  return data || [];
}

export async function submitQuiz(
  request: QuizSubmitRequest & { durationMs?: number },
): Promise<QuizResult> {
  const { data } = await api.post<QuizResult>(API_PATHS.quizSubmit, request);
  return data as QuizResult;
}

export async function fetchQuizHistory(
  topic?: string,
): Promise<QuizHistoryItem[]> {
  const path = topic
    ? API_PATHS.quizHistoryByTopic(topic)
    : API_PATHS.quizHistory;
  const { data } = await api.get<QuizHistoryItem[]>(path);
  return data || [];
}
