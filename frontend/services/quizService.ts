import useSWR from "swr";
import api from "@/lib/api";
import { API_PATHS } from "@/lib/constants";
import {
  QuizHistoryItem,
  QuizQuestion,
  QuizSessionPayload,
  QuizSubmitResult,
} from "@/types/quiz";
import { fetchQuizHistory as fetchHistoryLegacy } from "@/lib/quiz";

export interface QuizSubmitPayload {
  sessionId: string;
  topic: string;
  chapter: string;
  answers: Array<{ questionId: number; userAnswers: string[] }>;
  durationMs?: number;
}

export async function fetchQuizSession(
  topic: string,
  chapter: string,
): Promise<QuizSessionPayload> {
  const resp = (await api.get<QuizSessionPayload>(
    API_PATHS.quiz(topic, chapter),
  )) as unknown as QuizSessionPayload;
  const questions: QuizQuestion[] = Array.isArray(resp?.questions)
    ? resp.questions.map((q, optIdx) => ({
        id: q.id,
        type: q.type,
        difficulty: q.difficulty,
        question: q.question,
        options: (q.options ?? []).map((opt, idx) => ({
          id: opt.id ?? String.fromCharCode(65 + idx),
          label: opt.label ?? String(opt.id ?? idx),
        })),
        codeSnippet: q.codeSnippet ?? null,
      }))
    : [];
  return {
    sessionId: resp?.sessionId ?? "",
    topic,
    chapter,
    questions,
  };
}

export async function submitQuiz(
  payload: QuizSubmitPayload,
): Promise<QuizSubmitResult> {
  const resp = (await api.post<QuizSubmitResult>(
    API_PATHS.quizSubmit,
    payload,
  )) as unknown as QuizSubmitResult;
  return resp;
}

export function useQuizSession(topic?: string, chapter?: string) {
  const key = topic && chapter ? ["quiz-session", topic, chapter] : null;
  return useSWR<QuizSessionPayload | null>(
    key,
    () =>
      topic && chapter
        ? fetchQuizSession(topic, chapter)
        : Promise.resolve(null),
    {
      revalidateOnFocus: false,
    },
  );
}

export function useQuizHistory(topic?: string) {
  return useSWR<QuizHistoryItem[]>(
    ["quiz-history", topic || "all"],
    () => fetchHistoryLegacy(topic),
    { revalidateOnFocus: false },
  );
}

